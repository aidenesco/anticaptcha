package anticaptcha

import "context"

//ImageCaptchaResult is the api response from an image task
type ImageCaptchaResult struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

//ImageToText submits and retrieves an image task
func (c *Client) ImageToText(ctx context.Context, body string, opts ...OptionalValue) (result ImageCaptchaResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type": "ImageToTextTask",
		"body": body,
	}

	for _, v := range opts {
		v(data)
	}

	taskId, err = c.createTask(ctx, data)
	if err != nil {
		return
	}

	err = c.fetchTask(ctx, taskId, result)
	if err != nil {
		return
	}

	return
}
