package qorexGoHelpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int
	Status string
	Data []interface{}
	Rows int
	Error error
}

func (r *Response) SetStatus(statusCode int, statusText string, err error) {
	r.Code = statusCode
	r.Error = err
	r.Status = statusText
	return
}

func (r *Response) AppendData(entityList []interface{}) {
	r.Data = append(r.Data, entityList...)
	r.Rows = len(r.Data)
}

func (r *Response) GetJson() []byte {
	response, err := json.Marshal(r)
	if err != nil {
		r.SetStatus(http.StatusBadRequest, "Error converting results to JSON", err)
	}
	return response
}