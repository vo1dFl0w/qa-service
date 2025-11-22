package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/transport/utils"
)

type AnswerHandler struct {
	answerUsecase usecase.AnswerService
}

func NewAnswerHandler(answerUsecase usecase.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerUsecase: answerUsecase}
}

func (h *AnswerHandler) GetAnswer(answer_id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()
		
		if err := validateMethod(http.MethodGet, r.Method); err != nil {
			utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, err)
			return
		}

		res, err := h.answerUsecase.GetAnswerByID(ctx, answer_id)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}

func (h *AnswerHandler) AddAnswer(question_id int) http.HandlerFunc {
	type request struct {
		UserID uuid.UUID `json:"user_id,omitempty"`
		Text   string    `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		if err := validateMethod(http.MethodPost, r.Method); err != nil {
			utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, utils.ErrBadRequest)
			return
		}

		if req.Text == "" {
			utils.ErrorResponse(w, r, http.StatusBadRequest, fmt.Errorf("empty text: %s", utils.ErrBadRequest))
			return
		}

		res, err := h.answerUsecase.AddAnswerByID(ctx, question_id, req.UserID, req.Text)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]interface{}{"status": "success", "created": res})
	}
}

func (h *AnswerHandler) DeleteAnswer(answer_id int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		if err := validateMethod(http.MethodDelete, r.Method); err != nil {
			utils.ErrorResponse(w, r, http.StatusMethodNotAllowed, err)
			return
		}

		if err := h.answerUsecase.DeleteAnswerByID(ctx, answer_id); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]string{"status": "success"})
	}
}

