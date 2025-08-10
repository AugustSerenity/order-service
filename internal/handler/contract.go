package handler

type Service interface {
	GetOrderByID(id string) ([]byte, bool)
}
