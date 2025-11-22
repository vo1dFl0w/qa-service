package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
)

type AnswerRepositoryMock struct {
	mock.Mock
}

func (m *AnswerRepositoryMock) GetAnswer(ctx context.Context, answerID int) (*domain.Answer, error) {
	args := m.Called(ctx, answerID)

	if v := args.Get(0); v != nil {
		return v.(*domain.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AnswerRepositoryMock) SaveAnswer(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error) {
	args := m.Called(ctx, questionID, userID, text)

	if v := args.Get(0); v != nil {
		return v.(*domain.Answer), args.Error(1)
	}
	
	return nil, args.Error(1)
}

func (m *AnswerRepositoryMock) DeleteAnswer(ctx context.Context, answerID int) error {
	args := m.Called(ctx, answerID)

	return args.Error(0)
}