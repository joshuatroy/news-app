//go:generate mockgen -package=http -destination=./handler_mock_test.go . Service
package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"news-app/internal/domain"
)

const getArticlesByFeed = "/articles/feed"

// Service is an interface for our service layer for processing domain requests
type Service interface {
	GetArticles(context.Context, string) ([]domain.Article, error)
}

// handler is our internal representation of a http handler
type handler struct {
	service Service
}

// NewHandler is a constructor for a http handler
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type getArticlesRequest struct {
	FeedURL string `json:"feed_url,omitempty" validate:"required"`
}

func (h handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	var request getArticlesRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if err := validator.New().Struct(request); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	articles, err := h.service.GetArticles(r.Context(), request.FeedURL)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	h.writeSuccessResponse(w, articles)
}

func (h handler) writeSuccessResponse(w http.ResponseWriter, i interface{}) {
	body, _ := json.Marshal(i)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h handler) writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	response := struct {
		ErrorString string `json:"error"`
	}{
		ErrorString: err.Error(),
	}
	body, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}