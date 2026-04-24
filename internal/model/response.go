package model

type ResponseWrapper struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Page[E any] struct {
	Content []E `json:"content"`
	Total   int `json:"totalElements"`
}