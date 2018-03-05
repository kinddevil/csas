package model

import (
	"time"
)

type Calendar struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	SchoolId   int64     `json:"school_id"`
	Desc       string    `json:"description"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	Deleted    bool      `json:"deleted"`
	Type       string    `json:"type"`
	IsActive   bool      `json:"is_active"`
	End        time.Time `json:"end"`
	Start      time.Time `json:"start"`
	Events     string    `json:"events"`
	IsVisible  bool      `json:"is_visible"`
}
