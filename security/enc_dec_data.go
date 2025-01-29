package security

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/prachin77/pkr/models"
)

func EncryptData(sending_workspace *models.SendWorkSpaceFolder, host_publickey string) (string, error) {
	host_publickey_pemBlock := strings.TrimSpace(host_publickey)
	block, _ := pem.Decode([]byte(host_publickey_pemBlock))
	if block == nil {
		fmt.Println("error parsing host public key to pem bytes !")
		fmt.Println("check if it exists or path is correct  !")
		return "", errors.New("error parsing host public key to pem bytes")
	}

	host_public_key, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("error in parsing host public key to bytes")
		return "", err
	}

	label := []byte("")
	hash := sha256.New()

	result, err := rsa.EncryptOAEP(hash, rand.Reader, host_public_key, []byte(sending_workspace.Workspace_Password), label)
	if err != nil {
		fmt.Println("error encryting workspace password !")
		return "", err
	}

	base64_encrypted_password := base64.StdEncoding.EncodeToString(result)
	return base64_encrypted_password, nil
}

func DecryptData(encrypted_workspace_password string) (string, error) {
	host_private_key, host_privateKey_filepath, err := GetPrivateKeys()
	if err != nil || host_privateKey_filepath == "" || host_private_key == "" {
		fmt.Println("error retrieving host private key : ", err)
		return "", nil
	}

	block, _ := pem.Decode([]byte(host_private_key))
	if block == nil {
		fmt.Println("error parsing host private key to pem bytes !")
		fmt.Println("check if it exists or path is correct  !")
		return "", errors.New("error parsing host private key to pem bytes")
	}

	host_privatekey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("error in parsing host private key to bytes")
		return "", err
	}

	label := []byte("")
	hash := sha256.New()

	baseDecoded, _ := base64.StdEncoding.DecodeString(encrypted_workspace_password)
	plainText, err := rsa.DecryptOAEP(hash, rand.Reader, host_privatekey, []byte(baseDecoded), label)
	if err != nil {
		fmt.Println("error decrypting encrypted workspace password !")
		return "", err
	}

	return string(plainText), nil
}

func AESEncryptZipFile(zip_filePath string, encrypted_zip_FilePath string, AES_Key []byte, nonce []byte) error {
	inputFile, err := os.Open(zip_filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(encrypted_zip_FilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(AES_Key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, nonce)

	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

func EncryptZipFile(AES_Key string, client_publicKey string) (string, error) {
	client_publicKey_pemBlock := strings.TrimSpace(client_publicKey)
	block, _ := pem.Decode([]byte(client_publicKey_pemBlock))
	if block == nil {
		fmt.Println("error in parsing pem block !")
		return "", errors.New("error in parsing pem block")
	}

	publicKey_data, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("error in parsing public key !")
		return "", err
	}

	label := []byte("")
	hash := sha256.New()

	result, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey_data, []byte(AES_Key), label)
	if err != nil {
		fmt.Println("error encrypting file ")
		return "", err
	}

	base64Encrypted := base64.StdEncoding.EncodeToString(result)
	return base64Encrypted, nil
}

func AESDecryptZipFile(data_bytes []byte , AES_Key string , nonce string) ([]byte , error) {
	// create AES cipher text (a random text everytime)
	block , err := aes.NewCipher([]byte(AES_Key))
	if err != nil{
		return nil , err
	}

	stream := cipher.NewCTR(block , []byte(nonce))

	// Create a buffer to hold the decrypted data
	var decryptedData bytes.Buffer
	writer := &cipher.StreamWriter{S:stream , W:&decryptedData}

	if _ , err := writer.Write(data_bytes); err != nil{
		return nil ,err
	}

	return decryptedData.Bytes() , nil
}