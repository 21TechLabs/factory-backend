package dto

// CreateProductDto
type CreateProductDto struct {
	ProductId int `json:"product_id" validate:"required"`
	PlanIdx   int `json:"plan_idx"`
}
