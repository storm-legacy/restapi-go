package model

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Custom errors
var (
	errInvalidHash         = errors.New("Encoded hash is not correct")
	errIncompatibleVersion = errors.New("Incompatible version of Argon2")
)

// Structure for argon2id configuration
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var a2params = &argon2Params{
	memory:      65536,
	iterations:  3,
	parallelism: 2,
	saltLength:  16,
	keyLength:   32,
}

// Generate Argon2id standard encoded string from password
func encodePassword(password *string) (encodedHash string, err error) {

	// Generate random bytes for salt
	salt := make([]byte, a2params.saltLength)
	_, err = rand.Read(salt)
	if err != nil {
		return "", nil
	}

	// Encode password
	hash := argon2.IDKey(
		[]byte(*password),
		salt,
		a2params.iterations,
		a2params.memory,
		a2params.parallelism,
		a2params.keyLength,
	)

	// Create standardized representation for hashed password
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		a2params.memory,
		a2params.iterations,
		a2params.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// Check password with hashed one
func decodeHash(encodedHash *string) (p *argon2Params, salt, hash []byte, err error) {
	// Split string and check if is correct argon2 standard
	values := strings.Split(*encodedHash, "$")
	if len(values) != 6 {
		return nil, nil, nil, errInvalidHash
	}

	// Extract version from values[2]
	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err

	} else if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	// Extract params for argon2 hashing
	p = &argon2Params{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	// Extract and decode salt for password
	salt, err = base64.RawStdEncoding.Strict().DecodeString(values[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	// Extract and decode hash for password
	hash, err = base64.RawStdEncoding.Strict().DecodeString(values[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// Compare password with hashed one
func comparePasswordAndHash(password *string, encodedHash *string) (match bool, err error) {
	// Extract parameters from encoded password hash
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Hash plaintext password
	hash2 := argon2.IDKey([]byte(*password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Prevent timing attacks
	if subtle.ConstantTimeCompare(hash, hash2) == 1 {
		return true, nil
	}

	return false, nil
}
