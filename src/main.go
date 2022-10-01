package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
type UserResponse struct {
	Email     string `json:"email" gorm:"unique"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
type Response struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"msg"`
}

type SuccessResponse struct {
	Token string `json:"token"`
}

// Database handle
var db *gorm.DB

// Initailize datbase and connect to postgres db
func InitDB() error {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=mytestdb sslmode=disable password=dapper@123")

	if err != nil {
		fmt.Println("Failed to initialize database")
		return err
	}
	fmt.Println("Connected to database!")
	return nil
}

var connString = "user=postgres password=dapper@123 host=localhost port=5432 dbname=mytestdb sslmode=disable"

func main() {
	// Setup connection to our postgresql database
	InitDB()
	defer db.Close()
	// Remove existing table
	RemoveTables(db)
	// create user table in database
	db.AutoMigrate(&User{})

	// router
	router := gin.Default()
	router.POST("/signup", HandleSignup)
	router.POST("/login", HandleLogin)
	router.PUT("/users", HandleUpdateUser)
	router.GET("/users", HandleGetUser)

	router.Run("0.0.0.0:8080")

}

func HandleSignup(c *gin.Context) {
	fmt.Println("handle handleSignup")
	var newUser User
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	if err := db.Create(&newUser).Error; err != nil {
		resp := Response{Code: http.StatusInternalServerError, ErrorMsg: "create failed"}
		c.JSON(http.StatusInternalServerError, &resp)
		return
	}
	resp := SuccessResponse{Token: GetToken()}
	c.IndentedJSON(http.StatusOK, &resp)
}

func HandleLogin(c *gin.Context) {
	fmt.Println("handle Login")
	var newUser User
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newUser); err != nil {
		return
	}
	resp := SuccessResponse{Token: GetToken()}
	c.IndentedJSON(http.StatusCreated, &resp)
}

func HandleGetUser(c *gin.Context) {
	fmt.Println("handle handleGetUser")
	// fetch toke from header
	token := c.Request.Header.Get("x-authentication-token")
	fmt.Println(token)
	var users []User
	var resp []UserResponse

	// read users table from database
	db.Find(&users)
	// Form user response as reponse should not contain password field.
	if len(users) > 0 {
		for _, user := range users {
			var tmp UserResponse
			tmp.Email = user.Email
			tmp.Firstname = user.Firstname
			tmp.Lastname = user.Lastname
			resp = append(resp, tmp)
		}
	}
	c.IndentedJSON(http.StatusOK, &resp)
}

func HandleUpdateUser(c *gin.Context) {
	fmt.Println("handle handleUpdateUser")
	var newUser User
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newUser); err != nil {
		resp := Response{Code: http.StatusBadRequest, ErrorMsg: "Failed to parse request"}
		c.JSON(http.StatusInternalServerError, &resp)
		return
	}
	db.Model(&User{}).Updates(User{Firstname: newUser.Firstname, Lastname: newUser.Lastname})
	fmt.Println("User updated successfully")
	c.IndentedJSON(http.StatusOK, "200 OK")
}

func GetToken() string {
	// This method should call jwt library to get the jwt token
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

func RemoveTables(db *gorm.DB) error {
	db.Exec("DROP TABLE IF EXISTS users")
	return nil
}
