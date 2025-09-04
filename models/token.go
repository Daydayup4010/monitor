package models

import (
	"context"
	"uu/config"
)

type UUToken struct {
	Authorization string `json:"authorization" binding:"required"`
	Uk            string `json:"uk" binding:"required"`
}

type BuffToken struct {
	CsrfToken string `json:"csrf_token" binding:"required"`
	Session   string `json:"session" binding:"required"`
}

func (yp *UUToken) SetUUToken(ctx context.Context) error {
	err := config.RDB.HSet(ctx, "uu_token", "authorization", yp.Authorization, "uk", yp.Uk).Err()
	return err
}

func (yp *UUToken) GetUUToken(ctx context.Context) error {
	auth, err := config.RDB.HGet(ctx, "uu_token", "authorization").Result()
	uk, err := config.RDB.HGet(ctx, "uu_token", "uk").Result()
	if err == nil {
		yp.Authorization = auth
		yp.Uk = uk
	}

	return err
}

func (buff *BuffToken) SetBuffToken(ctx context.Context) error {
	err := config.RDB.HSet(ctx, "buff_token", "session", buff.Session, "csrf_token", buff.CsrfToken).Err()
	return err
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
