package helper

import "encoding/json"

type JobResponse struct {
	Code    int         `json:"code;"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type JobBody struct {
	QueueName string `json:"queue_name;"`
	Data      []byte `json:"data"`
}

func JobSuccess(data interface{}) *JobResponse {
	return &JobResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func JobFail(err error, data interface{}) *JobResponse {
	return &JobResponse{
		Code:    1,
		Message: err.Error(),
		Data:    data,
	}
}

func ParseJobBody(body []byte) (*JobBody, error) {
	var jobBody JobBody
	if err := json.Unmarshal(body, &jobBody); err != nil {
		return nil, err
	}
	return &jobBody, nil
}
