package internal

type GcApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewGcResponse(d interface{}, e error) *GcApiResponse {

	var success = true
	var msg = ""

	if e != nil {
		success = false
		msg = e.Error()
	}

	return &GcApiResponse{
		Success: success,
		Data:    d,
		Error:   msg,
	}
}
