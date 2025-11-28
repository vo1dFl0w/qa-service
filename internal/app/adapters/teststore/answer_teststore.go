package teststore

import (
	"context"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type AnswerRepository struct {
	storage *Storage
	answers map[int]*domain.Answer
}

func (r *AnswerRepository) SaveAnswer(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error) {
	q, ok := r.storage.questionRepository.questions[questionID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	id := len(r.answers)+1
	a := &domain.Answer{
		ID:         id,
		QuestionID: q.ID,
		UserID:     userID,
		Text:       text,
	}

	r.answers[id] = a

	return a, nil
}

func (r *AnswerRepository) GetAnswer(ctx context.Context, answerID int) (*domain.Answer, error) {
	a, ok := r.answers[answerID]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return a, nil
}

func (r *AnswerRepository) DeleteAnswer(ctx context.Context, answerID int) error {
	_, ok := r.answers[answerID]
	if !ok {
		return repository.ErrNoRowDeleted
	}
	delete(r.answers, answerID)

	return nil
}
