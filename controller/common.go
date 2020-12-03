package controller

import "time"

// ResponseBody 响应信息结构体
type ResponseBody struct {
	Duration string
	Data     interface{}
	Msg      string
}

// TimeCost 计算方法执行耗时
func TimeCost(start time.Time, body *ResponseBody) {
	body.Duration = time.Since(start).String()
}
