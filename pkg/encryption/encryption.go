package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/scrypt"
)

func AESDecrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create gcm")
	}

	nonce := make([]byte, gcm.NonceSize())

	payload, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}

	return payload, nil
}

func AESEncrypt(key, payload []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("failed to create cipher")
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("failed to create gcm")
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, errors.New("failed to create nonce")
	}

	ciphertext := aesgcm.Seal(nil, nonce, payload, nil)
	return ciphertext, nil
}

func GeneratePrivateKey() (*ecdh.PrivateKey, error) {
	priv, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		return nil, errors.New("failed to generate key pair")
	}

	return priv, nil
}

func GenerateSharedKey(priv *ecdh.PrivateKey, remote *ecdh.PublicKey) ([]byte, error) {
	shared, err := priv.ECDH(remote)
	if err != nil {
		return nil, errors.New("invalid remote public key")
	}

	hash := sha256.New
	hkdf := hkdf.New(hash, shared, nil, []byte("aes-key"))
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(hkdf, key); err != nil {
		return nil, errors.New("failed to generate AES key")
	}

	return key, nil
}

func PasswordDecrypt(password, encrypted []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(encrypted))
	if err != nil {
		return nil, errors.New("failed to decode encrypted payload")
	}

	salt := decoded[:16]
	encrypted = decoded[16:]

	// generate AES key from password
	key, err := scrypt.Key(password, salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, errors.New("failed to generate AES key")
	}

	// decrypt data
	decrypted, err := AESDecrypt(key, encrypted)
	if err != nil {
		return nil, errors.New("failed to decrypt payload")
	}

	return decrypted, nil
}

func PasswordEncrypt(password, payload []byte) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errors.New("failed to create salt")
	}

	// generate AES key from password
	key, err := scrypt.Key(password, salt, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, errors.New("failed to generate AES key")
	}

	ciphertext, err := AESEncrypt(key, payload)
	if err != nil {
		return nil, errors.New("failed to encrypt data")
	}

	// store salt inside
	encrypted := append(salt, ciphertext...)

	// encode to base64
	var encoded []byte
	base64.StdEncoding.Encode(encoded, encrypted)

	return encoded, nil
}
