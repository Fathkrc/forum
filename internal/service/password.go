package service

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      = 64 * 1024   // 64 MB of RAM used while hashing
	argonIterations  = 3           // does the work 3 times (slowness)
	argonParallelism = 2           // 2 CPU threads
	argonSaltLength  = 16          // 16 random bytes of salt
	argonKeyLength   = 32          // output hash length: 32 bytes
)

// HashPassword hashes the given password using Argon2id and returns a
// self-contained string in the PHC format:
//
//	$argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonKeyLength,
	)

	encoded := base64.RawStdEncoding.EncodeToString(hash)
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argonMemory,
		argonIterations,
		argonParallelism,
		saltEncoded,
		encoded,
	), nil
}

// VerifyPassword checks whether the given password matches an encoded
// Argon2id hash produced by HashPassword.
func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return false, errors.New("unsupported hash algorithm")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	candidate := argon2.IDKey(
		[]byte(password),
		salt,
		argonIterations,
		argonMemory,
		argonParallelism,
		argonKeyLength,
	)

	return subtle.ConstantTimeCompare(hash, candidate) == 1, nil
}
