package respsoiories

import (
	"database/sql"
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"strconv"
)

type Iproduct interface {
	Conn() error
	Insert(*models.Product) (int64, error)
	Delete(int64) bool
	Update(*models.Product) error
	SelectByKey(int64) (*models.Product, error)
	SelectAll() ([]*models.Product, error)
	SubProductNum(productID int64) error
}

type ProductManger struct {
	table     string
	mysqlConn *sql.DB
}

func (p *ProductManger) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return nil
}

func (p *ProductManger) Insert(product *models.Product) (int64, error) {
	if err := p.Conn(); err != nil {
		return 0, err
	}
	sql := "INSERT product set productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (p *ProductManger) Delete(productID int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	defer stmt.Close()
	return true
}

func (p *ProductManger) Update(product *models.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "UPDATE product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	return err
}

func (p *ProductManger) SelectByKey(productID int64) (*models.Product, error) {
	if err := p.Conn(); err != nil {
		return &models.Product{}, err
	}
	sql := "SELECT * FROM product where ID=" + strconv.FormatInt(productID, 10)
	row, err := p.mysqlConn.Query(sql)
	if err != nil {
		return &models.Product{}, err
	}
	defer row.Close()
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.Product{}, nil
	}
	productResult := &models.Product{}
	common.DataToStructByTagSql(result, productResult)
	return productResult, nil
}

func (p *ProductManger) SelectAll() ([]*models.Product, error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "SELECT * FROM " + p.table
	rows, err := p.mysqlConn.Query(sql)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	result := common.GetResultRows(rows)
	var productArray []*models.Product
	for _, v := range result {
		product := &models.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return productArray, nil
}

func (p *ProductManger) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "UPDATE " + p.table + " SET productNum=productNum-1 where ID =" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(productID)
	return err
}

func NewProductManager(table string, db *sql.DB) Iproduct {
	return &ProductManger{table: table, mysqlConn: db}
}
