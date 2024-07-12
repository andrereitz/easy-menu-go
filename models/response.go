package models

type GenericReponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ImageAddResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    int64  `json:"data"`
}

type DataReponse struct {
	Data any `json:"data"`
}
