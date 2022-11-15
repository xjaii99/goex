package common

import (
	"errors"
	"fmt"
	"github.com/nntaoli-project/goex/v2"
	"github.com/nntaoli-project/goex/v2/internal/logger"
	"net/http"
	"net/url"
)

type Trade struct {
	*V5
	apiOpts goex.ApiOptions
}

func (t *Trade) CreateOrder(order goex.Order, opts ...goex.OptionParameter) (*goex.Order, error) {
	reqUrl := fmt.Sprintf("%s%s", t.uriOpts.Endpoint, t.uriOpts.NewOrderUri)
	params := url.Values{}

	params.Set("instId", order.Pair.Symbol)
	//params.Set("tdMode", "cash")
	params.Set("side", adaptOrderSideToSym(order.Side))
	//params.Set("posSide", "")
	params.Set("ordType", adaptOrderTypeToSym(order.OrderTy))
	params.Set("px", goex.FloatToString(order.Price, order.Pair.PricePrecision))
	params.Set("sz", goex.FloatToString(order.Qty, order.Pair.QtyPrecision))
	goex.MergeOptionParams(&params, opts...)

	data, err := t.DoAuthRequest(http.MethodPost, reqUrl, &params, nil)
	if err != nil {
		logger.Errorf("[CreateOrder] err=%s, response=%s", err.Error(), string(data))
		return nil, err
	}

	ord, err := t.unmarshalOpts.CreateOrderResponseUnmarshaler(data)
	if err != nil {
		return nil, err
	}

	ord.Pair = order.Pair
	ord.Price = order.Price
	ord.Qty = order.Qty
	ord.Side = order.Side
	ord.OrderTy = order.OrderTy
	ord.Status = goex.OrderStatus_Pending

	return ord, err
}

func (t *Trade) GetOrderInfo(pair goex.CurrencyPair, id string, opt ...goex.OptionParameter) (*goex.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Trade) GetPendingOrders(pair goex.CurrencyPair, opt ...goex.OptionParameter) ([]goex.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Trade) GetHistoryOrders(pair goex.CurrencyPair, opt ...goex.OptionParameter) ([]goex.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Trade) CancelOrder(pair goex.CurrencyPair, id string, opt ...goex.OptionParameter) error {
	//TODO implement me
	panic("implement me")
}

func (t *Trade) DoAuthRequest(httpMethod, reqUrl string, params *url.Values, headers map[string]string) ([]byte, error) {
	reqBody, _ := goex.ValuesToJson(*params)
	reqBodyStr := string(reqBody)
	_url, _ := url.Parse(reqUrl)
	signStr, timestamp := SignParams(httpMethod, _url.RequestURI(), t.apiOpts.Secret, reqBodyStr)
	logger.Debugf("[DoAuthRequest] sign base64: %s, timestamp: %s", signStr, timestamp)

	headers = map[string]string{
		"Content-Type": "application/json; charset=UTF-8",
		//"Accept":               "application/json",
		"OK-ACCESS-KEY":        t.apiOpts.Key,
		"OK-ACCESS-PASSPHRASE": t.apiOpts.Passphrase,
		"OK-ACCESS-SIGN":       signStr,
		"OK-ACCESS-TIMESTAMP":  timestamp}

	respBody, err := goex.GetHttpCli().DoRequest(httpMethod, reqUrl, reqBodyStr, headers)
	if err != nil {
		return respBody, err
	}

	var baseResp BaseResp
	err = t.unmarshalOpts.ResponseUnmarshaler(respBody, &baseResp)
	if err != nil {
		return respBody, err
	}

	if baseResp.Code != 0 {
		return nil, errors.New(baseResp.Msg)
	}

	return baseResp.Data, nil
}

func NewTrade(opts ...goex.ApiOption) *Trade {
	var api = new(Trade)
	for _, opt := range opts {
		opt(&api.apiOpts)
	}
	return api
}
