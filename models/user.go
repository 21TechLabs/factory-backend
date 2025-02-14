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

	"github.com/21TechLabs/factory-be/dto"
	"github.com/21TechLabs/factory-be/models/payments"
	"github.com/21TechLabs/factory-be/notifications"
	"github.com/21TechLabs/factory-be/notifications/templates"
	"github.com/21TechLabs/factory-be/utils"
	"github.com/kamva/mgm/v3"
	"github.com/kataras/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type _Role struct {
	Admin   string
	Student string
	Client  string
}

var Roles = _Role{
	Admin:   "admin",
	Student: "student",
	Client:  "client",
}

type User struct {
	mgm.DefaultModel       `bson:",inline"`
	Name                   string    `bson:"name" json:"name"`                                     // User's name
	Role                   string    `bson:"role" json:"role"`                                     // User's role (e.g. "admin", "student")
	Email                  string    `bson:"email" json:"email"`                                   // User's email address
	ProfilePicURI          string    `bson:"profilePic" json:"profilePic"`                         // URI for the user's profile picture
	EmailVerified          bool      `bson:"emailVerified" json:"emailVerified"`                   // Whether the user's email has been verified
	EmailVerificationToken string    `bson:"emailVerificationToken" json:"emailVerificationToken"` // Token for email verification
	Password               string    `bson:"password" json:"password"`                             // User's password (hashed)
	PasswordResetToken     string    `bson:"passwordResetToken" json:"passwordResetToken"`         // Token for password reset
	PasswordTries          int       `bson:"passwordTries" json:"passwordTries"`                   // Number of tries for password reset
	OptedInForEmail        bool      `bson:"optedInForEmail" json:"optedInForEmail"`               // Whether the user has opted in for email alerts
	AccountSuspended       bool      `bson:"accountSuspended" json:"accountSuspended"`             // Whether the account is suspended
	AccountBlocked         bool      `bson:"accountBlocked" json:"accountBlocked"`                 // Whether the account is blocked
	MarkedForDeletion      bool      `bson:"markedForDeletion" json:"markedForDeletion"`           // Whether the account is marked for deletion
	DeleteAccountAfter     time.Time `bson:"deleteAccountAfter" json:"deleteAccountAfter"`         // Date/time after which the account will be deleted
	AccountDeleted         bool      `bson:"accountDeleted" json:"accountDeleted"`                 // Whether the account is deleted4
	AccountCreated         bool      `bson:"accountCreated" json:"accountCreated"`                 // Date/time when the account was created
	CoinBalance            int64     `bson:"coinBalance" json:"coinBalance"`                       // User's coin balance
}

func UserCreate(user dto.UserCreateDto, role string) (User, error) {

	var ctx = mgm.Ctx()
	userCount, err := mgm.Coll(&User{}).CountDocuments(ctx, bson.M{
		"email": user.Email,
		"role":  Roles.Admin,
	})

	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return User{}, err
		}
	}

	if userCount > 0 {
		return User{}, fmt.Errorf("user with this email already exists")
	}

	var newUser = User{
		Name:            user.Name,
		Role:            role,
		Email:           user.Email,
		OptedInForEmail: true,
		AccountCreated:  true,
	}

	newUser.Password, err = SaltPassword(user.Password, "")

	if err != nil {
		return User{}, err
	}

	// send email verification msg
	err = newUser.SendEmailVerifyEmail()
	if err != nil {
		return User{}, err
	}

	err = mgm.Coll(&newUser).Create(&newUser)

	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func JwtTokenVerifyAndGetUser(token string, secretKey []byte) (User, error) {
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
	user, err = UserGetByEmail(userPD.Email)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *User) GetDetails(allowPasswordResetToken bool) User {
	usr := User{
		Name:            u.Name,
		Role:            u.Role,
		Email:           u.Email,
		ProfilePicURI:   u.ProfilePicURI,
		EmailVerified:   u.EmailVerified,
		OptedInForEmail: u.OptedInForEmail,
		CoinBalance:     u.CoinBalance,
	}

	usr.ID = u.ID

	if allowPasswordResetToken {
		usr.Password = u.Password
		usr.PasswordResetToken = u.PasswordResetToken
	}

	return usr
}

func (u *User) SendEmailVerifyEmail() error {
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

func (u *User) sendPasswordResetEmail(token string) error {
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

func UserGetById(UserId string) (User, error) {
	var userFromDb User = User{}

	cur := mgm.Coll(&User{}).First(bson.M{
		"_id": UserId,
	}, &userFromDb)

	if cur.Error() != "" {
		if cur.Error() == "mongo: no documents in result" {
			return User{}, nil
		}

		return User{}, fmt.Errorf("failed to find the user user with id: %s \n %s", UserId, cur.Error())
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

	if stringContents == "alpha" {
		alphaNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	} else if stringContents == "num" {
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

func (user *User) ComparePassword(password string) bool {
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

func (user *User) GeneratePasswordResetToken(sendEmail bool) (token string, err error) {
	resetToken, err := GetAlphaNumString(56, "alphanum")

	if err != nil {
		return "", err
	}

	user.PasswordResetToken, err = SaltPassword(resetToken, "")

	if err != nil {
		return "", err
	}

	ctx := mgm.Ctx()
	_, err = mgm.Coll(user).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

	if err != nil {
		return "", err
	}

	if sendEmail {
		err = user.sendPasswordResetEmail(resetToken)
		if err != nil {
			return "", err
		}
	}

	return resetToken, err
}

func (user *User) CompareAndUpdatePasswordWithToken(token string, password string) error {
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

	ctx := mgm.Ctx()
	_, err = mgm.Coll(user).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

	if err != nil {
		return err
	}

	return nil
}

func UserGetByEmail(email string) (User, error) {
	ctx := mgm.Ctx()
	cursor, err := mgm.Coll(&User{}).Find(ctx, bson.M{"email": email, "accountDeleted": false, "markedForDeletion": false})

	if err != nil {
		return User{}, err
	}

	var users []User

	err = cursor.All(ctx, &users)

	if err != nil {
		return User{}, err
	}

	if len(users) == 0 {
		return User{}, fmt.Errorf("not found")
	}

	return users[0], nil

}

func (cdu *User) JwtTokenGet(expiryTime time.Time, secretKey []byte) (string, error) {
	claim := *cdu

	claim.PasswordResetToken = ""
	claim.EmailVerificationToken = ""
	claim.Password = ""

	token, err := jwt.Sign(jwt.HS256, secretKey, claim, jwt.MaxAge(time.Hour*24*7))

	if err != nil {
		return "", err
	}

	return string(token), nil
}

func UserVerifyEmailToken(email string, token string) (User, error) {
	user, err := UserGetByEmail(email)

	if err != nil {
		return user, err
	}

	if user.EmailVerificationToken != token {
		return user, fmt.Errorf("invalid token")
	}

	user.EmailVerified = true
	user.EmailVerificationToken = ""

	ctx := mgm.Ctx()
	_, err = mgm.Coll(&User{}).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

	if err != nil {
		return user, err
	}

	return user, nil
}

func (user *User) MarkAccountForDeletion() error {
	user.MarkedForDeletion = true

	ctx := mgm.Ctx()
	_, err := mgm.Coll(user).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{
			"markedForDeletion":  true,
			"deleteAccountAfter": time.Now().Add(time.Hour * 24 * 30),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (user *User) UploadFile(data []FileUpload) ([]File, error) {

	// step 1: Upload to S3 and create a store data in files collection
	// step 2: Insert everything to the mongodb mapping it to the user
	var buckerName = utils.GetEnv("S3_BUCKET_NAME", false)
	var endPointUrl = utils.GetEnv("AWS_ENDPOINT_URL", false)

	var files []File

	for _, file := range data {
		uploadOutput, err := utils.S3UploadFile(buckerName, file.Title, &file.File)

		if err != nil {
			return files, err
		}

		var location = uploadOutput.Location

		location = strings.Replace(location, "https://", "", 1)
		location = strings.Replace(location, endPointUrl, "", 1)

		var fileObj = File{
			UserId:      user.ID.Hex(),
			Name:        file.Title,
			Type:        file.File.Header.Get("Content-Type"),
			RelativeURL: location,
		}

		err = mgm.Coll(&fileObj).Create(&fileObj)

		if err != nil {
			return files, err
		}

		files = append(files, fileObj)
	}

	return files, nil
}

func UserLogin(loginDto dto.UserLoginDto) (User, error) {
	user, err := UserGetByEmail(loginDto.Email)

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

	isCorrectPassword := user.ComparePassword(loginDto.Password)

	if !isCorrectPassword {

		user.PasswordTries++

		if user.PasswordTries >= 5 {
			user.AccountBlocked = true
		}

		ctx := mgm.Ctx()
		_, err = mgm.Coll(&user).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

		if err != nil {
			fmt.Fprintf(os.Stdout, "UserLogin Error: Failed to update user\n%v", err)
			return user, err
		}

		fmt.Fprintln(os.Stdout, "UserLogin Error: Incorrect Password!")
		return user, fmt.Errorf("incorrect password")
	}

	user.PasswordTries = 0

	ctx := mgm.Ctx()
	_, err = mgm.Coll(&user).UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})

	if err != nil {
		fmt.Fprintf(os.Stdout, "UserLogin Error: Failed to update user\n%v", err)
		return user, err
	}

	return user, nil
}

func (user *User) GetActiveAppSubscriptionByAppCode(appCode string) (payments.UserSubscription, error) {
	var subscriptionColl = mgm.Coll(&payments.UserSubscription{})

	res := subscriptionColl.FindOne(mgm.Ctx(), bson.M{
		"appCode": appCode,
		"userId":  user.ID.Hex(),
		"status": bson.M{
			"$in": []string{
				string(payments.SubscriptionStatusActive),
				string(payments.SubscriptionStatusCharged),
				string(payments.SubscriptionStatusCompleted),
			},
		},
	}, &options.FindOneOptions{})

	if res.Err() != nil {
		return payments.UserSubscription{}, res.Err()
	}

	var subscription payments.UserSubscription

	err := res.Decode(&subscription)

	if err != nil {
		return payments.UserSubscription{}, err
	}

	return subscription, nil
}
