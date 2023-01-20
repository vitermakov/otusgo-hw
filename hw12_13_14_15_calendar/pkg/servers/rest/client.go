package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/pkg/servers/rest/rqres"
)

type ClientOpt func(client *Client)

type ClientAuth interface {
	Authorize(request *http.Request) error
}

type Client struct {
	client  *http.Client
	baseURL string
	auth    ClientAuth
}

func NewClient(baseURL string, opts ...ClientOpt) *Client {
	api := &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
	for _, optFunc := range opts {
		optFunc(api)
	}
	return api
}

func isGetMethod(meth string) bool {
	return strings.ToUpper(meth) == http.MethodGet
}

func isPostMethod(meth string) bool {
	meth = strings.ToUpper(meth)
	return meth == http.MethodPost || meth == http.MethodPut || meth == http.MethodDelete
}

func makeJSONBody(params interface{}) (io.Reader, error) {
	if params == nil {
		return nil, nil
	}
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bs), nil
}

func makeQueryURL(baseURL, resource string, params interface{}) (string, error) {
	var queryURL strings.Builder
	var query = make(url.Values)
	queryURL.WriteString(baseURL)
	queryURL.WriteString(resource)
	if params != nil {
		paramsMap, ok := params.(map[string]interface{})
		if !ok {
			return "", errors.New("wrong data for query params")
		}
		for key, value := range paramsMap {
			query.Set(key, fmt.Sprintf("%v", value))
		}
	}
	if len(query) > 0 {
		queryURL.WriteByte('?')
		queryURL.WriteString(query.Encode())
	}
	return queryURL.String(), nil
}

func WithAuth(auth ClientAuth) ClientOpt {
	return func(api *Client) {
		api.auth = auth
	}
}

func (ac *Client) authorize(request *http.Request) error {
	if ac.auth != nil {
		return ac.auth.Authorize(request)
	}
	return nil
}

func (ac *Client) doQuery(ctx context.Context, method, resource string, params interface{}) (*http.Response, error) {
	var (
		requestURL  string
		requestBody io.Reader
		err         error
	)
	switch {
	case isPostMethod(method):
		if requestURL, err = makeQueryURL(ac.baseURL, resource, nil); err != nil {
			return nil, err
		}
		if requestBody, err = makeJSONBody(params); err != nil {
			return nil, err
		}
	case isGetMethod(method):
		if requestURL, err = makeQueryURL(ac.baseURL, resource, params); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("method is not allowed")
	}
	request, err := http.NewRequestWithContext(ctx, method, requestURL, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	if err = ac.authorize(request); err != nil {
		return nil, err
	}
	response, err := ac.client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (ac *Client) Get(ctx context.Context, resource string, query map[string]interface{}) (*http.Response, error) {
	return ac.doQuery(ctx, http.MethodGet, resource, query)
}

func (ac *Client) Post(ctx context.Context, resource string, data interface{}) (*http.Response, error) {
	return ac.doQuery(ctx, http.MethodPost, resource, data)
}

func (ac *Client) Put(ctx context.Context, resource string, data interface{}) (*http.Response, error) {
	return ac.doQuery(ctx, http.MethodPut, resource, data)
}

func (ac *Client) Delete(ctx context.Context, resource string, data interface{}) (*http.Response, error) {
	return ac.doQuery(ctx, http.MethodDelete, resource, data)
}

func EncodeResponse(resp *http.Response, dataObj interface{}, reqResp bool) error {
	if resp == nil {
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return err
	}
	return rqres.ParseResponse(resp.StatusCode, data, dataObj, reqResp)
}
