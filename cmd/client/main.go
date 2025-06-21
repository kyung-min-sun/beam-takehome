package main

import (
	"log"
	"time"

	client "slai.io/takehome/pkg/client"
	"slai.io/takehome/pkg/common"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./client-files")
	if err != nil {
		log.Fatal(err)
	}


	someMessage := "hello there"
	lastChecked := time.Now()
	existingFiles := []common.FileWatchInfo{}
	for {

		log.Printf("Sending: '%s'", someMessage)

		echoValue, echoErr := c.Echo(someMessage)
		files, filesErr := c.FileWatch(lastChecked, &existingFiles)
		existingFiles = files
		lastChecked = time.Now()

		if echoErr != nil || filesErr != nil {
			log.Fatal("Unable to send request.")
		}

		log.Printf("Received: '%s'", echoValue)
		log.Printf("Received: '%v'", files)

		time.Sleep(time.Second)
	}

}
