package dto

type OTPCreateDto struct {
	CountryCode string `json:"country_code" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
