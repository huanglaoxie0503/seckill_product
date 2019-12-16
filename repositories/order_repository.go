package repositories

import (
	"database/sql"
	"seckill_product/common"
	"seckill_product/datamodels"
	"strconv"
)

// 定义订单接口
type IOrderRepository interface {
	Connect() (err error)
	// 插入
	Insert(order *datamodels.Order) (orderID int64, err error)
	// 删除
	Delete(productID int64) bool
	// 更新
	Update(order *datamodels.Order) (err error)
	// 根据商品ID 查询商品信息
	SelectByKey(orderID int64) (order *datamodels.Order, err error)
	// 查询所有商品
	SelectAll() (orderArray []*datamodels.Order, err error)
	// 查询订单相关信息
	SelectAllWithInfo() (map[int]map[string]string, error)
}

// 为实现接口创建打结构体
type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

// 实现接口方法的构造函数
func NewOrderManagerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderManagerRepository{
		table:     table,
		mysqlConn: sql,
	}
}

// 数据库连接
func (o *OrderManagerRepository) Connect() (err error) {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

// 订单插入
func (o *OrderManagerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Connect(); err != nil {
		return
	}

	insertSql := "insert" + o.table + "set userID=?,productID=?,orderStatus=?"
	stmt, errStmt := o.mysqlConn.Prepare(insertSql)
	if errStmt != nil {
		return productID, errStmt
	}
	result, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

// 订单删除
func (o *OrderManagerRepository) Delete(productID int64) (isOK bool) {
	if err := o.Connect(); err != nil {
		return
	}
	deleteSql := "delete from" + o.table + "where ID=?"
	stmt, errStmt := o.mysqlConn.Prepare(deleteSql)
	if errStmt != nil {
		return
	}
	_, err := stmt.Exec(productID)
	if err != nil {
		return
	}
	return true
}

// 更新
func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Connect(); errConn != nil {
		return errConn
	}
	updateSql := "update" + o.table + "set userID=?, productID=?,orderStatus=? where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, errStmt := o.mysqlConn.Prepare(updateSql)
	if errStmt != nil {
		return errStmt
	}
	_, errResult := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if errResult != nil {
		return errResult
	}
	return
}

// 根据订单ID，查询订单
func (o *OrderManagerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if errConn := o.Connect(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	querySql := "select * from" + o.table + "where ID=" + strconv.FormatInt(orderID, 10)
	row, errRow := o.mysqlConn.Query(querySql)
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

// 查询所有订单
func (o *OrderManagerRepository) SelectAll() (orderArray []*datamodels.Order, err error) {
	if errConn := o.Connect(); errConn != nil {
		return nil, errConn
	}
	selectSql := "select * from " + o.table
	rows, errRows := o.mysqlConn.Query(selectSql)
	if errRows != nil {
		return nil, errRows
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, err
	}
	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

// 查询订单相关信息
func (o *OrderManagerRepository) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	if errConn := o.Connect(); errConn != nil {
		return nil, errConn
	}
	//orderSql := "select o.ID, p.productName, o.orderStatus from seckill_product.order as o left join product as p on o.product=p.ID"
	orderSql := "Select o.ID,p.productName,o.orderStatus From seckill_product.order as o left join product as p on o.productID=p.ID"
	rows, errRows := o.mysqlConn.Query(orderSql)
	if errRows != nil {
		return nil, errRows
	}
	return common.GetResultRows(rows), err
}
