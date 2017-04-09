package main

import (
	"github.com/dgoodwin/gophoto/client"
	"github.com/dgoodwin/gophoto/server"

	"github.com/spf13/cobra"
)

func main() {

	var cmdServe = &cobra.Command{
		Use:   "serve [path to config file]",
		Short: "Launch the GoPhoto server",
		Run:   server.RunServer,
	}

	var cmdSync = &cobra.Command{
		Use:   "sync [path to photos directory]",
		Short: "Upload all new files in given directory",
		Run:   client.RunSync,
	}

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdServe)
	rootCmd.AddCommand(cmdSync)
	rootCmd.Execute()
}
