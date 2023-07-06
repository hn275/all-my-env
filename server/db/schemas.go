package db

var (
	TableUsers = "users"

	VendorGithub = "github"
)

type User struct {
	ID        int    `db:"id"`
	CreatedAt string `db:"created_at"`
	Vendor    string `db:"vendor"`
	UserName  string `db:"username"`
}
