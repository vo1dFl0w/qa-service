package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/vo1dFl0w/qa-service/internal/app/domain"
	"github.com/vo1dFl0w/qa-service/internal/mocks/usecase_mocks"
	"github.com/vo1dFl0w/qa-service/internal/transport/http_transport/handlers"
)

func TestQuestionHandler_GetQuestionsList(t *testing.T) {
	testCases := []struct {
		name       string
		mockResult []*domain.Question
		mockErr    error
		expStatus  int
		expCount   int
	}{
		{
			name:       "empty list",
			mockResult: []*domain.Question{},
			mockErr:    nil,
			expStatus:  http.StatusOK,
			expCount:   0,
		},
		{
			name: "non-empty list",
			mockResult: []*domain.Question{
				{ID: 1, Text: "q1", CreatedAt: time.Now()},
				{ID: 2, Text: "q2", CreatedAt: time.Now()},
			},
			mockErr:   nil,
			expStatus: http.StatusOK,
			expCount:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			questionUsecase := &mocks.QuestionUsecaseMock{}
			h := handlers.NewQuestionHandler(questionUsecase)

			questionUsecase.On("GetAllQuestions", mock.Anything).Return(tc.mockResult, tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /questions/", h.GetQuestionsList())

			req := httptest.NewRequest(http.MethodGet, "/questions/", nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expStatus, res.StatusCode)

			var got []*domain.Question
			err := json.NewDecoder(res.Body).Decode(&got)
			assert.NoError(t, err)
			assert.Len(t, got, tc.expCount)

			questionUsecase.AssertExpectations(t)
		})
	}
}

func TestQuestionHandler_CreateQuestion(t *testing.T) {
	testCases := []struct {
		name      string
		payload   map[string]interface{}
		mockErr   error
		expStatus int
	}{
		{
			name:      "success",
			payload:   map[string]interface{}{"text": "question_1"},
			mockErr:   nil,
			expStatus: http.StatusCreated,
		},
		{
			name:      "empty text",
			payload:   map[string]interface{}{"text": ""},
			mockErr:   nil,
			expStatus: http.StatusBadRequest,
		},
		{
			name:      "bad json",
			payload:   nil,
			mockErr:   nil,
			expStatus: http.StatusBadRequest,
		},
		{
			name:      "repo error",
			payload:   map[string]interface{}{"text": "question_1"},
			mockErr:   domain.ErrNoRowDeleted,
			expStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			questionUsecase := &mocks.QuestionUsecaseMock{}
			h := handlers.NewQuestionHandler(questionUsecase)

			var body []byte
			if tc.payload != nil {
				body, _ = json.Marshal(tc.payload)
			} else {
				body = []byte(`{invalid json`)
			}

			if tc.payload != nil && tc.payload["text"] != "" && tc.mockErr != nil {
				questionUsecase.On("CreateQuestion", mock.Anything, tc.payload["text"]).Return(nil, tc.mockErr)
			} else if tc.payload != nil && tc.payload["text"] != "" && tc.mockErr == nil {
				questionUsecase.On("CreateQuestion", mock.Anything, tc.payload["text"]).Return(&domain.Question{ID: 1, Text: "question_1"}, nil)
			}

			mux := http.NewServeMux()
			mux.HandleFunc("POST /questions/", h.CreateQuestion())

			req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(body))
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expStatus, res.StatusCode)
			questionUsecase.AssertExpectations(t)
		})
	}
}

func TestQuestionHandler_GetQuestionWithAnswers(t *testing.T) {
	testCases := []struct {
		name      string
		questionID int
		mockQ     *domain.Question
		mockAs    []*domain.Answer
		mockErr   error
		expStatus int
		expText   string
	}{
		{
			name:       "success",
			questionID: 1,
			mockQ:      &domain.Question{ID: 1, Text: "q1", CreatedAt: time.Now()},
			mockAs: []*domain.Answer{
				{ID: 1, QuestionID: 1, Text: "a1", UserID: uuid.New()},
			},
			mockErr:   nil,
			expStatus: http.StatusOK,
			expText:   "q1",
		},
		{
			name:       "not found",
			questionID: 100,
			mockQ:      nil,
			mockAs:     nil,
			mockErr:    domain.ErrNotFound,
			expStatus:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			questionUsecase := &mocks.QuestionUsecaseMock{}
			h := handlers.NewQuestionHandler(questionUsecase)

			questionUsecase.On("GetQuestionWithAnswers", mock.Anything, tc.questionID).Return(tc.mockQ, tc.mockAs, tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /questions/{id}", h.GetQuestionWithAnswers())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/questions/%d", tc.questionID), nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expStatus, res.StatusCode)
			if tc.mockQ != nil {
				var out struct {
					Question *domain.Question  `json:"question"`
					Answers  []*domain.Answer `json:"answers"`
				}
				err := json.NewDecoder(res.Body).Decode(&out)
				assert.NoError(t, err)
				assert.Equal(t, tc.expText, out.Question.Text)
				assert.Len(t, out.Answers, len(tc.mockAs))
			}

			questionUsecase.AssertExpectations(t)
		})
	}
}

func TestQuestionHandler_DeleteQuestion(t *testing.T) {
	testCases := []struct {
		name       string
		questionID int
		mockErr    error
		expStatus  int
	}{
		{
			name:       "success",
			questionID: 1,
			mockErr:    nil,
			expStatus:  http.StatusCreated,
		},
		{
			name:       "not found",
			questionID: 100,
			mockErr:    domain.ErrNotFound,
			expStatus:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			questionUsecase := &mocks.QuestionUsecaseMock{}
			h := handlers.NewQuestionHandler(questionUsecase)

			questionUsecase.On("DeleteQuestionByID", mock.Anything, tc.questionID).Return(tc.mockErr)

			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /questions/{id}", h.DeleteQuestion())

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/questions/%d", tc.questionID), nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tc.expStatus, res.StatusCode)
			questionUsecase.AssertExpectations(t)
		})
	}
}