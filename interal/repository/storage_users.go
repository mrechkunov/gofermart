package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
)

type StorageUsers struct {
	DBconnection *sql.DB
}

// создаем новый сторадж для работы с таблицей пользователей
func NewUsersStorage(DBconn *sql.DB) StorageUsers {
	return StorageUsers{DBconnection: DBconn}
}

// запрос данных по логину
func (su *StorageUsers) GetUserByLogin(ctx context.Context, login string) model.Users {
	var result model.Users
	err := su.DBconnection.QueryRowContext(ctx, "SELECT u_login, u_password, u_bearer FROM users WHERE u_login=$1", login).Scan(&result.Login, &result.Password, &result.Bearer)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("user is not exist in DB")
	}
	return result
}

// запрос данных по token
func (su *StorageUsers) GetUserByToken(ctx context.Context, token string) model.Users {
	var result model.Users
	err := su.DBconnection.QueryRowContext(ctx, "SELECT u_login, u_password, u_bearer FROM users WHERE u_bearer=$1", token).Scan(&result.Login, &result.Password, &result.Bearer)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("user is not exist in DB")
	}
	return result
}

// обновление данных в БД
func (su *StorageUsers) UpdateUser(ctx context.Context, user model.Users) error {
	sqlStatement := `UPDATE users 
				SET u_bearer = $1
				WHERE u_login = $2;`
	_, err := su.DBconnection.ExecContext(ctx, sqlStatement, user.Bearer, user.Login)
	if err != nil {
		logger.Log.Errorln("error while update user token in DB", err)
		return err
	}
	return nil
}

// добавление данных в БД
func (su *StorageUsers) InsertUser(ctx context.Context, user model.Users) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	sqlStatement := `INSERT INTO users (u_login, u_password, u_bearer) 
				VALUES ($1, $2, $3)`
	_, err := su.DBconnection.ExecContext(ctxWithTimeout, sqlStatement, user.Login, user.Password, user.Bearer)
	if err != nil {
		logger.Log.Errorln("error while insert user to DB", err)
		return err
	}
	return nil
}
