package services

import (
	"context"
	"fmt"
	"uu/config"
	"uu/models"
	"uu/utils"
)

type BuffResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		PageNum    int                `json:"page_num"`
		PageSize   int                `json:"page_size"`
		TotalCount int                `json:"total_count"`
		TotalPage  int                `json:"total_page"`
		Items      []*models.BuffItem `json:"items"`
	} `json:"data"`
}

type BuffInventoryResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		TotalCount int                     `json:"total_count"`
		TotalPage  int                     `json:"total_page"`
		Items      []*models.BuffInventory `json:"items"`
	} `json:"data"`
}

var BuffClient = utils.CreateClient("https://buff.163.com")

func GetBuffHeaders() map[string]string {
	ctx := context.Background()
	var buff models.BuffToken
	err := buff.GetBuffToken(ctx)
	if err != nil {
		config.Log.Errorf("get buff token error: %s", err)
	}
	cookie := fmt.Sprintf("session=%s; game=csgo", buff.Session)

	return map[string]string{
		"Accept":              utils.Header.Buff.Accept,
		"Accept-Language":     utils.Header.Buff.AcceptLanguage,
		"Accept-Encoding":     utils.Header.Buff.AcceptEncoding,
		"User-Agent":          utils.Header.Buff.UserAgent,
		"Timezone-Offset":     utils.Header.Buff.TimezoneOffset,
		"DeviceName":          utils.Header.Buff.DeviceName,
		"Device-Id-Weak":      utils.Header.Buff.DeviceIdWeak,
		"Screen-Scale":        utils.Header.Buff.ScreenScale,
		"Resolution":          utils.Header.Buff.Resolution,
		"Locale":              utils.Header.Buff.Locale,
		"Device-Id":           utils.Header.Buff.DeviceId,
		"Connection":          utils.Header.Buff.Connection,
		"Locale-Supported":    utils.Header.Buff.LocaleSupported,
		"Timezone":            utils.Header.Buff.Timezone,
		"Network":             utils.Header.Buff.Network,
		"Product":             utils.Header.Buff.Product,
		"Timezone-Offset-DST": utils.Header.Buff.TimezoneOffsetDst,
		"Model":               utils.Header.Buff.Model,
		"App-Version":         utils.Header.Buff.AppVersion,
		"Screen-Size":         utils.Header.Buff.ScreenSize,
		"App-Version-Code":    utils.Header.Buff.AppVersionCode,
		"System-Version":      utils.Header.Buff.SystemVersion,
		"Cookie":              cookie,
	}
}

func GetBuffItems(pageSize, pageNum string) ([]*models.BuffItem, int, error) {
	header := GetBuffHeaders()
	var buffResponse BuffResponse
	var opt = utils.RequestOptions{
		QueryParams: map[string]string{
			"page_size": pageSize,
			"page_num":  pageNum,
			"game":      "csgo",
			"tab":       "selling",
		},
		Headers: header,
		Result:  &buffResponse,
	}
	res, err := BuffClient.DoRequest("GET", "api/market/goods", opt)
	if err != nil || res.StatusCode() != 200 {
		config.Log.Errorf("request buff api error : %s", err)
	}
	return buffResponse.Data.Items, buffResponse.Data.TotalCount, err
}

// VerifyBuffToken 校验buff token
func VerifyBuffToken() {
	header := GetBuffHeaders()
	var buffResponse BuffResponse
	var buffToken models.BuffToken
	var opt = utils.RequestOptions{
		QueryParams: map[string]string{
			"page_size": "80",
			"page_num":  "2",
			"game":      "csgo",
		},
		Headers: header,
		Result:  &buffResponse,
	}
	res, err := BuffClient.DoRequest("GET", "api/market/goods", opt)
	if err != nil || res.StatusCode() != 200 || buffResponse.Code != "OK" {
		config.Log.Warn("buff token expired")
		buffToken.Expired = "yes"
	} else {
		buffToken.Expired = "no"
	}
	err = buffToken.UpdateBuffExpired()
	if err != nil {
		config.Log.Errorf("update buff expired error: %v", err)
	}
}

func GetBuffInventory() []*models.BuffInventory {
	header := GetBuffHeaders()
	var buffResponse BuffInventoryResponse
	var opt = utils.RequestOptions{
		QueryParams: map[string]string{
			"page_size": "999",
			"page_num":  "1",
			"game":      "csgo",
			"fold":      "true",
		},
		Headers: header,
		Result:  &buffResponse,
	}
	res, err := BuffClient.DoRequest("GET", "api/market/steam_inventory", opt)
	if err != nil || res.StatusCode() != 200 {
		config.Log.Errorf("request buff inventory api error : %s", err)
	}
	return buffResponse.Data.Items
}
