package domain

//go:generate mockgen -destination=../mock/domain/user.go -source=user.go

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	Username string
	Password string
	Role     Role
}

type UserRepository interface {
	Insert(user User) error
	FindAll() ([]User, error)
	FindByUsername(username string) (User, error)
	ExistsByRole(role Role) (bool, error)
	DeleteByUsername(username string) error
}
