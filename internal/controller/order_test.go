package controller

import (
	"bytes"
	"context"
	mocks2 "github.com/gam6itko/go-musthave-diploma/internal/accrual/mock"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	mocks3 "github.com/gam6itko/go-musthave-diploma/internal/jwt/mock"
	mocks "github.com/gam6itko/go-musthave-diploma/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderController_PostUserOrders(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("successful get orders", func(t *testing.T) {
		mJwtParser := mocks3.NewMockIParser(ctrl)
		mJwtParser.
			EXPECT().
			Parse("foo.bar.baz").
			Return(uint64(12), nil).
			MaxTimes(1)

		mAccural := mocks2.NewMockIClient(ctrl)
		mAccural.
			EXPECT().
			Get(uint64(1917)).
			Return(
				&diploma.Accrual{
					OrderNumber: "1917",
					Accrual:     100,
					Status:      "PROCESSED",
				},
				nil,
			).
			MaxTimes(1)

		mOrderRepo := mocks.NewMockIOrderRepository(ctrl)
		mOrderRepo.
			EXPECT().
			FindByID(context.Background(), uint64(1917)).
			Return(nil, nil).
			MaxTimes(1)
		mOrderRepo.
			EXPECT().
			InsertNew(
				gomock.Any(),
				&diploma.Order{
					ID:      1917,
					UserID:  12,
					Accrual: 100,
					Status:  diploma.StatusProcessed,
				},
			).
			Return(nil).
			MaxTimes(1)

		body := bytes.NewBufferString("1917")
		r := httptest.NewRequest("GET", "http://site.com", body)
		r.Header.Set("Content-Type", "text/plain")
		r.Header.Set("Authorization", "Bearer foo.bar.baz")
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}

		controller := NewOrderController(mJwtParser, mAccural, mOrderRepo)
		controller.PostUserOrders(w, r)

		require.Equal(t, http.StatusAccepted, w.Code)
	})
}
