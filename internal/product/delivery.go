package product

// HttpDelivery http delivery
type HttpDelivery interface {
	CreateProduct() echo.HandlerFunc
	UpdateProduct() echo.HandlerFunc
	GetByIDProduct() echo.HandlerFunc
	SearchProduct() echo.HandlerFunc
}
