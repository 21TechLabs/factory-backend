package dto

type DtoMapKey string

const (
	DtoMapKeyPaymentPlanCreate            DtoMapKey = "PaymentPlanCreate"
	DtoMapKeyOTPCreateDto                 DtoMapKey = "OTPCreateDto"
	DtoMapKeyUserCreateDto                DtoMapKey = "UserCreateDto"
	DtoMapKeyUserCreateStep1Dto           DtoMapKey = "UserCreateStep1Dto"
	DtoMapKeyUserUpdateDto                DtoMapKey = "UserUpdateDto"
	DtoMapKeyUserPasswordUpdateDto        DtoMapKey = "UserPasswordUpdateDto"
	DtoMapKeyUserLoginDto                 DtoMapKey = "UserLoginDto"
	DtoMapKeyUserRequestPasswordResetLink DtoMapKey = "UserRequestPasswordResetLink"
	DtoMapKeyDiscordTokenExchangeResponse DtoMapKey = "DiscordTokenExchangeResponse"
	DtoMapKeyDiscordGetExchangeTokenBody  DtoMapKey = "DiscordGetExchangeTokenBody"
	DtoMapKeyDiscordUserLoginBody         DtoMapKey = "DiscordUserLoginBody"
	DtoMapKeyDiscordUserWeb               DtoMapKey = "DiscordUserWeb"
)

var DTOMap = map[DtoMapKey]func() interface{}{
	"PaymentPlanCreate":            dtoMapToRef[ProductPlanCreate](),
	"OTPCreateDto":                 dtoMapToRef[OTPCreateDto](),
	"UserCreateDto":                dtoMapToRef[UserCreateDto](),
	"UserCreateStep1Dto":           dtoMapToRef[UserCreateStep1Dto](),
	"UserUpdateDto":                dtoMapToRef[UserUpdateDto](),
	"UserPasswordUpdateDto":        dtoMapToRef[UserPasswordUpdateDto](),
	"UserLoginDto":                 dtoMapToRef[UserLoginDto](),
	"UserRequestPasswordResetLink": dtoMapToRef[UserRequestPasswordResetLink](),
	"DiscordTokenExchangeResponse": dtoMapToRef[DiscordTokenExchangeResponse](),
	"DiscordGetExchangeTokenBody":  dtoMapToRef[DiscordGetExchangeTokenBody](),
	"DiscordUserLoginBody":         dtoMapToRef[DiscordUserLoginBody](),
	"DiscordUserWeb":               dtoMapToRef[DiscordUserWeb](),
}

func dtoMapToRef[T any]() func() interface{} {
	return func() interface{} {
		var x T
		return &x
	}
}
