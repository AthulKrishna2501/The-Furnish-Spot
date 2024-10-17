package route

import (
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin"
	adminuser "github.com/AthulKrishna2501/The-Furniture-Spot/admin/adminUser"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/category"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/product"
	"github.com/AthulKrishna2501/The-Furniture-Spot/captcha"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
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
	router.GET("/products", middleware.AuthMiddleware("user"), user.ListProducts)
	router.POST("/search-products", middleware.AuthMiddleware("user"), user.SearchProduct)

	//Admin
	router.POST("/adminlogin", admin.AdminLogin)
	router.GET("/viewcategories", category.ViewCategory)
	router.POST("/addcategory", category.AddCategory)
	router.PUT("/updatecategory/:id", category.EditCategory)
	router.DELETE("/deletecategory/:id", category.DeleteCategory)

	router.GET("/viewproducts", product.ViewProducts)
	router.POST("/addproducts", product.AddProducts)
	router.PUT("/updateproduct/:id", product.UpdateProduct)
	router.DELETE("/deleteproduct/:id", product.DeleteProduct)

	router.GET("/listusers", adminuser.ListUsers)
	router.POST("blockuser/:id", adminuser.BlockUser)
	router.POST("/unblockuser/:id", adminuser.UnblockUser)
}
