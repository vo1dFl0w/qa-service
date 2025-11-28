package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type AnswerService interface {
	AddAnswerByID(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error)
	GetAnswerByID(ctx context.Context, answerID int) (*domain.Answer, error)
	DeleteAnswerByID(ctx context.Context, answerID int) error
}

type answerService struct {
	repoAnswer   repository.AnswerRepository
	repoQuestion repository.QuestionRepository
}

func NewAnswerService(repoAnswer repository.AnswerRepository, repoQuestion repository.QuestionRepository) AnswerService {
	return &answerService{repoAnswer: repoAnswer, repoQuestion: repoQuestion}
}

func (s *answerService) AddAnswerByID(ctx context.Context, questionID int, userID uuid.UUID, text string) (*domain.Answer, error) {
	if err := validateID(questionID); err != nil {
		return nil, fmt.Errorf("invalid question id: %w", err)
	}

	if userID != uuid.Nil {
		_, err := uuid.Parse(userID.String())
		if err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
	}

	if text == "" {
		return nil, domain.ErrEmptyText
	}

	_, err := s.repoQuestion.FindQuestionByID(ctx, questionID)
	if err != nil {
		return nil, fmt.Errorf("find question by id: %w", err)
	}

	return s.repoAnswer.SaveAnswer(ctx, questionID, userID, text)
}

func (s *answerService) GetAnswerByID(ctx context.Context, answerID int) (*domain.Answer, error) {
	if err := validateID(answerID); err != nil {
		return nil, fmt.Errorf("invalid answer id: %w", err)
	}

	return s.repoAnswer.GetAnswer(ctx, answerID)
}
func (s *answerService) DeleteAnswerByID(ctx context.Context, answerID int) error {
	if err := validateID(answerID); err != nil {
		return fmt.Errorf("invalid answer id: %w", err)
	}

	return s.repoAnswer.DeleteAnswer(ctx, answerID)
}
