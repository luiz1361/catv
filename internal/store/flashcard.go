package store

// Flashcard represents a flashcard
type Flashcard struct {
	ID        int
	File      string
	Question  string
	Answer    string
	RevisitIn int // number of days until next review (<=0 means due)
}
