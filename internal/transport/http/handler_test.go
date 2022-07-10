package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"news-app/internal/domain"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handler_GetArticles(t *testing.T) {
	var (
		someFeedURL     = "https://some-feed-url"
		someTitle       = "some-title"
		someDescription = "some-description"
		someConent      = "some-content"
		someImage       = domain.Image{
			URL:   someFeedURL,
			Title: someTitle,
		}
		someArticles = []domain.Article{
			{
				Title:       someTitle,
				Description: someDescription,
				Content:     someConent,
				Image:       someImage,
			},
		}
	)

	t.Run("should return articles if service is successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockService(ctrl)
		handler := NewHandler(mockService)

		mockService.EXPECT().GetArticles(gomock.Any(), someFeedURL).Return(someArticles, nil)

		body := []byte(`{"feed_url":"https://some-feed-url"}`)
		req, err := http.NewRequest(http.MethodGet, getArticlesByFeed, bytes.NewReader(body))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		handler.GetArticles(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)

		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		defer res.Body.Close()

		var articles []domain.Article
		err = json.Unmarshal(bytes, &articles)
		require.NoError(t, err)

		assert.Equal(t, someArticles, articles)
	})

	t.Run("should return a internal server error if service layer fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockService(ctrl)
		handler := NewHandler(mockService)

		mockService.EXPECT().GetArticles(gomock.Any(), someFeedURL).Return([]domain.Article{}, assert.AnError)

		body := []byte(`{"feed_url":"https://some-feed-url"}`)
		req, err := http.NewRequest(http.MethodGet, getArticlesByFeed, bytes.NewReader(body))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		handler.GetArticles(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should return a bad request if json if request body is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockService(ctrl)
		handler := NewHandler(mockService)

		body := []byte(`{"invalid"}`)
		req, err := http.NewRequest(http.MethodGet, getArticlesByFeed, bytes.NewReader(body))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		handler.GetArticles(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return a bad request if url is missing from body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := NewMockService(ctrl)
		handler := NewHandler(mockService)

		body := []byte(`{"other-data":"some-other-data"}`)
		req, err := http.NewRequest(http.MethodGet, getArticlesByFeed, bytes.NewReader(body))
		require.NoError(t, err)

		w := httptest.NewRecorder()
		handler.GetArticles(w, req)

		res := w.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
