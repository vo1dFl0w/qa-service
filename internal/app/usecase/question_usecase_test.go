package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/mocks/repository_mocks"
)

func TestQuestionService_FindQuestionByID(t *testing.T) {
	questionRepo := &mocks.QuestionRepositoryMock{}
	questionUsecase := usecase.NewQuestionService(questionRepo)

	q := &domain.Question{
		ID:   1,
		Text: "question_1",
	}

	questionRepo.On("FindQuestionByID", mock.Anything, q.ID).Return(q, nil).Once()

	res, err := questionUsecase.FindQuestionByID(context.Background(), q.ID)
	assert.NoError(t, err)
	assert.Equal(t, q, res)

	questionRepo.On("FindQuestionByID", mock.Anything, 100).Return(nil, domain.ErrNotFound).Once()

	_, err = questionUsecase.FindQuestionByID(context.Background(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	questionRepo.AssertExpectations(t)
}

func TestQuestionService_GetAllQuestions(t *testing.T) {
	questionRepo := &mocks.QuestionRepositoryMock{}
	questionUsecase := usecase.NewQuestionService(questionRepo)

	q := []*domain.Question{
		{ID: 1, Text: "question_1"},
		{ID: 2, Text: "question_2"},
	}

	questionRepo.On("GetAllQuestions", mock.Anything).Return(q, nil).Once()

	res, err := questionUsecase.GetAllQuestions(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, q, res)

	questionRepo.AssertExpectations(t)
}

func TestQuestionService_CreateQuestion(t *testing.T) {
	questionRepo := &mocks.QuestionRepositoryMock{}
	questionUsecase := usecase.NewQuestionService(questionRepo)

	text := "question_1"
	q := &domain.Question{
		ID:   1,
		Text: text,
	}

	questionRepo.On("SaveQuestion", mock.Anything, text).Return(q, nil).Once()

	res, err := questionUsecase.CreateQuestion(context.Background(), text)
	assert.NoError(t, err)
	assert.Equal(t, q, res)

	_, err = questionUsecase.CreateQuestion(context.Background(), "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEmptyText)

	questionRepo.AssertExpectations(t)
}

func TestQuestionService_GetQuestionWithAnswers(t *testing.T) {
	questionRepo := &mocks.QuestionRepositoryMock{}
	answerRepo := &mocks.AnswerRepositoryMock{}
	questionUsecase := usecase.NewQuestionService(questionRepo)

	q := &domain.Question{ID: 1, Text: "question_1"}
	answers := []*domain.Answer{
		{ID: 1, QuestionID: q.ID, Text: "answer_1"},
		{ID: 2, QuestionID: q.ID, Text: "answer_2"},
	}

	questionRepo.On("GetQuestionWithAnswers", mock.Anything, q.ID).Return(q, answers, nil).Once()

	qs, as, err := questionUsecase.GetQuestionWithAnswers(context.Background(), q.ID)
	assert.NoError(t, err)
	assert.Equal(t, q, qs)
	assert.Equal(t, answers[0], as[0])
	assert.Equal(t, answers[1], as[1])

	questionRepo.AssertExpectations(t)
	answerRepo.AssertExpectations(t)
}

func TestQuestionService_DeleteQuestionByID(t *testing.T) {
	questionRepo := &mocks.QuestionRepositoryMock{}
	questionUsecase := usecase.NewQuestionService(questionRepo)

	questionRepo.On("DeleteQuestion", mock.Anything, 1).Return(nil).Once()
	err := questionUsecase.DeleteQuestionByID(context.Background(), 1)
	assert.NoError(t, err)

	questionRepo.On("DeleteQuestion", mock.Anything, 100).Return(domain.ErrNotFound).Once()
	err = questionUsecase.DeleteQuestionByID(context.Background(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	questionRepo.AssertExpectations(t)
}
