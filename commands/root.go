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
		dbName := "flashcards.db"
		var err error
		Store, err = store.NewStore(dbName)
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
