package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id        int    `gorm:"AUTO_INCREMENT" form:"id"  json:"id"`
	FirstName string `gorm:"not null" form:"firstname"  json:"firstname"`
	LastName  string `gorm:"not null" form:"lastname"  json:"lastname"`
}

func IniitDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./data.db")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}
	// create tables
	if !db.HasTable(&User{}) {
		db.CreateTable(&User{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&User{})
	}
	return db
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func PostUser(c *gin.Context) {
	db := IniitDb()
	defer db.Close()
	var user User
	c.Bind(&user)
	if user.FirstName != "" && user.LastName != "" {
		//  save the  user
		db.Create(&user)
		c.JSON(201, gin.H{"Success": user})
	} else {
		c.JSON(422, gin.H{"error": "Fields are empty"})
	}
}

func GetUsers(c *gin.Context) {
	// var users = []User{
	// 	{Id: 1, FirstName: "Sonu", LastName: "kumar"},
	// 	{Id: 1, FirstName: "Suresh", LastName: "kumar"},
	// }
	db := IniitDb()
	defer db.Close()
	var users []User
	db.Find(&users)
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	user_id := c.Params.ByName("id")
	db := IniitDb()
	defer db.Close()
	var user User
	db.Find(&user, user_id)
	if user.Id != 0 {
		c.JSON(200, user)
	} else {
		content := gin.H{"error": "user with  id#" + user_id + "  not found"}
		c.JSON(404, content)
	}

	// if user_id == 1 {
	// 	content := gin.H{
	// 		"id":        user_id,
	// 		"firstname": "Sonu",
	// 		"lastname":  "Kumar",
	// 	}
	// 	c.JSON(200, content)
	// } else if user_id == 2 {
	// 	content := gin.H{
	// 		"id":        user_id,
	// 		"firstname": "Suresh",
	// 		"lastname":  "Kumar",
	// 	}
	// 	c.JSON(200, content)
	// } else {
	// 	content := gin.H{"error": "user with  id#" + id + "  not found"}
	// 	c.JSON(404, content)
	// }
}
func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user_id, _ := strconv.ParseInt(id, 0, 64)
	db := IniitDb()
	defer db.Close()
	var user User
	db.First(&user, user_id)
	if user.FirstName != "" && user.LastName != "" {
		if user.Id != 0 {
			var newUser User
			c.Bind(&newUser)
			result := User{
				Id:        user.Id,
				FirstName: newUser.FirstName,
				LastName:  newUser.LastName,
			}
			db.Save(&result)
			c.JSON(200, gin.H{"success": result})
		} else {
			c.JSON(404, gin.H{"error": "User not found"})
		}
	} else {
		c.JSON(422, gin.H{"error": "Fields  are empty"})
	}

}
func DeleteUser(c *gin.Context) {
	db := IniitDb()
	defer db.Close()
	user_id := c.Params.ByName("id")
	// user_id, _ := strconv.ParseInt(id, 0, 64)
	var user User
	db.First(&user, user_id)
	if user.Id != 0 {
		db.Delete(&user)
		c.JSON(200, gin.H{"success": "User #" + user_id + " deleted"})
	} else {
		c.JSON(404, gin.H{"error": "Usernot found"})
	}
}
func main() {
	r := gin.Default()
	r.Use(Cors())
	v1 := r.Group("api/v1")
	v1.POST("/users", PostUser)
	v1.GET("/users", GetUsers)
	v1.GET("/users/:id", GetUser)
	v1.PUT("/users/:id", UpdateUser)
	v1.DELETE("/users/:id", DeleteUser)
	r.Run(":8080")
}
