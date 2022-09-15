package ClipboardHelper 

import (
	"golang.design/x/clipboard"	
	"context"
) 

var ctxCancelFunc context.CancelFunc
var ClipboardSet bool

func init() {
	Initialize()
}

/* Check if the module is ready to work */
func Initialize() {

	err := clipboard.Init()
	if err != nil {
	      panic(err)
	}
}

/* Read string from the clipboard */
func ReadClipboard() string {
	return string(clipboard.Read(clipboard.FmtText))
}

/* Set clipboard value */
func SetClipboard(value string) {
	clipboard.Write(clipboard.FmtText, []byte(value))
	ClipboardSet = true
}

/* Get a channel of clipboard values */
func GetClipboardChannel() <-chan []byte {

	ctx, cancel := context.WithCancel(context.Background())
	ctxCancelFunc = cancel

	ch := clipboard.Watch(ctx, clipboard.FmtText)
	return ch
}

/* Stop the "watch" of clipboard channel */
func CancelClipboardWatch() {
	ctxCancelFunc()
}


func main() {	
	ch := GetClipboardChannel()

	for data := range ch {	      
	      println(string(data))
	}    
}
