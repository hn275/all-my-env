package db

import "time"

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
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt TimeStamp

	// ie: hn275/envhub
	FullName string `gorm:"not null"`
	// ie: https://github.com/hn275/envhub
	Url string `gorm:"not null"`

	// relation
	User   User
	UserID int `gorm:"foreignKey"`

	Variables  []Variable   `gorm:"constraint:OnDelete:CASCADE"`
	Permission []Permission `gorm:"constraint:OnDelete:CASCADE"`
}

type Variable struct {
	ID        Base64EncodedID `gorm:"primaryKey"`
	CreatedAt TimeStamp       `gorm:"not null"`
	UpdatedAt TimeStamp       `gorm:"not null"`
	Key       string          `gorm:"not null;uniqueIndex:unique_key_repo"`
	Value     string          `gorm:"not null"`

	// relation
	Repository   Repository
	RepositoryID uint32 `gorm:"foreignKey;uniqueIndex:unique_key_repo"`
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
