/*
Yet another powerful customizable ActivityPub relay server written in Go.

Run Activity-Relay

API Server
	./Activity-Relay -c <Path of config file> server
Job Worker
	./Activity-Relay -c <Path of config file> worker
CLI Management Utility
	./Activity-Relay -c <Path of config file> control

Config

YAML Format
	ACTOR_PEM: actor.pem
	REDIS_URL: redis://localhost:6379

	RELAY_BIND: 0.0.0.0:8080
	RELAY_DOMAIN: relay.toot.yukimochi.jp
	RELAY_SERVICENAME: YUKIMOCHI Toot Relay Service
	RELAY_SUMMARY: |
		YUKIMOCHI Toot Relay Service is Running by Activity-Relay
	RELAY_ICON: https://example.com/example_icon.png
	RELAY_IMAGE: https://example.com/example_image.png
Environment Variable

This is Optional : When config file not exist, use environment variables.
	- ACTOR_PEM
	- REDIS_URL
	- RELAY_BIND
	- RELAY_DOMAIN
	- RELAY_SERVICENAME
	- RELAY_SUMMARY
	- RELAY_ICON
	- RELAY_IMAGE

*/
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
			fmt.Println(globalConfig.DumpWelcomeMessage("API Server"))
			err := api.Entrypoint(globalConfig)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			return nil
		},
	}

	var worker = &cobra.Command{
		Use:   "worker",
		Short: "Activity-Relay Job Worker",
		Long:  "Activity-Relay Job Worker is providing ActivityPub Activity deliverer",
		RunE: func(cmd *cobra.Command, args []string) error {
			initConfig(cmd, args)
			fmt.Println(globalConfig.DumpWelcomeMessage("Job Worker"))
			err := deliver.Entrypoint(globalConfig)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			return nil
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
