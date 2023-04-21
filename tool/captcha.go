package tool

import (
	"fmt"
	"github.com/dchest/captcha"
	"github.com/kataras/iris"
)

const (
	StdWidth  = 80
	StdHeight = 40
)

func GetCaptchaID(ctx iris.Context) {
	captchaMap := make(map[string]interface{}, 0)
	captchaMap["error_code"] = 0
	captchaMap["msg"] = "Get Captcha successfully"
	captchaMap["id"] = captcha.NewLen(4)
	ctx.JSON(captchaMap)
	fmt.Println(captchaMap)
	return
}

func GetCaptchaImg(ctx iris.Context) {
	captcha.Server(StdWidth, StdHeight).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}
