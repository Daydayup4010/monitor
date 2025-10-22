package models

import (
	"context"
	"uu/config"
	"uu/utils"
)

type UUToken struct {
	Authorization string `json:"authorization" binding:"required"`
	Uk            string `json:"uk" binding:"required"`
	Expired       string `json:"expired"`
}

type BuffToken struct {
	CsrfToken string `json:"csrf_token" binding:"required"`
	Session   string `json:"session" binding:"required"`
	Expired   string `json:"expired"`
}

func (uu *UUToken) SetUUToken(ctx context.Context) int {
	err := config.RDB.HSet(ctx, "uu_token", "authorization", uu.Authorization, "uk", uu.Uk).Err()
	if err != nil {
		config.Log.Errorf("update youpin token error: %s", err)
		return utils.ErrCodeUpdateUUToken
	}
	return utils.SUCCESS
}

func (uu *UUToken) GetUUToken(ctx context.Context) error {
	auth, err := config.RDB.HGet(ctx, "uu_token", "authorization").Result()
	uk, err := config.RDB.HGet(ctx, "uu_token", "uk").Result()
	if err == nil {
		uu.Authorization = auth
		uu.Uk = uk
	}

	return err
}

func (uu *UUToken) UpdateUUExpired() error {
	err := config.RDB.HSet(context.Background(), "uu_token", "expired", uu.Expired).Err()
	return err
}

func (uu *UUToken) GetUUExpired() int {
	ex, err := config.RDB.HGet(context.Background(), "uu_token", "expired").Result()
	if err == nil {
		uu.Expired = ex
	}
	if err != nil {
		config.Log.Errorf("Get token expired error : %v", err)
		return utils.ErrCodeGetTokenExpired
	}
	return utils.SUCCESS
}

func (buff *BuffToken) SetBuffToken(ctx context.Context) int {
	err := config.RDB.HSet(ctx, "buff_token", "session", buff.Session, "csrf_token", buff.CsrfToken).Err()
	if err != nil {
		config.Log.Errorf("update buff token error: %s", err)
		return utils.ErrCodeUpdateBuffToken
	}
	return utils.SUCCESS
}

func (buff *BuffToken) GetBuffToken(ctx context.Context) error {
	session, err := config.RDB.HGet(ctx, "buff_token", "session").Result()
	csrf, err := config.RDB.HGet(ctx, "buff_token", "csrf_token").Result()
	if err == nil {
		buff.Session = session
		buff.CsrfToken = csrf
	}
	return err
}

func (buff *BuffToken) UpdateBuffExpired() error {
	err := config.RDB.HSet(context.Background(), "buff_token", "expired", buff.Expired).Err()
	return err
}

func (buff *BuffToken) GetBuffExpired() int {
	ex, err := config.RDB.HGet(context.Background(), "buff_token", "expired").Result()
	if err == nil {
		buff.Expired = ex
	}
	if err != nil {
		config.Log.Errorf("Get token expired error : %v", err)
		return utils.ErrCodeGetTokenExpired
	}
	return utils.SUCCESS
}
