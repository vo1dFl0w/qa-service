package httptransport

import (
	"net/http"

	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/middlewares"
)

func (h *Handler) Router() http.Handler {
	h.root = middlewares.LoggerMiddleware(h.log)(h.router)

	h.router.HandleFunc("GET /questions/", h.questionHandler.GetQuestionsList())
	h.router.HandleFunc("POST /questions/", h.questionHandler.CreateQuestion())
	h.router.HandleFunc("GET /questions/{id}", h.questionHandler.GetQuestionWithAnswers())
	h.router.HandleFunc("DELETE /questions/{id}", h.questionHandler.DeleteQuestion())

	h.router.HandleFunc("POST /questions/{id}/answers/", h.answerHandler.AddAnswer())
	h.router.HandleFunc("GET /answers/{id}", h.answerHandler.GetAnswer())
	h.router.HandleFunc("DELETE /answers/{id}", h.answerHandler.DeleteAnswer())

	return h.router
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.root.ServeHTTP(w, r)
}
