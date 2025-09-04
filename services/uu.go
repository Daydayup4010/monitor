package services

import (
	"context"
	"uu/config"
	"uu/models"
	"uu/utils"
)

type UUResponse struct {
	Code       int             `json:"Code"`
	Msg        string          `json:"Msg"`
	Data       []*models.UItem `json:"Data"`
	TotalCount int             `json:"TotalCount"`
}

var client = utils.CreateClient("https://api.youpin898.com")

func GetHeaders() map[string]string {
	ctx := context.Background()
	var yp models.UUToken
	err := yp.GetUUToken(ctx)
	if err != nil {
		config.Log.Errorf("get youpin token error: %s", err)
	}
	return map[string]string{
		"accept":          utils.Header.UU.Accept,
		"accept-language": utils.Header.UU.AcceptLanguage,
		"app-version":     utils.Header.UU.AppVersion,
		"apptype":         utils.Header.UU.AppType,
		"content-type":    utils.Header.UU.ContentType,
		"platform":        utils.Header.UU.Platform,
		"user-agent":      utils.Header.UU.UserAgent,
		"uk":              yp.Uk,
		"authorization":   yp.Authorization,
	}
}

func GetUUItems(pageSize, PageNum int) ([]*models.UItem, int, error) {
	var header = GetHeaders()
	var uuResp UUResponse
	var opts = utils.RequestOptions{
		Body: map[string]int{
			"listSortType": 0,
			"sortType":     0,
			"pageSize":     pageSize,
			"pageIndex":    PageNum},
		Headers: header,
		Result:  &uuResp,
	}
	res, err := client.DoRequest("POST", "api/homepage/pc/goods/market/querySaleTemplate", opts)
	if err != nil || res.StatusCode() != 200 {
		config.Log.Warnf("request uu list api error: %s, code: %d", err, res.StatusCode())
	}
	return uuResp.Data, uuResp.TotalCount, err

}
