package api

import (
	//"go/token"
	"fmt"
	"main/db"
	"main/interceptor"
	"main/model"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
	//"path/filepath"
	"strconv"


	"github.com/gin-gonic/gin"
)


func SetupProductAPI(router *gin.Engine) {
	productAPI := router.Group("/api/v2")
	{
		productAPI.GET("/product",interceptor.JwtVerify, getProduct)
		productAPI.POST("/product",interceptor.JwtVerify, createProduct)
		productAPI.PUT("/product",interceptor.JwtVerify, editProduct)
		productAPI.GET("/product/:id", interceptor.JwtVerify, getProductByID)
		productAPI.DELETE("/product/:id" ,interceptor.JwtVerify, deleteProduct)
	}
}

// get product by search and all product
func getProduct(c *gin.Context) {
	var product []model.Product
	keyword := c.Query("keyword")
	if keyword != ""{
		keyword = fmt.Sprintf("%%%s%%",keyword)
		 db.GetDB().Where("name like ?", keyword).Find(&product)
	}else {
		db.GetDB().Order("created_at DESC").Find(&product)
	}
	c.JSON(http.StatusOK, product)
}

// get product by id
func getProductByID(c *gin.Context) {
	var product model.Product
	db.GetDB().Where("id = ?", c.Param("id")).First(&product)
	c.JSON(200, product)
}

//create product insert product to database 
func createProduct(c *gin.Context) {
	product := model.Product{}
	product.Name = c.PostForm("name")
	product.Stock,_ = strconv.ParseInt(c.PostForm("stock"),10,64)
	product.Price,_ = strconv.ParseFloat(c.PostForm("price"),64)
	product.CreatedAt = time.Now()
	db.GetDB().Create(&product)
	image,_ := c.FormFile("image")
    saveImage(image, &product, c)
	
	c.JSON(http.StatusOK, gin.H{"result": product})
}

// save image to database 
func saveImage(image *multipart.FileHeader, product *model.Product, c *gin.Context) {
	if image != nil {
		runningDir, _ := os.Getwd()
		product.Image = image.Filename
		extension := filepath.Ext(image.Filename)
		fileName := fmt.Sprintf("%d%s", product.ID, extension)
		filePath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)

		if fileExists(filePath) {
			os.Remove(filePath)
		}
		c.SaveUploadedFile(image, filePath)
		db.GetDB().Model(&product).Update("image", fileName)
	}
}
//check name in folder uploaded have image is name
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

//edit update Product
func editProduct(c *gin.Context) {
	var product model.Product
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	product.ID = uint(id)
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

	db.GetDB().Updates(&product)
	image, _ := c.FormFile("image")
	saveImage(image, &product, c)
	c.JSON(http.StatusOK, gin.H{"result": product})

}
// delete product
func deleteProduct(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	db.GetDB().Delete(&model.Product{}, id)
	c.JSON(http.StatusOK, gin.H{"result": "ok"})
}
