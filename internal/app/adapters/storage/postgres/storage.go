package postgres

import (
	"github.com/vo1dFl0w/qa-service/internal/app/repository"
	"gorm.io/gorm"
)

type Storage struct {
	DB                 *gorm.DB
	questionRepository repository.QuestionRepository
	answerRepository   repository.AnswerRepository
}

func New(db *gorm.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Question() repository.QuestionRepository {
	if s.questionRepository != nil {
		return s.questionRepository
	}

	s.questionRepository = &QuestionRepository{
		DB: s.DB,
	}

	return s.questionRepository
}

func (s *Storage) Answer() repository.AnswerRepository {
	if s.answerRepository != nil {
		return s.answerRepository
	}

	s.answerRepository = &AnswerRepository{
		DB: s.DB,
	}

	return s.answerRepository
}
