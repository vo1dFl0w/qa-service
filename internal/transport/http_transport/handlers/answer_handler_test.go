package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/mocks/usecase_mocks"
	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/handlers"
)

func TestAnswerHandler_GetAnswer(t *testing.T) {
	testCases := []struct {
		name       string
		answerID   int
		mockResult *domain.Answer
		mockErr    error
		expStatus  int
		expText    string
	}{
		{
			name:       "success",
			answerID:   1,
			mockResult: &domain.Answer{ID: 1, QuestionID: 10, UserID: uuid.New(), Text: "answer_1"},
			mockErr:    nil,
			expStatus:  http.StatusOK,
			expText:    "answer_1",
		},
		{
			name:       "not found",
			answerID:   100,
			mockResult: nil,
			mockErr:    domain.ErrNotFound,
			expStatus:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			answerUsecase := &mocks.AnswerUsecaseMock{}
			handler := handlers.NewAnswerHandler(answerUsecase)

			answerUsecase.On("GetAnswerByID", mock.Anything, tc.answerID).Return(tc.mockResult, tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /answers/{id}", handler.GetAnswer())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/answers/%d", tc.answerID), nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expStatus, res.StatusCode)
			if tc.mockResult != nil {
				var a domain.Answer
				json.NewDecoder(res.Body).Decode(&a)
				assert.Equal(t, tc.expText, a.Text)
			}
		})
	}
}

func TestAnswerHandler_AddAnswer(t *testing.T) {
	userID := uuid.New()

	testCases := []struct {
		name       string
		questionID int
		text       string
		mockErr    error
		expStatus  int
	}{
		{
			name:       "success",
			questionID: 1,
			text:       "answer_1",
			mockErr:    nil,
			expStatus:  http.StatusCreated,
		},
		{
			name:       "empty text",
			questionID: 1,
			text:       "",
			mockErr:    nil,
			expStatus:  http.StatusBadRequest,
		},
		{
			name:       "question not found",
			questionID: 100,
			text:       "answer_1",
			mockErr:    domain.ErrNotFound,
			expStatus:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			answerUsecase := &mocks.AnswerUsecaseMock{}
			handler := handlers.NewAnswerHandler(answerUsecase)

			reqBody, _ := json.Marshal(map[string]interface{}{
				"user_id": userID,
				"text":    tc.text,
			})

			if tc.text != "" && tc.mockErr != domain.ErrNotFound {
				answerUsecase.On("AddAnswerByID", mock.Anything, tc.questionID, userID, tc.text).Return(&domain.Answer{ID: 1, QuestionID: tc.questionID, UserID: userID, Text: tc.text}, tc.mockErr).Once()
			} else if tc.mockErr == domain.ErrNotFound {
				answerUsecase.On("AddAnswerByID", mock.Anything, tc.questionID, userID, tc.text).Return(nil, tc.mockErr).Once()
			}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /questions/{id}/answers/", handler.AddAnswer())

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/questions/%d/answers/", tc.questionID), bytes.NewReader(reqBody))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tc.expStatus, resp.StatusCode)
		})
	}
}

func TestDeleteHandler_DeleteAnswer(t *testing.T) {
	testCases := []struct {
		name      string
		answerID  int
		mockErr   error
		expStatus int
	}{
		{
			name:      "success",
			answerID:  1,
			mockErr:   nil,
			expStatus: http.StatusCreated,
		},
		{
			name:      "not found",
			answerID:  100,
			mockErr:   domain.ErrNotFound,
			expStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			answerUsecase := &mocks.AnswerUsecaseMock{}
			handler := handlers.NewAnswerHandler(answerUsecase)

			answerUsecase.On("DeleteAnswerByID", mock.Anything, tc.answerID).Return(tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /answers/{id}", handler.DeleteAnswer())

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/answers/%d", tc.answerID), nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			resp := w.Result()
			defer resp.Body.Close()
			assert.Equal(t, tc.expStatus, resp.StatusCode)
		})
	}
}
