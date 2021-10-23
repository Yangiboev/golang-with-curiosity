package product

import (
	"context"

	"github.com/Yangiboev/golang-with-curiosity/internal/models"
	"github.com/Yangiboev/golang-with-curiosity/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UseCase Product
type UseCase interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetByID(ctx context.Context, productID primitive.ObjectID) (*models.Product, error)
	Search(ctx context.Context, search string, pagination *utils.PaginationUC) (*models.ProductsList, error)
	PublishCreate(ctx context.Context, product *models.Product) error
	PublishUpdate(ctx context.Context, product *models.Product) error
}
