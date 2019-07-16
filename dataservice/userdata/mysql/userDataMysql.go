// Package mysql represents the mySQL implementation of the user data persistence layer
package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jfeng45/servicetmpl/container/logger"
	"github.com/jfeng45/servicetmpl/model"
	"github.com/jfeng45/servicetmpl/tools"
	"github.com/jfeng45/servicetmpl/tools/gdbc"
	"github.com/pkg/errors"
	"time"
)

const (
	DELETE_USER string ="delete from userinfo where username=?"
	QUERY_USER_BY_ID string ="SELECT * FROM userinfo where uid =?"
	QUERY_USER_BY_NAME = "SELECT * FROM userinfo where username =?"
	QUERY_USER = "SELECT * FROM userinfo "
	UPDATE_USER = "update userinfo set username=?, department=?, created=? where uid=?"
    INSERT_USER = "INSERT userinfo SET username=?,department=?,created=?"
)
// UserDataMySql is the MySQL implementation of UserDatainterface
type UserDataMySql struct {
	DB gdbc.Gdbc
}

func (userData *UserDataMySql) Remove(username string) (int64, error) {

	stmt, err := userData.DB.Prepare(DELETE_USER)
	if err!=nil {
		return 0, errors.Wrap(err, "")
	}
	defer stmt.Close()

	res, err := stmt.Exec(username)
	if err!=nil {
		return 0, errors.Wrap(err, "")
	}
	rowsAffected, err := res.RowsAffected()
	if err!=nil {
		return 0, errors.Wrap(err, "")
	}

	logger.Log.Debug("remove:row affected ", rowsAffected)
	return rowsAffected, nil
}

func (userData *UserDataMySql) Find(id int) (*model.User, error) {
	rows, err := userData.DB.Query(QUERY_USER_BY_ID, id)
	if err !=nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrieveUser(rows)
}
func retrieveUser(rows *sql.Rows) (*model.User, error) {
	if rows.Next() {
		return rowsToUser(rows)
	}
	return nil, nil
}
func rowsToUser(rows *sql.Rows) (*model.User, error) {
	var ds string
	user := &model.User{}
	err := rows.Scan(&user.Id, &user.Name, &user.Department, &ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	created, err := time.Parse(tools.FORMAT_ISO8601_DATE, ds)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	user.Created = created

	logger.Log.Debug("rows to User:", user)
	return user, nil
}
func (userData *UserDataMySql) FindByName(name string) (*model.User, error) {
	//logger.Log.Debug("call FindByName() and name is:", name)
	rows, err := userData.DB.Query(QUERY_USER_BY_NAME, name)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	return retrieveUser(rows)
}

func (userData *UserDataMySql) FindAll() ([]model.User, error) {

	rows, err := userData.DB.Query(QUERY_USER)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer rows.Close()
	users := []model.User{}

	//var ds string
	for rows.Next() {
		user, err :=rowsToUser(rows)
		if err != nil {
			return users, errors.Wrap(err,"")
		}
		users = append(users, *user)

	}
	//need to check error for rows.Next()
	if err =rows.Err(); err!= nil {
		return nil, errors.Wrap(err, "")
	}
	logger.Log.Debug("find user list:", users)
	return users, nil
}

func (userData *UserDataMySql) Update(user *model.User) (int64, error) {

	stmt, err := userData.DB.Prepare(UPDATE_USER)

	if err!=nil {
		return 0, errors.Wrap(err, "")
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.Department, user.Created, user.Id)
	if err!=nil {
		return 0, errors.Wrap(err, "")
	}
	rowsAffected, err := res.RowsAffected()

	if err!=nil {
		return 0, errors.Wrap(err, "")
	}
	logger.Log.Debug("update: rows affected: ", rowsAffected)

	return rowsAffected, nil
}

func (userData *UserDataMySql) Insert(user *model.User) (*model.User, error) {

	stmt, err := userData.DB.Prepare(INSERT_USER)
	if err!=nil {
		return nil, errors.Wrap(err, "")
	}
	defer stmt.Close()
	res, err := stmt.Exec(user.Name, user.Department, user.Created)
	if err!=nil {
		return nil, errors.Wrap(err, "")
	}
	id, err := res.LastInsertId()
	if err!=nil {
		return nil, errors.Wrap(err, "")
	}
	user.Id = int(id)
	logger.Log.Debug("user inserted:", user)
	return user, nil
}

