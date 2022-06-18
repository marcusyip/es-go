package main

import (
	"fmt"
	"os"

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
			// config := es.NewConfig()
			// db := database.Connect()
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
