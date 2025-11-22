package httptransport

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/middlewares"
	"github.com/vo1dFl0w/qa-service/internal/transport/utils"
)

var (
	questions         = "/questions/"
	questionsID       = "/questions/{id}"
	answersID         = "/answers/{id}"
	questionsIDanswers = "/questions/{id}/answers/"
)

func (h *Handler) Router() http.Handler {
	h.root = middlewares.LoggerMiddleware(h.log)(h.router)

	h.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		id, handlerType, err := parseURL(r.URL.Path)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		switch handlerType {
		case questions:
			if r.Method == http.MethodGet {
				h.questionHandler.GetQuestionsList()(w, r)
			} else {
				h.questionHandler.CreateQuestion()(w, r)
			}
		case questionsID:
			if r.Method == http.MethodGet {
				h.questionHandler.GetQuestionWithAnswers(id)(w, r)
			} else {
				h.questionHandler.DeleteQuestion(id)(w, r)
			}
		case answersID:
			if r.Method == http.MethodGet {
				h.answerHandler.GetAnswer(id)(w, r)
			} else {
				h.answerHandler.DeleteAnswer(id)(w, r)
			}
		case questionsIDanswers:
			if r.Method == http.MethodPost {
				h.answerHandler.AddAnswer(id)(w, r)
			}
		default:
			utils.ErrorResponse(w, r, http.StatusNotFound, fmt.Errorf("resource not found"))
			return
		}
	})

	return h.router
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.root.ServeHTTP(w, r)
}

func parseURL(url string) (int, string, error) {
	var (
		idReq   int
		pathReq string
	)

	path := strings.Trim(url, "/")
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == ' '
	})

	switch {
	case len(parts) == 0:
		return 0, "", fmt.Errorf("invalid url request")
	case len(parts) == 1 && parts[0] == "questions":
		pathReq = questions
		return 0, pathReq, nil
	case len(parts) == 2 && (parts[0] == "answers" || parts[0] == "questions"):
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, "", fmt.Errorf("invalid answer_id: %w", err)
		}
		idReq = id

		switch parts[0] {
		case "answers":
			pathReq = answersID
			return id, pathReq, nil
		case "questions":
			pathReq = questionsID
			return id, pathReq, nil

		}
	case len(parts) == 3 && parts[0] == "questions" && parts[2] == "answers":
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, "", fmt.Errorf("invalid answer_id")
		}
		idReq = id
		pathReq = questionsIDanswers
		return idReq, pathReq, nil
	default:
		return 0, "", fmt.Errorf("invalid url")
	}

	return idReq, pathReq, nil
}
