package domain

type UserRepository interface {
	Upsert(u *User) error
}
