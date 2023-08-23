package routes

import (
	"errors"
	"fiber_gorm_rest/database"
	"fiber_gorm_rest/models"

	"github.com/gofiber/fiber/v2"
)

type ProductSerializer struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) ProductSerializer {
	return ProductSerializer{
		ID: productModel.ID,
		Name: productModel.Name,
		SerialNumber: productModel.SerialNumber,
	}
}

func CreateProduct (c *fiber.Ctx) error {
	var product models.Product

	if err:=c.BodyParser(&product);err!=nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(product);

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

func GetProducts (c *fiber.Ctx) error {
	
	products := []models.Product{}
	database.Database.Db.Find(&products)
	responseProducts := []ProductSerializer{}
	
	for _, product := range products {
		responseProduct := CreateResponseProduct(product)
		responseProducts = append(responseProducts, responseProduct)
	}

	return c.Status(200).JSON(responseProducts)
}

//find product does not take a pointer to fiber.ctx
//find product is a helper function not a route
func FindProduct (id int, product * models.Product) error {
	database.Database.Db.Find(&product, "id=?", id)
	if product.ID==0 {
		return errors.New("product does not exist")
	}

	return nil
}

func GetProduct (c * fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err!=nil {
		return c.Status(400).JSON("ID must be an integer")
	}

	var product models.Product

	if err:=FindProduct(id, &product); err!=nil {
		return c.Status(400).JSON(err.Error())
	}

	responseProuct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProuct)
}

func UpdateProduct (c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err!=nil {
		return c.Status(400).JSON("ID should be an integer")
	}

	if err:=FindProduct(id, &product); err!=nil {
		return c.Status(400).JSON("Product does not exist")
	}

	type updateProduct struct {
		Name string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}
	
	var updateData updateProduct

	if err:=c.BodyParser(&updateData); err!=nil {
		return c.Status(500).JSON(err.Error())
	}

	if updateData.Name!="" {
		product.Name=updateData.Name
	}

	if updateData.SerialNumber!=""{
		product.SerialNumber = updateData.SerialNumber
	}

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

func DeleteProduct (c *fiber.Ctx) error {
	var product models.Product

	id, err := c.ParamsInt("id")

	if err!=nil {
		return c.Status(400).JSON("ID must be an integer")
	}

	if err:=FindProduct(id, &product); err!=nil {
		return c.Status(400).JSON("Product not found")
	}

	if err:=database.Database.Db.Delete(&product).Error; err!=nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully deleted product")

	
}
