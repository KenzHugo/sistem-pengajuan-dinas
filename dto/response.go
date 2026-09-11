package dto

// Response adalah pembungkus standar semua respons API
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse untuk endpoint dengan pagination
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

// OK membuat respons sukses
func OK(message string, data interface{}) Response {
	return Response{Success: true, Message: message, Data: data}
}

// Fail membuat respons gagal
func Fail(message string) Response {
	return Response{Success: false, Message: message}
}
