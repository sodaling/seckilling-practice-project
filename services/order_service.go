package services

import (
	"seckilling-practice-project/common"
	"seckilling-practice-project/models"
	"seckilling-practice-project/respsoiories"
)

type IOrderService interface {
	GetOrderByID(int64) (*models.Order, error)
	DeleteOrderByID(int64) bool
	UpdateOrder(*models.Order) error
	InsertOrder(*models.Order) (int64, error)
	GetAllOrder() ([]*models.Order, error)
	GetAllOrderInfo() (map[int]map[string]string, error)
	InsertOrderByMessage(*models.Message) (int64, error)
}

type OrderService struct {
	OrderRepository respsoiories.IOrderRepository
}

func (o *OrderService) GetOrderByID(orderID int64) (*models.Order, error) {
	return o.OrderRepository.SelectByKey(orderID)
}

func (o *OrderService) DeleteOrderByID(orderID int64) bool {
	return o.OrderRepository.Delete(orderID)
}

func (o *OrderService) UpdateOrder(order *models.Order) error {
	return o.OrderRepository.Update(order)
}

func (o *OrderService) InsertOrder(order *models.Order) (int64, error) {
	return o.OrderRepository.Insert(order)
}

func (o *OrderService) GetAllOrder() ([]*models.Order, error) {
	return o.OrderRepository.SelectAll()
}

func (o *OrderService) GetAllOrderInfo() (map[int]map[string]string, error) {
	return o.OrderRepository.SelectAllWithInfo()
}

func (o *OrderService) InsertOrderByMessage(message *models.Message) (int64, error) {
	order := &models.Order{
		UserId:      message.UserID,
		ProductId:   message.ProductID,
		OrderStatus: models.OrderSucess,
	}
	return o.InsertOrder(order)
}

func NewOrderService(orderRepository respsoiories.IOrderRepository) IOrderService {
	return &OrderService{OrderRepository: orderRepository}
}

func DefaultOrderService() IOrderService {
	mysqlCon, err := common.DefaultDb()
	if err != nil {
		panic(err)
	}
	orderRepo := respsoiories.NewOrderMangerRepository("order", mysqlCon)
	return NewOrderService(orderRepo)
}
