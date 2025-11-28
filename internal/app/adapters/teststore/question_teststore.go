package teststore

import (
	"context"
	"time"

	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type QuestionRepository struct {
	storage   *Storage
	questions map[int]*domain.Question
}

func (r *QuestionRepository) FindQuestionByID(ctx context.Context, questionID int) (*domain.Question, error) {
	q, ok := r.questions[questionID]
	if !ok {
		return nil, repository.ErrNotFound
	}

	return &domain.Question{
		ID:        q.ID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
	}, nil
}

func (r *QuestionRepository) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	res := make([]*domain.Question, 0, len(r.questions))

	for _, q := range r.questions {
		res = append(res, q)
	}

	return res, nil
}

func (r *QuestionRepository) SaveQuestion(ctx context.Context, text string) (*domain.Question, error) {
	id := len(r.questions)+1
	r.questions[id] = &domain.Question{ID: id, Text: text, CreatedAt: time.Now()}

	return r.questions[id], nil
}

func (r *QuestionRepository) GetQuestionWithAnswers(ctx context.Context, questionID int) (*domain.Question, []*domain.Answer, error) {
	q, ok := r.questions[questionID]
	if !ok {
		return nil, nil, repository.ErrNotFound
	}

	as := make([]*domain.Answer, 0, len(r.storage.answerRepository.answers))
	for _, a := range r.storage.answerRepository.answers {
		if a.QuestionID == questionID {
			as = append(as, a)
		}
	}

	return q, as, nil
}

func (r *QuestionRepository) DeleteQuestion(ctx context.Context, questionID int) error {
	_, ok := r.questions[questionID]
	if !ok {
		return repository.ErrNoRowDeleted
	}

	delete(r.questions, questionID)
	for _, a := range r.storage.answerRepository.answers {
		if a.QuestionID == questionID {
			delete(r.storage.answerRepository.answers, a.ID)
		}
	}

	return nil
}
