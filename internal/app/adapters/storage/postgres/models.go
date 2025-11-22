package postgres

import (
	"time"

	"github.com/google/uuid"
)

type Questions struct {
	ID           int           `gorm:"primaryKey;autoIncrement"`
	QuestionText string        `gorm:"type:text;not null"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	Answers      []Answers `gorm:"constraint:OnDelete:CASCADE;foreignKey:QuestionID"`
}

type Answers struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	QuestionID int       `gorm:"not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	AnswerText string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
