package cmd

import (
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

var deleteCmd = &cobra.Command{
	Use: "delete",
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
				// Node doesn't exist → NOOP
				slog.Info("node not found, skipping delete", "name", node.ObjectMeta.Name, "namespace", node.ObjectMeta.Namespace)
			case http.StatusOK:
				// Node exists → DELETE
				slog.Info("node exists, deleting", "name", node.ObjectMeta.Name, "namespace", node.ObjectMeta.Namespace)
				if err := sendJSON("DELETE", url, node); err != nil {
					slog.Warn("failed to delete node", "error", err)
				}
			default:
				slog.Warn("unexpected GET response", "status", resp.StatusCode, "body", string(body))
			}
		}

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringP("file", "f", "", "Path to the configuration file")
	deleteCmd.MarkFlagRequired("file")

	rootCmd.AddCommand(deleteCmd)

}
