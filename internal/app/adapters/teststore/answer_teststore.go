package teststore

import (
	"context"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
)

type AnswerRepository struct {
	storage *Storage
	answers map[int]*domain.Answer
}

func (r *AnswerRepository) SaveAnswer(ctx context.Context, question_id int, user_id uuid.UUID, text string) (*domain.Answer, error) {
	q, ok := r.storage.questionRepository.questions[question_id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	id := len(r.answers)+1
	a := &domain.Answer{
		ID:         id,
		QuestionID: q.ID,
		UserID:     user_id,
		Text:       text,
	}

	r.answers[id] = a

	return a, nil
}

func (r *AnswerRepository) GetAnswer(ctx context.Context, answer_id int) (*domain.Answer, error) {
	a, ok := r.answers[answer_id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	return a, nil
}

func (r *AnswerRepository) DeleteAnswer(ctx context.Context, answer_id int) error {
	_, ok := r.answers[answer_id]
	if !ok {
		return domain.ErrNoRowDeleted
	}
	delete(r.answers, answer_id)

	return nil
}
