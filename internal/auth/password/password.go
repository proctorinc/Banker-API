package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Parameters struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type DecodedHash struct {
	Parameters Argon2Parameters
	Salt       []byte
	Hash       []byte
}

func HashPassword(password string) (string, error) {
	params := getParams()
	salt, err := generateSalt(params.SaltLength)

	if err != nil {
		return "", fmt.Errorf("Failed to generate salt")
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format to argon2 hash string format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func VerifyPassword(password string, encodedHash string) error {
	decoded, err := decodeHash(encodedHash)
	params := decoded.Parameters

	if err != nil {
		return err
	}

	otherHash := argon2.IDKey(
		[]byte(password),
		decoded.Salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	// Compare hashes with constant time to avoid timing attacks
	if subtle.ConstantTimeCompare(decoded.Hash, otherHash) == 1 {
		return nil
	}

	return fmt.Errorf("Passwords do not match")
}

func getParams() Argon2Parameters {
	// Argon2 parameters per OWASP recommendation
	// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	return Argon2Parameters{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func decodeHash(encodedHash string) (data *DecodedHash, err error) {
	var version int
	params := Argon2Parameters{}

	vals := strings.Split(encodedHash, "$")

	if len(vals) != 6 {
		return nil, fmt.Errorf("invalid hash")
	}

	_, err = fmt.Sscanf(vals[2], "v=%d", &version)

	if err != nil {
		return nil, err
	}

	if version != argon2.Version {
		return nil, fmt.Errorf("invalid argon2 version")
	}

	_, err = fmt.Sscanf(
		vals[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	)

	if err != nil {
		return nil, err
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])

	if err != nil {
		return nil, err
	}

	params.SaltLength = uint32(len(salt))

	hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])

	if err != nil {
		return nil, err
	}

	params.KeyLength = uint32(len(hash))

	decoded := &DecodedHash{
		Parameters: params,
		Salt:       salt,
		Hash:       hash,
	}

	return decoded, nil
}

func generateSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}
