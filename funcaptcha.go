package anticaptcha

import "net/url"

//FuncaptchaResult is the api response from a funcaptcha task
type FuncaptchaResult struct {
	Token string `json:"token"`
}

//Funcaptcha submits and retrieves a funcaptcha task
func (c *Client) Funcaptcha(siteURL, siteKey, userAgent string, proxy *url.URL, opts ...OptionalValue) (result FuncaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":             "FunCaptchaTask",
		"websiteURL":       siteURL,
		"websitePublicKey": siteKey,
		"userAgent":        userAgent,
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

//FuncaptchaProxyless submits and retrieves a funcaptcha task
func (c *Client) FuncaptchaProxyless(siteURL, siteKey string, opts ...OptionalValue) (result FuncaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":             "FunCaptchaTaskProxyless",
		"websiteURL":       siteURL,
		"websitePublicKey": siteKey,
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
