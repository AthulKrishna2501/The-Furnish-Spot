package salesreport

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	db "github.com/AthulKrishna2501/The-Furniture-Spot/DB"
	"github.com/AthulKrishna2501/The-Furniture-Spot/models"
	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	log "github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
)

func GenerateReport(c *gin.Context) {
	var report models.SalesReport
	var productSales []models.ProductDetails

	filter := c.Query("filter")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	format := c.Query("format")
	query := db.Db.Model(&models.Order{})

	if filter == "daily" {
		query = query.Where("date(order_date) = ?", time.Now().Format("2006-01-02"))
	} else if filter == "weekly" {
		query = query.Where("order_date >= ? AND order_date <= ?", time.Now().AddDate(0, 0, -7), time.Now())
	} else if filter == "monthly" {
		query = query.Where("order_date >= ? AND order_date <= ?", time.Now().AddDate(0, -1, 0), time.Now())
	} else if startDate != "" && endDate != "" {
		query = query.Where("order_date BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Select(`
		COUNT(*) as total_sales_count,
		SUM(total) as total_order_amount,
		SUM(discount) as total_discount,
		SUM(CASE WHEN coupon_id IS NOT NULL THEN discount ELSE 0 END) as coupons_deduction
	`).Scan(&report).Error; err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("error in database query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Db.Table("order_items AS oi").
		Select(`
        p.product_id,
        p.product_name,
		  oi.quantity,
        SUM(oi.quantity) AS total_quantity_sold,
        SUM(oi.price * oi.quantity) AS total_price
    `).
		Joins("JOIN products AS p ON p.product_id = oi.product_id").
		Group("p.product_id, p.product_name,oi.quantity").
		Scan(&productSales).Error; err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error in database query")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	fmt.Println("Product Sales Data:", productSales)

	report.ProductSales = productSales
	if format == "excel" {
		filePath, err := GenerateExcelReport(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate Excel report"})
			return
		}
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=sales_report.xlsx")
		c.File(filePath)
	} else {
		filePath, err := GeneratePDFReport(report)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate PDF report"})
			return
		}
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=sales_report.pdf")
		c.File(filePath)

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=sales_report.xlsx")
		c.File(filePath)
	}
}
func GeneratePDFReport(report models.SalesReport) (string, error) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err := pdf.AddTTFFont("LiberationSans", "/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf")
	if err != nil {
		return "", err
	}

	pdf.SetFont("LiberationSans", "", 14)
	pdf.AddPage()

	pdf.Cell(nil, "Sales Report")
	pdf.Br(30)
	pdf.Cell(nil, "Total Sales Count: "+fmt.Sprintf("%d", report.TotalSalesCount))
	pdf.Br(15)
	pdf.Cell(nil, "Total Order Amount: "+fmt.Sprintf("%.2f", report.TotalOrderAmount))
	pdf.Br(15)
	pdf.Cell(nil, "Total Discount: "+fmt.Sprintf("%.2f", report.TotalDiscount))
	pdf.Br(15)
	pdf.Cell(nil, "Coupons Deduction: "+fmt.Sprintf("%.2f", report.CouponsDeduction))
	pdf.Br(30)

	pdf.SetFont("LiberationSans", "", 12)
	pdf.Cell(nil, "Product Details")
	pdf.Br(20)

	pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, "ProductID", gopdf.CellOption{Align: gopdf.Middle})
	pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, "Quantity", gopdf.CellOption{Align: gopdf.Middle})
	pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, "Total Price", gopdf.CellOption{Align: gopdf.Middle})
	pdf.Br(20)

	for _, product := range report.ProductSales {
		pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, fmt.Sprintf("%d", product.ProductID), gopdf.CellOption{Align: gopdf.Middle})
		pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, fmt.Sprintf("%d", product.Quantity), gopdf.CellOption{Align: gopdf.Middle})
		pdf.CellWithOption(&gopdf.Rect{W: 60, H: 20}, fmt.Sprintf("%.2f", product.TotalPrice), gopdf.CellOption{Align: gopdf.Middle})
		pdf.Br(20)
	}

	err = pdf.WritePdf("sales_report.pdf")
	if err != nil {
		return "", err
	}
	return "sales_report.pdf", nil
}

func GenerateExcelReport(report models.SalesReport) (string, error) {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "Sales Report")
	f.SetCellValue("Sheet1", "A2", "Total Sales Count")
	f.SetCellValue("Sheet1", "B2", report.TotalSalesCount)
	f.SetCellValue("Sheet1", "A3", "Total Order Amount")
	f.SetCellValue("Sheet1", "B3", report.TotalOrderAmount)
	f.SetCellValue("Sheet1", "A4", "Total Discount")
	f.SetCellValue("Sheet1", "B4", report.TotalDiscount)
	f.SetCellValue("Sheet1", "A5", "Coupons Deduction")
	f.SetCellValue("Sheet1", "B5", report.CouponsDeduction)

	outputPath := "./sales_report.xlsx"
	err := f.SaveAs(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func GetSalesData(c *gin.Context) {
	filter := c.Query("filter")

	var dates []string
	var sales []float64

	switch filter {
	case "yearly":
		rows, err := db.Db.Raw(`
			  SELECT TO_CHAR(order_date, 'YYYY') AS date, SUM(total) AS sales
			  FROM orders
			  GROUP BY TO_CHAR(order_date, 'YYYY')
			  ORDER BY TO_CHAR(order_date, 'YYYY')
		 `).Rows()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error in database query")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var date string
			var sale float64
			if err := rows.Scan(&date, &sale); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error in scanning column")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			dates = append(dates, date)
			sales = append(sales, sale)
		}

	case "monthly":
		rows, err := db.Db.Raw(`
			  SELECT TO_CHAR(order_date, 'YYYY-MM') AS date, SUM(total) AS sales
			  FROM orders
			  GROUP BY TO_CHAR(order_date, 'YYYY-MM')
			  ORDER BY TO_CHAR(order_date, 'YYYY-MM')
		 `).Rows()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error in scanning column")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var date string
			var sale float64
			if err := rows.Scan(&date, &sale); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("error in scanning column")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			dates = append(dates, date)
			sales = append(sales, sale)
		}

	case "weekly":
		rows, err := db.Db.Raw(`
			  SELECT TO_CHAR(order_date, 'IYYY-IW') AS date, SUM(total) AS sales
			  FROM orders
			  GROUP BY TO_CHAR(order_date, 'IYYY-IW')
			  ORDER BY TO_CHAR(order_date, 'IYYY-IW')
		 `).Rows()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("error in scannning column")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var date string
			var sale float64
			if err := rows.Scan(&date, &sale); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error in scanning rows")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			dates = append(dates, date)
			sales = append(sales, sale)
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dates": dates,
		"sales": sales,
	})
}

func GetTopSellingProducts(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitParam)

	var products []map[string]interface{}

	if err := db.Db.Table("order_items").
		Select("product_id, SUM(quantity) as total_sold").
		Group("product_id").
		Order("total_sold DESC").
		Limit(limit).
		Find(&products); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error in querying order_items")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, products)
}

func GetTopSellingCategories(c *gin.Context) {
	limitParam := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitParam)

	var categories []map[string]interface{}

	if err := db.Db.Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Select("products.category_id, SUM(order_items.quantity) as total_sold").
		Group("products.category_id").
		Order("total_sold DESC").
		Limit(limit).
		Find(&categories); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error is querying order_items")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetLedgerBook(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var ledgerEntries []map[string]interface{}

	if err := db.Db.Table("orders").
		Select("order_date as date, 'Sale' as type, total as amount, order_id").
		Where("order_date BETWEEN ? AND ?", startDate, endDate).
		Order("order_date DESC").
		Find(&ledgerEntries); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("error in querying orders")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return

	}

	c.JSON(http.StatusOK, ledgerEntries)
}
