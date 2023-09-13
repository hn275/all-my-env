package database

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/hn275/envhub/server/crypto"
)

var (
	idCounterMap map[uint32]uint16
	m            sync.Mutex

	ErrRepoMissingRepoID = errors.New("repository id not set")
	ErrIDNotGenerated    = errors.New("id not generated")
	ErrValueNotFound     = errors.New("value not found")
)

func init() {
	idCounterMap = make(map[uint32]uint16)
	m = sync.Mutex{}
}

type User struct {
	ID           uint32 `db:"id"`
	Login        string `db:"login"`
	RefreshToken string `db:"refresh_token"`
	Email        string `db:"email"`
}

type Repository struct {
	ID        uint32    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	FullName  string    `db:"full_name" json:"full_name"` // ie: hn275/envhub
	UserID    uint32    `db:"user_id" json:"user_id"`
}

// This table contains all the environment variables.
//
// `Value`'s are never saved raw. always the base64 encoding of the ciphered text,
// and the `ad` is the base64 decoded value of it's ID
type Variable struct {
	ID           string     `db:"id" json:"id"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at"`
	Key          string     `db:"variable_key" json:"key"`
	Value        string     `db:"variable_value" json:"value"`
	RepositoryID uint32     `db:"repository_id" json:"repository_id,omitempty"`
}

func (v *Variable) DecryptValue() error {
	ad, err := base64.StdEncoding.DecodeString(v.ID)
	if err != nil {
		return err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(v.Value)
	if err != nil {
		return err
	}

	plaintext, err := crypto.Decrypt(crypto.VariableKey, ciphertext, ad)
	if err != nil {
		return err
	}

	v.Value = string(plaintext)
	return nil
}

// Generates an ID for the variable. Panics if `RepositoryID` is not set
//
// schema to generate id:
//   - repoID: 8 bytes
//   - time: 4 bytes
//   - counter var: 2 bytes, reset to 0 every second
//   - random var: 2 bytes
func (v *Variable) GenID() error {
	if v.RepositoryID == 0 {
		return ErrRepoMissingRepoID
	}

	var buf [16]byte

	binary.LittleEndian.PutUint32(buf[:4], v.RepositoryID)

	t := time.Now().UTC().Unix()
	binary.BigEndian.PutUint32(buf[4:8], uint32(t))

	counter := idCounterMap[v.RepositoryID]
	m.Lock()
	idCounterMap[v.RepositoryID] = counter + 1
	m.Unlock()
	binary.BigEndian.PutUint16(buf[8:10], counter)

	if _, err := io.ReadFull(rand.Reader, buf[10:]); err != nil {
		return err
	}

	v.ID = base64.StdEncoding.EncodeToString(buf[:])
	return nil
}

// Cipher value, will panic if `Variable.ID` is an empty value
func (v *Variable) EncryptValue() error {
	if v.ID == "" {
		return ErrIDNotGenerated
	}

	if v.Value == "" {
		return ErrValueNotFound
	}

	ad, err := base64.StdEncoding.DecodeString(v.ID)
	if err != nil {
		return err
	}

	ciphertext, err := crypto.Encrypt(crypto.VariableKey, []byte(v.Value), ad)
	if err != nil {
		return err
	}

	v.Value = base64.StdEncoding.EncodeToString(ciphertext)
	return nil
}

func RefreshVariableCounter() {
	for {
		m.Lock()
		idCounterMap = make(map[uint32]uint16)
		m.Unlock()
		time.Sleep(time.Second)
	}
}

// This table describes the type of access an user have for each repo.
// By default all users would have read-only access (that is if github api says so).
// This table only holds records for write access, including the repo owner, which
// the access entry should be written when they link up the repository
type Permission struct {
	ID           uint32 `db:"id"`
	RepositoryID uint32 `db:"repository_id"`
	UserID       uint32 `db:"user_id"`
}
