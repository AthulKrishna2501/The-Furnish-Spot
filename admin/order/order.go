package order

import (
	"fmt"
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func ListOrders(c *gin.Context) {
	var orders []models.Order

	if err := db.Db.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch orders"})
		return
	}

	var orderResponses []responsemodels.OrderResponse

	for _, order := range orders {
		fmt.Printf("Order ID: %d, Status: %s\n", order.OrderID, order.Status)

		var totalQuantity int
		var orderItems []models.OrderItem
		if err := db.Db.Where("order_id = ?", order.OrderID).Find(&orderItems).Error; err == nil {
			for _, item := range orderItems {
				totalQuantity += item.Quantity
			}
		} else {
			fmt.Printf("Error fetching items for Order ID %d: %v\n", order.OrderID, err)
		}

		orderResponses = append(orderResponses, responsemodels.OrderResponse{
			UserID:        order.UserID,
			OrderID:       order.OrderID,
			Total:         order.Total,
			Quantity:      totalQuantity,
			Status:        order.Status,
			Method:        order.Method,
			PaymentStatus: order.PaymentStatus,
			OrderDate:     order.OrderDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"orders": orderResponses})
}

func ChangeOrderStatus(c *gin.Context) {
	var order models.Order

	OrderID := c.Param("id")

	if err := db.Db.First(&order, OrderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if order.Status == "Delivered" || order.Status == "Shipped" || order.Status == "Failed" || order.Status == "Canceled" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid status"})
		return
	}

	if input.Status == "Canceled" && order.Status != "Canceled" {
		var product models.Product
		var item models.OrderItem
		if err := db.Db.First(&product, item.ProductID).Error; err == nil {
			product.Quantity += order.Quantity
			if err := db.Db.Save(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product quantity"})
				return
			}
		}
	}

	order.Status = input.Status
	if err := db.Db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "new_status": order.Status})
}
