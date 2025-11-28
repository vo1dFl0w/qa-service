package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	DB *gorm.DB
}

func (r *AnswerRepository) SaveAnswer(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error) {
	as := &Answers{
		QuestionID: questionID,
		UserID:     userID,
		AnswerText: text,
	}
	res := r.DB.WithContext(ctx).Create(as)
	
	if res.Error != nil {
		return nil, fmt.Errorf("save answer error: %w", res.Error)
	}

	return &domain.Answer{
		ID: as.ID,
		QuestionID: as.QuestionID,
		UserID: as.UserID,
		Text: as.AnswerText,
		CreatedAt: as.CreatedAt,
	}, nil
}

func (r *AnswerRepository) GetAnswer(ctx context.Context, answerID int) (*domain.Answer, error) {
	as := &Answers{}

	if err := r.DB.WithContext(ctx).First(as, answerID).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get answer: %w", err)
	}

	return &domain.Answer{
			ID:         as.ID,
			QuestionID: as.QuestionID,
			UserID:     as.UserID,
			Text:       as.AnswerText,
			CreatedAt:  as.CreatedAt},
		nil
}

func (r *AnswerRepository) DeleteAnswer(ctx context.Context, answerID int) error {
	res := r.DB.WithContext(ctx).Delete(&Answers{}, answerID)

	if res.Error != nil {
		return fmt.Errorf("delete answer: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return repository.ErrNoRowDeleted
	}

	return nil
}
