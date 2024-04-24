package controller

import (
	"fmt"
	"time"
	"webssh/core"

	"github.com/gin-gonic/gin"
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
	// 输出一下入参
	fmt.Println("CheckSSH in: ", c.Request.URL.Query())
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
	defer sshClient.Close()

	if err != nil {
		fmt.Println(err)
		responseBody.Msg = err.Error()
	}
	return &responseBody
}
