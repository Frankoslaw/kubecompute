package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kubecompute/internal/core/domain"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use: "apply",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := cmd.Flag("file").Value.String()

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		sections := strings.Split(string(content), "---")
		for _, section := range sections {
			section = strings.TrimSpace(section)
			if section == "" {
				continue
			}

			var node domain.Node
			if err := yaml.Unmarshal([]byte(section), &node); err != nil {
				slog.Warn("failed to unmarshal node", "error", err)
				continue
			}

			url := fmt.Sprintf("http://localhost:8080/ns/%s/nodes/%s", node.ObjectMeta.Namespace, node.ObjectMeta.Name)

			resp, err := http.Get(url)
			if err != nil {
				slog.Warn("failed to GET node", "error", err)
				continue
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)

			switch resp.StatusCode {
			case http.StatusNotFound:
				// Node doesn't exist → POST
				slog.Info("node not found, creating", "name", node.ObjectMeta.Name, "namespace", node.ObjectMeta.Namespace)
				if err := sendJSON("POST", fmt.Sprintf("http://localhost:8080/ns/%s/nodes", node.ObjectMeta.Namespace), node); err != nil {
					slog.Warn("failed to create node", "error", err)
				}
			case http.StatusOK:
				// Node exists → PUT (update)
				slog.Info("node exists, updating", "name", node.ObjectMeta.Name, "namespace", node.ObjectMeta.Namespace)
				if err := sendJSON("PUT", url, node); err != nil {
					slog.Warn("failed to update node", "error", err)
				}
			default:
				slog.Warn("unexpected GET response", "status", resp.StatusCode, "body", string(body))
			}
		}

		return nil
	},
}

func init() {
	applyCmd.Flags().StringP("file", "f", "", "Path to the configuration file")
	applyCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(applyCmd)
}

func sendJSON(method, url string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s %s failed: %s (%d)", method, url, string(body), resp.StatusCode)
	}
	return nil
}
