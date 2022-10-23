package controllers

import (
	"final_project/dto"
	"final_project/helpers"
	"final_project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userCtrl struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *userCtrl {
	return &userCtrl{
		db: db,
	}
}

func (ctrl *userCtrl) Register(c *gin.Context) {
	u := models.User{}
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "Invalid payload"})
		return
	}

	var isRegistered bool
	err := ctrl.db.Model(u).Select("count(*) > 0").Where("email=? OR username=?", u.Email, u.Username).Find(&isRegistered).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "Cannot find user"})
		return
	}

	validation := helpers.Validate(u)
	if len(validation) > 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:           "Validation Error",
			ValidationError: validation,
		})
		return
	}

	if isRegistered {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "User already registered"})
		return
	}

	if err := ctrl.db.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Error: "Cannot register user",
		})
		return
	}
	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.User{
			Age:      &u.Age,
			Email:    &u.Email,
			ID:       &u.ID,
			Username: &u.Username,
		},
		Message: "Register success",
	})
}

func (ctrl *userCtrl) Login(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "Invalid payload"})
		return
	}
	user := models.User{}
	if err := ctrl.db.First(&user, "email=?", payload.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "Failed find user"})
		return
	}

	if err := helpers.ComparePass(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "Wrong email or password"})
		return
	}

	token := helpers.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, dto.Response{
		Data:    gin.H{"token": token},
		Message: "Login success",
	})
}

func (ctrl *userCtrl) Update(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	var payload struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid payload"})
		return
	}

	u := models.User{}
	if err := ctrl.db.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Error: "failed update user"})
		return
	}

	u.Email = payload.Email
	u.Username = payload.Username

	if err := ctrl.db.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data: dto.User{
			Age:       &u.Age,
			Email:     &u.Email,
			ID:        &u.ID,
			Username:  &u.Username,
			UpdatedAt: &u.UpdatedAt,
		},
	})
}

func (ctrl *userCtrl) Delete(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	if err := ctrl.db.Unscoped().Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Data: gin.H{
			"message": "Your account has been successfully deleted",
		},
	})
}
