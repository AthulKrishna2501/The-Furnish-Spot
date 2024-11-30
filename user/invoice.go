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

	err := pdf.AddTTFFont("DejaVuSans", "./fonts/arial.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to add font: %w", err)
	}
	err = pdf.SetFont("DejaVuSans", "", 12)
	if err != nil {
		return nil, fmt.Errorf("cannot set font: %w", err)
	}

	leftMargin := 20.0
	rightMargin := 550.0

	pdf.SetFont("DejaVuSans", "B", 18)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "The Furniture Spot")
	pdf.Br(25)

	pdf.SetFont("DejaVuSans", "", 12)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Invoice ID: %s", invoice.InvoiceID))
	pdf.Br(15)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Date: %s", invoice.Date.Format("02-Jan-2006")))
	pdf.Br(15)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Customer ID: %d", invoice.UserID))
	pdf.Br(25)

	headers := []string{"Product ID", "Quantity", "Unit Price", "Discount", "Total"}
	columnWidths := []float64{100, 100, 100, 100, 100}
	pdf.SetFont("DejaVuSans", "B", 12)

	xPos := leftMargin
	for i, header := range headers {
		pdf.SetX(xPos)
		pdf.CellWithOption(&gopdf.Rect{W: columnWidths[i], H: 15}, header, gopdf.CellOption{Align: gopdf.Center})
		xPos += columnWidths[i]
	}
	pdf.Br(20)

	pdf.SetFont("DejaVuSans", "", 12)
	for _, item := range invoice.Items {
		itemTotal := (item.UnitPrice * float64(item.Quantity)) - item.Discount

		xPos := leftMargin
		row := []string{
			fmt.Sprintf("%d", item.ProductID),
			fmt.Sprintf("%d", item.Quantity),
			fmt.Sprintf("%.2f", item.UnitPrice),
			fmt.Sprintf("%.2f", item.Discount),
			fmt.Sprintf("%.2f", itemTotal),
		}

		for i, cell := range row {
			pdf.SetX(xPos)
			pdf.CellWithOption(&gopdf.Rect{W: columnWidths[i], H: 15}, cell, gopdf.CellOption{Align: gopdf.Center})
			xPos += columnWidths[i]
		}
		pdf.Br(15)
	}

	tableEndY := pdf.GetY()
	pdf.Line(leftMargin, tableEndY, rightMargin, tableEndY)
	pdf.Br(25)

	pdf.SetFont("DejaVuSans", "", 12)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Subtotal: %.2f", invoice.Subtotal))
	pdf.Br(15)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Discount Applied: %.2f", invoice.Discount))
	pdf.Br(15)
	pdf.SetFont("DejaVuSans", "B", 14)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, fmt.Sprintf("Total: %.2f", invoice.Total))
	pdf.Br(25)

	pdf.SetFont("DejaVuSans", "", 10)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Thank you for shopping with The Furniture Spot!")
	pdf.Br(15)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "For support, contact us at support@thefurniturespot.com.")

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
