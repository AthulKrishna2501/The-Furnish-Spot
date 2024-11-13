package order

import (
	"net/http"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func ListOrders(c *gin.Context) {
	status := c.Query("status")
	sort := c.Query("sort")
	order := c.Query("order")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var orders []models.Order
	db := db.Db.Model(&models.Order{})

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if startDate != "" && endDate != "" {
		db = db.Where("order_date BETWEEN ? AND ?", startDate, endDate)
	}

	switch sort {
	case "order_date":
		db = db.Order("order_date " + order)
	case "total":
		db = db.Order("total " + order)
	default:
		if sort != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort parameter"})
			return
		}
	}

	if err := db.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot fetch orders"})
		return
	}

	var orderResponses []responsemodels.OrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, responsemodels.OrderResponse{
			OrderID:        order.OrderID,
			UserID:         order.UserID,
			Quantity:       order.Quantity,
			DiscountAmount: order.Discount,
			Total:          order.Total,
			Method:         order.Method,
			Status:         order.Status,
			PaymentStatus:  order.PaymentStatus,
			OrderDate:      order.OrderDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": orderResponses})
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
	} else {
		if order.Status == "Canceled" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order is already canceled"})
			return
		} else if order.Status == "Delivered" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel a delivered order"})
			return
		}
	}

	order.Status = input.Status
	if err := db.Db.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully", "new_status": order.Status})
}
