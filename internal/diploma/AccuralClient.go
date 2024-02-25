package diploma

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AccuralClient struct {
	httpClient *http.Client
	host       string
}

func NewAccuralClient(httpClient *http.Client, host string) *AccuralClient {
	host = strings.TrimRight(host, "/")
	if !strings.Contains(host, "://") {
		host = fmt.Sprintf("http://%s", host)
	}

	return &AccuralClient{
		httpClient,
		host,
	}
}

func (ths AccuralClient) Get(orderID uint64) (acc *Accural, err error) {
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
	acc = new(Accural)
	err = decoder.Decode(acc)
	if err != nil {
		return
	}

	return
}
