package store

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Store manages the database connection and queries
type Store struct {
	DB *sql.DB
}

// NewStore creates a new Store
func NewStore(dbName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	createTable := `CREATE TABLE IF NOT EXISTS flashcards (
			  id INTEGER PRIMARY KEY AUTOINCREMENT,
			  file TEXT,
			  question TEXT,
			  answer TEXT,
			  next_review DATE
		  );`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

// GetFlashcardsForReview returns all flashcards that are due for review
func (s *Store) GetFlashcardsForReview() ([]Flashcard, error) {
	rows, err := s.DB.Query("SELECT id, question, answer FROM flashcards WHERE next_review <= date('now') ORDER BY next_review ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []Flashcard
	for rows.Next() {
		var fc Flashcard
		err := rows.Scan(&fc.ID, &fc.Question, &fc.Answer)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, fc)
	}
	return flashcards, nil
}

// UpdateFlashcard updates a flashcard's complexity and next review date
func (s *Store) UpdateFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("UPDATE flashcards SET next_review=? WHERE id=?", fc.NextReview.Format("2006-01-02"), fc.ID)
	return err
}

// IsFileProcessed checks if a file has already been processed
func (s *Store) IsFileProcessed(filePath string) (bool, error) {
	var count int
	row := s.DB.QueryRow("SELECT COUNT(*) FROM flashcards WHERE file = ?", filePath)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertFlashcard inserts a new flashcard into the database
func (s *Store) InsertFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("INSERT INTO flashcards (file, question, answer, next_review) VALUES (?, ?, ?, ?)", fc.File, fc.Question, fc.Answer, fc.NextReview.Format("2006-01-02"))
	return err
}

// Close closes the database connection
func (s *Store) Close() {
	s.DB.Close()
}
