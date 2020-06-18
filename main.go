package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)
import "Activity-Relay/api"
import "Activity-Relay/command"
import "Activity-Relay/deliver"
import "Activity-Relay/models"

var (
	version string

	globalConfig *models.RelayConfig
)

func main() {
	var app = buildCommand()
	app.PersistentFlags().StringP("config", "c", "config.yml", "Path of config file.")

	app.Execute()
}

func buildCommand() *cobra.Command {
	var server = &cobra.Command{
		Use:   "server",
		Short: "Activity-Relay API Server",
		Long:  "Activity-Relay API Server is providing WebFinger API, ActivityPub inbox",
		RunE: func(cmd *cobra.Command, args []string) error {
			initConfig(cmd, args)
			return api.Entrypoint(globalConfig)
		},
	}

	var worker = &cobra.Command{
		Use:   "worker",
		Short: "Activity-Relay Job Worker",
		Long:  "Activity-Relay Job Worker is providing ActivityPub Activity deliverer",
		RunE: func(cmd *cobra.Command, args []string) error {
			initConfig(cmd, args)
			return deliver.Entrypoint(globalConfig)
		},
	}

	var control = &cobra.Command{
		Use:   "control",
		Short: "Activity-Relay CLI",
		Long:  "Activity-Relay CLI Management Utility",
	}
	command.Registor(control)

	var app = &cobra.Command{
		Short: "YUKIMOCHI Activity-Relay",
		Long:  "YUKIMOCHI Activity-Relay - ActivityPub Relay Server",
	}
	app.AddCommand(server)
	app.AddCommand(worker)
	app.AddCommand(control)

	return app
}

func initConfig(cmd *cobra.Command, args []string) {
	configPath := cmd.Flag("config").Value.String()
	file, err := os.Open(configPath)
	defer file.Close()

	if err == nil {
		viper.SetConfigType("yaml")
		viper.ReadConfig(file)
	} else {
		fmt.Println("Config file not exist. Use environment variables.")

		viper.BindEnv("ACTOR_PEM")
		viper.BindEnv("REDIS_URL")
		viper.BindEnv("RELAY_BIND")
		viper.BindEnv("RELAY_DOMAIN")
		viper.BindEnv("RELAY_SERVICENAME")
		viper.BindEnv("RELAY_SUMMARY")
		viper.BindEnv("RELAY_ICON")
		viper.BindEnv("RELAY_IMAGE")
	}

	globalConfig, err = models.NewRelayConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
