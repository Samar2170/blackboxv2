package userservice

import (
	"blackbox-v2/internal/models"
	"blackbox-v2/pkg/config"
	"blackbox-v2/pkg/db"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var signingKey []byte

func init() {
	signingKey = []byte(config.Config.SigningKey)
	db.DB.AutoMigrate(&User{}, &UserSession{}, UserMetaData{})
}

type UserClaim struct {
	Username string `json:"username"`
	UserCid  string `json:"user_cid"`
	jwt.RegisteredClaims
}

func getCIDForUser() string {
	return uuid.New().String()
}

func customHash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), 14)
	return string(bytes), err
}
func checkPasswordHashed(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createToken(username string, userCID string) (string, error) {
	claims := UserClaim{
		username,
		userCID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cognitio",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func LoginUser(username string, password string) (string, error) {
	user, err := getUserByUsername(username)
	if err != nil || user.Username == "" {
		return "", err
	}

	if !checkPasswordHashed(user.Password, password) {
		return "", errors.New("invalid password")
	}
	token, err := createToken(username, user.UserCID)
	if err != nil {
		return "", err
	}
	us := UserSession{
		User:      *user,
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * 10 * time.Hour).Unix(),
	}
	err = models.CreateModelInstance(&us)
	if err != nil {
		return "", err
	}
	return token, nil
}

func SignupUser(email, username, password string) (string, error) {
	u, err := getUserByUsername(username)
	if err == nil || u.Username == username {
		return "", errors.New("user already exists with this username")
	}
	hashedPassword, err := customHash(password)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	user := User{
		Email:    email,
		Username: username,
		Password: hashedPassword,
		UserCID:  getCIDForUser(),
	}
	err = models.CreateModelInstance(&user)
	if err != nil {
		return "", err
	}
	usermd := UserMetaData{
		UserID: user.ID,
	}
	err = models.CreateModelInstance(&usermd)
	if err != nil {
		return user.UserCID, errors.New("error creating user metadata")
	}
	err = CreateDirForUser(&user)
	if err != nil {
		return user.UserCID, errors.New("error creating user directory")
	}
	return user.UserCID, nil
}

func VerifyToken(token string) (User, error) {
	empty := User{}
	claims := UserClaim{}
	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return empty, err
	}
	if !tkn.Valid {
		return empty, errors.New("invalid token")
	}
	user, err := GetUserByCID(claims.UserCid)
	if err != nil {
		return empty, err
	}
	return *user, nil
}

// func sendEmailVerification(user *User) error {
// 	otp := rand.Int()
// 	uv := UserSignupVerification{
// 		User:      *user,
// 		UserID:    user.ID,
// 		Email:     user.Email,
// 		OTP:       fmt.Sprintf("%d", otp),
// 		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
// 		Verified:  false,
// 	}
// 	err := createModelInstance(&uv)
// 	if err != nil {
// 		return err
// 	}
// 	template := fmt.Sprintf(`Hi`+user.Username+`,
// 	Welcome to Cognitio. Please verify your email by entering the following OTP:
// 	`+fmt.Sprintf("%d", otp)+`

// 	Thanks,
// 	`, user.Username, otp)
// 	message := []byte(template)
// 	auth := smtp.PlainAuth("", emailAccount, emailPassword, smtpHost)
// 	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, emailAccount, []string{user.Email}, message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func VerifyEmailOTP(userCID, otp string) error {
// 	user, err := getUserByCID(userCID)
// 	if err != nil {
// 		return err
// 	}
// 	uv, err := getUserVerificationByEmail(user.Email)
// 	if err != nil {
// 		return err
// 	}
// 	if uv.OTP != otp {
// 		return errors.New("invalid otp")
// 	}
// 	uv.Verified = true
// 	err = updateModelInstance(uv)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
