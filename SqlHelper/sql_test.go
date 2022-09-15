package SqlHelper

import (
	"testing"
	"reflect"
)


func TestMessage(t *testing.T) {
	InitializeDB()

	var message Message
	message.date = GetTimestamp()
	message.device = "test"
	message.checksum = "test"
	message.mode = 1
	message.message = "test message"

	/* Check insert new message functionality */
	InsertNewMessage(message)
	messageFromDB := QueryLastMessage()

	if !reflect.DeepEqual(message, messageFromDB) {
		 t.Errorf("Message is incorrect, got: %v, sent: %v", messageFromDB, message)
	}

	/* Check delete message functionality */
	DeleteMessageByTimestamp(message.date)
	messageFromDB = QueryLastMessage()

	if reflect.DeepEqual(message, messageFromDB) {
		 t.Errorf("Unable to delete the test message, last message: %v, test message: %v", messageFromDB, message)
	}
}

func TestConfig(t *testing.T) {
	InitializeDB()

	var config Config
	config.key = "testKey"
	config.value = "testValue"

	// delete old test config if exist
	DeleteConfigByKey(config.key)

	// check that no test config exist 
	configFromDB := QueryConfigByKey(config.key)

	if reflect.DeepEqual(config, configFromDB) { 
		t.Error("Not able to delete config from the Database")
	}

	InsertConfig(config)
	configFromDB = QueryConfigByKey(config.key)

	if !reflect.DeepEqual(config, configFromDB) {
		t.Error("Unable to add new config inside the DB")
	}

	// delete old test config if exist
	DeleteConfigByKey(config.key)

	// check that no test config exist 
	configFromDB = QueryConfigByKey(config.key)

	if reflect.DeepEqual(config, configFromDB) { 
		t.Error("Not able to delete config from the Database")
	}

}
