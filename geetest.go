package anticaptcha

import (
	"context"
	"net/url"
)

//GeeTestResult is the api response from a geetest task
type GeeTestResult struct {
	Challenge string `json:"challenge"`
	Validate  string `json:"validate"`
	SecCode   string `json:"seccode"`
}

//GeeTest submits and retrieves a geetest task
func (c *Client) GeeTest(ctx context.Context, siteURL, siteKey, challenge, userAgent string, proxy *url.URL, opts ...OptionalValue) (result GeeTestResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "GeeTestTask",
		"websiteURL": siteURL,
		"gt":         siteKey,
		"challenge":  challenge,
		"userAgent":  userAgent,
	}

	for _, v := range opts {
		v(data)
	}

	err = addProxyInfo(proxy, data)
	if err != nil {
		return
	}

	taskId, err = c.createTask(ctx, data)
	if err != nil {
		return
	}

	err = c.fetchTask(ctx, taskId, &result)
	if err != nil {
		return
	}

	return
}

//GeeTestProxyless submits and retrieves a geetest task
func (c *Client) GeeTestProxyless(ctx context.Context, siteURL, siteKey, challenge string, opts ...OptionalValue) (result GeeTestResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "GeeTestTaskProxyless",
		"websiteURL": siteURL,
		"gt":         siteKey,
		"challenge":  challenge,
	}

	for _, v := range opts {
		v(data)
	}

	taskId, err = c.createTask(ctx, data)
	if err != nil {
		return
	}

	err = c.fetchTask(ctx, taskId, &result)
	if err != nil {
		return
	}

	return
}
