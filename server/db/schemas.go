package db

import "time"

var (
	TableUsers = "users"

	VendorGithub = "github"
)

type User struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt string
	Vendor    string `gorm:"not null"`
	UserName  string `gorm:"not null,unique"`

	// relation
	Repositories []Repository `gorm:"constraint:OnDelete:CASCADE"`
}

type Repository struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	FullName  string `gorm:"not null"`
	Url       string `gorm:"not null"`

	// relation
	User   User
	UserID int `gorm:"foreignKey"`

	Variables []Variable `gorm:"constraint:OnDelete:CASCADE"`
}

type Variable struct {
	ID        int       `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Key       string    `gorm:"not null"`
	Value     string    `gorm:"not null"`
	Nonce     string    `gorm:"not null"`

	// relation
	Repository   Repository
	RepositoryID int `gorm:"foreignKey"`
}
