package usecase

import (
	"context"
	"fmt"

	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type QuestionService interface {
	FindQuestionByID(ctx context.Context, questionID int) (*domain.Question, error)
	GetAllQuestions(ctx context.Context) ([]*domain.Question, error)
	CreateQuestion(ctx context.Context, text string) (*domain.Question, error)
	GetQuestionWithAnswers(ctx context.Context, questionID int) (*domain.Question, []*domain.Answer, error)
	DeleteQuestionByID(ctx context.Context, questionID int) error
}

type questionService struct {
	repository repository.QuestionRepository
}

func NewQuestionService(repository repository.QuestionRepository) QuestionService {
	return &questionService{repository: repository}
}

func (s *questionService) FindQuestionByID(ctx context.Context, questionID int) (*domain.Question, error) {
	if err := validateID(questionID); err != nil {
		return nil, fmt.Errorf("invalid question id: %w", err)
	}

	return s.repository.FindQuestionByID(ctx, questionID)
}

func (s *questionService) GetAllQuestions(ctx context.Context) ([]*domain.Question, error) {
	return s.repository.GetAllQuestions(ctx)
}

func (s *questionService) CreateQuestion(ctx context.Context, text string) (*domain.Question, error) {
	if text == "" {
		return nil, domain.ErrEmptyText
	}

	return s.repository.SaveQuestion(ctx, text)
}

func (s *questionService) GetQuestionWithAnswers(ctx context.Context, questionID int) (*domain.Question, []*domain.Answer, error) {
	if err := validateID(questionID); err != nil {
		return nil, nil, fmt.Errorf("invalid question id: %w", err)
	}

	return s.repository.GetQuestionWithAnswers(ctx, questionID)
}

func (s *questionService) DeleteQuestionByID(ctx context.Context, questionID int) error {
	if err := validateID(questionID); err != nil {
		return fmt.Errorf("invalid question id: %w", err)
	}

	return s.repository.DeleteQuestion(ctx, questionID)
}
