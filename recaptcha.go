package anticaptcha

import "net/url"

//RecaptchaResult is the api response from a recaptcha task
type RecaptchaResult struct {
	GRecaptchaResponse string      `json:"gRecaptchaResponse"`
	Cookies            interface{} `json:"cookies"`
}

//Recaptcha submits and retrieves a recaptcha task
func (c *Client) Recaptcha(siteURL, siteKey, userAgent string, proxy *url.URL, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "NoCaptchaTask",
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

	taskId, err = c.createTask(data)
	if err != nil {
		return
	}

	err = c.fetchTask(taskId, &result)
	if err != nil {
		return
	}

	return
}

//RecaptchaProxyless submits and retrieves a recaptcha task
func (c *Client) RecaptchaProxyless(siteURL, siteKey string, opts ...OptionalValue) (result RecaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":       "NoCaptchaTaskProxyless",
		"websiteURL": siteURL,
		"websiteKey": siteKey,
	}

	for _, v := range opts {
		v(data)
	}

	taskId, err = c.createTask(data)
	if err != nil {
		return
	}

	err = c.fetchTask(taskId, &result)
	if err != nil {
		return
	}

	return
}

//RecaptchaV3 submits and retrieves a recaptcha v3 task
func (c *Client) RecaptchaV3(siteURL, siteKey string, minScore float64, opts ...OptionalValue) (result RecaptchaResult, err error) {
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

	taskId, err = c.createTask(data)
	if err != nil {
		return
	}

	err = c.fetchTask(taskId, &result)
	if err != nil {
		return
	}

	return
}
