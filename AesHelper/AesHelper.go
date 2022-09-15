/*
*	FILE			: AES_Example.go
*	PROJECT			: INFO-1340 - Block Ciphers
*	PROGRAMMER		: Daniel Pieczewski, ref: https://github.com/mickelsonm
*	FIRST VERSION		: 2020-04-12
*	DESCRIPTION		:
*		The function(s) in this file make up example code for encryption and decryption of a block of text
*		using the Golang standard library AES implementation using the Cipher Feedback mode of encryption (CFB). 
*		DISCLAIMER: There is no way that this a secure implementation of AES. This is only for my personal learning.
*		So help you God if this ends up in some commercial application.
 */
 package AesHelper

 import (	 
	 "crypto/aes"
	 "crypto/cipher"
	 "crypto/rand"
	 "encoding/base64"
	 "errors"
	 "io"
	
 )
 
 /*
  *	FUNCTION		: encrypt
  *	DESCRIPTION		:
  *		This function takes a string and a cipher key and uses AES to encrypt the message
  *
  *	PARAMETERS		:
  *		byte[] key	: Byte array containing the cipher key
  *		byte[] message	: byte[] containing the message to encrypt
  *
  *	RETURNS			:
  *		byte[] encoded	: byte[] containing the encoded user input
  *		error err	: Error message
  */
 func Encrypt(key []byte, plainText []byte) (encoded []byte, err error) {
		
	 //Create a new AES cipher using the key
	 block, err := aes.NewCipher(key)
 
	 //IF NewCipher failed, exit:
	 if err != nil {
		 return
	 }
 
	 //Make the cipher text a byte array of size BlockSize + the length of the message
	 cipherText := make([]byte, aes.BlockSize+len(plainText))
 
	 //iv is the ciphertext up to the blocksize (16)
	 iv := cipherText[:aes.BlockSize]
	 if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		 return
	 }
 
	 //Encrypt the data:
	 stream := cipher.NewCFBEncrypter(block, iv)
	 stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)
 	

 	 dst := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
 	 base64.StdEncoding.Encode(dst, cipherText)

	 //Return string encoded in base64
	 return dst, err
 }
 
 /*
  *	FUNCTION		: decrypt
  *	DESCRIPTION		:
  *		This function takes a string and a key and uses AES to decrypt the string into plain text
  *
  *	PARAMETERS		:
  *		byte[] key	: Byte array containing the cipher key
  *		byte[] secure	: byte[] containing an encrypted message
  *
  *	RETURNS			:
  *		byte[] decoded	: byte[] containing the decrypted equivalent of secure
  *		error err	: Error message
  */
 func Decrypt(key []byte, secure []byte) (decoded []byte, err error) {
	 //Remove base64 encoding:
	 cipherText := make([]byte, base64.StdEncoding.DecodedLen(len(secure)))
	 _, err = base64.StdEncoding.Decode(cipherText, []byte(secure))
	 
	 //IF DecodeString failed, exit:
	 if err != nil {
		 return
	 }
 
	 //Create a new AES cipher with the key and encrypted message
	 block, err := aes.NewCipher(key)
 
	 //IF NewCipher failed, exit:
	 if err != nil {
		 return
	 }
 
	 //IF the length of the cipherText is less than 16 Bytes:
	 if len(cipherText) < aes.BlockSize {
		 err = errors.New("Ciphertext block size is too short!")
		 return
	 }
 
	 iv := cipherText[:aes.BlockSize]
	 cipherText = cipherText[aes.BlockSize:]
 
	 //Decrypt the message
	 stream := cipher.NewCFBDecrypter(block, iv)
	 stream.XORKeyStream(cipherText, cipherText)
 
	 return cipherText, err
 }