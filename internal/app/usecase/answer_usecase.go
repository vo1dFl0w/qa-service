package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type AnswerService interface {
	AddAnswerByID(ctx context.Context, question_id int, user_id uuid.UUID, text string) (*domain.Answer, error)
	GetAnswerByID(ctx context.Context, answer_id int) (*domain.Answer, error)
	DeleteAnswerByID(ctx context.Context, answer_id int) error
}

type answerService struct {
	repoAnswer   repository.AnswerRepository
	repoQuestion repository.QuestionRepository
}

func NewAnswerService(repoAnswer repository.AnswerRepository, repoQuestion repository.QuestionRepository) AnswerService {
	return &answerService{repoAnswer: repoAnswer, repoQuestion: repoQuestion}
}

func (s *answerService) AddAnswerByID(ctx context.Context, question_id int, user_id uuid.UUID, text string) (*domain.Answer, error) {
	if err := validateID(question_id); err != nil {
		return nil, fmt.Errorf("invalid question id: %w", err)
	}

	if user_id != uuid.Nil {
		_, err := uuid.Parse(user_id.String())
		if err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
	}

	if text == "" {
		return nil, domain.ErrEmptyText
	}

	_, err := s.repoQuestion.FindQuestionByID(ctx, question_id)
	if err != nil {
		return nil, fmt.Errorf("find_question_by_id error: %w", err)
	}

	return s.repoAnswer.SaveAnswer(ctx, question_id, user_id, text)
}

func (s *answerService) GetAnswerByID(ctx context.Context, answer_id int) (*domain.Answer, error) {
	if err := validateID(answer_id); err != nil {
		return nil, fmt.Errorf("invalid answer id: %w", err)
	}

	return s.repoAnswer.GetAnswer(ctx, answer_id)
}
func (s *answerService) DeleteAnswerByID(ctx context.Context, answer_id int) error {
	if err := validateID(answer_id); err != nil {
		return fmt.Errorf("invalid answer id: %w", err)
	}

	return s.repoAnswer.DeleteAnswer(ctx, answer_id)
}
