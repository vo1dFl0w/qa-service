package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	DB *gorm.DB
}

func (r *QuestionRepository) FindQuestionByID(ctx context.Context, questionID int) (*domain.Question, error) {
	qs := &Questions{}

	if err := r.DB.WithContext(ctx).First(qs, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("find question by id: %w", err)
	}

	return &domain.Question{
		ID:        qs.ID,
		Text:      qs.QuestionText,
		CreatedAt: qs.CreatedAt,
	}, nil
}

func (r *QuestionRepository) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	var res []Questions

	if err := r.DB.WithContext(ctx).Find(&res).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get all questions: %w", err)
	}

	qs := make([]*domain.Question, 0, len(res))
	for _, q := range res {
		qs = append(qs, &domain.Question{
			ID:        q.ID,
			Text:      q.QuestionText,
			CreatedAt: q.CreatedAt,
		})
	}
	return qs, nil
}

func (r *QuestionRepository) SaveQuestion(ctx context.Context, text string) (*domain.Question, error) {
	qs := &Questions{
		QuestionText: text,
	}

	res := r.DB.WithContext(ctx).Create(qs)
	if res.Error != nil {
		return nil, fmt.Errorf("create error: %w", res.Error)
	}

	return &domain.Question{
		ID: qs.ID,
		Text: qs.QuestionText,
		CreatedAt: qs.CreatedAt,
	}, nil
}

func (r *QuestionRepository) GetQuestionWithAnswers(ctx context.Context, questionID int) (*domain.Question, []*domain.Answer, error) {
	var qs Questions

	if err := r.DB.WithContext(ctx).Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&qs, questionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, repository.ErrNotFound
		}
		return nil, nil, fmt.Errorf("get question and answers: %w", err)
	}

	q := &domain.Question{ID: qs.ID, Text: qs.QuestionText, CreatedAt: qs.CreatedAt}
	answers := make([]*domain.Answer, 0, len(qs.Answers))
	for _, a := range qs.Answers {
		answers = append(answers, &domain.Answer{
			ID:         a.ID,
			QuestionID: a.QuestionID,
			UserID:     a.UserID,
			Text:       a.AnswerText,
			CreatedAt:  a.CreatedAt,
		})
	}

	return q, answers, nil
}

func (r *QuestionRepository) DeleteQuestion(ctx context.Context, questionID int) error {
	res := r.DB.WithContext(ctx).Delete(&Questions{}, questionID)
	if res.Error != nil {
		return fmt.Errorf("delete_question error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return repository.ErrNoRowDeleted
	}

	return nil
}
