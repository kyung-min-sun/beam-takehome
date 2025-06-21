package server

import (
	"encoding/base64"
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
	
		if file.Deleted {
			// Remove the file or directory
			log.Printf("Removing file: %s", fullPath)
			err := os.RemoveAll(fullPath)
			if err != nil {
				log.Printf("Error removing %s: %v", fullPath, err)
				continue
			}
			
			// Clean up empty parent directories
			cleanupEmptyDirectories(directory, file.Path)
			continue
		}
	
		// Create directory if it doesn't exist
		dirPath := filepath.Dir(fullPath)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Printf("Error creating directory %s: %v", dirPath, err)
			continue
		}
		
		// Write file data
		data, err := base64.StdEncoding.DecodeString(file.Base64)
		if err != nil {
			log.Printf("Error decoding base64 string %s: %v", file.Base64, err)
			continue
		}
		err = os.WriteFile(fullPath, data, 0644)
		if err != nil {
			log.Printf("Error writing file %s: %v", fullPath, err)
			continue
		}
		
		log.Printf("Successfully wrote file: %s", fullPath)
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

// cleanupEmptyDirectories removes empty directories up the directory tree
func cleanupEmptyDirectories(baseDir, filePath string) {
	// Get the directory path of the deleted file
	dirPath := filepath.Dir(filePath)
	
	// Start from the immediate parent directory and work our way up
	for dirPath != "." && dirPath != "/" {
		fullDirPath := filepath.Join(baseDir, dirPath)
		
		// Check if directory is empty
		entries, err := os.ReadDir(fullDirPath)
		if err != nil {
			log.Printf("Error reading directory %s: %v", fullDirPath, err)
			break
		}
		
		// If directory is not empty, stop cleaning up
		if len(entries) > 0 {
			break
		}
		
		// Remove empty directory
		err = os.Remove(fullDirPath)
		if err != nil {
			log.Printf("Error removing empty directory %s: %v", fullDirPath, err)
			break
		}
		
		log.Printf("Removed empty directory: %s", fullDirPath)
		
		// Move up to parent directory
		dirPath = filepath.Dir(dirPath)
	}
}

