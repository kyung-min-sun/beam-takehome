package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

const maxConnectionAttempts = 100
const hostURL = "ws://localhost:5555/"

func init() {
}

type Client struct {
	Directory string
	SessionId string
	ws        *websocket.Conn
	connected bool
	hostURL   string
	channels  map[string]chan []byte
}

func NewClient(directory string) (*Client, error) {
	var client *Client = &Client{
		Directory: directory,
		hostURL:   hostURL,
	}

	err := client.connect()
	if err != nil {
		return nil, err
	}

	client.connected = true
	client.channels = make(map[string]chan []byte)

	return client, nil
}

func (c *Client) connect() error {
	connected := false
	attempts := 0

	for {
		log.Println("Connection attempt: ", attempts)

		if attempts > maxConnectionAttempts {
			break
		}

		ws, _, err := websocket.DefaultDialer.Dial(c.hostURL, nil)
		c.ws = ws

		if err != nil {
			attempts++
			continue
		}

		connected = true
		break
	}

	// We weren't able to connect to the host, bail
	if !connected {
		return nil
	}

	// Start receiving messages
	go c.rx()

	return nil
}

func (c *Client) rx() {
	for {
		_, message, err := c.ws.ReadMessage()
		if ce, ok := err.(*websocket.CloseError); ok {

			switch ce.Code {
			case websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived,
				websocket.CloseAbnormalClosure:
				return
			}
		}

		var msg common.BaseResponse

		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		} else {
			if _, ok := c.channels[msg.RequestId]; ok {
				c.channels[msg.RequestId] <- message
			} else {
				log.Println("channel not found")
			}
		}
	}
}

func (c *Client) tx(msg []byte) error {
	err := c.ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}

// Request implementations
func (r *Client) Echo(value string) (string, error) {
	requestId := uuid.NewString()

	var request *common.EchoRequest = &common.EchoRequest{
		BaseRequest: common.BaseRequest{
			RequestId:   requestId,
			RequestType: string(common.Echo),
		},
		Value: value,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return "", err
	}

	var response common.EchoResponse = common.EchoResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle echo response: ", err)
		return "", err
	}

	return response.Value, err
}

func findFile(files []common.FileWatchInfo, path string) (common.FileWatchInfo, bool) {
	for _, file := range files {
			if file.Path == path {
					return file, true
			}
	}
	return common.FileWatchInfo{}, false
}

func scanDirectory(root string) ([]common.FileWatchInfo, error) {
	var files []common.FileWatchInfo

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
				return err
		}
		
		// Skip directories
		if info.IsDir() {
				return nil
		}

		files = append(files, common.FileWatchInfo{
			FileInfo: info,
			Path: path,
		})
		
		return nil
	})

	return files, err
}

func getFileWatchPayload(root string, path string, info common.FileWatchInfo) *common.FileWatchPayload {
			// Read file data
		data, err := os.ReadFile(path)
		if err != nil {
				fmt.Printf("Error reading file %s: %v\n", path, err)
				return nil
		}
		
		// Get relative path
		relPath, err := filepath.Rel(root, path)
		if err != nil {
				relPath = path
		}

		fileInfo := common.FileWatchPayload{
				Path:   relPath,
				Base64: base64.StdEncoding.EncodeToString(data),
				Deleted: false,
		}

	return &fileInfo
}

func (r *Client) FileWatch(lastChecked time.Time, existingFiles *[]common.FileWatchInfo) ([]common.FileWatchInfo, error) {
	if existingFiles == nil {
		existingFiles = &[]common.FileWatchInfo{}
	}

	requestId := uuid.NewString()
	newFiles, err := scanDirectory(r.Directory)
	if err != nil {
		return newFiles, err
	}

	requestFiles := []common.FileWatchPayload{}

	for _, file := range newFiles {
		if file.ModTime().After(lastChecked) {
			payload := getFileWatchPayload(r.Directory, file.Path, file)
			if payload != nil {
				requestFiles = append(requestFiles, *payload)
			}
		}
	}

	for _, file := range *existingFiles {
		_, exists := findFile(newFiles, file.Path);
		if !exists {
			payload := common.FileWatchPayload{
				Path: file.Path,
				Deleted: true,
			}
			requestFiles = append(requestFiles, payload)
		}
	}

	fmt.Printf("requestFiles: %v\n", requestFiles)
	var request *common.FileWatchRequest = &common.FileWatchRequest{
		BaseRequest: common.BaseRequest{
			RequestId:   requestId,
			RequestType: string(common.FileWatch),
		},
		Files: requestFiles,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return newFiles, err
	}

	r.channels[requestId] = make(chan []byte)

	err = r.tx(payload)
	if err != nil {
		return newFiles, err
	}

	var response common.FileWatchResponse = common.FileWatchResponse{}

	msg := <-r.channels[requestId]
	err = json.Unmarshal(msg, &response)
	if err != nil {
		log.Println("Unable to handle file watch response: ", err)
		return newFiles, err
	}

	return newFiles, err
}
