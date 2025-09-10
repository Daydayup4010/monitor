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

var BuffClient = utils.CreateClient("https://buff.163.com")

func GetBuffHeaders() map[string]string {
	ctx := context.Background()
	var buff models.BuffToken
	err := buff.GetBuffToken(ctx)
	if err != nil {
		config.Log.Errorf("get buff token error: %s", err)
	}
	cookie := fmt.Sprintf("session=%s;csrf_token=%s", buff.Session, buff.CsrfToken)
	return map[string]string{
		"accept":          utils.Header.Buff.Accept,
		"accept-language": utils.Header.Buff.AcceptLanguage,
		//"content-type":    utils.Header.Buff.ContentType,
		"user-agent": utils.Header.Buff.UserAgent,
		"Cookie":     cookie,
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
