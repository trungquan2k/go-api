package users

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/quanht2k/golang_basic_training/app/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	secret = []byte("123123123")
)

type UserApi struct {
	DB *gorm.DB
}

func InitUserApi(db *gorm.DB) *UserApi {
	return &UserApi{DB: db}
}
type Claims struct {
	PhoneNumber string `json:"phonenumber"`
	jwt.StandardClaims
}

func generateToken(phoneNumber string) (string, error) {
	claims := Claims{
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func authenticateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.PhoneNumber, nil
	}
	return "", fmt.Errorf("invalid token")
}

func (u UserApi) GetListUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Authorization token missing"})
		}
		token := authHeader[len("Bearer "):]
		_, err := authenticateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
		}

		var users []*models.User
		if err := u.DB.Table(models.User{}.TableName()).
			Scan(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error fetching users"})
		}

		return c.JSON(fiber.Map{"status": fiber.StatusOK, "data": users, "message": "User list retrieved successfully"})
	}
}

func(u UserApi) SignIn() fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := models.SignInInput{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}


		token, err := generateToken(payload.PhoneNumber)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error generating token"})
		}
		if payload.Password == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid password"})
		}

		response := fiber.Map{
			"status":  fiber.StatusOK,
			"data":    nil, // Omit user data for security reasons
			"message": "Login successful",
			"token":   token,
		}

		return c.JSON(response)
	}

}

func (u UserApi) SignUp() fiber.Handler {
	return 	func(ctx *fiber.Ctx) error{
		payload := models.SignUpInput{}
		// Parse the request body into the User object
		if err := ctx.BodyParser(&payload); err != nil {
			return ctx.JSON(&fiber.Map{"status": http.StatusBadRequest, "error": err.Error()})
		}
		var existingUser *models.User
		if err := db.First(&existingUser, "phonenumber = ?", payload.PhoneNumber).Error; err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
		}

		if existingUser.Password != payload.Password {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid password"})
		}

		var user *models.User
		result := u.DB.Table(models.User{}.TableName()).Where("phonenumber = ?", payload.PhoneNumber).First(&user)
		
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		user = &models.User{
			PhoneNumber: payload.PhoneNumber,
			Password: string(hashedPassword),
		}
		result = u.DB.Table(models.User{}.TableName()).Create(&user)

		if result.Error != nil {
			return ctx.JSON(&fiber.Map{"status": http.StatusBadRequest, "error": result.Error})
		}
		// Return the created user
		return ctx.JSON(&fiber.Map{"status": http.StatusOK, "data": models.FilterUserRecord(user), "message": "Tạo mới user thành công"})
	}

}

// 
func (u UserApi) GetUserDetail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		var user *models.User
		result := u.DB.Table(models.User{}.TableName()).Where("id = ?", id).First(&user)
		if result.Error != nil {
			return ctx.JSON(&fiber.Map{"status": http.StatusNotFound, "error": "User not found"})
		}

		return ctx.JSON(&fiber.Map{"status": http.StatusOK, "data": models.FilterUserRecord(user), "message": "Lấy dữ liệu thành công"})
	}
}

func (u UserApi) UpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		id := c.Params("id")
		// var user *models.User
		payload := models.SignUpInput{}

		result := u.DB.Table(models.User{}.TableName()).Where("id = ?", id).First(&payload)
		if result.Error != nil {
			return c.JSON(&fiber.Map{"status": http.StatusNotFound, "error": "User not found"})
		}
		// Param data to request
		newUser := new(models.User)
		if err := c.BodyParser(newUser); err != nil {
			return c.JSON(&fiber.Map{"status": http.StatusBadRequest, "error": err.Error()})
		}
		// Update the user's fields
		payload.PhoneNumber = newUser.PhoneNumber
		payload.Password = newUser.Password

		// Save the updated user
		if err := u.DB.Save(&payload).Error; err != nil {
			return c.JSON(&fiber.Map{"status": http.StatusInternalServerError, "error": err.Error()})
		}
		user := &models.User{
			UserName: payload.PhoneNumber,
			Password: payload.Password,
		}
		// Return a success message
		return c.JSON(&fiber.Map{"status": http.StatusOK, "data": models.FilterUserRecord(user), "message": "Cập nhật user thành công"})
	}
}

func (u UserApi)DeleteOneUser() fiber.Handler{
	return func(c *fiber.Ctx) error {

	// Get id for each element
	id := c.Params("id")
	var user *models.User
	result := u.DB.Table(models.User{}.TableName()).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return c.JSON(&fiber.Map{"status": http.StatusNotFound, "error": "User not found"})
	}
	// Delete the user from the database
	if err := u.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(&fiber.Map{"status": http.StatusInternalServerError, "error": err.Error()})
	}
	// Return a success message
	return c.JSON(&fiber.Map{"status": http.StatusOK, "message": "Xoá tài khoản thành công"})
} 
}