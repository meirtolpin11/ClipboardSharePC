# Clipboard Share 

Welcome to "Clipboard Share" project. We all know that moment when we want to copy some link from one computer to another, and we send it by mail/telegram/whatsapp or even copying it on a piece of paper - as we all know, it's not the best way to do so. 

So I decided to solve this problem by writing a clipboard share cross platform application. And because I know that we all have trust issues, and no one is going to share with me his clipboard, I decided to go for a serverless structure.

## How it's working ?
The Idea is simple - 
* We start to monitor clipboard changes
* When there is any change to the clipboard - the data is encrypted and sent to the rest of the clients (*your clients only*) using `Google Cloud Pub/Sub` ~(We all trust Google, so it's okey)~
* Another clients are receiving the data, decrypted it, and updating the clipboard.

## Requirements 
* Operation System - for now only `Linux` and `Windows` systems are supported, someday it will go to mobile too.
* Google Cloud Pub/Sub Service account - you need to create your personal Google Cloud account and get Pub/Sub API key ([Tutorial](https://youtu.be/rWcLDax-VmM))

## So, How to use ?

* after getting the google api json creds file, firstly run the following command (it will initialize the creds)

  ``` ClipboardShare.exe -pubsubJson creds.json ```
* next, when the database is created - just run the application (you can also change the encryption key, check ```--help```)

## HELP 
```
Usage of ClipboardShare.exe:
  -enc string
        Encryption key for the messages (default "somemagichere123")
  -modifyClipboard
        Modify clipboard when recieving messages (default true)
  -pubsubJson string
        Pub Sub creds json
  -sendClipboard
        Send clipboard data when it's changed (default true)
```