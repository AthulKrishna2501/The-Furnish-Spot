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
	err = pdf.SetFont("DejaVuSans", "", 12)
	if err != nil {
		return nil, fmt.Errorf("cannot set font: %w", err)
	}

	pdf.Cell(nil, fmt.Sprintf("Invoice ID: %s", invoice.InvoiceID))
	pdf.Br(17)

	pdf.Cell(nil, fmt.Sprintf("Date: %s", invoice.Date.Format("02-Jan-2006")))
	pdf.Br(17)

	pdf.Cell(nil, fmt.Sprintf("Customer ID: %d", invoice.UserID))
	pdf.Br(17)

	pdf.SetFont("DejaVuSans", "", 10)

	pdf.Cell(nil, "ProductID")
	pdf.SetX(150)
	pdf.Cell(nil, "Qty")
	pdf.SetX(180)
	pdf.Cell(nil, "Unit Price")
	pdf.SetX(230)
	pdf.Cell(nil, "Discount")
	pdf.SetX(280)
	pdf.Cell(nil, "Total")
	pdf.Br(15)

	pdf.Line(5, pdf.GetY(), 400, pdf.GetY())
	pdf.Br(10)

	subtotal := 0.0
	for _, item := range invoice.Items {
		fmt.Printf("Item: %+v\n", item)

		itemTotal := (item.UnitPrice * float64(item.Quantity)) - item.Discount
		subtotal += itemTotal

		pdf.Cell(nil, fmt.Sprintf("%d", item.ProductID))
		pdf.SetX(150)
		pdf.Cell(nil, fmt.Sprintf("%d", item.Quantity))

		pdf.SetX(180)
		pdf.Cell(nil, fmt.Sprintf("%.2f", item.UnitPrice))

		pdf.SetX(230)
		pdf.Cell(nil, fmt.Sprintf("%.2f", item.Discount))

		pdf.SetX(280)
		pdf.Cell(nil, fmt.Sprintf("%.2f", itemTotal))

		pdf.Br(10)
	}

	discountApplied := invoice.Subtotal - subtotal
	total := subtotal - discountApplied

	pdf.Br(15)
	pdf.SetFont("DejaVuSans", "", 12)
	pdf.Cell(nil, fmt.Sprintf("Subtotal: %.2f", subtotal))
	pdf.Br(15)
	pdf.Cell(nil, fmt.Sprintf("Discount Applied: %.2f", discountApplied))
	pdf.Br(15)

	pdf.Cell(nil, fmt.Sprintf("Total: %.2f", total))
	pdf.Br(15)

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
	var subtotal float64
	for _, item := range items {
		totalPrice := (item.Price * float64(item.Quantity)) - item.Discount

		invoiceItems = append(invoiceItems, models.InvoiceItem{
			ProductID:  item.ProductID,
			Discount:   float64(item.Discount),
			Quantity:   item.Quantity,
			UnitPrice:  item.Price,
			TotalPrice: totalPrice, 
		})

		subtotal += totalPrice
	}

	invoice.InvoiceID = fmt.Sprintf("INV-%d", order.OrderID)
	invoice.Date = time.Now()
	invoice.UserID = order.UserID
	invoice.Subtotal = subtotal 
	invoice.Discount = order.Discount
	invoice.Total = subtotal - order.Discount
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
