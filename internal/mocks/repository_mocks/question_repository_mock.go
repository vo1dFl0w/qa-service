package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
)

type QuestionRepositoryMock struct {
	mock.Mock
}

func (m *QuestionRepositoryMock) FindQuestionByID(ctx context.Context, question_id int) (*domain.Question, error) {
	args := m.Called(ctx, question_id)

	if v := args.Get(0); v != nil {
		return v.(*domain.Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *QuestionRepositoryMock) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	args := m.Called(ctx)

	if v := args.Get(0); v != nil {
		return v.([]*domain.Question), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *QuestionRepositoryMock) SaveQuestion(ctx context.Context, text string) (*domain.Question, error) {
	args := m.Called(ctx, text)
	
	if v := args.Get(0); v != nil {
		return v.(*domain.Question), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *QuestionRepositoryMock) GetQuestionWithAnswers(ctx context.Context, question_id int) (*domain.Question, []*domain.Answer, error) {
	var q *domain.Question
	var as []*domain.Answer

	args := m.Called(ctx, question_id)
	if v := args.Get(0); v != nil {
		q = v.(*domain.Question)
	}
	if v := args.Get(1); v != nil {
		as = v.([]*domain.Answer)
	}
	return q, as, args.Error(2)
}

func (m *QuestionRepositoryMock) DeleteQuestion(ctx context.Context, question_id int) error {
	args := m.Called(ctx, question_id)
	return args.Error(0)
}