package dto

type UserCreateDto struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	CountryCode     string `json:"country_code" validate:"required"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
	PhoneNumberOtp  string `json:"phone_number_otp" validate:"required"`
}

// Admin DTOs
type UserCreateStep1Dto struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UserUpdateStep2Dto struct {
	CountryCode    string `json:"country_code"`
	PhoneNumber    string `json:"phone_number"`
	PhoneNumberOtp string `json:"phone_number_otp"`
}

type UserUpdateStep3Dto struct {
	UpdateToken     string `json:"update_token"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserUpdateDto struct {
	Name                     string `json:"name"`
	Email                    string `json:"email"`
	CountryCode              string `json:"country_code"`
	PhoneNumber              string `json:"phone_number"`
	Password                 string `json:"password"`
	OptedInForEmailAlerts    bool   `json:"opted_in_for_email_alerts"`
	OptedInForSmsAlerts      bool   `json:"opted_in_for_sms_alerts"`
	OptedInForWhatsappAlerts bool   `json:"opted_in_for_whatsapp_alerts"`
}

type UserPasswordUpdateDto struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type UserLoginDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRequestPasswordResetLink struct {
	Email string `json:"email" validate:"required,email"`
}
