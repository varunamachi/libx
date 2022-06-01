package iox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/varunamachi/libx/errx"
	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrKey    = errors.New("failed generate key")
	ErrCipher = errors.New("failed to create cipher")
	ErrFile   = errors.New("failed read/write file")
	ErrInput  = errors.New("invalid input")
)

const saltSize = 32
const magicSize = 4

var magic = []byte{0xE1, 0xEA, 0xE1, 0xA0}

type FileCrytor interface {
	Encrypt(in []byte) ([]byte, error)
	Decrypt(in []byte) ([]byte, error)
	IsEncrypted(in []byte) bool
}

type aesGCMCryptor struct {
	password string
}

func NewCryptor(password string) FileCrytor {
	return &aesGCMCryptor{
		password: password,
	}
}

func EncryptToBase64Str(in, password string) (string, error) {
	c := NewCryptor(password)
	out, err := c.Encrypt([]byte(in))
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(out), nil
}

func DecryptBase64Str(in, password string) (string, error) {
	ba, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return "", err
	}

	c := NewCryptor(password)
	out, err := c.Decrypt(ba)
	if err != nil {
		return "", err
	}
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func EncryptToFile(reader io.Reader, path, password string) error {

	in, err := ioutil.ReadAll(reader)
	if err != nil {
		return errx.Errf(err, "failed to read plaintext to file at %s", path)
	}

	c := NewCryptor(password)
	out, err := c.Encrypt(in)
	if err != nil {
		return err
	}

	if err = os.WriteFile(path, out, 0700); err != nil {
		return errx.Errf(err, "failed write encrypted data to file")
	}
	return nil
}

func DecryptFromFile(path, password string, writer io.Writer) error {
	in, err := os.ReadFile(path)
	if err != nil {
		return errx.Errf(
			err, "failed to read ciphertext from file at '%s'", path)
	}

	c := NewCryptor(password)
	out, err := c.Decrypt(in)
	if err != nil {
		return err
	}

	if _, err = writer.Write(out); err != nil {
		return errx.Errf(err, "failed write encrypted data to file")
	}
	return nil
}

func (c *aesGCMCryptor) Encrypt(in []byte) ([]byte, error) {
	if c.IsEncrypted(in) {
		return nil, errx.Errf(ErrInput, "the data is already encrypted")
	}

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, errx.Errf(err, "failed to create salt")
	}

	gcm, err := c.getGCM(salt)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errx.Errf(err, "failed to create nonce")
	}

	val := gcm.Seal(nonce, nonce, in, nil)
	out := make([]byte, 0, len(val)+saltSize+magicSize)
	out = append(out, magic...)
	out = append(out, salt...)
	out = append(out, val...)

	return out, nil
}

func (c *aesGCMCryptor) Decrypt(in []byte) ([]byte, error) {
	if !c.IsEncrypted(in) {
		return nil, errx.Errf(ErrInput, "the input is not properly encrypted")
	}

	in = in[magicSize:]
	salt := in[:saltSize]
	gcm, err := c.getGCM(salt)
	if err != nil {
		return nil, err
	}
	in = in[saltSize:]

	nonceSize := gcm.NonceSize()
	if len(in) < nonceSize {
		return nil, errx.Errf(err, "input data is too small to decrypt")
	}

	nonce, cipherText := in[:nonceSize], in[nonceSize:]
	out, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errx.Errf(err, "failed to decrypt")
	}

	return out, nil
}

func (c *aesGCMCryptor) getGCM(salt []byte) (cipher.AEAD, error) {

	key := pbkdf2.Key([]byte(c.password), salt, 65536, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errx.Errf(err, "failed to create AES cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errx.Errf(err, "failed to create GCM cipher")
	}

	return gcm, nil
}

func (c *aesGCMCryptor) IsEncrypted(in []byte) bool {
	if len(in) < magicSize+saltSize {
		return false
	}
	for i := 0; i < magicSize; i++ {
		if in[i] != magic[i] {
			return false
		}
	}
	return true
}
