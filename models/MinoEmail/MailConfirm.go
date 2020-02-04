package MinoEmail

import (
	"errors"
	"git.ntmc.tech/root/MinoIC-PE/models"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoDatabase"
	"github.com/jinzhu/gorm"
	"time"
)

func ConfirmKey(key string) bool {
	DB := MinoDatabase.GetDatabase()
	defer DB.Close()
	var keyInfo MinoDatabase.RegConfirmKey
	if !DB.Where("Key = ?", key).First(&keyInfo).RecordNotFound() {
		if keyInfo.ValidTime.After(time.Now()) {
			var user MinoDatabase.User
			if !DB.Where("ID = ?", keyInfo.ID).First(&user).RecordNotFound() {
				user.EmailConfirmed = true
				DB.Model(&user).Update(MinoDatabase.User{
					EmailConfirmed: true,
				})
				return true
			}
		}
	}
	return false
}

func ConfirmRegister(user MinoDatabase.User) error {
	if user.EmailConfirmed {
		return errors.New("User Already confirmed! ")
	}
	key := MinoDatabase.RegConfirmKey{
		UserName:  user.Name,
		UserEmail: user.Email,
		Model:     gorm.Model{},
		Key:       models.RandKey(15),
		UserID:    user.ID,
		ValidTime: time.Now().Add(30 * time.Minute),
	}
	DB := MinoDatabase.GetDatabase()
	defer DB.Close()
	DB.Create(&key)
	SendConfirmMail(key)
	return nil
}