package commands

import (
	"catv/internal/store"
	"catv/internal/tui"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review flashcards that are due",
	Run: func(cmd *cobra.Command, args []string) {
		flashcards, err := Store.GetFlashcardsForReview()
		if err != nil {
			tui.PrintError("DB query error:", err)
			return
		}

		if len(flashcards) == 0 {
			tui.PrintInfo("No flashcards due for review.")
			return
		}

		// Run Bubble Tea TUI for review
		model := tui.NewReviewModel(flashcards)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running review TUI:", err)
		}

		// After review, update DB only for flashcards that were actually answered
		for i, fc := range flashcards {
			// Only update if user interacted with the card (answered correct/incorrect or provided difficulty)
			if model.FlashcardWasCorrect(i) && model.FlashcardDifficulty(i) > 0 {
				fc.NextReview = store.NextReviewDate(model.FlashcardDifficulty(i))
				err := Store.UpdateFlashcard(fc)
				if err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Updated flashcard %d: next_review=%s", fc.ID, fc.NextReview.Format("2006-01-02")))
				}
			} else if !model.FlashcardWasCorrect(i) && model.FlashcardDifficulty(i) > 0 {
				// If explicitly marked incorrect, set next_review to tomorrow
				fc.NextReview = store.NextReviewDate(1)
				err := Store.UpdateFlashcard(fc)
				if err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Marked flashcard %d incorrect: next_review set to tomorrow", fc.ID))
				}
			}
			// If neither correct nor incorrect was set, do not update next_review
		}
	},
}

// ...existing code...
