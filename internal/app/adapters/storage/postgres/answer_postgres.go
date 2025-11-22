package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	DB *gorm.DB
}

func (r *AnswerRepository) SaveAnswer(ctx context.Context, question_id int, user_id uuid.UUID, text string) (*domain.Answer, error) {
	as := &Answers{
		QuestionID: question_id,
		UserID:     user_id,
		AnswerText: text,
	}
	res := r.DB.WithContext(ctx).Create(as)
	
	if res.Error != nil {
		return nil, fmt.Errorf("save_answer error: %w", res.Error)
	}

	return &domain.Answer{
		ID: as.ID,
		QuestionID: as.QuestionID,
		UserID: as.UserID,
		Text: as.AnswerText,
		CreatedAt: as.CreatedAt,
	}, nil
}

func (r *AnswerRepository) GetAnswer(ctx context.Context, answer_id int) (*domain.Answer, error) {
	as := &Answers{}

	if err := r.DB.WithContext(ctx).First(as, answer_id).Error; err != nil {
		return nil, fmt.Errorf("get_answer error: %w", err)
	}

	return &domain.Answer{
			ID:         as.ID,
			QuestionID: as.QuestionID,
			UserID:     as.UserID,
			Text:       as.AnswerText,
			CreatedAt:  as.CreatedAt},
		nil
}

func (r *AnswerRepository) DeleteAnswer(ctx context.Context, answer_id int) error {
	res := r.DB.WithContext(ctx).Delete(&Answers{}, answer_id)

	if res.Error != nil {
		return fmt.Errorf("delete_answer error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.ErrNoRowDeleted
	}

	return nil
}
