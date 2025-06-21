package common

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
	Files []FileWatchInfo
}

type FileWatchResponse struct {
	BaseResponse
	Value string
}

type FileWatchInfo struct {
	Path     string
	Name     string
	Size     int64
	Data     []byte
	Base64   string
}