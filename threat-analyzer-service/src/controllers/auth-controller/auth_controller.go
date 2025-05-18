package authcontroller

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	mysqlconfig "github.com/yatender-pareek/threat-analyzer-service/src/config/my-sql-config"
	authdto "github.com/yatender-pareek/threat-analyzer-service/src/dto/auth-dto"
	userentity "github.com/yatender-pareek/threat-analyzer-service/src/models/user-model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	validate = validator.New()
)

// Register handles user registration
// @Summary Register a new user
// @Description Creates a new user account with username, password, and email
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body authdto.RegisterRequestDTO true "User registration details"
// @Success 200 {object} authdto.RegisterResponseDTO "Registration success"
// @Failure 400 {object} map[string]string "error: Invalid input"
// @Failure 409 {object} map[string]string "error: Username or email already exists"
// @Failure 500 {object} map[string]string "error: Server error"
// @Router /api/register [post]
func Register(c *gin.Context) {
	var registerDTO authdto.RegisterRequestDTO
	if err := c.ShouldBindJSON(&registerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := validate.Struct(&registerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	user := userentity.User{
		Username: registerDTO.Username,
		Password: string(hashedPassword),
		Email:    registerDTO.Email,
	}

	if err := mysqlconfig.GetDB().Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") || err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	response := authdto.RegisterResponseDTO{
		Message: "User registered successfully",
	}
	c.JSON(http.StatusOK, response)
}

// Login handles user login and JWT generation
// @Summary User login
// @Description Authenticates a user and returns a JWT token using query parameters
// @Tags Auth
// @Accept json
// @Produce json
// @Param username query string true "User username"
// @Param password query string true "User password"
// @Success 200 {object} authdto.LoginResponseDTO "Login success with JWT token"
// @Failure 400 {object} map[string]string "error: Invalid input"
// @Failure 401 {object} map[string]string "error: Invalid credentials"
// @Failure 500 {object} map[string]string "error: Server error"
// @Router /api/login [get]
func Login(c *gin.Context) {
	var loginDTO authdto.LoginRequestDTO
	if err := c.ShouldBindQuery(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := validate.Struct(&loginDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	var storedUser userentity.User
	if err := mysqlconfig.GetDB().Where("username = ?", loginDTO.Username).First(&storedUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(loginDTO.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginDTO.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(JWTSecretKey())
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error" + err.Error()})
		return
	}

	log.Printf("Generated token for user %s: %s", loginDTO.Username, tokenString)

	response := authdto.LoginResponseDTO{
		Token: tokenString,
	}
	c.JSON(http.StatusOK, response)
}
