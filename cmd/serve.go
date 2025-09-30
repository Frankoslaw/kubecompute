package cmd

import (
	"database/sql"
	"kubecompute/internal/adapter/handler"
	"kubecompute/internal/adapter/provision"
	provision_terraform "kubecompute/internal/adapter/provision/terraform"
	"kubecompute/internal/adapter/repository/sqlc"
	"kubecompute/internal/core/port"
	"kubecompute/internal/core/service"
	"kubecompute/internal/core/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	_ "modernc.org/sqlite"

	_ "kubecompute/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Connect to DB
		conn, err := connectDB("./kubecompute.db?_journal_mode=WAL&_synchronous=NORMAL&_busy_timeout=5000")
		if err != nil {
			return err
		}
		defer conn.Close()

		// Setup repository and services
		repository := sqlc.NewSqlcNodeRepository(conn)

		dockerProvider, err := provision.NewDockerNodeProvider()
		if err != nil {
			return err
		}
		terraformVirtualboxProvider := provision_terraform.NewVirtualboxNodeProvider()

		reconciler := service.NewNodeReconciler(repository)
		reconciler.RegisterProvider("docker", dockerProvider)
		reconciler.RegisterProvider("terraform-virtualbox", terraformVirtualboxProvider)
		reconciler.RegisterProvider("virtualbox", terraformVirtualboxProvider) // alias, as virtualbox native provider is currently not implemented

		workQueue := util.NewWorkQueue[port.ReconcileRequest](32)
		controller := service.NewNodeController(repository, reconciler, workQueue)
		controller.Start(cmd.Context())

		nodeService := service.NewNodeService(repository, controller)

		// Setup Gin router and handlers
		router := gin.Default()
		nodeHandler := handler.NewNodeHandler(nodeService)
		nodeHandler.RegisterRoutes(router)

		// Swagger endpoint
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// Start server
		return router.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

}

func connectDB(sqliteDSN string) (*sql.DB, error) {
	conn, err := sql.Open("sqlite", sqliteDSN)
	if err != nil {
		return nil, err
	}
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)

	return conn, nil
}
