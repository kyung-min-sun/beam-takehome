package main

import (
	"log"
	"time"

	client "slai.io/takehome/pkg/client"
)

func main() {
	log.Println("Starting client...")

	c, err := client.NewClient("./")
	if err != nil {
		log.Fatal(err)
	}


	someMessage := "hello there"
	lastChecked := time.Now()
	for {

		log.Printf("Sending: '%s'", someMessage)

		echoValue, echoErr := c.Echo(someMessage)
		files, filesErr := c.FileWatch(lastChecked)
		lastChecked = time.Now()

		if echoErr != nil || filesErr != nil {
			log.Fatal("Unable to send request.")
		}

		log.Printf("Received: '%s'", echoValue)
		log.Printf("Received: '%s'", files)

		time.Sleep(time.Second)
	}

}
