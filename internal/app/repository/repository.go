package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
)

type AnswerRepository interface {
	SaveAnswer(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error)
	GetAnswer(ctx context.Context, answerID int) (*domain.Answer, error)
	DeleteAnswer(ctx context.Context, answerID int) error
}

type QuestionRepository interface {
	FindQuestionByID(ctx context.Context, questionID int) (*domain.Question, error)
	GetAllQuestions(ctx context.Context) ([]*domain.Question, error)
	SaveQuestion(ctx context.Context, text string) (*domain.Question, error)
	GetQuestionWithAnswers(ctx context.Context, questionID int) (*domain.Question, []*domain.Answer, error)
	DeleteQuestion(ctx context.Context, questionID int) error
}
