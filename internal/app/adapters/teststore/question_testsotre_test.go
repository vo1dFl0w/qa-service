package teststore_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/qa-service/internal/app/adapters/teststore"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

func TestQuestionRepository_FindQuestionByID(t *testing.T) {
	s := teststore.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)

	res, err := s.Question().FindQuestionByID(context.TODO(), q.ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, q, res)

	_, err = s.Question().FindQuestionByID(context.TODO(), 100)
	assert.Error(t, err)
}

func TestQuestionRepository_GetAllQuestions(t *testing.T) {
	s := teststore.New()

	q1, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q1)

	q2, err := s.Question().SaveQuestion(context.TODO(), "question_2")
	assert.NoError(t, err)
	assert.NotNil(t, q2)

	res, err := s.Question().GetAllQuestions(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, q1, res[0])
	assert.Equal(t, q2, res[1])
}

func TestQuestionRepository_SaveQuestion(t *testing.T) {
	s := teststore.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, 1, q.ID)
	assert.Equal(t, "question_1", q.Text)
}

func TestQuestionRepository_GetQuestionWithAnswers(t *testing.T) {
	s := teststore.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, 1, q.ID)
	assert.Equal(t, "question_1", q.Text)

	userID := uuid.New()

	a1, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_1")
	assert.NoError(t, err)
	assert.NotNil(t, a1)
	assert.Equal(t, 1, a1.ID)
	assert.Equal(t, q.ID, a1.QuestionID)
	assert.Equal(t, "answer_1", a1.Text)
	assert.Equal(t, a1.UserID, userID)

	a2, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_2")
	assert.NoError(t, err)
	assert.NotNil(t, a2)
	assert.Equal(t, 2, a2.ID)
	assert.Equal(t, q.ID, a2.QuestionID)
	assert.Equal(t, "answer_2", a2.Text)
	assert.Equal(t, a2.UserID, userID)

	qRes, aRes, err := s.Question().GetQuestionWithAnswers(context.TODO(), q.ID)
	assert.NoError(t, err)
	assert.NotNil(t, qRes)
	assert.NotNil(t, aRes)
	assert.Equal(t, q, qRes)
	assert.Equal(t, a1, aRes[0])
	assert.Equal(t, a2, aRes[1])

	_, _, err = s.Question().GetQuestionWithAnswers(context.TODO(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestQuestionRepository_DeleteQuestion(t *testing.T) {
	s := teststore.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, 1, q.ID)
	assert.Equal(t, "question_1", q.Text)

	userID := uuid.New()

	a1, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_1")
	assert.NoError(t, err)
	assert.NotNil(t, a1)
	assert.Equal(t, 1, a1.ID)
	assert.Equal(t, q.ID, a1.QuestionID)
	assert.Equal(t, "answer_1", a1.Text)
	assert.Equal(t, a1.UserID, userID)

	a2, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_2")
	assert.NoError(t, err)
	assert.NotNil(t, a2)
	assert.Equal(t, 2, a2.ID)
	assert.Equal(t, q.ID, a2.QuestionID)
	assert.Equal(t, "answer_2", a2.Text)
	assert.Equal(t, a2.UserID, userID)

	err = s.Question().DeleteQuestion(context.TODO(), q.ID)
	assert.NoError(t, err)

	_, err = s.Answer().GetAnswer(context.TODO(), a1.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)

	_, err = s.Answer().GetAnswer(context.TODO(), a2.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)

	err = s.Question().DeleteQuestion(context.TODO(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNoRowDeleted)
}
