package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/qa-service/internal/app/usecase"
	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/utils"
)

type AnswerHandler struct {
	answerUsecase usecase.AnswerService
}

func NewAnswerHandler(answerUsecase usecase.AnswerService) *AnswerHandler {
	return &AnswerHandler{answerUsecase: answerUsecase}
}

func (h *AnswerHandler) GetAnswer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		answerID, err := getIdFromURL(r)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		res, err := h.answerUsecase.GetAnswerByID(ctx, answerID)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusOK, res)
	}
}

func (h *AnswerHandler) AddAnswer() http.HandlerFunc {
	type request struct {
		UserID uuid.UUID `json:"user_id,omitempty"`
		Text   string    `json:"text"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		questionID, err := getIdFromURL(r)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, ErrBadRequest)
			return
		}

		if req.Text == "" {
			utils.ErrorResponse(w, r, http.StatusBadRequest, fmt.Errorf("empty text: %s", ErrBadRequest))
			return
		}

		res, err := h.answerUsecase.AddAnswerByID(ctx, questionID, req.UserID, req.Text)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]interface{}{"status": "success", "created": res})
	}
}

func (h *AnswerHandler) DeleteAnswer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), ctxDelay)
		defer cancel()

		answerID, err := getIdFromURL(r)
		if err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		if err := h.answerUsecase.DeleteAnswerByID(ctx, answerID); err != nil {
			utils.ErrorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		utils.Response(w, r, http.StatusCreated, map[string]string{"status": "success"})
	}
}
