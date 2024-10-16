package route

import (
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/category"
	"github.com/AthulKrishna2501/The-Furniture-Spot/captcha"
	"github.com/AthulKrishna2501/The-Furniture-Spot/user"

	"github.com/gin-gonic/gin"
)

func RegisterURL(router *gin.Engine) {
	//User
	router.GET("/getcaptcha", captcha.GetCaptcha)
	router.GET("/captcha/:captchaID", captcha.CaptchaHandler)
	router.POST("/signup", user.SignUp)
	router.POST("/verifyotp", user.VerifyOTP)
	router.POST("/resendotp/:email", user.ResendOTP)
	router.POST("/login", user.Login)
	router.GET("/products", user.ListProducts)
	router.POST("/search-products", user.SearchProduct)

	//Admin

	router.POST("/adminlogin", admin.AdminLogin)
	router.GET("/viewcategories", category.ViewCategory)
	router.POST("/addcategory", category.AddCategory)
	router.PUT("/updatecategory/:id", category.EditCategory)
	router.DELETE("/deletecategory/:id", category.DeleteCategory)

}
