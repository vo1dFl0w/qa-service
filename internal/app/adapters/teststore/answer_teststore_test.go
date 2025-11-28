package teststore_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/qa-service/internal/app/adapters/teststore"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

func TestAnswerRepository_SaveAnswer(t *testing.T) {
	s := teststore.New()
	userID := uuid.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)

	a, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_1")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ID)
	assert.Equal(t, q.ID, a.QuestionID)
	assert.Equal(t, "answer_1", a.Text)
	assert.Equal(t, a.UserID, userID)

	_, err = s.Answer().SaveAnswer(context.TODO(), 100, userID, "answer_1")
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestAnswerRepository_GetAnswer(t *testing.T) {
	s := teststore.New()

	userID := uuid.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)

	a, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_1")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ID)
	assert.Equal(t, q.ID, a.QuestionID)
	assert.Equal(t, "answer_1", a.Text)
	assert.Equal(t, a.UserID, userID)

	res, err := s.Answer().GetAnswer(context.TODO(), a.ID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, a, res)

	_, err = s.Answer().GetAnswer(context.TODO(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestAnswerRepository_DeleteAnswer(t *testing.T) {
	s := teststore.New()

	userID := uuid.New()

	q, err := s.Question().SaveQuestion(context.TODO(), "question_1")
	assert.NoError(t, err)
	assert.NotNil(t, q)

	a, err := s.Answer().SaveAnswer(context.TODO(), q.ID, userID, "answer_1")
	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, 1, a.ID)
	assert.Equal(t, q.ID, a.QuestionID)
	assert.Equal(t, "answer_1", a.Text)
	assert.Equal(t, a.UserID, userID)

	err = s.Answer().DeleteAnswer(context.TODO(), a.ID)
	assert.NoError(t, err)

	err = s.Answer().DeleteAnswer(context.TODO(), 100)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repository.ErrNoRowDeleted)
}
