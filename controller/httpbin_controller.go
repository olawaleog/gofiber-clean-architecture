package controller

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/service"
)

type HttpBinController struct {
	service.HttpService
}

func NewHttpBinController(httpBinService *service.HttpService) *HttpBinController {
	return &HttpBinController{HttpService: *httpBinService}
}
