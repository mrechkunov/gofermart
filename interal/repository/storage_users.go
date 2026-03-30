package repository

import (
	"database/sql"

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
	var su StorageUsers
	su.DBconnection = DBconn
	return su
}

// запрос данных по логину
func (su *StorageUsers) GetByLogin(uLogin string) model.Users {
	var result model.Users
	err := su.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	err = su.DBconnection.QueryRow("SELECT u_login, u_password, u_bearer FROM users WHERE u_login=$1", uLogin).Scan(&result.Login, &result.Password, &result.Bearer)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("user is not exist in DB")
	}
	return result
}

// запрос данных по token
func (su *StorageUsers) GetByToken(token string) model.Users {
	var result model.Users
	err := su.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	err = su.DBconnection.QueryRow("SELECT u_login, u_password, u_bearer FROM users WHERE u_bearer=$1", token).Scan(&result.Login, &result.Password, &result.Bearer)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("Запись не найдена")
	}
	return result
}

// обновление данных в БД
func (su *StorageUsers) UpdateUser(user model.Users) error {
	err := su.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `UPDATE users 
		SET u_bearer = $1
		WHERE u_login = $2;`
	_, err = su.DBconnection.Exec(sqlStatement, user.Bearer, user.Login)
	if err != nil {
		logger.Log.Errorln("error while update user token in DB", err)
		return err
	}
	return nil
}

// добавление данных в БД
func (su *StorageUsers) InsertUser(user model.Users) error {
	err := su.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `INSERT INTO users 
			(u_login, u_password, u_bearer) 
			VALUES ($1, $2, $3)`
	_, err = su.DBconnection.Exec(sqlStatement, user.Login, user.Password, user.Bearer)
	if err != nil {
		logger.Log.Errorln("error while insert user to db", err)
		return err
	}
	return nil
}
