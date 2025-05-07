package model

type GeneralResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(message string, data interface{}) *GeneralResponse {
	r := new(GeneralResponse)
	r.Success = true
	r.Code = 200
	r.Message = message
	r.Data = data
	return r
}

func FailedResponse(code int, message string, data interface{}) *GeneralResponse {
	r := new(GeneralResponse)
	r.Success = false
	r.Code = code
	r.Message = message
	r.Data = data
	return r
}
