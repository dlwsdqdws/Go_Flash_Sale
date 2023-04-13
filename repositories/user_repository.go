package repositories

import (
	"database/sql"
	"errors"
	"pro-iris/common"
	"pro-iris/datamodels"
	"strconv"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

func NewUserManagerRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{table, db}
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (u *UserManagerRepository) Conn() (err error) {
	if u.mysqlConn == nil {
		mysql, errSql := common.NewMysqlConn()
		if errSql != nil {
			return errSql
		}
		u.mysqlConn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return
}

func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("userName cannot be empty")
	}
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "Select * from " + u.table + " where userName=?"
	row, errRow := u.mysqlConn.Query(sql, userName)
	defer row.Close()
	if errRow != nil {
		return &datamodels.User{}, errRow
	}

	result := common.GetResultRow(row)

	if len(result) == 0 {
		return &datamodels.User{}, errors.New("user doesn't exist")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}

func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return
	}
	sql := "INSERT " + u.table + " SET nickName=?, userName=?, password=?"
	stmt, errStmt := u.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return userId, errStmt
	}
	result, errRes := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if errRes != nil {
		return userId, errRes
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	if err = u.Conn(); err != nil {
		return &datamodels.User{}, err
	}
	sql := "select * from " + u.table + " where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("user doesn't exist")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
