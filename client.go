package anticaptcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const softId = 948

//Client allows programmatic access to the anti-captcha api
type Client struct {
	key           string
	host          string
	client        *http.Client
	delay         time.Duration
	checkInterval time.Duration
}

//ClientOption is an option used to modify the anti-captcha client
type ClientOption func(c *Client)

//OptionalValue is an option used to add values to an api request
type OptionalValue func(map[string]interface{})

//WithOptional is an option that adds values to an api request
func WithOptional(key string, value interface{}) OptionalValue {
	return func(m map[string]interface{}) {
		m[key] = value
	}
}

//NewClient returns a new Client with the applied options
func NewClient(key string, opts ...ClientOption) (client *Client) {
	client = &Client{
		key:           key,
		client:        http.DefaultClient,
		delay:         time.Second * 10,
		checkInterval: time.Second * 3,
	}

	client.host = "api.anti-captcha.com"

	for _, v := range opts {
		v(client)
	}

	return
}

//WithDelay is an option that makes the Client use the provided delay
func WithDelay(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.delay = duration
	}
}

//WithHost is an option that makes the Client use the provided host
func WithHost(host string) ClientOption {
	return func(c *Client) {
		c.host = host
	}
}

//WithCheckInterval is an option that makes the Client use the provided check interval
func WithCheckInterval(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.checkInterval = duration
	}
}

func (c *Client) createTask(ctx context.Context, task interface{}) (taskId int64, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"softId":    softId,
		"task":      task,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := createRequest(ctx, http.MethodPost, fmt.Sprintf("https://%s/createTask", c.host), bytes.NewReader(sendBytes))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	var response struct {
		ErrorId          int64  `json:"errorId"`
		ErrorCode        string `json:"errorCode"`
		ErrorDescription string `json:"errorDescription"`
		TaskId           int64  `json:"taskId"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.ErrorId != 0 {
		err = errors.New("anticaptcha: " + response.ErrorDescription)
		return
	}

	taskId = response.TaskId
	return
}

func (c *Client) getTaskResult(ctx context.Context, taskId int64, dst interface{}) (ready bool, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := createRequest(ctx, http.MethodPost, fmt.Sprintf("https://%s/getTaskResult", c.host), bytes.NewReader(sendBytes))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	var response struct {
		ErrorId          int64           `json:"errorId"`
		ErrorCode        string          `json:"errorCode"`
		ErrorDescription string          `json:"errorDescription"`
		Status           string          `json:"status"`
		Solution         json.RawMessage `json:"solution"`
		Cost             string          `json:"cost"`
		IP               string          `json:"ip"`
		CreateTime       int64           `json:"createTime"`
		EndTime          int64           `json:"endTime"`
		SolveCount       int64           `json:"solveCount"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.ErrorId != 0 {
		ready = false
		err = errors.New("anticaptcha: " + response.ErrorDescription)
		return
	}

	switch response.Status {
	case "ready":
		ready = true
		err = json.Unmarshal(response.Solution, dst)
		return
	case "processing":
		ready = false
		return
	}
	return
}

func (c *Client) fetchTask(ctx context.Context, taskId int64, dst interface{}) (err error) {
	ticker := time.NewTicker(c.checkInterval)
	time.Sleep(c.delay)

	for {
		select {
		case <-ticker.C:
			ready, err := c.getTaskResult(ctx, taskId, dst)
			if err != nil {
				return err
			}

			if ready {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

//GetBalance retrieves the current account balance
func (c *Client) GetBalance(ctx context.Context) (balance float64, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := createRequest(ctx, http.MethodPost, fmt.Sprintf("https://%s/getBalance", c.host), bytes.NewReader(sendBytes))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	var response struct {
		ErrorId          int64   `json:"errorId"`
		ErrorCode        string  `json:"errorCode"`
		ErrorDescription string  `json:"errorDescription"`
		Balance          float64 `json:"balance"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.ErrorId != 0 {
		err = errors.New("anticaptcha: " + response.ErrorDescription)
		return
	}

	balance = response.Balance

	return
}

//ReportIncorrectImageCaptcha reports an incorrect captcha for a refund
func (c *Client) ReportIncorrectImageCaptcha(ctx context.Context, taskId int64) (err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := createRequest(ctx, http.MethodPost, fmt.Sprintf("https://%s/reportIncorrectImageCaptcha", c.host), bytes.NewReader(sendBytes))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	var response struct {
		ErrorId          int64  `json:"errorId"`
		ErrorCode        string `json:"errorCode"`
		ErrorDescription string `json:"errorDescription"`
		Status           string `json:"status"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.ErrorId != 0 {
		err = errors.New("anticaptcha: " + response.ErrorDescription)
		return
	}

	return
}

//ReportIncorrectRecaptcha reports an incorrect captcha for a refund
func (c *Client) ReportIncorrectRecaptcha(ctx context.Context, taskId int64) (err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	req, err := createRequest(ctx, http.MethodPost, fmt.Sprintf("https://%s/reportIncorrectRecaptcha", c.host), bytes.NewReader(sendBytes))
	if err != nil {
		return
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = resp.Body.Close()
	if err != nil {
		return
	}

	var response struct {
		ErrorId          int64  `json:"errorId"`
		ErrorCode        string `json:"errorCode"`
		ErrorDescription string `json:"errorDescription"`
		Status           string `json:"status"`
	}

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return
	}

	if response.ErrorId != 0 {
		err = errors.New("anticaptcha: " + response.ErrorDescription)
		return
	}

	return
}

func addProxyInfo(proxy *url.URL, to map[string]interface{}) error {
	p, err := strconv.Atoi(proxy.Port())
	if err != nil {
		return err
	}

	to["proxyPort"] = p
	to["proxyType"] = proxy.Scheme
	to["proxyAddress"] = proxy.Hostname()

	pp, hasPassword := proxy.User.Password()
	pu := proxy.User.Username()

	if pu != "" {
		to["proxyLogin"] = pu
		if hasPassword {
			to["proxyPassword"] = pp
		}
	}

	return nil
}

func createRequest(ctx context.Context, method, url string, body io.Reader) (r *http.Request, err error) {
	r, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	r.Header.Set("User-Agent", "anticaptcha (github.com/aidenesco/anticaptcha)")
	r.Header.Set("Content-Type", "application/json")

	return
}
