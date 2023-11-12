package userservice

// accept login requests and provide token,
// accept token and provide user info

import (
	"blackbox-v2/pkg/db"

	"gorm.io/gorm"
)

type User struct {
	ID       uint
	UserCID  string `gorm:"uniqueIndex:idx_cid"`
	Email    string `gorm:"unique"`
	Username string `gorm:"uniqueIndex:idx_username"`
	Password string
}

func (u *User) Create() error {
	err := db.DB.Create(u).Error
	return err
}

func (u *User) Update() error {
	err := db.DB.Save(u).Error
	return err
}
func (u *User) toJson() map[string]interface{} {
	return map[string]interface{}{
		"email":    u.Email,
		"username": u.Username,
		"user_id":  u.ID,
	}
}

func getUserByUsername(username string) (*User, error) {
	var user User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func GetUserByCID(cid string) (*User, error) {
	var user User
	err := db.DB.Where("user_c_id = ?", cid).First(&user).Error
	return &user, err
}
func GetUserByID(id uint) (User, error) {
	var user User
	err := db.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

type UserMetaData struct {
	*gorm.Model
	ID             uint
	User           User
	UserID         uint  `gorm:"index"`
	DirCreated     bool  `gorm:"default:false"`
	NoteDirCreated bool  `gorm:"default:false"`
	StorageUsed    int64 `gorm:"default:0"`
	DirName        string
}

func (u *UserMetaData) Create() error {
	err := db.DB.Create(u).Error
	return err
}

func (u *UserMetaData) Update() error {
	err := db.DB.Save(u).Error
	return err
}

func GetUserMetaDataByUserID(userID uint) (*UserMetaData, error) {
	var user UserMetaData
	err := db.DB.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

type UserSession struct {
	*gorm.Model
	ID        uint   `gorm:"PrimaryIndex"`
	User      User   `gorm:"foreignKey:UserID"`
	SessionID string `gorm:"index"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"index"`
	ExpiresAt int64  `gorm:"index"`
	LoggedOut bool   `gorm:"index"`
}

func (u *UserSession) Create() error {
	err := db.DB.Create(u).Error
	return err
}

func (u *UserSession) Update() error {
	err := db.DB.Save(u).Error
	return err
}
func GetUserSessionBySessionID(sessionID string) (*UserSession, error) {
	var user UserSession
	err := db.DB.Preload("User").Where("session_id = ?", sessionID).First(&user).Error
	return &user, err
}

// type UserSignupVerification struct {
// 	*gorm.Model
// 	ID        uint   `gorm:"PrimaryIndex"`
// 	User      User   `gorm:"foreignKey:UserID"`
// 	UserID    uint   `gorm:"index"`
// 	Email     string `gorm:"index"`
// 	OTP       string `gorm:"index"`
// 	ExpiresAt int64  `gorm:"index"`
// 	Verified  bool   `gorm:"index"`
// }

// func (u *UserSignupVerification) create() error {
// 	err := db.DB.Create(u).Error
// 	return err
// }
// func (u *UserSignupVerification) update() error {
// 	err := db.DB.Save(u).Error
// 	return err
// }
// func getUserVerificationByEmail(email string) (*UserSignupVerification, error) {
// 	var user UserSignupVerification
// 	err := db.DB.Where("email = ?", email).First(&user).Error
// 	return &user, err
// }
