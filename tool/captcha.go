package tool

import (
	"github.com/dchest/captcha"
	"github.com/kataras/iris"
	"time"
)

const (
	StdWidth  = 80
	StdHeight = 40
)

//func GetCaptchaID(ctx iris.Context) {
//	captchaMap := make(map[string]interface{}, 0)
//	captchaID := captcha.NewLen(4)
//	captchaMap["error_code"] = 0
//	captchaMap["msg"] = "Get Captcha successfully"
//	captchaMap["id"] = captchaID
//	ctx.SetCookieKV("captcha_id", captchaID, iris.CookieExpires(time.Minute*5))
//}

func GetCaptchaID(ctx iris.Context) {
	captchaMap := make(map[string]interface{}, 0)
	captchaID := captcha.NewLen(4)
	captchaMap["error_code"] = 0
	captchaMap["msg"] = "Get Captcha successfully"
	captchaMap["id"] = captchaID
	ctx.SetCookieKV("captcha_id", captchaID, iris.CookieExpires(time.Minute*5))
	ctx.ContentType("image/png")
	captcha.WriteImage(ctx.ResponseWriter(), captchaID, StdWidth, StdHeight)
}

func GetCaptchaImg(ctx iris.Context) {
	captcha.Server(StdWidth, StdHeight).ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}

func VerifyCaptcha(ctx iris.Context) {
	captchaID := ctx.GetCookie("captcha_id")
	captchaInput := ctx.URLParam("code")
	//fmt.Println(captchaID, captchaInput)

	if captcha.VerifyString(captchaID, captchaInput) {
		captchaMap := make(map[string]interface{}, 0)
		captchaMap["error_code"] = 0
		captchaMap["msg"] = "Pass verification"
		ctx.JSON(captchaMap)
	} else {
		captchaMap := make(map[string]interface{}, 0)
		captchaMap["error_code"] = 1
		captchaMap["msg"] = "Fail to pass verification"
		ctx.JSON(captchaMap)
	}
}
