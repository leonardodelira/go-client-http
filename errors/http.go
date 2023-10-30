package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ResponseError struct {
	StatusCode int
	Message    string `json:"message"`
	Code       int    `json:"code"`
}

type InvalidResource struct {
	ResponseError

	FailedFields []string `json:"failed_fields"`
}

type response struct {
	Message string              `json:"message"`
	Code    int                 `json:"code"`
	Extras  map[string][]string `json:"extras"`
}

func FromBadRequest(resp *http.Response) InvalidResource {
	var r response
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return InvalidResource{}
	}

	extras := []string{}
	if failedFields, ok := r.Extras["failed_fields"]; ok {
		extras = append(extras, failedFields...)
	}

	return InvalidResource{
		ResponseError: ResponseError{
			StatusCode: resp.StatusCode,
			Message:    r.Message,
			Code:       r.Code,
		},
		FailedFields: extras,
	}
}

func FromHTTPResponse(resp *http.Response) ResponseError {
	var r response
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &r); err != nil {
		return ResponseError{}
	}

	return ResponseError{
		StatusCode: resp.StatusCode,
		Message:    r.Message,
		Code:       r.Code,
	}
}

func (e InvalidResource) Error() string {
	return fmt.Sprintf("StatusCode=%d, Message=%s, FailedFields=%v", e.StatusCode, e.Message, e.FailedFields)
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("StatusCode=%d, Message=%s", e.StatusCode, e.Message)
}
