package notionapi

type ErrorCode string

type Error struct {
	Object  ObjectType `json:"object"`
	Status  int        `json:"status"`
	Code    ErrorCode  `json:"code"`
	Message string     `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

type RateLimitedError struct {
	Message string
}

func (e *RateLimitedError) Error() string {
	return e.Message
}
