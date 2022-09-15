package MessageHelper

import (
	"ClipboardShare/SqlHelper"
	"ClipboardShare/AesHelper"
	"crypto/sha1"
	"bytes"
	"encoding/gob"
	"os"
	"io"
	"fmt"
)

// creating Message object from the input clipboard data
// Input - the <string> content of the message 
// Output - a SqlHelper.Message object 
func CreateMessage(data string) *SqlHelper.Message {

	message := new(SqlHelper.Message)
	h := sha1.New()

	message.Date = SqlHelper.GetTimestamp()
	message.Message = data
	message.Mode = 1 // sent mode
	message.Device, _ = os.Hostname()

	// calculate checksum for the message
	io.WriteString(h, fmt.Sprintf("%d", message.Date))
	io.WriteString(h, message.Message)
	io.WriteString(h, message.Device)
	message.Checksum = string(h.Sum(nil))

	return message
}

// Converting SqlHelper.Message object to Encrypted byte array
// encryption key is stored in the local database
// input - SqlHelper.Message object
// output - AES encrypted bytes array
func SerializeMessage(message *SqlHelper.Message) []byte {
	var serializedMessage bytes.Buffer
	encryptionKey := SqlHelper.QueryConfigByKey("encryptionKey").Value

	enc := gob.NewEncoder(&serializedMessage)
	_ = enc.Encode(message)

	serializedData, _ := AesHelper.Encrypt([]byte(encryptionKey), serializedMessage.Bytes())
	return serializedData
}

// Converting Encrypted byte array to SqlHelper.Message object
// encryption key is stored in the local database
// input - AES encrypted bytes array
// output - SqlHelper.Message object
func DeserializeMessage(serializedMessage []byte) (*SqlHelper.Message, int) {
	var message *SqlHelper.Message
	encryptionKey := SqlHelper.QueryConfigByKey("encryptionKey").Value

	serializedMessage, _ = AesHelper.Decrypt([]byte(encryptionKey), serializedMessage)

	reader := bytes.NewReader(serializedMessage)
	dec := gob.NewDecoder(reader)
	_ = dec.Decode(&message)

	message.Mode = 0

	// validate checksum
	h := sha1.New()

	io.WriteString(h, fmt.Sprintf("%d", message.Date))
	io.WriteString(h, message.Message)
	io.WriteString(h, message.Device)

	if ( string(h.Sum(nil)) != message.Checksum ) {
		// in case the message checksum is wrong - the data is diffrent from the original one
		return nil, -1
	}

	return message, 0 
}
