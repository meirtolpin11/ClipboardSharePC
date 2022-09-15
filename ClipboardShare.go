package main

import (
	"ClipboardShare/ClipboardHelper"
	"ClipboardShare/PubsubHelper"
	"ClipboardShare/MessageHelper"
	"ClipboardShare/SqlHelper"
	"flag"
	"fmt"
	"os"
	"log"
)

func ClipboardReciever(modifyClipboard bool) {
	// creating messages channel, all the new messages will be added to here
	messages := make(chan []byte)

	log.Println("Subscribing to the PubSub Topic")
	go PubsubHelper.DefaultSubscribe(messages)

	log.Println("Starting PubSub listener")
	for data := range messages {
		deserializedMessage, err := MessageHelper.DeserializeMessage(data)
		
		if err != 0 {
			// something with this message gone wrong - ignoring it
			// for example - wrong checksum
			continue
		}

		// some debug
		log.Println(fmt.Sprintf("Received: %s, From Device: %s, Date: %d", 
			deserializedMessage.Message, deserializedMessage.Device, deserializedMessage.Date))

		// update the computer clipboard to the received data
		if modifyClipboard {
			ClipboardHelper.SetClipboard(deserializedMessage.Message)
		}
		
	}
}

func ClipboardSender(sendClipboard bool) {

	log.Println("Getting Clipboard Listener")
	clipChan := ClipboardHelper.GetClipboardChannel()

	log.Println("Starting Clipboard Listener")
	for data := range clipChan {
		
		if ClipboardHelper.ClipboardSet {
			ClipboardHelper.ClipboardSet = false
			continue
		}

		// public the new clipboard content to the Google Cloud PubSub
		if (sendClipboard) {
			log.Println("Sending " + string(data))
			PubsubHelper.DefaultPublish(MessageHelper.SerializeMessage(MessageHelper.CreateMessage(string(data))))
		}
	}
}

func main() {

	encryptionKey :=  flag.String("enc", "somemagichere123", "Encryption key for the messages")
	modifyClipboard := flag.Bool("modifyClipboard", true, "Modify clipboard when recieving messages")
	sendClipboard := flag.Bool("sendClipboard", true, "Send clipboard data when it's changed")
	pubsubJson := flag.String("pubsubJson", "", "Pub Sub creds json")
	flag.Parse()

	if *pubsubJson != "" {
		PubsubHelper.SetPubSubConfig(*pubsubJson)
		os.Exit(0)
	}

	// insert encryption key to the database
	SqlHelper.InsertConfig(SqlHelper.Config{Key: "encryptionKey", Value: *encryptionKey})

	go ClipboardReciever(*modifyClipboard)	
	ClipboardSender(*sendClipboard)
}