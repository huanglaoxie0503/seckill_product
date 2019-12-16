package repositories

import (
	"database/sql"
	"seckill_product/datamodels"
)

// 定义用户接口
type IUserRepository interface {
	Connect() (err error)
	Select(userName string) (user *datamodels.User, err error)
	Insert(user datamodels.User) (userId int64, err error)
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

func (u UserManagerRepository) Connect() (err error) {
	panic("implement me")
}

func (u UserManagerRepository) Select(userName string) (user *datamodels.User, err error) {
	panic("implement me")
}

func (u UserManagerRepository) Insert(user datamodels.User) (userId int64, err error) {
	panic("implement me")
}
