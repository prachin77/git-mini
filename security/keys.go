package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/prachin77/pkr/models"
)

var (
	RSAKeySize = 4096
)

// GenerateRSAKeys generates an RSA key pair (private & public keys) and returns them.
func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		return nil, nil, fmt.Errorf("error generating RSA keys: %v", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

func ParsePublicKeyToBytes(public_key *rsa.PublicKey) []byte {
	public_key_bytes := x509.MarshalPKCS1PublicKey(public_key)
	publickey_pem_block := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: public_key_bytes,
		},
	)

	return publickey_pem_block
}

func ParsePrivateKeyToBytes(private_key *rsa.PrivateKey) []byte {
	private_key_bytes := x509.MarshalPKCS1PrivateKey(private_key)
	privatekey_pem_block := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: private_key_bytes,
		},
	)

	return privatekey_pem_block
}

func StorePrivateKeys(privateKey *rsa.PrivateKey, privateKeyFilePath string) error {
	private_pem_key := ParsePrivateKeyToBytes(privateKey)
	if private_pem_key == nil {
		return errors.New("error converting private key to []byte")
	}

	fmt.Println("private key successfully stored into config folder ...")
	return os.WriteFile(privateKeyFilePath, private_pem_key, 0777)
}

func StorePublicKeys(publicKey *rsa.PublicKey, publicKeyFilePath string) error {
	public_pem_block := ParsePublicKeyToBytes(publicKey)
	if public_pem_block == nil {
		return errors.New("error converting public key to []byte")
	}

	fmt.Println("public key successfully stored into config folder ...")
	return os.WriteFile(publicKeyFilePath, public_pem_block, 0777)
}

func GetPublicKey() (string, string, error) {
	public_key_data, err := os.ReadFile(models.PUBLIC_KEY_FILE)
	if err != nil {
		return "", "", err
	}

	path, err := filepath.Abs(models.PUBLIC_KEY_FILE)
	if err != nil {
		fmt.Println("error retrieving host public key file path : ", err)
		return "", "", err
	}
	return string(public_key_data), path, nil
}

func GetPrivateKeys() (string, string, error) {
	private_key_data, err := os.ReadFile(models.PRIVATE_KEY_FILE)
	if err != nil {
		return "", "", err
	}

	path, err := filepath.Abs(models.PRIVATE_KEY_FILE)
	if err != nil {
		fmt.Println("error retrieving host public key file path : ", err)
		return "", "", err
	}

	return string(private_key_data), path, nil
}

// ParseBytesToPublicKey converts a PEM-encoded public key (as []byte) into an *rsa.PublicKey
func ParseBytesToPublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA public key: %v", err)
	}

	return publicKey, nil
}

func GenerateAESKeys() ([]byte, error) {
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}
	return aesKey, nil
}

func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}
