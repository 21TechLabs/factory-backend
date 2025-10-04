package models

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/21TechLabs/factory-backend/dto"
	"github.com/21TechLabs/factory-backend/notifications"
	"github.com/21TechLabs/factory-backend/notifications/templates"
	"github.com/21TechLabs/factory-backend/utils"
	"github.com/kataras/jwt"
	"gorm.io/gorm"
)

type UserRole int

const (
	UserRoleAdmin  UserRole = iota + 1 // 1
	UserRoleClient UserRole = iota + 1 // 2
)

type UserStore struct {
	DB        *gorm.DB
	FileStore *FileStore
}

func NewUserStore(db *gorm.DB, fs *FileStore) *UserStore {
	return &UserStore{DB: db, FileStore: fs}
}

type User struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name                   string    `gorm:"column:name" json:"name"`
	Role                   UserRole  `gorm:"column:role" json:"role"`
	Email                  string    `gorm:"column:email;uniqueIndex" json:"email"`
	ProfilePicURI          string    `gorm:"column:profile_picture_url" json:"profilePicURI"`
	EmailVerified          bool      `gorm:"column:email_verified" json:"emailVerified"`
	EmailVerificationToken string    `gorm:"column:email_verification_token" json:"-"`
	Password               string    `gorm:"column:password" json:"-"`
	PasswordResetToken     string    `gorm:"column:password_reset_token" json:"-"`
	PasswordTries          int       `gorm:"column:password_tries" json:"passwordTries"`
	OptedInForEmail        bool      `gorm:"column:opted_in_for_email" json:"optedInForEmail"`
	AccountSuspended       bool      `gorm:"column:account_suspended" json:"accountSuspended"`
	AccountBlocked         bool      `gorm:"column:account_blocked" json:"accountBlocked"`
	MarkedForDeletion      bool      `gorm:"column:marked_for_deletion" json:"markedForDeletion"`
	DeleteAccountAfter     time.Time `gorm:"column:delete_account_after" json:"deleteAccountAfter"`
	AccountDeleted         bool      `gorm:"column:account_deleted" json:"accountDeleted"`
	AccountCreated         bool      `gorm:"column:account_created" json:"accountCreated"`
	Tokens                 int64     `gorm:"column:tokens" json:"tokens"`
	Files                  []File    `gorm:"foreignKey:UserID;references:ID" json:"files"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}

func (us *UserStore) UserCreate(user dto.UserCreateDto) (User, error) {

	var userCount int64
	result := us.DB.Model(&User{}).Where("email = ?", user.Email).Count(&userCount)

	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return User{}, result.Error
		}
	}

	if userCount > 0 {
		return User{}, fmt.Errorf("user with this email already exists")
	}

	var newUser = User{
		Name:            user.Name,
		Role:            UserRoleClient,
		Email:           user.Email,
		OptedInForEmail: true,
		AccountCreated:  true,
	}
	var err error
	newUser.Password, err = SaltPassword(user.Password, "")

	if err != nil {
		return User{}, err
	}

	result = us.DB.Create(&newUser)

	if result.Error != nil {
		return User{}, result.Error
	}

	// send email verification msg
	err = us.SendEmailVerifyEmail(&newUser)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (user *User) UserIsAdmin() bool {
	return user.Role == UserRoleAdmin
}

func (user *User) UserIsClient() bool {
	return user.Role == UserRoleClient
}

func (uc *UserStore) JwtTokenVerifyAndGetUser(token string, secretKey []byte) (User, error) {
	verifiedToken, err := jwt.Verify(jwt.HS256, secretKey, []byte(token))
	if err != nil {
		return User{}, err
	}

	var userPD User
	err = verifiedToken.Claims(&userPD)
	if err != nil {
		return User{}, err
	}

	var user User
	user, err = uc.UserGetByEmail(userPD.Email)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (us *UserStore) GetDetails(u *User, allowPasswordResetToken bool) User {
	return *u
}

func (us *UserStore) SendEmailVerifyEmail(u *User) error {
	// generate token
	token, err := GetAlphaNumString(64, "alnum")
	if err != nil {
		return err
	}
	// save token
	u.EmailVerificationToken = token
	// send email
	var frontendURL = utils.GetEnv("FRONTEND_URL", false)
	var req = notifications.NewRequest([]string{u.Email}, fmt.Sprintf("Welcome to the family %s.", u.Name), fmt.Sprintf("Hey %s, we are glad that you joined our family, to verify your email visit %s/verify-email?email=%s&&token=%s", u.Name, frontendURL, u.Email, token))

	var template templates.WelcomeMessage = templates.WelcomeMessage{
		Name:      u.Name,
		BrandName: "",
		Link:      fmt.Sprintf("%s/verify-email?email=%s&&token=%s", frontendURL, u.Email, token),
	}
	err = template.ParseAsHTML(req)

	if err == nil {
		go func() {
			_, err := req.SendEmail()
			if err != nil {
				log.Printf("Failed to send email to %s", u.Email)
			}
		}()
	}

	return nil
}

func (us *UserStore) Update(u *User) error {

	var model = us.DB.Model(&User{})

	result := model.Where("id = ?", u.ID).Updates(u)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return fmt.Errorf("user with id %d not found", u.ID)
		}
		return fmt.Errorf("failed to update user with id %d: %v", u.ID, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", u.ID)
	}

	return nil
}

func (us *UserStore) sendPasswordResetEmail(u *User, token string) error {
	var frontendURL = utils.GetEnv("FRONTEND_URL", false)

	var req = notifications.NewRequest([]string{u.Email}, fmt.Sprintf("%s, your password reset email.", u.Name), fmt.Sprintf("to reset the password please visit %s/reset-password?email=%s&&token=%s", frontendURL, u.Email, token))

	var template templates.ResetPasswordMessage = templates.ResetPasswordMessage{
		Name:      u.Name,
		BrandName: "",
		Link:      fmt.Sprintf("%s/reset-password?email=%s&&token=%s", frontendURL, u.Email, token),
	}
	template.ParseAsHTML(req)
	go func() {
		_, err := req.SendEmail()

		if err != nil {
			log.Default().Panicf("Failed to send email to %s", u.Email)
		}
	}()
	return nil
}

func (us *UserStore) UserGetById(id uint) (User, error) {
	var userFromDb User = User{}

	result := us.DB.Model(&User{}).Where("id = ?", id).First(&userFromDb)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return User{}, nil
		}
		return User{}, fmt.Errorf("failed to find the user with id: %d \n %s", id, result.Error.Error())
	}

	return userFromDb, nil
}

func GetRandomNumber(limit int64) (int64, error) {
	v, err := rand.Int(rand.Reader, big.NewInt(limit))
	if err != nil {
		return 0, err
	}
	return v.Int64(), nil
}

func GetAlphaNumString(length int, stringContents string) (string, error) {
	var alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	switch stringContents {
	case "alphanum":
		alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "alpha":
		alphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case "num":
		alphaNum = "0123456789"
	}

	var alphaNumLen = len(alphaNum)
	var result string

	for i := 0; i < length; i++ {
		randomNum, err := GetRandomNumber(int64(alphaNumLen))

		if err != nil {
			return "", err
		}

		result += string(alphaNum[randomNum])
	}

	return result, nil
}

func SaltPassword(password string, salt string) (string, error) {
	var err error

	if len(salt) == 0 {
		salt, err = GetAlphaNumString(36, "alphanum")
		if err != nil {
			return "", err
		}
	}

	password = strings.ReplaceAll(password, "", salt)
	sh512 := sha512.New()
	_, err = sh512.Write([]byte(password))

	if err != nil {
		return "", err
	}

	var sha512_hash = hex.EncodeToString(sh512.Sum(nil))

	password = fmt.Sprintf("%s.%s", sha512_hash, salt)
	return password, nil
}

func (us *UserStore) ComparePassword(user *User, password string) bool {
	var password_and_salt []string = strings.Split(user.Password, ".")

	var current_passwd = password_and_salt[0]
	var salt = password_and_salt[1]

	hased_password, err := SaltPassword(password, salt)

	if err != nil {
		return false
	}

	var password_and_salt2 []string = strings.Split(hased_password, ".")

	if len(password_and_salt2) != 2 {
		return false
	}

	var to_compare = password_and_salt2[0]

	return current_passwd == to_compare
}

func (us *UserStore) GeneratePasswordResetToken(user *User, sendEmail bool) (token string, err error) {
	resetToken, err := GetAlphaNumString(56, "alphanum")

	if err != nil {
		return "", err
	}

	user.PasswordResetToken, err = SaltPassword(resetToken, "")

	if err != nil {
		return "", err
	}

	err = us.Update(user)

	if err != nil {
		return "", err
	}

	if sendEmail {
		err = us.sendPasswordResetEmail(user, resetToken)
		if err != nil {
			return "", err
		}
	}

	return resetToken, err
}

func (us *UserStore) CompareAndUpdatePasswordWithToken(user *User, token string, password string) error {
	var token_and_salt []string = strings.Split(user.PasswordResetToken, ".")

	var salt = token_and_salt[1]

	salted_token, err := SaltPassword(token, salt)

	if err != nil {
		return err
	}

	if user.PasswordResetToken != salted_token {
		return fmt.Errorf("invalid token")
	}

	user.Password, err = SaltPassword(password, "")

	if err != nil {
		return err
	}

	user.PasswordResetToken = ""

	err = us.Update(user)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) UserGetByEmail(email string) (User, error) {
	// cursor, err := mgm.Coll(&User{}).Find(ctx, map[string]interface{}{"email": email, "accountDeleted": false, "markedForDeletion": false})
	var user User

	result := us.DB.Model(&User{}).Where("email = ?", email).First(&user)

	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil

}

func (user *User) JwtTokenGet(expiryTime time.Time, secretKey []byte) (string, error) {
	claim := *user

	claim.PasswordResetToken = ""
	claim.EmailVerificationToken = ""
	claim.Password = ""

	token, err := jwt.Sign(jwt.HS256, secretKey, claim, jwt.MaxAge(time.Hour*24*7))

	if err != nil {
		return "", err
	}

	return string(token), nil
}

func (us *UserStore) UserVerifyEmailToken(email string, token string) (User, error) {
	user, err := us.UserGetByEmail(email)

	if err != nil {
		return user, err
	}

	if user.EmailVerificationToken != token {
		return user, fmt.Errorf("invalid token")
	}

	user.EmailVerified = true
	user.EmailVerificationToken = ""

	err = us.Update(&user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserStore) MarkAccountForDeletion(user *User) error {
	user.MarkedForDeletion = true

	user.DeleteAccountAfter = time.Now().Add(time.Hour * 24 * 30)
	user.MarkedForDeletion = true

	err := us.Update(user)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) UploadFile(user *User, data []FileUpload) ([]File, error) {
	return us.FileStore.UploadFile(data, user.ID)
}

func (us *UserStore) UserLogin(loginDto dto.UserLoginDto) (User, error) {
	user, err := us.UserGetByEmail(loginDto.Email)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: Failed to find user\n%v", err)
		return user, err
	}
	if user.AccountBlocked {
		fmt.Fprintln(os.Stdout, "UserLogin Error: Account Blocked!")
		return user, fmt.Errorf("account blocked")
	}

	if user.AccountDeleted {
		fmt.Fprintln(os.Stdout, "UserLogin Error: Account Deleted!")
		return user, fmt.Errorf("account deleted")
	}

	if user.AccountSuspended {
		fmt.Fprintln(os.Stdout, "UserLogin Error: Account Suspended!")
		return user, fmt.Errorf("account suspended")
	}

	if user.PasswordTries >= 5 {
		fmt.Fprintln(os.Stdout, "UserLogin Error: Account Blocked!")
		return user, fmt.Errorf("account blocked")
	}

	isCorrectPassword := us.ComparePassword(&user, loginDto.Password)

	if !isCorrectPassword {

		user.PasswordTries++

		if user.PasswordTries >= 5 {
			user.AccountBlocked = true
		}

		err = us.Update(&user)

		if err != nil {
			fmt.Fprintf(os.Stdout, "UserLogin Error: Failed to update user\n%v", err)
			return user, err
		}

		fmt.Fprintln(os.Stdout, "UserLogin Error: Incorrect Password!")
		return user, fmt.Errorf("incorrect password")
	}

	user.PasswordTries = 0

	err = us.Update(&user)

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: Failed to update user\n%v", err)
		return user, err
	}

	return user, nil
}

func (us *UserStore) UserGetAllBy(filter map[string]interface{}, start, limit int) ([]User, error) {
	var users []User

	result := us.DB.Model(&User{}).Where(filter).Offset(start).Limit(limit).Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to find users: %v", result.Error)
	}

	return users, nil
}
