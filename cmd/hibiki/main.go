package main

import (
	"github.com/rl404/hibiki/internal/utils"
	"github.com/spf13/cobra"
)

// @title Hibiki API
// @description Hibiki API.
// @BasePath /
// @schemes http https
func main() {
	cmd := cobra.Command{
		Use:   "hibiki",
		Short: "Hibiki API",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Run API server",
		RunE: func(*cobra.Command, []string) error {
			return server()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "consumer",
		Short: "Run message consumer",
		RunE: func(*cobra.Command, []string) error {
			return consumer()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
		RunE: func(*cobra.Command, []string) error {
			return migrate()
		},
	})

	cronCmd := cobra.Command{
		Use:   "cron",
		Short: "Cron",
	}

	cronCmd.AddCommand(&cobra.Command{
		Use:   "update",
		Short: "Update old data",
		RunE: func(*cobra.Command, []string) error {
			return cronUpdate()
		},
	})

	cronCmd.AddCommand(&cobra.Command{
		Use:   "fill",
		Short: "Fill missing data",
		RunE: func(*cobra.Command, []string) error {
			return cronFill()
		},
	})

	cmd.AddCommand(&cronCmd)

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
