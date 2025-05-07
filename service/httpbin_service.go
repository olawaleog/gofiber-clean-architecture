package service

import "context"

type HttpService interface {
	PostMethod(ctx context.Context, url, method string, body *map[string]interface{}, header *map[string]interface{}, isForm bool) map[string]interface{}
}
