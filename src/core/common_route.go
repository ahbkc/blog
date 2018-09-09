package core

import (
	"net/http"
	"github.com/mojocn/base64Captcha"
	"path/filepath"
	"strings"
	"utils"
)

//file path handel
func GetFilePath(name string) string {
	if strings.HasSuffix(name,".tmpl") {
		return filepath.Join(utils.Dir, "/src/resource", utils.AdminTmplHtmlPath, name)
	}
	return filepath.Join(utils.Dir, "/src/resource", utils.AdminHtmlPath, name)
}

//output verifyCode picture
func VerifyCodeGenerate(w http.ResponseWriter, r *http.Request) {
	var config = base64Captcha.ConfigDigit{
		Height:     38,
		Width:      120,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	var capD base64Captcha.CaptchaInterface
	utils.IdKeyD, capD = base64Captcha.GenerateCaptcha("", config)
	capD.WriteTo(w) //write to response
}
