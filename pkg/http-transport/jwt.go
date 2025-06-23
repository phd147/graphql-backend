package http_transport

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"os"
)

var keyID = uuid.Nil.String()

const privateKeyFilePath = ".keys/private_key.pem"
const publicKeyFilePath = ".keys/public_key.pem"

//go:embed .keys/private_key.pem
var privateKeyBts []byte

//go:embed .keys/public_key.pem
var publicKeyBts []byte

func LoadRSAKeys() (KeyPair, error) {
	privateKeyBlock, _ := pem.Decode(privateKeyBts)
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return KeyPair{}, errors.New("invalid private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return KeyPair{}, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyBts)
	if publicKeyBlock == nil || publicKeyBlock.Type != "PUBLIC KEY" {
		return KeyPair{}, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return KeyPair{}, err
	}

	return KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil

}

// Helper to generate keys and store them in the global maps
func generateAndStoreKeys(keyID string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate key pair for %s: %w", keyID, err)
	}

	publicKey := &priv.PublicKey

	// Save private key to PEM file
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}

	if err := os.WriteFile(privateKeyFilePath, pem.EncodeToMemory(privateKeyPEM), 0600); err != nil {
		return fmt.Errorf("failed to write private key to file: %w", err)
	}

	// Save public key to PEM file
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := os.WriteFile(publicKeyFilePath, pem.EncodeToMemory(publicKeyPEM), 0644); err != nil {
		return fmt.Errorf("failed to write public key to file: %w", err)
	}
	return nil
}

type KeyPair struct {
	PrivateKey any
	PublicKey  any
}

type JwtHandler interface {
	GenerateToken(ctx context.Context, userClaims UserClaims) (string, error)
	GetPublicKey(ctx context.Context, keyID string) (any, error)
}

type jwtHandler struct {
	Key KeyPair
}

func (j jwtHandler) GetPublicKey(ctx context.Context, keyID string) (any, error) {
	// TODO: Implement this method to return the public key based on the keyID
	return j.Key.PublicKey, nil
}

func (j jwtHandler) GenerateToken(ctx context.Context, userClaims UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, userClaims)
	token.Header["kid"] = keyID

	signedToken, err := token.SignedString(j.Key.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func NewJWTHandler(key KeyPair) JwtHandler {
	return &jwtHandler{
		Key: key,
	}
}
