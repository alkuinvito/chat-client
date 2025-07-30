package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
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

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("invalid nonce")
	}

	// retrieve nonce from ciphertext
	nonce := ciphertext[:nonceSize]
	payload := ciphertext[nonceSize:]

	decrypted, err := gcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, errors.New("failed to decrypt data")
	}

	return decrypted, nil
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

	// store nonce inside ciphertext
	ciphertext = append(nonce, ciphertext...)

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
	salt := encrypted[:16]
	encrypted = encrypted[16:]

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
	return encrypted, nil
}
