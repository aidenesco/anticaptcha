package anticaptcha

import (
	"context"
	"net/url"
)

//RecaptchaResult is the api response from a recaptcha task
type RecaptchaResult struct {
	GRecaptchaResponse string      `json:"gRecaptchaResponse"`
	Cookies            interface{} `json:"cookies"`
}

//RecaptchaV2 submits and retrieves a recaptcha v2 task
func (c *Client) RecaptchaV2(ctx context.Context, siteURL, siteKey, userAgent string, proxy *url.URL, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "RecaptchaV2Task",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
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

//RecaptchaV2Proxyless submits and retrieves a recaptcha v2 task
func (c *Client) RecaptchaV2Proxyless(ctx context.Context, siteURL, siteKey string, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "RecaptchaV2TaskProxyless",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
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

//RecaptchaV2Enterprise submits and retrieves a recaptcha v2 enterprise task
func (c *Client) RecaptchaV2Enterprise(ctx context.Context, siteURL, siteKey, userAgent string, proxy *url.URL, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "RecaptchaV2EnterpriseTask",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
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

//RecaptchaV2EnterpriseProxyless submits and retrieves a recaptcha v2 enterprise task
func (c *Client) RecaptchaV2EnterpriseProxyless(ctx context.Context, siteURL, siteKey string, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "RecaptchaV2EnterpriseTaskProxyless",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
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

//RecaptchaV3Proxyless submits and retrieves a recaptcha v3 task
func (c *Client) RecaptchaV3Proxyless(ctx context.Context, siteURL, siteKey string, minScore float64, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "RecaptchaV3TaskProxyless",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
		"minScore":   minScore,
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

//RecaptchaV3Enterprise submits and retrieves a recaptcha v3 enterprise task
func (c *Client) RecaptchaV3Enterprise(ctx context.Context, siteURL, siteKey string, minScore float64, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":         "RecaptchaV3TaskProxyless",
		"websiteURL":   siteURL,
		"websiteKey":   siteKey,
		"minScore":     minScore,
		"isEnterprise": true,
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
