package handler

const (
	RESPONSE_INVALID_ID = "Invalid ID"
	RESPONSE_DUPLICATED = "Duplicated Name"
)

type Data struct {
	Status int `json:"status,omitempty"`
}

type Response struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

func MakeResponse(kind string, model string) Response {
	var data Data
	var resp Response

	switch kind {
	case RESPONSE_INVALID_ID:
		data.Status = 404
		resp.Code = "rest_" + model + "_invalid_id"
		resp.Message = "Invalid ID."
	case RESPONSE_DUPLICATED:
		data.Status = 409
		resp.Code = "rest_" + model + "_duplicated_name"
		resp.Message = "Duplicated Name."
	}

	resp.Data = data

	return resp
}

// import model "example.com/logos106/saroop-api/models"

// JsonError is a generic error in JSON format
//
// swagger:response jsonError
// type jsonError struct {
// 	// in: body
// 	Message string `json:"message"`
// }
