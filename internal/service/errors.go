package service

import "errors"

// Ошибки сервиса
var (
	ErrEmptyIP         = errors.New("IP address is required")
	ErrInvalidIPFormat = errors.New("Invalid IP address format")
	ErrNoDataFromAPIs  = errors.New("Unable to get city information from any API")
)
