package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/8run0/otel/backend/pkg/otel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrCreateUser   = fmt.Errorf("database failed to create user")
	ErrDeleteUser   = fmt.Errorf("database failed to delete user")
	ErrUpdateUser   = fmt.Errorf("database failed to update user")
	ErrNoUserFound  = fmt.Errorf("database failed to get user")
	ErrNoUsersFound = fmt.Errorf("database failed to get users")
)

type userDatabaseImpl interface {
	GetUsers(ctx context.Context) (users []User, err error)
	GetUserByID(ctx context.Context, id int64) (user User, err error)
	CreateUser(ctx context.Context, req *CreateUserRequest) (id int64, err error)
	DeleteUser(ctx context.Context, req *DeleteUserRequest) (err error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest) (err error)
}

type User struct {
	ID           uint64    `db:"id"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	CreatedOn    time.Time `db:"created_on"`
	LastEdited   time.Time `db:"last_edited"`
}

type CreateUserRequest struct {
	Name         string
	PasswordHash string
}

type DeleteUserRequest struct {
	ID int64
}

type UpdateUserRequest struct {
	ID      int64
	Updates User
}

type UserDatabase struct {
	userDatabaseImpl
	*otel.Tools
}

func NewUserDatabase() *UserDatabase {
	db, _ := sqlx.Connect("sqlite3", "./users.db")
	userDB := &userDatabase{
		DB: db,
	}
	userDB.initalise()
	return &UserDatabase{
		userDatabaseImpl: userDatabaseSpanner{
			next: userDB,
		},
	}
}
func (db *userDatabase) initalise() {
	createTable :=
		`CREATE TABLE IF NOT EXISTS 'users' ('id' INTEGER PRIMARY KEY AUTOINCREMENT, 'name' VARCHAR(64), 'password_hash' VARCHAR(255),'created_on' TIMESTAMP , 'last_edited' TIMESTAMP);`
	_, err := db.Exec(createTable)
	if err != nil {
		fmt.Println("user table initalised failed", err)
		return
	}
	fmt.Println("user table initalised")
}

var _ userDatabaseImpl = userDatabase{}

type userDatabase struct {
	*sqlx.DB
}

func (db userDatabase) GetUsers(ctx context.Context) (users []User, err error) {
	//GetUsers business logic goes here
	users, err = db.getUsers(ctx)
	if err != nil {
		return nil, ErrNoUsersFound
	}
	return users, nil
}
func (db userDatabase) getUsers(ctx context.Context) (users []User, err error) {
	users, err = db.getUsersQuery(ctx)
	if err != nil {
		return nil, err
	}
	return users, err
}
func (db userDatabase) getUsersQuery(ctx context.Context) (users []User, err error) {
	users = []User{}
	err = db.SelectContext(ctx, &users, "SELECT * FROM users;")
	return users, err
}

func (db userDatabase) GetUserByID(ctx context.Context, id int64) (user User, err error) {
	//GetUserByID business logic goes here
	user, err = db.getUserByID(ctx, id)
	if err != nil {
		return User{}, ErrNoUserFound
	}
	return user, nil
}
func (db userDatabase) getUserByID(ctx context.Context, id int64) (user User, err error) {
	user, err = db.getUserByIDQuery(ctx, id)
	return user, err
}
func (db userDatabase) getUserByIDQuery(ctx context.Context, id int64) (user User, err error) {
	user = User{}
	err = db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = ?;", id)
	return user, err
}

func (db userDatabase) UpdateUser(ctx context.Context, req *UpdateUserRequest) (err error) {
	//UpdateUser business logic goes here
	if err := db.updateUser(ctx, req); err != nil {
		return ErrUpdateUser
	}
	return nil
}
func (db userDatabase) updateUser(ctx context.Context, req *UpdateUserRequest) error {
	stmt, err := db.updateUserStatement(ctx)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(req.Updates.Name, req.Updates.PasswordHash, time.Now(), req.ID)
	if err != nil {
		return err
	}
	return nil
}
func (db userDatabase) updateUserStatement(ctx context.Context) (*sql.Stmt, error) {
	tx, err := db.prepareTX(ctx)
	if err != nil {
		return nil, err
	}
	stmt, err := db.prepareStatement(ctx, tx, "UPDATE users SET (name, password_hash, last_edited) VALUES (?,?,?) WHERE id = ?;")
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (db userDatabase) CreateUser(ctx context.Context, req *CreateUserRequest) (id int64, err error) {
	dbUser := &User{
		Name:         req.Name,
		PasswordHash: req.PasswordHash,
		CreatedOn:    time.Now(),
		LastEdited:   time.Now(),
	}
	newUserID, err := db.createNewUser(ctx, dbUser)
	if err != nil {
		return 0, ErrCreateUser
	}
	return newUserID, nil
}
func (db userDatabase) createNewUser(ctx context.Context, dbUser *User) (int64, error) {
	tx, stmt, err := db.createUserStatement(ctx)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(dbUser.Name, dbUser.PasswordHash, dbUser.CreatedOn, dbUser.LastEdited)
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return res.LastInsertId()
}
func (db userDatabase) createUserStatement(ctx context.Context) (*sql.Tx, *sql.Stmt, error) {
	tx, err := db.prepareTX(ctx)
	if err != nil {
		return nil, nil, err
	}
	stmt, err := db.prepareStatement(ctx, tx, "INSERT INTO users (name, password_hash, created_on, last_edited) VALUES (?,?,?,?);")
	if err != nil {
		return nil, nil, err
	}
	return tx, stmt, nil
}

func (db userDatabase) DeleteUser(ctx context.Context, req *DeleteUserRequest) (err error) {
	if err = db.deleteUser(ctx, req); err != nil {
		return ErrDeleteUser
	}
	return nil
}
func (db userDatabase) deleteUser(ctx context.Context, req *DeleteUserRequest) error {
	stmt, err := db.deleteUserStatement(ctx)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(req.ID)
	if err != nil {
		return err
	}
	return nil
}
func (db userDatabase) deleteUserStatement(ctx context.Context) (*sql.Stmt, error) {
	tx, err := db.prepareTX(ctx)
	if err != nil {
		return nil, err
	}
	stmt, err := db.prepareStatement(ctx, tx, "DELETE FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (db userDatabase) prepareStatement(ctx context.Context, tx *sql.Tx, sql string) (*sql.Stmt, error) {
	stmt, err := tx.PrepareContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
func (db userDatabase) prepareTX(ctx context.Context) (*sql.Tx, error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}
