package main

import (
	"fmt"
	"os"

	"github.com/jasonli0226/ssh-connection-manager/pkg/app"
	"github.com/jasonli0226/ssh-connection-manager/pkg/cli"
	"github.com/jasonli0226/ssh-connection-manager/pkg/infra/storage"
	"github.com/jasonli0226/ssh-connection-manager/pkg/logging"
)

func main() {
	logging.Init()

	repo := storage.NewFileConnectionRepository()
	sshManager := app.NewSSHManager(repo)
	rootCmd := cli.NewRootCommand(sshManager)

	if err := rootCmd.Execute(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to execute command")
		fmt.Println(err)
		os.Exit(1)
	}
}
