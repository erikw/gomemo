package commands

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	debugFlag bool
	logger    *slog.Logger
)

var rootCmd = &cobra.Command{
	Use:   "gomemo",
	Short: "Gomemo - a personal note taking application",
	Long: `Gomemo is a personal note taking application with a web interface.
Use subcommands to interact with the application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogger(debugFlag)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show help when no subcommand is provided
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func GetLogger() *slog.Logger {
	return logger
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Enable debug logging")
	rootCmd.AddCommand(seedCmd)
	rootCmd.AddCommand(serveCmd)
	// Remove completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initLogger(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}

	logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}),
	)

	if debug {
		logger.Debug("Debug logging enabled.")
	}
}
