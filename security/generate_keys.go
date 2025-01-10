package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
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
