package models

type JSONErrorInfo struct {
	Code    int    `json:"code,ompitempty"`
	Message string `json:"message,ompitempty"`
}

type JSONError struct {
	Error JSONErrorInfo `json:"error"`
}

type JSONData struct {
	Data interface{} `json:"data"`
}
