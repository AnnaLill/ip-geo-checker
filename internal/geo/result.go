package geo

// Result результат запроса к одному геолокационному API
type Result struct {
	APIName string // Название API для идентификации
	City    string
	Error   error
}
