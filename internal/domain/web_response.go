package domain

type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omnitempty"`
	Errors    interface{} `json:"erros,omnitempty"`
}
