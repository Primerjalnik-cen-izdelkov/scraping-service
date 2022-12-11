package models

type JSONErrorInfo struct {
    Code    int    `json:"code,ompitempty" example:"404"`
    Message string `json:"message,ompitempty" example:"not found"`
}

type JSONError struct {
    Error JSONErrorInfo `json:"error"`
}

type JSONData struct {
    Data interface{} `json:"data"`
}
