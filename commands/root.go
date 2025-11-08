package commands

import (
	"fmt"
	"os"

	"catv/store"
	"catv/tui"

	"github.com/spf13/cobra"
)

var Store *store.Store
var Model string

var RootCmd = &cobra.Command{
	Use:   "catv",
	Short: "Ollama-powered spaced repetition flashcards CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		homeDir, err := os.UserHomeDir()
		if err != nil {
			tui.PrintError("Could not determine home directory:", err)
			os.Exit(1)
		}
		catvDir := homeDir + string(os.PathSeparator) + ".catv"
		if err := os.MkdirAll(catvDir, 0700); err != nil {
			tui.PrintError("Could not create .catv directory:", err)
			os.Exit(1)
		}
		dbPath := catvDir + string(os.PathSeparator) + "flashcards.db"
		Store, err = store.NewStore(dbPath)
		if err != nil {
			tui.PrintError("DB error:", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Default to review command
		ReviewCmd.Run(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if Store != nil {
			Store.Close()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&Model, "model", "llama3.1", "Ollama model to use for flashcard generation")
	RootCmd.AddCommand(GenerateCmd)
	RootCmd.AddCommand(ReviewCmd)
}
