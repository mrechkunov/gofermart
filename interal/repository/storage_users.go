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
	err = su.DBconnection.QueryRow("SELECT * FROM users WHERE uLogin=$1", uLogin).Scan(&result.Login, &result.Password, &result.Bearer)
	if err != nil {
		logger.Log.Infoln(err)
	}
	return result
}

// обновление данных в БЛ
func (su *StorageUsers) UpdateUser(user model.Users) error {
	err := su.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `UPDATE users 
		SET ubearer = $1
		WHERE ulogin = $2;`
	_, err = su.DBconnection.Exec(sqlStatement, user.Bearer, user.Login)
	if err != nil {
		logger.Log.Errorln("error while update token in DB", err)
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
			(uLogin, upassword, ubearer) 
			VALUES ($1, $2, $3)`
	_, err = su.DBconnection.Exec(sqlStatement, user.Login, user.Password, user.Bearer)
	if err != nil {
		logger.Log.Errorln("error while insert to db", err)
		return err
	}
	return nil
}

func (su *StorageUsers) Close() error {
	return su.DBconnection.Close()
}

// func (d *sql.DB) SetData() {

// }

// func (d *DB) SetData(shortURL string, originalURL string, cookie string) error {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Fatal(err)
// 	}
// 	// проверяем есть ли такой URL в DB
// 	var shortURLFromDB string
// 	d.dbconn.QueryRow("SELECT shorturl FROM storage WHERE shorturl=$1", shortURL).Scan(&shortURLFromDB)
// 	if shortURLFromDB == shortURL {
// 		logger.Log.Infoln("shortURL already exist in DB")
// 		return errors.New("409 Conflict")
// 	} else {
// 		uid, _ := cryptoauth.GetIDFromCookie(cookie)
// 		sqlStatement := `INSERT INTO storage
// 			(count, uuid, originalurl, shorturl, cookie, isdeleted)
// 			VALUES ($1, $2, $3, $4, $5, $6)`
// 		_, err := d.dbconn.Exec(sqlStatement, d.counter, uid, originalURL, shortURL, cookie, false)
// 		if err != nil {
// 			logger.Log.Errorln("error while insert to db", err)
// 			return err
// 		}
// 		d.counter++
// 	}
// 	return nil
// }

// func (d *DB) GetData(shortURL string) (string, bool) {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Warnln(err)
// 	}
// 	var res string
// 	err = d.dbconn.QueryRow("SELECT originalurl FROM storage WHERE shorturl=$1", shortURL).Scan(&res)
// 	isFound := true
// 	if err != nil {
// 		isFound = false
// 	}
// 	return res, isFound
// }

// func (d *DB) IsCookieExist(cookie string) bool {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Warnln(err)
// 	}
// 	var isFound bool
// 	var queryres string
// 	err = d.dbconn.QueryRow("SELECT cookie FROM storage WHERE cookie=$1", cookie).Scan(&queryres)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			isFound = false
// 		} else {
// 			logger.Log.Infoln(err)
// 		}
// 	} else {
// 		if queryres == cookie {
// 			logger.Log.Infoln("cookie is exist in DB", queryres)
// 			isFound = true
// 		}
// 	}
// 	return isFound
// }

// func (d *DB) GetDataByUID(uid uint32) []model.ResponseDataBatchByCookie {
// 	var result []model.ResponseDataBatchByCookie

// 	rows, err := d.dbconn.Query("SELECT shortURL, originalURL FROM storage WHERE uuid=$1", uid)
// 	if err != nil {
// 		logger.Log.Fatalln(err)
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var r model.ResponseDataBatchByCookie
// 		if err := rows.Scan(&r.ShortURL, &r.OriginalURL); err != nil {
// 			logger.Log.Fatalln(err)
// 		}
// 		result = append(result, r)
// 	}
// 	if err := rows.Err(); err != nil {
// 		logger.Log.Fatalln(err)
// 	}
// 	return result
// }

// func (d *DB) IsDeleted(shortURL string) bool {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Warnln(err)
// 	}
// 	var result bool
// 	err = d.dbconn.QueryRow("SELECT isdeleted FROM storage WHERE shorturl=$1", shortURL).Scan(&result)
// 	return result
// }

// func (d *DB) IsCreator(shortURL string, cookie string) bool {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Warnln(err)
// 	}
// 	var resultCookie string
// 	err = d.dbconn.QueryRow("select cookie from storage where shorturl=$1", shortURL).Scan(&resultCookie)
// 	if err != nil {
// 		logger.Log.Warnln("error whele select cookie from DB", err)
// 	}
// 	if resultCookie == cookie {
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func (d *DB) SetIsDeleted(shortURL []string) {
// 	err := d.dbconn.Ping()
// 	if err != nil {
// 		logger.Log.Warnln(err)
// 	}
// 	sqlStatement := `UPDATE storage AS s
// 		SET isdeleted = true
// 		FROM (SELECT * FROM UNNEST($1::text[]) AS t(shorturl)) AS data
// 		WHERE s.shorturl = data.shorturl;`
// 	_, err = d.dbconn.Exec(sqlStatement, shortURL)
// 	if err != nil {
// 		logger.Log.Errorln("error while UPDATE isdeleted fiels in db", err)
// 	}
// }
