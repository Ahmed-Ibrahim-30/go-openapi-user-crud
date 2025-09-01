package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Person struct {
	ID       uint `gorm:"primarykey"`
	Name     string
	Email    string `gorm:"unique"`
	Age      int
	Password string
}

func DBGormAccess(dsn string) {
	db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Turn Off Gorm Self Naming
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})

	err := db.AutoMigrate(&Person{})
	if err != nil {
		return
	}

	res := db.Create(&Person{Name: "Salem", Email: "b@gmail.com", Age: 26, Password: "2389"})
	if res.Error != nil {
		fmt.Println("error:", res.Error)
	}
	db.Delete(&Person{ID: 4})

	db.Model(&Person{}).Where("id = ?", 9).Updates(Person{
		Name: "Mohamed",
	})

	var people []Person
	result := db.Find(&people)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	for _, person := range people {
		fmt.Println("ID: ", person.ID, "Name: ", person.Name, "Email: ", person.Email, "Age: ", person.Age, "Password: ", person.Password)
	}
}
