package db

import (
	"encoding/base64"
	"time"

	"github.com/hn275/envhub/server/crypto"
)

var (
	TableUsers       = "users"
	TableRepos       = "repositories"
	TablePermissions = "permissions"

	VendorGithub = "github"
)

// UTC time stamp RFC3339
type TimeStamp = time.Time
type Base64EncodedID = string

type User struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt TimeStamp
	Vendor    string `gorm:"not null"`
	UserName  string `gorm:"not null,unique"`

	// relation
	Repositories []Repository `gorm:"constraint:OnDelete:CASCADE"`
	Permission   []Permission `gorm:"constraint:OnDelete:CASCADE"`
}

type Repository struct {
	ID        uint32    `gorm:"primaryKey" json:"id"`
	CreatedAt TimeStamp `json:"created_at"`

	// ie: hn275/envhub
	FullName string `gorm:"not null" json:"full_name"`
	// ie: https://github.com/hn275/envhub
	Url string `gorm:"not null" json:"url"`

	// relation
	User   User `json:"-"`
	UserID int  `gorm:"foreignKey" json:"-"`

	Variables  []Variable   `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Permission []Permission `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

// This table contains all the environment variables.
//
// `Value`'s are never saved raw. always the base64 encoding of the ciphered text,
// and the `ad` is the base64 decoded value of it's ID
type Variable struct {
	ID        Base64EncodedID `gorm:"primaryKey" json:"id"`
	CreatedAt TimeStamp       `gorm:"not null" json:"created_at"`
	UpdatedAt TimeStamp       `gorm:"not null" json:"updated_at"`
	Key       string          `gorm:"not null;uniqueIndex:unique_key_repo" json:"key"`
	Value     string          `gorm:"not null" json:"value"`

	// relation
	Repository   Repository `json:"-"`
	RepositoryID uint32     `gorm:"foreignKey;uniqueIndex:unique_key_repo" json:"-"`
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
	RepositoryID uint `gorm:"foreignKey;uniqueIndex:unique_user_repo"`

	User   User
	UserID uint `gorm:"foreignKey;uniqueIndex:unique_user_repo"`
}
