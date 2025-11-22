package httptransport

import (
	"log/slog"
	"net/http"

	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/handlers"
)

type Handler struct {
	router          *http.ServeMux
	root            http.Handler
	log             *slog.Logger
	questionHandler *handlers.QuestionHandler
	answerHandler   *handlers.AnswerHandler
}

func NewHandler(log *slog.Logger, questionUsecase usecase.QuestionService, answerUsecase usecase.AnswerService) *Handler {
	h := &Handler{
		router: http.NewServeMux(),
		log:    log,
		questionHandler: handlers.NewQuestionHandler(questionUsecase),
		answerHandler: handlers.NewAnswerHandler(answerUsecase),
	}

	h.Router()
	return h
}
