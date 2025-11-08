package commands

import (
	"fmt"
	"os"

	"catv/config"
	"catv/store"
	"catv/tui"

	"github.com/spf13/cobra"
)

var Cfg *config.Config
var Store *store.Store

var RootCmd = &cobra.Command{
	Use:   "catv",
	Short: "Ollama-powered spaced repetition flashcards CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		Cfg, err = config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		Store, err = store.NewStore(Cfg.Database.Name)
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
	RootCmd.AddCommand(GenerateCmd)
	RootCmd.AddCommand(ReviewCmd)
}
