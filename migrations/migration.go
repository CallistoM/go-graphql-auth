package migrations

import (
	// standard libraries
	"log"
	// custom structs
	"github.com/callistom/go-graphql-auth/structs"
	// db
	"github.com/jinzhu/gorm"
	// encryption library
	"golang.org/x/crypto/bcrypt"
)

var (
	err error
)

// CreateMigrations create all structs
func CreateMigrations(db *gorm.DB) (bool, error) {

	if err != nil {
		log.Fatal(err)
	}

	// migrate user struct
	if err := db.AutoMigrate(&structs.User{}).Error; err != nil {
		return false, err
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	// test data
	user := structs.User{Name: "TestUser1", Mail: "test@user.com", Password: string(hash)}

	// create first user
	if err := db.FirstOrCreate(&user).Error; err != nil {
		return false, err
	}

	return true, nil
}
