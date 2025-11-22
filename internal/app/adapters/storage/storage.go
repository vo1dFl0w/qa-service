package postgres

import "github.com/vo1dFl0w/qa-service/internal/app/repository"

type Storage interface {
	Question() repository.QuestionRepository
	Answer() repository.AnswerRepository
}
