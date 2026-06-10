package entity

import "time"

type User struct {
	Id          int64      `json:"id"`
	Uid         string     `json:"uid"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	Phone       string     `json:"phone"`
	Password    string     `json:"-"`
	Gender      int        `json:"gender"`
	Status      int        `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type UserWallet struct {
	Id          int64 `json:"id"`
	UserId      int64 `json:"user_id"`
	Chips       int64 `json:"chips"`
	Gold        int64 `json:"gold"`
	Diamonds    int   `json:"diamonds"`
	FrozenChips int64 `json:"frozen_chips"`
	Version     int   `json:"version"`
}
