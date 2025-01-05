package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

const (
	RSA_key_size = 2048
	AES_key_size = 32
)

func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSA_key_size)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating RSA keys: %v", err)
	}
	fmt.Println("private & public keys generated successfully ...")
	return privateKey, &privateKey.PublicKey, nil
}

// EncryptWithAES encrypts data using AES-GCM.
func EncryptWithAES(aesKey []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating AES cipher: %v", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating GCM: %v", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating nonce: %v", err)
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// StorePrivateKeys encrypts and stores the private key securely using AES.
func StorePrivateKeys(privateKey *rsa.PrivateKey, privateKeyFilePath string) error {
	// Generate a random AES key
	aesKey := make([]byte, AES_key_size)
	_, err := rand.Read(aesKey)
	if err != nil {
		return fmt.Errorf("error generating AES key: %v", err)
	}

	// Marshal private key to PEM format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Encrypt the private key using AES
	encryptedKey, nonce, err := EncryptWithAES(aesKey, privateKeyBytes)
	if err != nil {
		return fmt.Errorf("error encrypting private key: %v", err)
	}

	// Append nonce to encrypted key
	data := append(nonce, encryptedKey...)

	// Save encrypted private key to file
	err = ioutil.WriteFile(privateKeyFilePath, data, 0600)
	if err != nil {
		return fmt.Errorf("error saving private key file: %v", err)
	}

	fmt.Println("private key successfully saved inside file ...")
	return nil
}

// StorePublicKeys stores the public key in PEM format.
func StorePublicKeys(publicKey *rsa.PublicKey, publicKeyFilePath string) error {
	// Marshal public key to PEM format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("error marshalling public key: %v", err)
	}

	// Encode the public key in PEM
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Save the public key to a file
	err = ioutil.WriteFile(publicKeyFilePath, publicKeyPEM, 0644)
	if err != nil {
		return fmt.Errorf("error saving public key file: %v", err)
	}

	fmt.Println("public key successfully saved inside file ...")
	return nil
}
