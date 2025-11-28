package domain

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         int       `json:"id"`
	QuestionID int       `json:"question_id"`
	UserID     uuid.UUID `json:"user_id"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"created_at"`
}