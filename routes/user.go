package routes

import (
	"errors"
	"fiber_gorm_rest/database"
	"fiber_gorm_rest/models"

	"github.com/gofiber/fiber/v2"
)

//all user routes go here

type UserSerializer struct {
	// this is not the model user, this is a serializer
	ID uint `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

func CreateResponseUser (userModel models.User) UserSerializer {
	return UserSerializer{
		ID: userModel.ID,
		FirstName: userModel.FirstName,
		LastName: userModel.LastName,
	}
}

// endpoints take a * fiber.Ctx as input
func CreateUser (c *fiber.Ctx) error {
	var user models.User

	if err:=c.BodyParser(&user); err!=nil {
		return c.Status(400).JSON(err.Error())
		//we also have to look at validation
	}

	database.Database.Db.Create(&user)
	//creates the insert value
	//since user was created from models.User
	//the mapped data from the body is automatically inserted there
	//we don't have to explicitly mention the table

	responseUser := CreateResponseUser(user)
	//the user struct that we will return

	return c.Status(200).JSON(responseUser)
}

func GetUsers (c *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)
	responseUsers := []UserSerializer{}

	for _, user:=range(users) {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

//findUser is a helper function not a route
func findUser (id int, user *models.User) error {
	database.Database.Db.Find(&user, "id=?", id)
	//store the result of the query in user

	if user.ID==0 {
		return errors.New("user does not exist")
	}

	return nil
}

func GetUser (c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err!=nil {
		return c.Status(400).JSON("User ID should be an integer")
	}

	if err:=findUser(id, &user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseUser:= CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func UpdateUser (c * fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err!=nil {
		return c.Status(400).JSON("ID should be an integer")
	}

	if err:=findUser(id, &user); err!=nil {
		return c.Status(400).JSON("User does not exist")
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	var updateData UpdateUser

	if err:=c.BodyParser(&updateData); err!=nil {
		return c.Status(500).JSON(err.Error())
	}

	if updateData.FirstName!="" {
		user.FirstName = updateData.FirstName
	}

	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}

	// user.FirstName = updateData.FirstName
	// user.LastName = updateData.LastName

	database.Database.Db.Save(&user)

	responseUser := CreateResponseUser(user)

	return c.Status(200).JSON(responseUser)
}

func DeleteUser (c * fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var user models.User

	if err!=nil {
		return c.Status(400).JSON("ID should be an integer")
	}

	if err:=findUser(id, &user); err!=nil {
		return c.Status(400).JSON("User does not exist")
	}

	if err:= database.Database.Db.Delete(&user).Error; err!=nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully deleted user")
}