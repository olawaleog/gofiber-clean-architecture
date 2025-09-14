package client

import (
	"context"
)

type HttpClient interface {
	Send(ctx context.Context, url string, method string, requestBody *map[string]interface{}, response *map[string]interface{}, isForm bool) (map[string]interface{}, error)
}
