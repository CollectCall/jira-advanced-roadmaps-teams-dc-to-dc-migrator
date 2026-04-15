package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type SecretStore map[string]SavedSecrets

type SavedSecrets struct {
	SourceAuthToken string `json:"sourceAuthToken,omitempty"`
	SourceUsername  string `json:"sourceUsername,omitempty"`
	SourcePassword  string `json:"sourcePassword,omitempty"`
	TargetAuthToken string `json:"targetAuthToken,omitempty"`
	TargetUsername  string `json:"targetUsername,omitempty"`
	TargetPassword  string `json:"targetPassword,omitempty"`
}

type encryptedSecretStore struct {
	Salt       string `json:"salt"`
	Nonce      string `json:"nonce"`
	Ciphertext string `json:"ciphertext"`
}

func defaultSecretStorePath(configPath string) string {
	dir := filepath.Dir(configPath)
	return filepath.Join(dir, "secrets.enc.json")
}

func secretStoreHasProfile(path, profile string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var wrapper encryptedSecretStore
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return false
	}
	return wrapper.Ciphertext != ""
}

func loadSecretStore(path, passphrase string) (SecretStore, error) {
	store := SecretStore{}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}
	var wrapper encryptedSecretStore
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, err
	}
	plaintext, err := decryptSecretPayload(wrapper, passphrase)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(plaintext, &store); err != nil {
		return nil, err
	}
	return store, nil
}

func saveSecretStore(path, passphrase string, store SecretStore) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	plaintext, err := json.Marshal(store)
	if err != nil {
		return err
	}
	wrapper, err := encryptSecretPayload(plaintext, passphrase)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(wrapper, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func encryptSecretPayload(plaintext []byte, passphrase string) (encryptedSecretStore, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return encryptedSecretStore{}, err
	}
	key := deriveSecretKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return encryptedSecretStore{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return encryptedSecretStore{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return encryptedSecretStore{}, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return encryptedSecretStore{
		Salt:       base64.StdEncoding.EncodeToString(salt),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

func decryptSecretPayload(wrapper encryptedSecretStore, passphrase string) ([]byte, error) {
	salt, err := base64.StdEncoding.DecodeString(wrapper.Salt)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(wrapper.Nonce)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(wrapper.Ciphertext)
	if err != nil {
		return nil, err
	}
	key := deriveSecretKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("could not unlock secret store: %w", err)
	}
	return plaintext, nil
}

func deriveSecretKey(passphrase string, salt []byte) []byte {
	sum := sha256.Sum256(append([]byte(passphrase), salt...))
	key := sum[:]
	for i := 0; i < 100000; i++ {
		next := sha256.Sum256(append(key, salt...))
		key = next[:]
	}
	out := make([]byte, len(key))
	copy(out, key)
	return out
}

func loadSecretsIntoConfig(cfg *Config) error {
	if !cfg.HasSavedSecrets || cfg.Profile == "" {
		return nil
	}
	if cfg.SourceAuthToken != "" || cfg.SourceUsername != "" || cfg.SourcePassword != "" || cfg.TargetAuthToken != "" || cfg.TargetUsername != "" || cfg.TargetPassword != "" {
		return nil
	}

	passphrase, err := promptSecretValue("Secret store passphrase", "Enter the passphrase to unlock saved credentials for this profile.")
	if err != nil {
		return err
	}
	store, err := loadSecretStore(cfg.SecretStorePath, passphrase)
	if err != nil {
		return err
	}
	secrets, ok := store[cfg.Profile]
	if !ok {
		return nil
	}
	cfg.SourceAuthToken = secrets.SourceAuthToken
	cfg.SourceUsername = secrets.SourceUsername
	cfg.SourcePassword = secrets.SourcePassword
	cfg.TargetAuthToken = secrets.TargetAuthToken
	cfg.TargetUsername = secrets.TargetUsername
	cfg.TargetPassword = secrets.TargetPassword
	return nil
}
