package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
)

type AnswerUsecaseMock struct {
	mock.Mock
}

func (m *AnswerUsecaseMock) GetAnswerByID(ctx context.Context, answerID int) (*domain.Answer, error) {
	args := m.Called(ctx, answerID)

	if v := args.Get(0); v != nil {
		return v.(*domain.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AnswerUsecaseMock) AddAnswerByID(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error) {
	args := m.Called(ctx, questionID, userID, text)

	if v := args.Get(0); v != nil {
		return v.(*domain.Answer), args.Error(1)
	}
	
	return nil, args.Error(1)
}

func (m *AnswerUsecaseMock) DeleteAnswerByID(ctx context.Context, answerID int) error {
	args := m.Called(ctx, answerID)

	return args.Error(0)
}
