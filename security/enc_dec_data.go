package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/prachin77/pkr/models"
	"github.com/prachin77/pkr/utils"
)

func EncryptWorkspacePassword(sending_workspace *models.SendWorkSpaceFolder, host_publickey string) (string, error) {
	host_publickey = strings.TrimSpace(host_publickey)
	block, _ := pem.Decode([]byte(host_publickey))
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

func DecryptWorkspacePassword(encrypted_workspace_password string) (string, error) {
	host_private_key, host_privateKey_filepath, err := utils.GetHostPrivateKeys()
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
