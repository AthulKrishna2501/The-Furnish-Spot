package route

import (
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin"
	adminuser "github.com/AthulKrishna2501/The-Furniture-Spot/admin/adminUser"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/category"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/coupon"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/offer"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/order"
	"github.com/AthulKrishna2501/The-Furniture-Spot/admin/product"
	salesreport "github.com/AthulKrishna2501/The-Furniture-Spot/admin/salesReport"
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
	router.GET("/googlelogin", user.HandleGoogleLogin)
	router.GET("/auth/google/callback", user.HandleGoogleCallback)
	router.POST("/verifyotp", user.VerifyOTP)
	router.POST("/resendotp/:email", user.ResendOTP)
	router.POST("/login", user.Login)
	router.PUT("/forgotpassword", middleware.AuthMiddleware("user"), user.ForgotPassword)

	//Products
	router.GET("/products", user.ViewProducts)
	router.GET("/search-products", user.SearchProducts)

	//Profile
	router.GET("/viewprofile", middleware.AuthMiddleware("user"), user.UserProfile)
	router.POST("/editprofile", middleware.AuthMiddleware("user"), user.EditProfile)
	router.GET("/viewaddress", middleware.AuthMiddleware("user"), user.ViewAddress)
	router.POST("profile/addaddress", middleware.AuthMiddleware("user"), user.AddAddress)
	router.PUT("profile/updateaddress/:id", middleware.AuthMiddleware("user"), user.EditAddress)
	router.DELETE("/profile/deleteaddress/:id", middleware.AuthMiddleware("user"), user.DeleteAddress)
	router.GET("/user/wallet", middleware.AuthMiddleware("user"), user.ViewWallet)
	router.GET("/wallet/transactions", middleware.AuthMiddleware("user"), user.GetWalletTransactions)

	//Orders
	router.GET("/vieworders", middleware.AuthMiddleware("user"), user.ViewOrders)
	router.DELETE("/orders/:id/delete", middleware.AuthMiddleware("user"), user.CancelOrders)
	router.POST("/users/order", middleware.AuthMiddleware("user"), user.Orders)
	router.GET("/paypal/confirmpayment", user.CapturePayPalOrder)
	router.GET("/paypal/cancel-payment", user.CapturePayPalOrder)
	router.POST("/user/returnorder", middleware.AuthMiddleware("user"), user.ReturnOrder)
	router.POST("/user/generate-invoice/:id", middleware.AuthMiddleware("user"), user.GenerateInvoiceHandler)
	//Cart
	router.GET("/user/cart", middleware.AuthMiddleware("user"), user.Cart)
	router.POST("/user/addtocart", middleware.AuthMiddleware("user"), user.AddToCart)
	router.DELETE("user/removeitem/:id", middleware.AuthMiddleware("user"), user.RemoveItem)

	//Whishlist
	router.GET("/user/viewwhishlist", middleware.AuthMiddleware("user"), user.ViewWhishlist)
	router.POST("/user/addtowhishlist", middleware.AuthMiddleware("user"), user.AddToWhishlist)
	router.DELETE("/user/removeitem", middleware.AuthMiddleware("user"), user.WishlistRemoveItem)
	router.DELETE("/user/clearwishlist", middleware.AuthMiddleware("user"), user.ClearWishlist)

	//Coupons
	router.GET("/coupons", middleware.AuthMiddleware("user"), user.ViewCoupons)

	//ReivewRatings
	router.POST("/user/review", middleware.AuthMiddleware("user"), user.AddReviews)
	router.PUT("/user/editreview", middleware.AuthMiddleware("user"), user.EditReview)
	router.DELETE("/user/deletereview/:id", middleware.AuthMiddleware("user"), user.DeleteReview)

	//Admin
	router.POST("/adminlogin", admin.AdminLogin)
	router.GET("/viewcategories", middleware.AuthMiddleware("admin"), category.ViewCategory)
	router.POST("/addcategory", middleware.AuthMiddleware("admin"), category.AddCategory)
	router.PUT("/updatecategory/:id", middleware.AuthMiddleware("admin"), category.EditCategory)
	router.DELETE("/deletecategory/:id", middleware.AuthMiddleware("admin"), category.DeleteCategory)

	router.GET("/viewproducts", middleware.AuthMiddleware("admin"), product.ViewProducts)
	router.POST("/addproducts", middleware.AuthMiddleware("admin"), product.AddProducts)
	router.PUT("/updateproduct/:id", middleware.AuthMiddleware("admin"), product.UpdateProduct)
	router.DELETE("/deleteproduct/:id", middleware.AuthMiddleware("admin"), product.DeleteProduct)
	router.PUT("/admin/updatestock/:id", middleware.AuthMiddleware("admin"), product.UpdateProductStock)

	router.GET("/listusers", middleware.AuthMiddleware("admin"), adminuser.ListUsers)
	router.POST("blockuser/:id", middleware.AuthMiddleware("admin"), adminuser.BlockUser)
	router.POST("/unblockuser/:id", middleware.AuthMiddleware("admin"), adminuser.UnblockUser)

	router.GET("/admin/listorders", middleware.AuthMiddleware("admin"), order.ListOrders)
	router.PUT("/admin/changeorderstatus/:id", middleware.AuthMiddleware("admin"), order.ChangeOrderStatus)

	router.GET("/admin/viewcoupons", middleware.AuthMiddleware("admin"), coupon.ViewCoupons)
	router.POST("/admin/addcoupon", middleware.AuthMiddleware("admin"), coupon.AddCoupon)
	router.DELETE("/admin/deletecoupon/:id", middleware.AuthMiddleware("admin"), coupon.DeleteCoupon)

	router.GET("/admin/viewoffers", middleware.AuthMiddleware("admin"), offer.ViewOffers)
	router.POST("/admin/addoffer", middleware.AuthMiddleware("admin"), offer.AddOffer)
	router.PUT("/admin/updateoffer", middleware.AuthMiddleware("admin"), offer.UpdateOffer)

	router.GET("/generate-report", middleware.AuthMiddleware("admin"), salesreport.GenerateReport)
	router.GET("/get-sales-data", middleware.AuthMiddleware("admin"), salesreport.GetSalesData)
	router.GET("/top-selling-product", middleware.AuthMiddleware("admin"), salesreport.GetTopSellingProducts)
	router.GET("/top-selling-category", middleware.AuthMiddleware("admin"), salesreport.GetTopSellingCategories)
	router.GET("/ledger-book", middleware.AuthMiddleware("admin"), salesreport.GetLedgerBook)

}
