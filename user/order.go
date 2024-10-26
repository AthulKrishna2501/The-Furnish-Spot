package user

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/middleware"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models/responsemodels"
	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	var input models.OrderInput
	var address models.Address
	var cart []models.Cart

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	claims, _ := c.Get("claims")
	customClaims, ok := claims.(*middleware.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := customClaims.ID

	if err := db.Db.Where("user_id=? AND address_id =?", userID, input.AddressID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	if err := db.Db.Table("carts").
		Select("carts.*, products.product_id, products.price").
		Joins("left join products on carts.product_id = products.product_id").
		Where("carts.user_id = ?", userID).
		Scan(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not fetch cart", "details": err.Error()})
		return
	}

	fmt.Printf("Fetched Cart Contents: %+v\n", cart)

	if len(cart) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	var cartLock sync.Mutex
	cartLock.Lock()
	defer cartLock.Unlock()

	var totalAmount float64
	var totalQuantity int
	var orderItems []models.OrderItem

	for _, item := range cart {
		productID := item.ProductID
		product := models.Product{}

		if err := db.Db.First(&product, productID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "product_id": productID})
			return
		}

		fmt.Printf("Product ID: %d, Price: %.2f, Cart Quantity: %d\n", product.ProductID, product.Price, item.Quantity)

		if product.Quantity < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product", "product_id": product.ProductID})
			return
		}

		itemPrice := float64(item.Quantity) * product.Price
		totalAmount += itemPrice      
		totalQuantity += item.Quantity

		fmt.Printf("Item Price for Product ID %d: %.2f (Quantity: %d)\n", productID, itemPrice, item.Quantity)

		orderItem := models.OrderItem{
			ProductID: productID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
		}
		orderItems = append(orderItems, orderItem)

		product.Quantity -= item.Quantity
		if err := db.Db.Save(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}
	}

	fmt.Printf("Total Amount Calculated: %.2f\n", totalAmount)
	fmt.Printf("Total Quantity Calculated: %d\n", totalQuantity)

	order := models.Order{
		UserID:        int(userID),
		Total:         totalAmount,
		OrderDate:     time.Now(),
		Status:        "Pending",
		Method:        "COD",
		PaymentStatus: "Processing",
	}

	if err := db.Db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not place order", "details": err.Error()})
		return
	}

	for _, item := range orderItems {
		item.OrderID = order.OrderID 
		if err := db.Db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save order items", "details": err.Error()})
			return
		}
	}

	if err := db.Db.Where("user_id=?", userID).Delete(&models.Cart{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear the cart"})
		return
	}

	orderResponse := responsemodels.OrderResponse{
		UserID:        int(userID),
		OrderID:       order.OrderID,
		Total:         totalAmount,
		Quantity:      totalQuantity,
		Status:        order.Status,
		Method:        order.Method,
		PaymentStatus: order.PaymentStatus,
		OrderDate:     order.OrderDate,
	}

	fmt.Printf("Order Response: %+v\n", orderResponse)

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": orderResponse})
}
