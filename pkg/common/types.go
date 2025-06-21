package common

import "os"

type RequestType string

const (
	Echo RequestType = "ECHO"
	FileWatch RequestType = "FILE_WATCH"
)

type BaseRequest struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type BaseResponse struct {
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
}

type EchoRequest struct {
	BaseRequest
	Value string
}

type EchoResponse struct {
	BaseResponse
	Value string
}

type FileWatchRequest struct {
	BaseRequest
	Files []FileWatchPayload
}

type FileWatchResponse struct {
	BaseResponse
	Value string
}

type FileWatchPayload struct {
	Path     string
	Base64   string
	Deleted  bool
}

type FileWatchInfo struct {
	os.FileInfo
	Path string
}