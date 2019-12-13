package datamodels

// 商品模式
type Product struct {
	ID           int64  `json:"id" sql:"ID" oscar:"ID"`
	ProductName  string `json:"productName" sql:"productName" oscar:"ProductName"`
	ProductNum   int64  `json:"productNum" sql:"productNum" oscar:"ProductNum"`
	ProductImage string `json:"productImage" sql:"productImage" oscar:"ProductImage"`
	ProductUrl   string `json:"productUrl" sql:"productUrl" oscar:"ProductUrl"`
}
