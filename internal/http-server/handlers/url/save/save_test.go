package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YLysov0017/go_pet_n1/internal/http-server/handlers/url/save"
	"github.com/YLysov0017/go_pet_n1/internal/http-server/handlers/url/save/mocks"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			url:   "https://google.com",
			alias: "test_alias",
		},
		{
			name:  "Empty alias",
			url:   "https://google.com",
			alias: "",
		},
		{
			name:      "Empty URL",
			url:       "",
			alias:     "some_alias",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			alias:     "some_alias",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			url:       "https://google.com",
			alias:     "test_alias",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc // Исключение проблем с параллельным запуском тестов

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Параллельный запуск тестов

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")). // Мок будет ожидать конкретного url
													Return(int64(1), tc.mockError). // Вернуть случайный id и ошибку
													Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock, 6)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
