package main

import (
	"github.com/dgoodwin/gophoto/server"

	"github.com/spf13/cobra"
)

func main() {

	var cmdServe = &cobra.Command{
		Use:   "serve [path to config file]",
		Short: "Launch the GoPhoto server",
		Run:   server.RunServer,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdServe)
	rootCmd.Execute()
}
