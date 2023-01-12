package entity

import (
	"blog-admin-api/pkg/db"
)

type FriendLink struct {
	Base
	Name  string `json:"name"`
	Email string `json:"email"`
	Url   string `json:"url"`
}

func (*FriendLink) TableName() string {
	return "friend_link"
}

func (h *FriendLink) Create() error {
	return db.Orm.Create(h).Error
}

func (*FriendLink) Delete(id int) error {
	h := FriendLink{}
	h.Id = id

	return db.Orm.Delete(h).Error
}

func (h *FriendLink) Update(id int) error {
	h.Id = id

	return db.Orm.Model(h).Updates(h).Error
}

func (*FriendLink) FindWithPage(page, limit int) ([]*FriendLink, int64, error) {
	if limit == 0 {
		limit = 6
	}

	friend := make([]*FriendLink, 0)
	var count int64

	if err := db.Orm.Model(&FriendLink{}).Count(&count).Error; err != nil {
		return friend, count, err
	}
	err := db.Orm.Offset((page - 1) * limit).Limit(limit).Order("id desc").Find(&friend).Error

	return friend, count, err
}

func (*FriendLink) GetById(id int) (*FriendLink, error) {
	harems := &FriendLink{}
	if err := db.Orm.First(&harems, id).Error; err != nil {
		return harems, err
	}

	return harems, nil
}
