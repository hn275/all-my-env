package database

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/hn275/envhub/server/crypto"
)

var (
	TableUsers       = "users"
	TableRepos       = "repositories"
	TablePermissions = "permissions"
	TableVariables   = "variables"

	VendorGithub = "github"

	idCounterMap map[uint64]uint16
	m            sync.Mutex

	ErrRepoMissingRepoID = errors.New("repository id not set")
)

func init() {
	idCounterMap = make(map[uint64]uint16)
	m = sync.Mutex{}
}

// UTC time stamp RFC3339
type TimeStamp = time.Time
type HexEncodedID = string

type User struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt TimeStamp
	Vendor    string `gorm:"not null"`
	UserName  string `gorm:"not null,unique"`

	// relation
	Repositories []Repository `gorm:"constraint:OnDelete:CASCADE"`
	Permission   []Permission `gorm:"constraint:OnDelete:CASCADE"`
}

type Repository struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt TimeStamp `json:"created_at"`

	// ie: hn275/envhub
	FullName string `gorm:"not null" json:"full_name"`
	// ie: https://github.com/hn275/envhub
	Url string `gorm:"not null" json:"url"`

	// relation
	User   User   `json:"-"`
	UserID uint64 `gorm:"foreignKey" json:"-"`

	Variables  []Variable   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Permission []Permission `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

// This table contains all the environment variables.
//
// `Value`'s are never saved raw. always the base64 encoding of the ciphered text,
// and the `ad` is the base64 decoded value of it's ID
type Variable struct {
	ID        HexEncodedID `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt TimeStamp    `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt TimeStamp    `gorm:"not null" json:"updated_at,omitempty"`
	Key       string       `gorm:"not null;uniqueIndex:unique_key_repo" json:"key,omitempty"`
	Value     string       `gorm:"not null" json:"value,omitempty"`

	// relation
	Repository   Repository `json:"-"`
	RepositoryID uint64     `gorm:"foreignKey;uniqueIndex:unique_key_repo" json:"-"`
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

	plaintext, err := crypto.Decrypt(ciphertext, ad)
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

	binary.LittleEndian.PutUint64(buf[:8], v.RepositoryID)

	t := time.Now().UTC().Unix()
	binary.BigEndian.PutUint32(buf[8:12], uint32(t))

	counter := idCounterMap[v.RepositoryID]
	m.Lock()
	idCounterMap[v.RepositoryID] = counter + 1
	m.Unlock()
	binary.BigEndian.PutUint16(buf[12:14], counter)

	if _, err := io.ReadFull(rand.Reader, buf[14:]); err != nil {
		return err
	}

	v.ID = hex.EncodeToString(buf[:])
	return nil
}

// Cipher value, will panic if `Variable.ID` is an empty value
func (v *Variable) EncryptValue() error {
	if v.ID == "" {
		panic("variable id not generated")
	}

	ad, err := base64.StdEncoding.DecodeString(v.ID)
	if err != nil {
		return err
	}

	ciphertext, err := crypto.Encrypt(v.Value, ad)
	if err != nil {
		return err
	}

	v.Value = base64.StdEncoding.EncodeToString(ciphertext)
	return nil
}

func RefreshVariableCounter() {
	for {
		m.Lock()
		idCounterMap = make(map[uint64]uint16)
		m.Unlock()
		time.Sleep(time.Second)
	}
}

// This table describes the type of access an user have for each repo.
//
// By default all users would have read-only access (that is if github api says so).
//
// This table only holds records for write access.
// If gorm query returns a `gorm.ErrRecordNotFound`, user doesn't have read access.
type Permission struct {
	ID uint `gorm:"primaryKey"`

	// relation
	Repository   Repository
	RepositoryID uint64 `gorm:"foreignKey;uniqueIndex:unique_user_repo"`

	User   User
	UserID uint64 `gorm:"foreignKey;uniqueIndex:unique_user_repo"`
}
