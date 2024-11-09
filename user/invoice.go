package user

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
)

func GeneratePDF(invoice models.Invoice) ([]byte, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("DejaVuSans", "/home/athul/Documents/The Furnish spot/fonts/arial.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to add font: %w", err)
	}
	err = pdf.SetFont("DejaVuSans", "", 14)
	if err != nil {
		return nil, fmt.Errorf("cannot set font:%w", err)
	}

	err = pdf.Cell(nil, fmt.Sprintf("Invoice %s", invoice.InvoiceID))
	if err != nil {
		return nil, fmt.Errorf("failed to add InvoiceID: %w", err)
	}
	pdf.Br(20)

	err = pdf.Cell(nil, fmt.Sprintf("Date: %s", invoice.Date.Format("02-Jan-2006")))
	if err != nil {
		return nil, fmt.Errorf("failed to add date: %w", err)
	}
	pdf.Br(20)

	err = pdf.Cell(nil, fmt.Sprintf("Customer ID: %d", invoice.UserID))
	if err != nil {
		return nil, fmt.Errorf("failed to add customer name: %w", err)
	}
	pdf.Br(15)

	err = pdf.Cell(nil, fmt.Sprintf("Billing Address: %s", invoice.BillingAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to add billing address: %w", err)
	}
	pdf.Br(20)

	for _, item := range invoice.Items {
		err = pdf.Cell(nil, fmt.Sprintf("%s - Qty: %d - Price: %.2f - Total: %.2f",
			item.Description, item.Quantity, item.UnitPrice, item.TotalPrice))
		if err != nil {
			return nil, fmt.Errorf("failed to add item: %w", err)
		}
		pdf.Br(15)
	}

	err = pdf.Cell(nil, fmt.Sprintf("Subtotal: %.2f", invoice.Subtotal))
	if err != nil {
		return nil, fmt.Errorf("failed to add subtotal: %w", err)
	}
	pdf.Br(15)

	err = pdf.Cell(nil, fmt.Sprintf("Tax: %.2f", invoice.Tax))
	if err != nil {
		return nil, fmt.Errorf("failed to add tax: %w", err)
	}
	pdf.Br(15)

	err = pdf.Cell(nil, fmt.Sprintf("Discount: %d", invoice.Discount))
	if err != nil {
		return nil, fmt.Errorf("failed to add discount: %w", err)
	}
	pdf.Br(15)

	err = pdf.Cell(nil, fmt.Sprintf("Total: %.2f", invoice.Total))
	if err != nil {
		return nil, fmt.Errorf("failed to add total: %w", err)
	}

	var buffer bytes.Buffer
	_, err = pdf.WriteTo(&buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buffer.Bytes(), nil
}

func GenerateInvoiceHandler(c *gin.Context) {
	var invoice models.Invoice
	var order models.Order
	var items []models.OrderItem
	OrderID := c.Param("id")

	if err := db.Db.Where("order_id = ?", OrderID).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find order"})
		return
	}

	if err := db.Db.Where("order_id = ?", OrderID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot find order items"})
		return
	}

	var invoiceItems []models.InvoiceItem
	for _, item := range items {
		invoiceItems = append(invoiceItems, models.InvoiceItem{
			Quantity:  item.Quantity,
			UnitPrice: item.Price,
		})
	}

	invoice.InvoiceID = fmt.Sprintf("INV-%d", order.OrderID)
	invoice.Date = time.Now()
	invoice.UserID = order.UserID
	invoice.Subtotal = order.Total
	invoice.Discount = order.Discount
	invoice.Total = order.Total
	invoice.Items = invoiceItems

	pdfBytes, err := GeneratePDF(invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate invoice: %v", err)})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", invoice.InvoiceID))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
