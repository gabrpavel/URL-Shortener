package delete_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"URL-Shortener/internal/http-server/handlers/url/delete"
	"URL-Shortener/internal/http-server/handlers/url/delete/mocks"
	"URL-Shortener/internal/lib/logger/handlers/slogdiscard"
	"URL-Shortener/internal/storage"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		mockError error
		respCode  int
		respError string
	}{
		{
			name:     "Success",
			alias:    "test_alias",
			respCode: http.StatusOK,
		},
		{
			name:      "Alias not found",
			alias:     "missing_alias",
			mockError: storage.ErrURLNotFound,
			respCode:  http.StatusNotFound,
			respError: "not found",
		},
		{
			name:      "DeleteURL internal error",
			alias:     "test_alias",
			mockError: errors.New("unexpected error"),
			respCode:  http.StatusInternalServerError,
			respError: "internal error",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respCode:  http.StatusBadRequest,
			respError: "invalid request",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.alias != "" && tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
			} else if tc.alias != "" {
				urlDeleterMock.On("DeleteURL", tc.alias).Return(nil).Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock)

			req, err := http.NewRequest(http.MethodDelete, "/delete/"+tc.alias, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.respCode, rr.Code)

			body := rr.Body.String()
			if tc.respError != "" {
				require.Contains(t, body, tc.respError)
			} else {
				require.NotContains(t, body, "error")
			}
		})
	}
}
