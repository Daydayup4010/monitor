package models

import (
	"context"
	"uu/config"
)

type UUToken struct {
	Authorization string `json:"authorization" binding:"required"`
	Uk            string `json:"uk" binding:"required"`
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
