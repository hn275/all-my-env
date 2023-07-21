package db

import "time"

var (
	TableUsers = "users"
	TableRepos = "repositories"

	VendorGithub = "github"
)

// UTC time stamp RFC3339
type TimeStamp = time.Time

type User struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt TimeStamp
	Vendor    string `gorm:"not null"`
	UserName  string `gorm:"not null,unique"`

	// relation
	Repositories []Repository `gorm:"constraint:OnDelete:CASCADE"`
}

type Repository struct {
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt TimeStamp
	FullName  string `gorm:"not null"`
	Url       string `gorm:"not null"`

	// relation
	User   User
	UserID int `gorm:"foreignKey"`

	Variables []Variable `gorm:"constraint:OnDelete:CASCADE"`
}

type Variable struct {
	ID        uint64    `gorm:"primaryKey"`
	CreatedAt TimeStamp `gorm:"not null"`
	UpdatedAt TimeStamp `gorm:"not null"`
	Key       string    `gorm:"not null"`
	Value     string    `gorm:"not null"`
	Nonce     string    `gorm:"not null"`

	// relation
	Repository   Repository
	RepositoryID int `gorm:"foreignKey"`
}
