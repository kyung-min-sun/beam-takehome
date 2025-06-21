package server

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"slai.io/takehome/pkg/common"
)

func HandleEcho(msg []byte, client *Client) error {
	log.Println("Received ECHO request.")

	var request common.EchoRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	response := &common.EchoResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: request.Value,
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}

func HandleFileWatch(msg []byte, client *Client, directory string) error {
	log.Println("Received FILE_WATCH request.")

	var request common.FileWatchRequest
	err := json.Unmarshal(msg, &request)

	if err != nil {
		log.Fatal("Invalid echo request.")
	}

	files := request.Files

	// Write files to filesystem
	for _, file := range files {
		// Create full file path
		fullPath := directory + "/" + file.Path 
		
		// Create directory if it doesn't exist
		dirPath := filepath.Dir(fullPath)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Printf("Error creating directory %s: %v", dirPath, err)
			continue
		}
		
		// Write file data
		err = os.WriteFile(fullPath, file.Data, 0644)
		if err != nil {
			log.Printf("Error writing file %s: %v", fullPath, err)
			continue
		}
		
		log.Printf("Successfully wrote file: %s (size: %d bytes)", fullPath, file.Size)
	}

	response := &common.FileWatchResponse{
		BaseResponse: common.BaseResponse{
			RequestId:   request.RequestId,
			RequestType: request.RequestType,
		},
		Value: "Received file watch request.",
	}

	responsePayload, err := json.Marshal(response)
	if err != nil {
		return err
	}

	err = client.ws.WriteMessage(websocket.TextMessage, responsePayload)
	if err != nil {
		return err
	}

	return nil
}

