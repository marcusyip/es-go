package main

import (
	"context"
	"fmt"
	"os"

	"github.com/es-go/es-go/es"
	"github.com/es-go/es-go/es/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd    *cobra.Command
	ConfigFile string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			config := es.NewConfig()
			db := database.Connect()

			sql := fmt.Sprintf(
				"INSERT INTO %s (aggregate_id, version, event_type, payload, created_at) VALUES ($1, $2, $3, $4, $5)",
				config.TableName)

			_, err := db.Exec(context.Background(), sql, "test-id", 1, "created_event", []byte("{\"amount\":100,\"currency\":\"HKD\"}"), "NOW()")
			if err != nil {
				panic(err)
			}

			// conf := config.ProvideConfig()
			// app, err := app.BuildApp(conf)
			// if err != nil {
			// 	panic(err)
			// }
			// app.Start()
		},
	}
	cmd.AddCommand(
	// migrate.NewMigrateCommand(),
	// alerter.NewAlerterCommand(),
	)
	return cmd
}

func init() {
	cobra.OnInitialize(onInitialize)
	rootCmd = newRootCommand()
	// rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file")
}

func onInitialize() {
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}
}

func main() {
	Execute()
}
