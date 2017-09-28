package httputils

import "net/http"

type NopRoundTripper struct{}

func (rt NopRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
	}, nil
}
