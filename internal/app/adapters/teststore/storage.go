package teststore

import (
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
)

type Storage struct{
	questionRepository *QuestionRepository
	answerRepository *AnswerRepository
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Question() repository.QuestionRepository {
	if s.questionRepository != nil {
		return s.questionRepository
	}

	s.questionRepository = &QuestionRepository{
		storage: s,
		questions: make(map[int]*domain.Question),
	}

	return s.questionRepository
}

func (s *Storage) Answer() repository.AnswerRepository {
	if s.answerRepository != nil {
		return s.answerRepository
	}

	s.answerRepository = &AnswerRepository{
		storage: s,
		answers: make(map[int]*domain.Answer),
	}

	return s.answerRepository
}