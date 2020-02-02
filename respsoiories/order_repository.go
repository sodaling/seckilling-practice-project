package respsoiories

import (
	"database/sql"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*models.Order) (int64, error)
	Delete(int64) bool
	Update(*models.Order) error
	SelectByKey(int64) (*models.Order, error)
	SelectAll() ([]*models.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func (o *OrderMangerRepository) SelectByKey(orderID int64) (*models.Order, error) {
	if err := o.Conn(); err != nil {
		return &models.Order{}, nil
	}

	sql := "Select * from " + o.table + " where ID = " + strconv.FormatInt(orderID, 10)
	row, err := o.mysqlConn.Query(sql)
	defer row.Close()
	if err != nil {
		return &models.Order{}, nil
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.Order{}, nil
	}
	order := &models.Order{}
	common.DataToStructByTagSql(result, order)
	return order, nil
}

func (o *OrderMangerRepository) SelectAll() ([]*models.Order, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}
	var orderArray = make([]*models.Order, len(result))
	for _, v := range result {
		order := &models.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return orderArray, nil
}

func (o *OrderMangerRepository) Conn() error {
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

func (o *OrderMangerRepository) Insert(order *models.Order) (int64, error) {
	if err := o.Conn(); err != nil {
		return 0, nil
	}
	sql := "INSERT 'order' SET userID=?,productID=?,orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return 0, nil
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return 0, nil
	}
	return result.LastInsertId()
}

func (o *OrderMangerRepository) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}
	sql := "DELETE FROM " + o.table + " WHERE orderID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		return false
	}
	_, err = stmt.Exec(orderID)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderMangerRepository) Update(order *models.Order) error {
	if err := o.Conn(); err != nil {
		return err
	}
	sql := "update " + o.table + "set userID=?,productID=?,orderStatus=? where ID =" + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.ID)
	return nil
}

func (o *OrderMangerRepository) SelectAllWithInfo() (map[int]map[string]string, error) {
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "Select o.ID,p.productName,o.orderStatus From miaosha.order as o left join product as p on o.productID=p.ID"
	rows,err := o.mysqlConn.Query(sql)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	return common.GetResultRows(rows),nil
}

func NewOrderMangerRepository(table string, mysqlConn *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: mysqlConn}
}
