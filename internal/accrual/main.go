package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gam6itko/go-musthave-diploma/internal/diploma"
	"net/http"
	"strings"
)

type IClient interface {
	Get(orderID uint64) (acc *diploma.Accrual, err error)
}

type Client struct {
	httpClient *http.Client
	host       string
}

func NewAccrualClient(httpClient *http.Client, host string) *Client {
	host = strings.TrimRight(host, "/")
	if !strings.Contains(host, "://") {
		host = fmt.Sprintf("http://%s", host)
	}

	return &Client{
		httpClient,
		host,
	}
}

func (ths Client) Get(orderID uint64) (acc *diploma.Accrual, err error) {
	resp, err := ths.httpClient.Get(
		fmt.Sprintf("%s/api/orders/%d", ths.host, orderID),
	)
	if err != nil {
		return
	}

	if resp.StatusCode == http.StatusNoContent {
		err = errors.New("unregistered")
		return
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		err = errors.New("too many requests")
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		err = errors.New("internal server error")
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status %d", resp.StatusCode)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	acc = new(diploma.Accrual)
	err = decoder.Decode(acc)
	if err != nil {
		return
	}

	return
}
