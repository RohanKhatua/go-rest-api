package database

import (
	"fiber_gorm_rest/models"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

//handles the database connection
func ConnectDb () {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	//the name of the database is "db"

	if err!=nil {
		log.Fatal("Failed to connect to DB")
		os.Exit(2)
	}

	log.Println("Connected to DB successfully")

	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migration")

	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})
	//sort of like npx prisma db push


	Database = DbInstance{Db:db}
}