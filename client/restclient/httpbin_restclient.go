package restclient

import (
	"context"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/client"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/common"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
)

func NewHttpRestClient(config configuration.Config) client.HttpClient {
	return &RestClient{config}
}

type RestClient struct {
	configuration.Config
}

func (h RestClient) Send(ctx context.Context, url string, method string, requestBody *map[string]interface{}, hd *map[string]interface{}, isForm bool) map[string]interface{} {
	var headers []common.HttpHeader
	for key, value := range *hd {
		headers = append(headers, common.HttpHeader{Key: key, Value: value.(string)})
	}
	response := make(map[string]interface{})
	//headers = append(headers, common.HttpHeader{Key: "Authorization", Value: "Bearer " + h.Config.Get("PAYSTACK_SECRET_KEY")})
	httpClient := common.ClientComponent[map[string]interface{}, map[string]interface{}]{
		HttpMethod:     method,
		UrlApi:         url,
		RequestBody:    requestBody,
		ResponseBody:   &response,
		Headers:        headers,
		ConnectTimeout: 60000,
		ActiveTimeout:  60000,
		IsFormData:     isForm,
	}
	res := httpClient.Execute(ctx)
	common.Logger.Error(res)
	//exception.PanicLogging(err)
	return *httpClient.ResponseBody
}
