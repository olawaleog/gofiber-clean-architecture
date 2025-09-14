package impl

import (
	"context"

	"github.com/RizkiMufrizal/gofiber-clean-architecture/client"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

func NewHttpBinServiceImpl(httpBinClient *client.HttpClient) service.HttpService {
	return &httpBinServiceImpl{HttpClient: *httpBinClient}
}

type httpBinServiceImpl struct {
	client.HttpClient
}

//
//func (h httpBinServiceImpl) PostMethod(ctx context.Context, url, method string, body, response *map[string]interface{}) map[string]interface{} {
//	//TODO implement me
//	panic("implement me")
//}

func (h *httpBinServiceImpl) PostMethod(ctx context.Context, url string, method string, body *map[string]interface{}, header *map[string]interface{}, isForm bool) (map[string]interface{}, error) {

	response, err := h.HttpClient.Send(ctx, url, method, body, header, isForm)
	common.NewLogger().Info("log response service ", response)
	return response, err
}
