package product

import "github.com/gin-gonic/gin"

// HttpDelivery http delivery
type HttpDelivery interface {
	CreateProduct() gin.HandlerFunc
	UpdateProduct() gin.HandlerFunc
	GetByIDProduct() gin.HandlerFunc
	SearchProduct() gin.HandlerFunc
}
