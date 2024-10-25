package order

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
)

func ListOrders(c *gin.Context) {
	var orders []models.Order

	if err := db.Db.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": orders})
}

func ChangeOrderStatus(c *gin.Context) {
	var orders models.Order

	OrderId := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := db.Db.First(&orders, OrderId).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	if input.Status == "Canceled" && orders.Status != "Canceled" {
		var products models.Product
		if err := db.Db.First(&products, orders.ProductID).Error; err == nil {
			products.Quantity += orders.Quantity
			db.Db.Save(&products)
		}
	}
	orders.Status = input.Status

	db.Db.Save(&orders)

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func CancelOrders(c *gin.Context) {
	var order models.Order

	orderID := c.Param("id")

	if err := db.Db.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == "Devilered" || order.Status == "Shipped" || order.Status == "Failed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
		return
	}

	order.Status = "Canceled"

	db.Db.Save(&order)

	var product models.Product

	if err := db.Db.First(&product, order.ProductID).Error; err == nil {
		product.Quantity += order.Quantity
		db.Db.Save(&product)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled successfully"})
}
