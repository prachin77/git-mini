package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/prachin77/pkr/models"
)

func EncryptWorkspacePassword(sending_workspace *models.SendWorkSpaceFolder , host_publickey string) (string , error) {
	// Decode the PEM-encoded public key
	block , _ := pem.Decode([]byte(host_publickey))
	if block == nil {
		return "" , errors.New("invalid public key format")
	}	

	// parse public key 	
	publickey , err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	
	// confirm that key is RSA public key 
	rsaPublicKey , ok := publickey.(*rsa.PublicKey)
	if !ok{
		return "" , errors.New("sorry , not an RSA public key")
	}
	
	// encrypt password using RSA-OAEP
	label := []byte("")
	hash := sha256.New()
	encrypted_bytes , err := rsa.EncryptOAEP(hash , rand.Reader , rsaPublicKey , []byte(sending_workspace.Workspace_Password),label)
	if err != nil{
		return "" , err
	}

	return base64.StdEncoding.EncodeToString(encrypted_bytes) , nil
}