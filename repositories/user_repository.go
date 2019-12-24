package repositories

import (
	"database/sql"
	"errors"
	"seckill_product/common"
	"seckill_product/datamodels"
	"strconv"
)

// 定义用户接口
type IUserRepository interface {
	Connect() (err error)
	Select(userName string) (user *datamodels.User, err error)
	Insert(user *datamodels.User) (userId int64, err error)
}

// 定义用户结构体
type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

// 定义构造方法
func NewUserManagerRepository(table string, db *sql.DB) IUserRepository {
	return &UserManagerRepository{
		table:     table,
		mysqlConn: db,
	}
}

// 实现接口方法

// 数据库连接
func (u *UserManagerRepository) Connect() (err error) {
	if u.mysqlConn == nil {
		// 此处数据库连接方法需要修改
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

// 查询
func (u *UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("条件不能为空！")
	}
	if err = u.Connect(); err != nil {
		return &datamodels.User{}, err
	}
	sqlSelect := "select * from" + u.table + "where userName=?"
	rows, errRows := u.mysqlConn.Query(sqlSelect, userName)
	defer rows.Close()
	if errRows != nil {
		return &datamodels.User{}, errRows
	}
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	// 返回 user  err
	return
}

// 插入
func (u *UserManagerRepository) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Connect(); err != nil {
		return
	}
	sqlInsert := "insert" + u.table + "set nickName=?, userName=?, passWord=?"
	stmt, errStmt := u.mysqlConn.Prepare(sqlInsert)
	if errStmt != nil {
		return userId, errStmt
	}
	result, errResult := stmt.Exec(user.UserName, user.NickName, user.HashPassword)
	if errResult != nil {
		return userId, errResult
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(userId int64) (user *datamodels.User, err error) {
	sql := "select * from" + u.table + "where ID=" + strconv.FormatInt(userId, 10)
	row, errRow := u.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.User{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.User{}, errors.New("用户不存在！")
	}
	user = &datamodels.User{}
	common.DataToStructByTagSql(result, user)
	return
}
