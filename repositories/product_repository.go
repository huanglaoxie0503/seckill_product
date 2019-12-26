package repositories

import (
	"database/sql"
	"log"
	"seckill_product/common"
	"seckill_product/datamodels"
	"strconv"
)

// 第一步：先开发对应接口
// 第二步：实现定义的接口

// 定义接口
type IProduct interface {
	// 连接数据库
	Connect() error
	// 插入数据
	Insert(*datamodels.Product) (int64, error)
	// 删除数据
	Delete(int64) bool
	// 更新数据
	Update(*datamodels.Product) error
	// 查询数据
	SelectByKey(int64) (*datamodels.Product, error)
	// 查询全部数据
	SelectAll() ([]*datamodels.Product, error)
	// 商品数量
	SubProductNum(productID int64) error
}

// 定义结构体
type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

// 定义构造函数
func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{
		table:     table,
		mysqlConn: db,
	}
}

// 以下为实现接口定义的方法

// 数据库连接
func (p *ProductManager) Connect() (err error) {
	// mysql 连接
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}

	// 判断表是否存在
	if p.table == "" {
		p.table = "product"
	}

	return
}

// 商品插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	// 判断连接是否存在
	if err = p.Connect(); err != nil {
		return
	}
	// sql语句
	sqlInsert := "insert product set productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sqlInsert)
	if errSql != nil {
		return 0, errSql
	}
	// 准备传入参数
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

// 商品删除
func (p *ProductManager) Delete(productID int64) (isOk bool) {
	// 判断连接是否存在
	if err := p.Connect(); err != nil {
		return false
	}
	sqlDel := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sqlDel)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(productID)
	if err != nil {
		return false
	}
	return true
}

// 商品更新
func (p *ProductManager) Update(product *datamodels.Product) (err error) {
	// 判断连接是否存在
	if err := p.Connect(); err != nil {
		return err
	}
	sqlUpdate := "update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sqlUpdate)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return
}

// 根据商品ID查询商品
func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	// 判断连接是否存在
	if err = p.Connect(); err != nil {
		return &datamodels.Product{}, err
	}
	sqlSelect := "select * from " + p.table + "where ID=" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sqlSelect)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	common.DataToStructByTagSql(result, productResult)
	return
}

// 获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, err error) {
	// 判断连接是否存在
	if err = p.Connect(); err != nil {
		return nil, err
	}
	log.Print("数据库连接成功")

	sqlSelect := "select * from " + p.table
	rows, err := p.mysqlConn.Query(sqlSelect)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	results := common.GetResultRows(rows)
	if len(results) == 0 {
		return nil, nil
	}

	for _, v := range results {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}

// 商品数量扣除
func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Connect(); err != nil {
		return err
	}
	sqlUpdate := "update " + p.table + " set " + " productNum=productNum-1 where ID =" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sqlUpdate)
	defer stmt.Close()
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
