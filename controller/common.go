package controller

import "time"

type ResponseBody struct {
	Duration string
	Data     interface{}
	Msg      string
}

func TimeCost(start time.Time, body *ResponseBody) {
	body.Duration = time.Since(start).String()
}
