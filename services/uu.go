package services

import (
	"context"
	"uu/config"
	"uu/models"
	"uu/utils"
)
import "encoding/json"

type UUResponse struct {
	Code       int             `json:"Code"`
	Msg        string          `json:"Msg"`
	Data       []*models.UItem `json:"Data"`
	TotalCount int             `json:"TotalCount"`
}

type UUInventory struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ItemCount  int                  `json:"itemCount"`
		ItemsInfos []*models.UItemsInfo `json:"itemsInfos"`
		TotalCount int                  `json:"totalCount"`
		Valuation  string               `json:"valuation"`
	} `json:"data"`
}

type OpenResponse struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data []*Item `json:"data"`
}

type Item struct {
	SaleTemplateResponse  *SaleTemplateResponse `json:"saleTemplateResponse"`
	SaleCommodityResponse *SaleCommodity        `json:"saleCommodityResponse"`
}

type SaleTemplateResponse struct {
	TemplateId       int    `json:"templateId"`
	TemplateHashName string `json:"templateHashName"`
	IconUrl          string `json:"iconUrl"`
	ExteriorName     string `json:"exteriorName"`
	RarityName       string `json:"rarityName"`
	QualityName      string `json:"qualityName"`
}

type SaleCommodity struct {
	MinSellPrice             string `json:"minSellPrice"`
	FastShippingMinSellPrice string `json:"fastshippingminSellPrice"`
	ReferencePrice           string `json:"referencePrice"`
	SellNum                  int64  `json:"sellNum"`
}

var client = utils.CreateClient("https://gw-openapi.youpin898.com")

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

// BuildBodyWithSign 构建带签名的请求体
// body: 原始请求体参数
// 返回: 添加了 timestamp, sign, appKey 的请求体
func BuildBodyWithSign(body map[string]interface{}) (map[string]interface{}, error) {
	// 获取签名参数
	signParams, err := GetSignParams(body)
	if err != nil {
		return nil, err
	}

	// 将签名参数添加到请求体中
	body["timestamp"] = signParams.Timestamp
	body["sign"] = signParams.Sign
	body["appKey"] = signParams.AppKey

	return body, nil
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

func VerifyUUToken() {
	var header = GetHeaders()
	var uuResp UUResponse
	var uuToken models.UUToken
	var opts = utils.RequestOptions{
		Body: map[string]int{
			"listSortType": 0,
			"sortType":     0,
			"pageSize":     100,
			"pageIndex":    3},
		Headers: header,
		Result:  &uuResp,
	}
	res, err := client.DoRequest("POST", "api/homepage/pc/goods/market/querySaleTemplate", opts)
	if err != nil || res.StatusCode() != 200 || uuResp.Code != 0 {
		config.Log.Warn("uu token expired")
		uuToken.Expired = "yes"
	} else {
		uuToken.Expired = "no"
	}
	err = uuToken.UpdateUUExpired()
	if err != nil {
		config.Log.Errorf("uu token expired update error: %v", err)
	}
}

func GetUUInventory() []*models.UItemsInfo {
	var header = GetHeaders()
	var uuInventory UUInventory
	var opts = utils.RequestOptions{
		Body: map[string]int{
			"isMerge":   1,
			"pageSize":  999,
			"pageIndex": 1},
		Headers: header,
		Result:  &uuInventory,
	}
	res, err := client.DoRequest("POST", "api/youpin/pc/inventory/list", opts)
	if err != nil || res.StatusCode() != 200 {
		config.Log.Warnf("request uu inventory list api error: %s, code: %d", err, res.StatusCode())
	}
	return uuInventory.Data.ItemsInfos
}

// RequestItem 请求列表项
type RequestItem struct {
	TemplateHashName string `json:"templateHashName"`
}

func GetUUGoods(hashNames []string) []*models.UBaseInfo {
	var uuResp OpenResponse
	var requestList []RequestItem
	var infos []*models.UBaseInfo
	for _, name := range hashNames {
		requestList = append(requestList, RequestItem{
			TemplateHashName: name,
		})
	}

	body := map[string]interface{}{
		"requestList": requestList,
	}

	// 添加签名参数到请求体
	signedBody, err := BuildBodyWithSign(body)
	if err != nil {
		config.Log.Errorf("build signed body error: %s", err)
		return infos
	}

	var opts = utils.RequestOptions{
		Body: signedBody,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	res, err := client.DoRequest("POST", "open/v1/api/batchGetOnSaleCommodityInfo", opts)
	if err != nil || res.StatusCode() != 200 {
		config.Log.Errorf("request uu goods api error: %s", err)
		return infos
	}

	// 手动解析 JSON

	if err := json.Unmarshal(res.Body(), &uuResp); err != nil {
		config.Log.Errorf("json unmarshal error: %s", err)
		return infos
	}

	for _, item := range uuResp.Data {
		infos = append(infos, &models.UBaseInfo{
			Id:          item.SaleTemplateResponse.TemplateId,
			HashName:    item.SaleTemplateResponse.TemplateHashName,
			IconUrl:     item.SaleTemplateResponse.IconUrl,
			QualityName: item.SaleTemplateResponse.QualityName,
			RarityName:  item.SaleTemplateResponse.RarityName,
		})
	}
	return infos
}
