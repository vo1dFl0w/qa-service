package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/utils"
)

type QuestionHandler struct {
	questionUsecase usecase.QuestionService
}

func NewQuestionHandler(questionUsecase usecase.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionUsecase: questionUsecase,
	}
}

func (h *QuestionHandler) GetQuestionsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		res, err := h.questionUsecase.GetAllQuestions(ctx)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}

func (h *QuestionHandler) CreateQuestion() http.HandlerFunc {
	type request struct {
		Text string `json:"text"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, utils.ErrBadRequest)
			return
		}

		if strings.TrimSpace(req.Text) == "" {
			utils.ErrorResponse(w, r, http.StatusBadRequest, fmt.Errorf("empty text"))
			return
		}

		res, err := h.questionUsecase.CreateQuestion(ctx, req.Text)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]interface{}{"status": "success", "created": res})
	}
}

func (h *QuestionHandler) GetQuestionWithAnswers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		questionID, err := getIdFromURL(r)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		q, a, err := h.questionUsecase.GetQuestionWithAnswers(ctx, questionID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusOK, map[string]interface{}{"question": q, "answers": a})
	}
}

func (h *QuestionHandler) DeleteQuestion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		questionID, err := getIdFromURL(r)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.questionUsecase.DeleteQuestionByID(ctx, questionID); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]string{"status": "success"})
	}
}
