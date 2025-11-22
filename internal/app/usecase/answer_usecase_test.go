package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/mocks/repository_mocks"
)

func TestAnswerService_AddAnswerByID(t *testing.T) {
	answerRepo := &mocks.AnswerRepositoryMock{}
	questionRepo := &mocks.QuestionRepositoryMock{}
	answerUsecase := usecase.NewAnswerService(answerRepo, questionRepo)

	userID := uuid.New()
	text := "answer_1"

	q := &domain.Question{
		ID:   1,
		Text: "question_1",
	}

	answerExp := &domain.Answer{
		ID:         1,
		QuestionID: q.ID,
		UserID:     userID,
		Text:       "answer_1"}

	questionRepo.On("FindQuestionByID", mock.Anything, q.ID).Return(q, nil).Once()
	answerRepo.On("SaveAnswer", mock.Anything, q.ID, userID, text).Return(answerExp, nil).Once()

	a1, err := answerUsecase.AddAnswerByID(context.Background(), q.ID, userID, text)
	assert.NoError(t, err)
	assert.NotNil(t, a1)

	questionRepo.On("FindQuestionByID", mock.Anything, 100).Return(nil, domain.ErrNotFound).Once()

	_, err = answerUsecase.AddAnswerByID(context.Background(), 100, userID, text)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	_, err = answerUsecase.AddAnswerByID(context.Background(), q.ID, userID, "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEmptyText)

	questionRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
}

func TestAnswerService_GetAnswerByID(t *testing.T) {
	answerRepo := &mocks.AnswerRepositoryMock{}
	answerUsecase := usecase.NewAnswerService(answerRepo, nil)

	answerID := 1
	answerExp := &domain.Answer{
		ID:         answerID,
		QuestionID: 1,
		UserID:     uuid.New(),
		Text:       "answer_1",
	}

	answerRepo.On("GetAnswer", mock.Anything, answerID).Return(answerExp, nil).Once()

	res, err := answerUsecase.GetAnswerByID(context.Background(), answerID)
	assert.NoError(t, err)
	assert.Equal(t, answerExp, res)

	answerRepo.On("GetAnswer", mock.Anything, 100).Return(nil, domain.ErrNotFound).Once()

	_, err = answerUsecase.GetAnswerByID(context.Background(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	answerRepo.AssertExpectations(t)
}

func TestAnswerService_DeleteAnswerByID(t *testing.T) {
	answerRepo := &mocks.AnswerRepositoryMock{}
	answerUsecase := usecase.NewAnswerService(answerRepo, nil)

	answerRepo.On("DeleteAnswer", mock.Anything, 5).Return(nil).Once()

	err := answerUsecase.DeleteAnswerByID(context.Background(), 5)
	assert.NoError(t, err)

	answerRepo.On("DeleteAnswer", mock.Anything, 100).Return(domain.ErrNotFound).Once()

	err = answerUsecase.DeleteAnswerByID(context.Background(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	answerRepo.AssertExpectations(t)
}
