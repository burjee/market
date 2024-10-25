package structure

type Order struct {
	Type     string `json:"type"     binding:"required,oneof=buy sell"`
	Price    int    `json:"price"    binding:"required,number,gte=1,lte=99999"`
	Quantity int    `json:"quantity" binding:"required,number,gte=1,lte=99999"`
}
