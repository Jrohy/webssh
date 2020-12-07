package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"webssh/core"
)

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

// CheckSSH 检查ssh连接是否能连接
func CheckSSH(c *gin.Context) *ResponseBody {
	responseBody := ResponseBody{Msg: "success"}
	defer TimeCost(time.Now(), &responseBody)
	sshInfo := c.DefaultQuery("sshInfo", "")
	sshClient, err := core.DecodedMsgToSSHClient(sshInfo)
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
		return &responseBody
	}
	err = sshClient.GenerateClient()
	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
	}
	return &responseBody
}
