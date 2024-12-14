package redirect_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YLysov0017/go_pet_n1/internal/http-server/handlers/url/redirect"
	"github.com/YLysov0017/go_pet_n1/internal/http-server/handlers/url/redirect/mocks"
	"github.com/YLysov0017/go_pet_n1/internal/lib/api"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/handlers/slogdiscard"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {

	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
		respCode  int
	}{
		{
			name:     "Success",
			alias:    "test_Alias_google",
			url:      "https://www.google.com",
			respCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}
			router := chi.NewRouter()
			router.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetMock))
			ts := httptest.NewServer(router)
			defer ts.Close()
			redirected, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			require.NoError(t, err)
			require.Equal(t, tc.url, redirected)
		})
	}
}
