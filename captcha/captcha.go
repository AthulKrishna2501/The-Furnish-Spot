package captcha

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func GetCaptcha(c *gin.Context) {
	captchaID := captcha.New()
	c.JSON(http.StatusOK, gin.H{
		"captcha_id":  captchaID,
		"captcha_url": "/captcha/" + captchaID,
	})
}

func CaptchaHandler(c *gin.Context) {
	captchaID := c.Param("captchaID")
	c.Header("Content-Type", "image/png")
	captcha.WriteImage(c.Writer, captchaID, captcha.StdWidth, captcha.StdHeight)
}
