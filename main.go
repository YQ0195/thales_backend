package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Connect to Postgres
	db, err := initDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// 2. Start Gin
	r := gin.Default()

	r.Static("/images", "./public/images")

	// GET /products with optional query: ?search=term&sortBy=price&sortOrder=asc&page=1&limit=6
	r.GET("/products", func(c *gin.Context) {
		search := c.Query("search")
		productType := c.Query("type") // new filter
		sortBy := c.DefaultQuery("sortBy", "id")
		sortOrder := strings.ToLower(c.DefaultQuery("sortOrder", "asc"))
		page := c.DefaultQuery("page", "1")
		limit := c.DefaultQuery("limit", "10")

		pg, _ := strconv.Atoi(page)
		lim, _ := strconv.Atoi(limit)
		offset := (pg - 1) * lim

		query := `SELECT id, name, type, price, description, picture_url FROM products WHERE 1=1`
		args := []interface{}{}
		paramIdx := 1

		if search != "" {
			query += " AND LOWER(name) LIKE LOWER($" + strconv.Itoa(paramIdx) + ")"
			args = append(args, "%"+search+"%")
			paramIdx++
		}

		if productType != "" {
			query += " AND LOWER(type) = LOWER($" + strconv.Itoa(paramIdx) + ")"
			args = append(args, productType)
			paramIdx++
		}

		// Sort logic
		switch sortBy {
		case "price":
			query += " ORDER BY price"
		case "name":
			query += " ORDER BY name"
		default:
			query += " ORDER BY id"
		}

		if sortOrder == "desc" {
			query += " DESC"
		} else {
			query += " ASC"
		}

		query += " LIMIT $" + strconv.Itoa(paramIdx) + " OFFSET $" + strconv.Itoa(paramIdx+1)
		args = append(args, lim, offset)

		var products []Product
		err := db.Select(&products, query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	// POST /products (multipart form: fields + image)
	r.POST("/products", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
			return
		}

		// Save image to disk
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		path := filepath.Join("public/images", filename)
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL := "/images/" + filename

		name := c.PostForm("name")
		typeVal := c.PostForm("type")
		priceStr := c.PostForm("price")
		description := c.PostForm("description")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		query := `INSERT INTO products (name, type, price, description, picture_url) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		var id int
		err = db.QueryRow(query, name, typeVal, price, description, imageURL).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id, "picture_url": imageURL})
	})

	// PUT /products/:id
	r.PUT("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var updatedProd Product
		if err := c.ShouldBindJSON(&updatedProd); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		query := `UPDATE products SET name=$1, type=$2, price=$3, description=$4, picture_url=$5 WHERE id=$6`
		_, err := db.Exec(query, updatedProd.Name, updatedProd.Type, updatedProd.Price, updatedProd.Description, updatedProd.PictureURL, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
	})

	// GET /products/:id
	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var prod Product
		err := db.Get(&prod, `SELECT id, name, type, price, description, picture_url FROM products WHERE id=$1`, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, prod)
	})

	// DELETE /products/:id
	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec(`DELETE FROM products WHERE id=$1`, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "deleted"})
	})
	// PUT /products/:id/with-image - update with new image
	r.PUT("/products/:id/with-image", func(c *gin.Context) {
		id := c.Param("id")
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
			return
		}

		// Save new image
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		path := filepath.Join("public/images", filename)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imageURL := "/images/" + filename

		// Collect form fields
		name := c.PostForm("name")
		typeVal := c.PostForm("type")
		priceStr := c.PostForm("price")
		description := c.PostForm("description")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		// Update DB
		query := `UPDATE products SET name=$1, type=$2, price=$3, description=$4, picture_url=$5 WHERE id=$6`
		_, err = db.Exec(query, name, typeVal, price, description, imageURL, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated with image"})
	})

	// PUT /products/:id/no-image - update without changing image
	r.PUT("/products/:id/no-image", func(c *gin.Context) {
		id := c.Param("id")

		name := c.PostForm("name")
		typeVal := c.PostForm("type")
		priceStr := c.PostForm("price")
		description := c.PostForm("description")

		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
			return
		}

		query := `UPDATE products SET name=$1, type=$2, price=$3, description=$4 WHERE id=$5`
		_, err = db.Exec(query, name, typeVal, price, description, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Product updated without image"})
	})

	// 4. Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// 5. Run server
	r.Run(":8080")
}
