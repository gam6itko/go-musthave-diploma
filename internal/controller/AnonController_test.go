package controller

import (
	"bytes"
	"context"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	jwt_mocks "github.com/gam6itko/go-musthave-diploma/internal/jwt/mock"
	mocks "github.com/gam6itko/go-musthave-diploma/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnonController_PostUserRegister(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Content-Type is not application/json", func(t *testing.T) {
		mIssuer := jwt_mocks.NewMockIIssuer(ctrl)

		mUserRepo := mocks.NewMockIUserRepository(ctrl)

		r := &http.Request{}
		w := &httptest.ResponseRecorder{}

		controller := NewAnonController(mIssuer, mUserRepo)
		controller.PostUserRegister(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("json parser fail", func(t *testing.T) {
		mIssuer := jwt_mocks.NewMockIIssuer(ctrl)

		mUserRepo := mocks.NewMockIUserRepository(ctrl)

		r := httptest.NewRequest("GET", "http://site.com", nil)
		r.Header.Set("Content-Type", "application/json")
		w := &httptest.ResponseRecorder{}

		controller := NewAnonController(mIssuer, mUserRepo)
		controller.PostUserRegister(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("bad credentials", func(t *testing.T) {
		mIssuer := jwt_mocks.NewMockIIssuer(ctrl)

		mUserRepo := mocks.NewMockIUserRepository(ctrl)

		json := `{}`
		r := httptest.NewRequest("GET", "http://site.com", bytes.NewBufferString(json))
		r.Header.Set("Content-Type", "application/json")
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}

		controller := NewAnonController(mIssuer, mUserRepo)
		controller.PostUserRegister(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Equal(t, "missing required credentials\n", w.Body.String())
	})

	t.Run("successful registration", func(t *testing.T) {
		mIssuer := jwt_mocks.NewMockIIssuer(ctrl)
		mIssuer.
			EXPECT().
			Issue(uint64(12)).
			Return("foo.bar.baz", nil).
			MaxTimes(1)

		mUserRepo := mocks.NewMockIUserRepository(ctrl)
		mUserRepo.
			EXPECT().
			FindByLogin(context.Background(), "username").
			Return(nil, nil).
			MaxTimes(1)

		mUserRepo.
			EXPECT().
			InsertNew(context.Background(), "username", gomock.Any()).
			Return(uint64(12), nil).
			MaxTimes(1)

		json := `{"login": "username", "password": "qwerty"}`
		r := httptest.NewRequest("GET", "http://site.com", bytes.NewBufferString(json))
		r.Header.Set("Content-Type", "application/json")
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}

		controller := NewAnonController(mIssuer, mUserRepo)
		controller.PostUserRegister(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "Bearer foo.bar.baz", w.Header().Get("Authorization"))
	})
}

func TestAnonController_PostUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("successful login", func(t *testing.T) {
		mIssuer := jwt_mocks.NewMockIIssuer(ctrl)
		mIssuer.
			EXPECT().
			Issue(uint64(12)).
			Return("foo.bar.baz", nil).
			MaxTimes(1)

		hash, err := bcrypt.GenerateFromPassword([]byte("qwerty"), bcrypt.MinCost)
		require.NoError(t, err)

		user := &diploma.User{
			ID:           12,
			PasswordHash: hash,
		}

		mUserRepo := mocks.NewMockIUserRepository(ctrl)
		mUserRepo.
			EXPECT().
			FindByLogin(context.Background(), "username").
			Return(user, nil).
			MaxTimes(1)

		mUserRepo.
			EXPECT().
			InsertNew(context.Background(), "username", gomock.Any()).
			Return(uint64(12), nil).
			MaxTimes(1)

		json := `{"login": "username", "password": "qwerty"}`
		r := httptest.NewRequest("GET", "http://site.com", bytes.NewBufferString(json))
		r.Header.Set("Content-Type", "application/json")
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}

		controller := NewAnonController(mIssuer, mUserRepo)
		controller.PostUserLogin(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "Bearer foo.bar.baz", w.Header().Get("Authorization"))
	})
}
