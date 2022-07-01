package svc

import (
	"context"
	"fmt"

	"github.com/8run0/otel/backend/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCreateUser   = fmt.Errorf("service failed to create user")
	ErrDeleteUser   = fmt.Errorf("service failed to delete user")
	ErrUpdateUser   = fmt.Errorf("service failed to update user")
	ErrNoUserFound  = fmt.Errorf("service failed to get user")
	ErrNoUsersFound = fmt.Errorf("service failed to get users")
)

type userServiceImpl interface {
	GetUsers(ctx context.Context) (users []User, err error)
	GetUserByID(ctx context.Context, id int64) (user User, err error)
	CreateUser(ctx context.Context, req CreateUserRequest) (id int64, err error)
	DeleteUser(ctx context.Context, req DeleteUserRequest) (err error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) (err error)
}

type User struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type CreateUserRequest struct {
	Name     string
	Password string
}

type DeleteUserRequest struct {
	ID int64
}

type UpdateUserRequest struct {
	ID     int64
	Update CreateUserRequest
}

type UserService struct {
	userServiceImpl
}

func NewUserService() UserService {
	userService := userServiceSpanner{
		next: &userService{
			UserDatabase: db.NewUserDatabase(),
		},
	}
	return UserService{
		userServiceImpl: &userService,
	}
}

var _ userServiceImpl = &userService{}

type userService struct {
	*db.UserDatabase
}

func (u userService) GetUsers(ctx context.Context) (users []User, err error) {
	//GetUsers business logic goes here
	dbUsers, err := u.UserDatabase.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	users = make([]User, len(dbUsers))
	for pos, dbUser := range dbUsers {
		users[pos] = User{
			ID:   dbUser.ID,
			Name: dbUser.Name,
		}
	}
	return users, nil
}

func (u userService) GetUserByID(ctx context.Context, id int64) (user User, err error) {
	//GetUserByID business logic goes here
	dbUser, err := u.UserDatabase.GetUserByID(ctx, id)
	if err != nil {
		return User{}, err
	}
	user = User{
		ID:   dbUser.ID,
		Name: dbUser.Name,
	}
	return user, nil
}

func (u userService) CreateUser(ctx context.Context, req CreateUserRequest) (id int64, err error) {
	name := req.Name
	password := req.Password
	passwordHash, err := u.hashPassword(password)
	if err != nil {
		return 0, err
	}
	return u.UserDatabase.CreateUser(ctx, &db.CreateUserRequest{
		Name:         name,
		PasswordHash: passwordHash,
	})
}

func (u userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u userService) DeleteUser(ctx context.Context, req DeleteUserRequest) (err error) {
	//DeleteUser business logic goes here
	return u.UserDatabase.DeleteUser(ctx, &db.DeleteUserRequest{
		ID: req.ID,
	})
}

func (u userService) UpdateUser(ctx context.Context, req UpdateUserRequest) (err error) {
	//UpdateUser business logic goes here
	name := req.Update.Name
	password := req.Update.Password
	passwordHash, err := u.hashPassword(password)
	if err != nil {
		return err
	}
	return u.UserDatabase.UpdateUser(ctx, &db.UpdateUserRequest{
		ID:      req.ID,
		Updates: db.User{Name: name, PasswordHash: passwordHash},
	})
}
