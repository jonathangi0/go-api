package users

import (
	"errors"
	"time"

	"github.com/gothinkster/golang-gin-realworld-example-app/common"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Models should only be concerned with database schema, more strict checking should be put in validator.
//
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
//
// HINT: If you want to split null and "", you should use *string instead of string.
type UserModel struct {
	ID              uint       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID          *uint64    `gorm:"column:user_id"`
	Name            string     `gorm:"column:name;size:255;not null"`
	PayrollNumber   *string    `gorm:"column:payroll_number;size:255"`
	AreaName        *string    `gorm:"column:area_name;size:255"`
	Branch          *string    `gorm:"column:branch;size:255"`
	BranchAccess    *string    `gorm:"column:branch_access"`
	Email           string     `gorm:"column:email;size:255;uniqueIndex;not null"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at"`
	PasswordHash    string     `gorm:"column:password;not null"`
	RememberToken   *string    `gorm:"column:remember_token;size:100"`
	SessionID       *string    `gorm:"column:session_id"`
	Active          bool       `gorm:"column:active;not null"`
	CreatedAt       *time.Time `gorm:"column:created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at"`

	// API compatibility fields; they are not columns in public.users.
	Username string  `gorm:"-"`
	Bio      string  `gorm:"-"`
	Image    *string `gorm:"-"`
}

func (UserModel) TableName() string {
	return "users"
}

// A hack way to save ManyToMany relationship,
// gorm will build the alias as FollowingBy <-> FollowingByID <-> "following_by_id".
//
// DB schema looks like: id, created_at, updated_at, deleted_at, following_id, followed_by_id.
//
// Retrieve them by:
//
//	db.Where(FollowModel{ FollowingID:  v.ID, FollowedByID: u.ID, }).First(&follow)
//	db.Where(FollowModel{ FollowedByID: u.ID, }).Find(&follows)
//
// More details about gorm.Model: http://jinzhu.me/gorm/models.html#conventions
type FollowModel struct {
	gorm.Model
	Following    UserModel
	FollowingID  uint
	FollowedBy   UserModel
	FollowedByID uint
}

// Migrate the schema of database if needed
func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&UserModel{})
	db.AutoMigrate(&FollowModel{})
}

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
//
//	err := userModel.setPassword("password0")
func (u *UserModel) setPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
//
//	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *UserModel) checkPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// You could input the conditions and it will return an UserModel in database with error info.
//
//	userModel, err := FindOneUser(&UserModel{Username: "username0"})
func FindOneUser(condition interface{}) (UserModel, error) {
	db := common.GetDB()
	var model UserModel
	if user, ok := condition.(UserModel); ok {
		if user.Name == "" && user.Username != "" {
			user.Name = user.Username
		}
		condition = user
	}
	err := db.Where(condition).First(&model).Error
	return model, err
}

// You could input an UserModel which will be saved in database returning with error info
//
//	if err := SaveOne(&userModel); err != nil { ... }
func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}

// You could update properties of an UserModel to database returning with error info.
//
//	err := db.Model(userModel).Updates(UserModel{Username: "wangzitian0"}).Error
func (model *UserModel) Update(data interface{}) error {
	db := common.GetDB()
	err := db.Model(model).Updates(data).Error
	return err
}

// You could add a following relationship as userModel1 following userModel2
//
//	err = userModel1.following(userModel2)
func (u UserModel) following(v UserModel) error {
	db := common.GetDB()
	var follow FollowModel
	err := db.FirstOrCreate(&follow, &FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).Error
	return err
}

// You could check whether  userModel1 following userModel2
//
//	followingBool = myUserModel.isFollowing(self.UserModel)
func (u UserModel) isFollowing(v UserModel) bool {
	db := common.GetDB()
	var follow FollowModel
	db.Where(FollowModel{
		FollowingID:  v.ID,
		FollowedByID: u.ID,
	}).First(&follow)
	return follow.ID != 0
}

// You could delete a following relationship as userModel1 following userModel2
//
//	err = userModel1.unFollowing(userModel2)
func (u UserModel) unFollowing(v UserModel) error {
	db := common.GetDB()
	err := db.Where("following_id = ? AND followed_by_id = ?", v.ID, u.ID).Delete(&FollowModel{}).Error
	return err
}

// You could get a following list of userModel
//
//	followings := userModel.GetFollowings()
func (u UserModel) GetFollowings() []UserModel {
	db := common.GetDB()
	var follows []FollowModel
	var followings []UserModel
	db.Preload("Following").Where(FollowModel{
		FollowedByID: u.ID,
	}).Find(&follows)
	for _, follow := range follows {
		followings = append(followings, follow.Following)
	}
	return followings
}
