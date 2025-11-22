package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"gorm.io/gorm"
)

type QuestionRepository struct {
	DB *gorm.DB
}

func (r *QuestionRepository) FindQuestionByID(ctx context.Context, question_id int) (*domain.Question, error) {
	qs := &Questions{}

	if err := r.DB.WithContext(ctx).First(qs, question_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find_question_by_id error: %w", err)
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
		return nil, fmt.Errorf("get_all_questions error: %w", err)
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

func (r *QuestionRepository) GetQuestionWithAnswers(ctx context.Context, question_id int) (*domain.Question, []*domain.Answer, error) {
	var qs Questions

	if err := r.DB.WithContext(ctx).Preload("Answers", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&qs, question_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, domain.ErrNotFound
		}
		return nil, nil, fmt.Errorf("get_question_and_answers error: %w", err)
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

func (r *QuestionRepository) DeleteQuestion(ctx context.Context, question_id int) error {
	res := r.DB.WithContext(ctx).Delete(&Questions{}, question_id)
	if res.Error != nil {
		return fmt.Errorf("delete_question error: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return domain.ErrNoRowDeleted
	}

	return nil
}
