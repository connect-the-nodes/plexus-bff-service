package apiresponse

type Response[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func OK[T any](data T) Response[T] {
	return Response[T]{
		Success: true,
		Data:    data,
	}
}

func Error(message string) Response[any] {
	return Response[any]{
		Success: false,
		Message: message,
	}
}
