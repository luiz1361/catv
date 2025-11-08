package store

import (
	"time"
)

// Flashcard represents a flashcard
type Flashcard struct {
	ID         int
	File       string
	Question   string
	Answer     string
	NextReview time.Time
}

// NextReviewDate calculates the next review date based on the revisit days
func NextReviewDate(days int) time.Time {
	return time.Now().AddDate(0, 0, days)
}
