package utils

import "github.com/go-playground/validator/v10"

func ValidateStruct(structToValidate interface{}) error {
	return validator.New().Struct(structToValidate)
}

type DTOGtLt struct {
	Gte int64 `json:"gte" validate:"omitempty,gte=0"`
	Lte int64 `json:"lte" validate:"omitempty,gte=0"`
}

type SortBy struct {
	Field     string `json:"field" validate:"required"`
	Direction string `json:"direction" validate:"required,oneof=asc desc"`
}
