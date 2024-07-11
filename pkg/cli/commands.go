package cli

import (
	"fmt"
	"strconv"

	"github.com/jasonli0226/ssh-connection-manager/pkg/app"
	"github.com/jasonli0226/ssh-connection-manager/pkg/logging"
	"github.com/spf13/cobra"
)

func NewRootCommand(sshManager *app.SSHManager) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ssh-manager",
		Short: "SSH Connection Manager CLI",
		Long:  `A CLI tool to manage SSH connections and connect to remote servers using aliases.`,
	}

	rootCmd.AddCommand(newAddCommand(sshManager))
	rootCmd.AddCommand(newListCommand(sshManager))
	rootCmd.AddCommand(newConnectCommand(sshManager))
	rootCmd.AddCommand(newDeleteCommand(sshManager))

	return rootCmd
}

func newAddCommand(sshManager *app.SSHManager) *cobra.Command {
	return &cobra.Command{
		Use:   "add <alias> <host> <user> <port>",
		Short: "Add a new SSH connection",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			alias := args[0]
			host := args[1]
			user := args[2]
			port, err := strconv.Atoi(args[3])
			if err != nil {
				logging.Log.Error().Err(err).Msg("Invalid port number")
				fmt.Println("Invalid port number")
				return
			}

			fmt.Print("Enter password: ")
			var password string
			fmt.Scanln(&password)

			err = sshManager.AddConnection(alias, host, user, password, port)
			if err != nil {
				logging.Log.Error().Err(err).Msg("Failed to add connection")
				fmt.Printf("Error adding connection: %v\n", err)
			} else {
				logging.Log.Info().Str("alias", alias).Msg("Connection added successfully")
				fmt.Println("Connection added successfully")
			}
		},
	}
}

func newListCommand(sshManager *app.SSHManager) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all saved SSH connections",
		Run: func(cmd *cobra.Command, args []string) {
			connections, err := sshManager.ListConnections()
			if err != nil {
				logging.Log.Error().Err(err).Msg("Failed to list connections")
				fmt.Printf("Error listing connections: %v\n", err)
				return
			}

			if len(connections) == 0 {
				fmt.Println("No connections found")
				return
			}

			fmt.Println("Saved connections:")
			for _, conn := range connections {
				fmt.Printf("- %s: %s@%s:%d\n", conn.Alias, conn.User, conn.Host, conn.Port)
			}
			logging.Log.Info().Int("count", len(connections)).Msg("Listed connections")
		},
	}
}

func newConnectCommand(sshManager *app.SSHManager) *cobra.Command {
	return &cobra.Command{
		Use:   "connect <alias>",
		Short: "Connect to a saved SSH server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			alias := args[0]
			logging.Log.Info().Str("alias", alias).Msg("Attempting to connect")
			err := sshManager.Connect(alias)
			if err != nil {
				logging.Log.Error().Err(err).Str("alias", alias).Msg("Failed to connect")
				fmt.Printf("Error connecting to %s: %v\n", alias, err)
			} else {
				logging.Log.Info().Str("alias", alias).Msg("Connection closed")
			}
		},
	}
}

func newDeleteCommand(sshManager *app.SSHManager) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <alias>",
		Short: "Delete a saved SSH connection",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			alias := args[0]
			err := sshManager.DeleteConnection(alias)
			if err != nil {
				logging.Log.Error().Err(err).Str("alias", alias).Msg("Failed to delete connection")
				fmt.Printf("Error deleting connection %s: %v\n", alias, err)
			} else {
				logging.Log.Info().Str("alias", alias).Msg("Connection deleted successfully")
				fmt.Printf("Connection %s deleted successfully\n", alias)
			}
		},
	}
}
