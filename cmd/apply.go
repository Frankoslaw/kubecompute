package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"kubecompute/internal/adapter/provision"
	"kubecompute/internal/adapter/repository/sqlc"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"kubecompute/internal/core/service"
	"kubecompute/internal/core/util"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	_ "modernc.org/sqlite"
)

var (
	filePath string
	delete   bool = false
)

var applyCmd = &cobra.Command{
	Use: "apply",
	RunE: func(cmd *cobra.Command, args []string) error {
		// temporary serve bootstrap
		conn, err := sql.Open("sqlite", "./kubecompute.db?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000")
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		conn.SetMaxIdleConns(1)
		conn.SetMaxOpenConns(1)

		repository := sqlc.NewSqlcNodeRepository(conn)
		dockerProvider, err := provision.NewDockerNodeProvider()
		if err != nil {
			return err
		}
		reconciler := service.NewNodeReconciler(repository, dockerProvider)

		workQueue := util.NewWorkQueue[port.ReconcileRequest](64)
		controller := service.NewNodeController(repository, reconciler, workQueue)
		controller.Start(cmd.Context())

		nodeService := service.NewNodeService(repository, controller)

		// true apply
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

			var obj map[string]interface{}
			if err := yaml.Unmarshal([]byte(section), &obj); err != nil {
				slog.Warn("failed to unmarshal yaml", "error", err)
				continue
			}

			apiVersion, ok1 := obj["apiVersion"].(string)
			kind, ok2 := obj["kind"].(string)
			if !ok1 || !ok2 {
				continue
			}

			if apiVersion == "kubecompute.io/v1alpha1" && kind == "Node" {
				var node domain.Node
				if err := yaml.Unmarshal([]byte(section), &node); err != nil {
					slog.Warn("failed to unmarshal node", "error", err)
					continue
				}

				var err error
				if !delete {
					err = nodeService.CreateNode(cmd.Context(), &node)
				} else {
					err = nodeService.DeleteNode(cmd.Context(), node.ObjectMeta.NamespacedName())
				}

				if err != nil {
					slog.Warn("failed to create node", "error", err)
					continue
				}
			}
		}

		time.Sleep(5 * time.Second) // wait for reconciliation

		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the configuration file")
	applyCmd.MarkFlagRequired("file")
	applyCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete the resources in the configuration file")

	rootCmd.AddCommand(applyCmd)
}

func setupDatabase(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(ctx); err != nil {
		return nil, err
	}

	return dbpool, nil
}
