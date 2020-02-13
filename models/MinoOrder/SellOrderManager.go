package MinoOrder

import (
	"errors"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoDatabase"
	"git.ntmc.tech/root/MinoIC-PE/models/MinoMessage"
	"git.ntmc.tech/root/MinoIC-PE/models/PterodactylAPI"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

func SellCreate(SpecID uint, userID uint) uint {
	DB := MinoDatabase.GetDatabase()
	var (
		wareSpec    MinoDatabase.WareSpec
		finalPrice  uint
		originPrice uint
	)
	DB.Where("id = ?", SpecID).First(&wareSpec)
	switch wareSpec.ValidDuration {
	case 3 * 24 * time.Hour:
		originPrice = 1
		finalPrice = 1
	case 30 * 24 * time.Hour:
		originPrice = uint(wareSpec.PricePerMonth)
		finalPrice = uint(0.01 * float32(100-wareSpec.Discount) * wareSpec.PricePerMonth)
	case 90 * 24 * time.Hour:
		originPrice = uint(wareSpec.PricePerMonth * 3)
		finalPrice = uint(0.01 * float32(100-wareSpec.Discount) * wareSpec.PricePerMonth * 3)
	}
	beego.Debug(originPrice, finalPrice)
	order := MinoDatabase.Order{
		Model:       gorm.Model{},
		SpecID:      SpecID,
		UserID:      userID,
		OriginPrice: originPrice,
		FinalPrice:  finalPrice,
		Paid:        false,
		Confirmed:   false,
	}
	DB.Create(&order)
	return order.ID
}

func SellGet(orderID uint) (MinoDatabase.Order, error) {
	var order MinoDatabase.Order
	DB := MinoDatabase.GetDatabase()
	if !DB.Where("id = ?", orderID).First(&order).RecordNotFound() {
		return order, nil
	}
	return MinoDatabase.Order{}, errors.New("cant find order by id " + strconv.Itoa(int(orderID)))
}

func SellPaymentCheck(orderID uint, keyString string, selectedIP int) error {
	DB := MinoDatabase.GetDatabase()
	var (
		order MinoDatabase.Order
		key   MinoDatabase.WareKey
		user  MinoDatabase.User
		spec  MinoDatabase.WareSpec
		exp   string
	)
	if DB.Where("id = ?", orderID).First(&order).RecordNotFound() {
		return errors.New("cant find order by id: " + strconv.Itoa(int(orderID)))
	}
	if order.Paid {
		return errors.New("order is already paid: " + strconv.Itoa(int(orderID)))
	}
	if DB.Where("key = ?", keyString).First(&key).RecordNotFound() {
		return errors.New("cant find key: " + keyString)
	}
	if key.Exp.Before(time.Now()) || key.SpecID != order.SpecID {
		return errors.New("invalid key: " + keyString)
	}
	if DB.Where("id = ?", order.UserID).First(&user).RecordNotFound() {
		return errors.New("cant find order`s owner by id: " + strconv.Itoa(int(order.UserID)))
	}
	if DB.Where("id = ?", order.SpecID).First(&spec).RecordNotFound() {
		return errors.New("cant find wareSpec by id: " + strconv.Itoa(int(order.SpecID)))
	}
	pteUser, ok := PterodactylAPI.GetUser(PterodactylAPI.ConfGetParams(), user.Name, true)
	if !ok {
		return errors.New("cant find pte user: " + user.Name)
	}
	pteUserID := pteUser.Uid
	switch spec.ValidDuration {
	case 3 * 24 * time.Hour:
		exp = (time.Now().AddDate(0, 0, 3)).Format("2006-01-02 15:04:05")
	case 30 * 24 * time.Hour:
		exp = (time.Now().AddDate(0, 30, 0)).Format("2006-01-02 15:04:05")
	case 90 * 24 * time.Hour:
		exp = (time.Now().AddDate(0, 90, 0)).Format("2006-01-02 15:04:05")
	}
	go func() {
		entity := MinoDatabase.WareEntity{
			Model:            gorm.Model{},
			UserID:           orderID,
			ServerExternalID: user.Name + strconv.Itoa(int(orderID)),
			UserExternalID:   user.Name,
			DeleteStatus:     0,
			ValidDate:        time.Now().Add(spec.ValidDuration),
		}
		DB.Create(&entity)
		err := PterodactylAPI.PterodactylCreateServer(PterodactylAPI.ConfGetParams(), PterodactylAPI.PterodactylServer{
			UserId:      pteUserID,
			ExternalId:  user.Name + strconv.Itoa(int(orderID)),
			Name:        user.Name + strconv.Itoa(int(orderID)),
			Description: "到期时间：" + exp,
			Suspended:   false,
			Limits: PterodactylAPI.PterodactylServerLimit{
				Memory: spec.Memory,
				Swap:   spec.Swap,
				Disk:   spec.Disk,
				IO:     spec.Io,
				CPU:    spec.Cpu,
			},
			Allocation: selectedIP,
			NestId:     spec.Nest,
			EggId:      spec.Egg,
			PackId:     0,
		})
		if err == nil {
			DB.Model(&order).Update("allocation_id", selectedIP)
			DB.Model(&order).Update("confirmed", true)
			DB.Model(&order).Update("paid", true)
			DB.Delete(&key)
			beego.Info("Key used: " + key.Key)
			MinoMessage.Send("ADMIN", user.ID, "您的订单 #"+strconv.Itoa(int(order.ID))+" 已成功创建对应服务器，请前往控制台确认")
			beego.Info("order id confirmed: " + strconv.Itoa(int(orderID)))
		} else {
			beego.Error("cant create server for order id: " + strconv.Itoa(int(orderID)) + "with error: " + err.Error())
		}
	}()
	return nil
}
