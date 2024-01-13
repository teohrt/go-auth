package repositories

type DBUser struct {
	UserID     string
	Username   string
	AuthID     string
	Email      string
	Created_at string
	Updated_at string
}

type CreateUserInput struct {
	AuthID   string
	Username string
}
