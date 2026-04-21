package users_repository_postgres

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}
