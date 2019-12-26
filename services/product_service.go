package services

import (
	"seckill_product/datamodels"
	"seckill_product/repositories"
)

// 定义商品服务接口
type IProductService interface {
	GetProductByID(productID int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(productID int64) (isOK bool)
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) (err error)
	SubNumberOne(productID int64) error
}

// 定义商品服务结构体
type ProductService struct {
	productRepository repositories.IProduct
}

// 初始化函数
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{productRepository: repository}
}

// 按商品ID查询数据
func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productID)
}

// 查询所有商品
func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

// 按商品 ID 删除数据
func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

// 插入数据
func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

// 更新数据
func (p *ProductService) UpdateProduct(product *datamodels.Product) (err error) {
	return p.productRepository.Update(product)
}

// 商品数量扣除
func (p *ProductService) SubNumberOne(productID int64) error {
	return p.productRepository.SubProductNum(productID)
}
