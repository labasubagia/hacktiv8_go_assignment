package controllers

import (
	"errors"
	"final_project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type commentCtrl struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *commentCtrl {
	return &commentCtrl{db: db}
}

func (ctrl *commentCtrl) Create(c *gin.Context) {

	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid").Error())
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid").Error())
		return
	}

	cm := models.Comment{}
	if err := c.ShouldBindJSON(&cm); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cm.UserID = uint(userID)

	if err := ctrl.db.Create(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, comment{
		ID:        &cm.ID,
		Message:   &cm.Message,
		PhotoID:   &cm.PhotoID,
		UserID:    &cm.UserID,
		CreatedAt: &cm.CreatedAt,
	})
}

func (ctrl *commentCtrl) List(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid").Error())
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid").Error())
		return
	}

	cms := []models.Comment{}
	if err := ctrl.db.Preload("User").Preload("Photo").Where("user_id=?", userID).Find(&cms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := []comment{}
	for _, cm := range cms {
		result = append(result, comment{
			ID:        &cm.ID,
			UserID:    &cm.UserID,
			CreatedAt: &cm.CreatedAt,
			UpdatedAt: &cm.UpdatedAt,
			Message:   &cm.Message,
			PhotoID:   &cm.PhotoID,
			User: &user{
				ID:       &cm.User.ID,
				Username: &cm.User.Username,
				Email:    &cm.User.Email,
			},
			Photo: &photo{
				ID:      &cm.Photo.ID,
				Title:   &cm.Photo.Title,
				Caption: &cm.Photo.Caption,
				URL:     &cm.Photo.URL,
				UserID:  &cm.Photo.UserID,
			},
		})
	}

	c.JSON(http.StatusCreated, result)
}

func (ctrl *commentCtrl) Update(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid").Error())
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid").Error())
		return
	}

	IDStr := c.Param("commentId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("comment id invalid").Error())
		return
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	cm := models.Comment{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&cm).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	cm.Message = payload.Message

	if err := ctrl.db.Save(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, comment{
		ID:        &cm.ID,
		UserID:    &cm.UserID,
		Message:   &cm.Message,
		UpdatedAt: &cm.UpdatedAt,
	})
}

func (ctrl *commentCtrl) Delete(c *gin.Context) {
	userData, ok := c.MustGet("userData").(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user data invalid").Error())
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, errors.New("user id invalid").Error())
		return
	}

	IDStr := c.Param("commentId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("comment id invalid").Error())
		return
	}

	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).Delete(&models.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
