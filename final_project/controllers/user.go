package controllers

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := ctrl.db.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user{
		Age:      &u.Age,
		Email:    &u.Email,
		ID:       &u.ID,
		Username: &u.Username,
	})
}

func (ctrl *userCtrl) Login(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := models.User{}
	if err := ctrl.db.First(&user, "email=?", payload.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.ComparePass(user.Password, payload.Password); err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("invalid username or password").Error())
		return
	}

	token := helpers.GenerateToken(user.ID, user.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (ctrl *userCtrl) Update(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid"))
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid"))
		return
	}

	var payload struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	u := models.User{}
	if err := ctrl.db.First(&u, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	u.Email = payload.Email
	u.Username = payload.Username

	if err := ctrl.db.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user{
		Age:       &u.Age,
		Email:     &u.Email,
		ID:        &u.ID,
		Username:  &u.Username,
		UpdatedAt: &u.UpdatedAt,
	})
}

func (ctrl *userCtrl) Delete(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid"))
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid"))
		return
	}

	if err := ctrl.db.Debug().Delete(&models.User{}, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
