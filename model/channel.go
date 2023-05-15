package model

import (
	_ "gorm.io/driver/sqlite"
	"one-api/common"
)

type Channel struct {
	Id           int    `json:"id"`
	Type         int    `json:"type" gorm:"default:0"`
	Key          string `json:"key" gorm:"not null"`
	Status       int    `json:"status" gorm:"default:1"`
	Name         string `json:"name" gorm:"index"`
	Weight       int    `json:"weight"`
	CreatedTime  int64  `json:"created_time" gorm:"bigint"`
	TestTime     int64  `json:"test_time" gorm:"bigint"`
	ResponseTime int    `json:"response_time"` // in milliseconds
	BaseURL      string `json:"base_url" gorm:"column:base_url"`
	Other        string `json:"other"`
}

func GetAllChannels(startIdx int, num int) ([]*Channel, error) {
	var channels []*Channel
	var err error
	err = DB.Order("id desc").Limit(num).Offset(startIdx).Omit("key").Find(&channels).Error
	return channels, err
}

func SearchChannels(keyword string) (channels []*Channel, err error) {
	err = DB.Omit("key").Where("id = ? or name LIKE ?", keyword, keyword+"%").Find(&channels).Error
	return channels, err
}

func GetChannelById(id int, selectAll bool) (*Channel, error) {
	channel := Channel{Id: id}
	var err error = nil
	if selectAll {
		err = DB.First(&channel, "id = ?", id).Error
	} else {
		err = DB.Omit("key").First(&channel, "id = ?", id).Error
	}
	return &channel, err
}

func GetRandomChannel() (*Channel, error) {
	// TODO: consider weight
	channel := Channel{}
	var err error = nil
	if common.UsingSQLite {
		err = DB.Where("status = ?", common.ChannelStatusEnabled).Order("RANDOM()").Limit(1).First(&channel).Error
	} else {
		err = DB.Where("status = ?", common.ChannelStatusEnabled).Order("RAND()").Limit(1).First(&channel).Error
	}
	return &channel, err
}

func BatchInsertChannels(channels []Channel) error {
	var err error
	err = DB.Create(&channels).Error
	return err
}

func (channel *Channel) Insert() error {
	var err error
	err = DB.Create(channel).Error
	return err
}

func (channel *Channel) Update() error {
	var err error
	err = DB.Model(channel).Updates(channel).Error
	return err
}

func (channel *Channel) UpdateResponseTime(responseTime int64) {
	err := DB.Model(channel).Select("response_time", "test_time").Updates(Channel{
		TestTime:     common.GetTimestamp(),
		ResponseTime: int(responseTime),
	}).Error
	if err != nil {
		common.SysError("failed to update response time: " + err.Error())
	}
}

func (channel *Channel) Delete() error {
	var err error
	err = DB.Delete(channel).Error
	return err
}
