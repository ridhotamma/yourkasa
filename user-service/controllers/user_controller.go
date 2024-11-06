package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ridhotamma/yourkasa/user-service/dto"
	"github.com/ridhotamma/yourkasa/user-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (c *UserController) Create(ctx *gin.Context) {
	var input dto.CreateUserDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		FirstName:         input.FirstName,
		LastName:          input.LastName,
		Email:             input.Email,
		PasswordHash:      string(hashedPassword),
		ProfilePictureUrl: input.ProfilePictureUrl,
		Role:              models.Role(input.Role),
	}

	if err := c.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "id": user.ID})
}

func (c *UserController) Update(ctx *gin.Context) {
	var input dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	var user models.User
	if err := c.db.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := map[string]interface{}{}
	if input.FirstName != "" {
		updates["first_name"] = input.FirstName
	}
	if input.LastName != "" {
		updates["last_name"] = input.LastName
	}
	if input.ProfilePictureUrl != "" {
		updates["profile_picture_url"] = input.ProfilePictureUrl
	}
	if input.Role != "" {
		updates["role"] = input.Role
	}

	if err := c.db.Model(&user).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UserController) List(ctx *gin.Context) {
	var users []models.User
	if err := c.db.Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var userList []dto.UserListDTO
	for _, user := range users {
		userList = append(userList, dto.UserListDTO{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      string(user.Role),
		})
	}

	ctx.JSON(http.StatusOK, userList)
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var user models.User
	if err := c.db.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userDetail := dto.UserDetailDTO{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		ProfilePictureUrl: user.ProfilePictureUrl,
		Role:              string(user.Role),
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		LastLoggedIn:      user.LastLoggedIn,
	}

	ctx.JSON(http.StatusOK, userDetail)
}

func (c *UserController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.db.Delete(&models.User{}, id).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (c *UserController) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var user models.User
	if err := c.db.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	userDetail := dto.UserDetailDTO{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		ProfilePictureUrl: user.ProfilePictureUrl,
		Role:              string(user.Role),
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		LastLoggedIn:      user.LastLoggedIn,
	}

	ctx.JSON(http.StatusOK, userDetail)
}
