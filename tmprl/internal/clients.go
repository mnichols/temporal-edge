package internal

import (
	"context"
	"github.com/mnichols/temporal-edge/tmprl/generated/tmprl/v1/tmprlv1connect"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewClient(ctx context.Context, httpClient Doer) (tmprlv1connect.TmprlServiceClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	
	client := tmprlv1connect.NewTmprlServiceClient(httpClient, "http://localhost:8081")
	return client, nil
}
