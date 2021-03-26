package anticaptcha

import (
	"bytes"
	"encoding/json"
	"errors"
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
	client        *http.Client
	timeout       time.Duration
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
		timeout:       time.Minute,
		delay:         time.Second * 10,
		checkInterval: time.Second * 3,
	}

	for _, v := range opts {
		v(client)
	}

	return
}

//WithTimeout is an option that makes the Client use the provided timeout
func WithTimeout(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = duration
	}
}

//WithDelay is an option that makes the Client use the provided delay
func WithDelay(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.delay = duration
	}
}

//WithCheckInterval is an option that makes the Client use the provided check interval
func WithCheckInterval(duration time.Duration) ClientOption {
	return func(c *Client) {
		c.checkInterval = duration
	}
}

func (c *Client) createTask(task interface{}) (taskId int64, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"softId":    softId,
		"task":      task,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := c.client.Post("https://api.anti-captcha.com/createTask", "application/json", bytes.NewReader(sendBytes))
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
		err = errors.New("anticaptcha: " + response.ErrorCode)
		return
	}

	taskId = response.TaskId
	return
}

func (c *Client) getTaskResult(taskId int64, dst interface{}) (ready bool, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := c.client.Post("https://api.anti-captcha.com/getTaskResult", "application/json", bytes.NewReader(sendBytes))
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

	if response.Status != "ready" {
		ready = false
		return
	}

	ready = true

	err = json.Unmarshal(response.Solution, dst)
	if err != nil {
		return
	}

	return
}

func (c *Client) fetchTask(taskId int64, dst interface{}) (err error) {
	time.Sleep(c.delay)
	ticker := time.Tick(c.checkInterval)
	timeout := time.After(c.timeout)

	for {
		select {
		case <-ticker:
			ready, err := c.getTaskResult(taskId, dst)
			if err != nil {
				return err
			}

			if ready {
				return nil
			}
		case <-timeout:
			return errors.New("anticaptcha: timeout exceeded fetching task")
		}
	}

}

//GetBalance retrieves the current account balance
func (c *Client) GetBalance() (balance float64, err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := c.client.Post("https://api.anti-captcha.com/getBalance", "application/json", bytes.NewReader(sendBytes))
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
		err = errors.New("anticaptcha: " + response.ErrorCode)
		return
	}

	balance = response.Balance

	return
}

//ReportIncorrectImageCaptcha reports an incorrect captcha for a refund
func (c *Client) ReportIncorrectImageCaptcha(taskId int64) (err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := c.client.Post("https://api.anti-captcha.com/reportIncorrectImageCaptcha", "application/json", bytes.NewReader(sendBytes))
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
		err = errors.New("anticaptcha: " + response.ErrorCode)
		return
	}

	return
}

//ReportIncorrectRecaptcha reports an incorrect captcha for a refund
func (c *Client) ReportIncorrectRecaptcha(taskId int64) (err error) {
	data := map[string]interface{}{
		"clientKey": c.key,
		"taskId":    taskId,
	}

	sendBytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	resp, err := c.client.Post("https://api.anti-captcha.com/reportIncorrectRecaptcha", "application/json", bytes.NewReader(sendBytes))
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
		err = errors.New("anticaptcha: " + response.ErrorCode)
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
