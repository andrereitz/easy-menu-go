package models

type GenericReponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type DataReponse struct {
	Data any `json:"data"`
}
