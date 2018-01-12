package model

import (
	"time"
)

type Assets struct {
	Id int64 `json:"id"`

	Key        string    `json:"name"`
	Value      string    `json:"value"`
	Desc       string    `json:"description"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	IsDeleted  bool      `json:"is_deleted"`
	Type       string    `json:"type"`
	Status     int       `json:"status"`
}
