package controller

import (
	"bytes"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	mocks3 "github.com/gam6itko/go-musthave-diploma/internal/jwt/mock"
	mocks "github.com/gam6itko/go-musthave-diploma/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserController_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("successful get balance", func(t *testing.T) {
		mJwtParser := mocks3.NewMockIParser(ctrl)
		mJwtParser.
			EXPECT().
			Parse("foo.bar.baz").
			Return(uint64(12), nil).
			MaxTimes(1)

		mUserRepo := mocks.NewMockIUserRepository(ctrl)
		mUserRepo.
			EXPECT().
			FindByID(gomock.Any(), uint64(12)).
			Return(
				&diploma.User{
					ID:              12,
					BalanceCurrent:  100,
					BalanceWithdraw: 20,
				},
				nil,
			).
			MaxTimes(1)

		body := bytes.NewBufferString("1917")
		r := httptest.NewRequest("GET", "http://site.com", body)
		r.Header.Set("Content-Type", "text/plain")
		r.Header.Set("Authorization", "Bearer foo.bar.baz")
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}

		controller := NewUserController(mJwtParser, mUserRepo)
		controller.GetUserBalance(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.JSONEq(t, `{"current": 100, "withdrawn": 20}`, w.Body.String())
	})
}
