package anticaptcha

import "context"

//SquareNetResult is the api response from a squarenet task
type SquareNetResult struct {
	CellNumbers []int64 `json:"cellNumbers"`
}

//SquareNet submits and retrieves a squarenet task
func (c *Client) SquareNet(ctx context.Context, body, object string, rows, columns int64) (result SquareNetResult, err error) {
	var taskId int64
	data := map[string]interface{}{
		"type":         "SquareNetTextTask",
		"body":         body,
		"objectName":   object,
		"rowsCount":    rows,
		"columnsCount": columns,
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
