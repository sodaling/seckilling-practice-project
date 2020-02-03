package respsoiories

import (
	"database/sql"
	"errors"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
)

type IUserRepository interface {
	Conn() error
	Select(userName string) (*models.User, error)
	Insert(user *models.User) (userID int64, err error)
}

type UserManagerRepository struct {
	table     string
	mysqlconn *sql.DB
}

func (u *UserManagerRepository) Conn() error {
	if u.mysqlconn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlconn = mysql
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManagerRepository) Select(userName string) (*models.User, error) {
	if userName == "" {
		return &models.User{}, errors.New("username can not be empty")
	}
	if err := u.Conn(); err != nil {
		return &models.User{}, err
	}
	sql := "SELECT * FROM " + u.table + " where userName=?"
	row, err := u.mysqlconn.Query(sql, userName)
	if err != nil {
		return &models.User{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.User{}, errors.New("user is not exist")
	}
	retUser := &models.User{}
	common.DataToStructByTagSql(result, retUser)
	return retUser, nil
}

func (u *UserManagerRepository) Insert(user *models.User) (userID int64, err error) {
	if err := u.Conn(); err != nil {
		return 0, err
	}
	sql := "INSERT " + u.table + " SET nickName=?,userName=?,passWord=?"
	stmt, err := u.mysqlconn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	result, err := stmt.Exec(user.NickName, user.UserName, user.HashPassword)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func NewUserManagerRepository(table string, mysqlconn *sql.DB) IUserRepository {
	return &UserManagerRepository{table: table, mysqlconn: mysqlconn}
}
