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

type photoCtrl struct {
	db *gorm.DB
}

func NewPhotoController(db *gorm.DB) *photoCtrl {
	return &photoCtrl{
		db: db,
	}
}

func (ctrl *photoCtrl) Create(c *gin.Context) {
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

	p := models.Photo{}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	p.UserID = uint(userID)

	if err := ctrl.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, photo{
		ID:        &p.ID,
		Title:     &p.Title,
		Caption:   &p.Caption,
		URL:       &p.URL,
		UserID:    &p.UserID,
		CreatedAt: &p.CreatedAt,
	})
}

func (ctrl *photoCtrl) List(c *gin.Context) {
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

	photos := []models.Photo{}
	if err := ctrl.db.Preload("User").Where("user_id=?", userID).Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := []photo{}
	for _, p := range photos {
		result = append(result, photo{
			ID:        &p.ID,
			Title:     &p.Title,
			URL:       &p.URL,
			UserID:    &p.UserID,
			CreatedAt: &p.CreatedAt,
			UpdatedAt: &p.UpdatedAt,
			Caption:   &p.Caption,
			User: &user{
				Username: &p.User.Username,
				Email:    &p.User.Email,
			},
		})
	}

	c.JSON(http.StatusCreated, result)
}

func (ctrl *photoCtrl) Update(c *gin.Context) {
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

	IDStr := c.Param("photoId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("photo id invalid").Error())
		return
	}

	var payload struct {
		Title   string `json:"title"`
		Caption string `json:"caption"`
		URL     string `json:"photo_url"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	p := models.Photo{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	p.Caption = payload.Caption
	p.Title = payload.Title
	p.URL = payload.URL

	if err := ctrl.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, photo{
		ID:        &p.ID,
		Title:     &p.Title,
		URL:       &p.URL,
		UserID:    &p.UserID,
		UpdatedAt: &p.UpdatedAt,
		Caption:   &p.Caption,
	})
}

func (ctrl *photoCtrl) Delete(c *gin.Context) {
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

	IDStr := c.Param("photoId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("photo id invalid").Error())
		return
	}

	if err := ctrl.db.Select("Comment").Where("id=? AND user_id=?", ID, userID).Delete(&models.Photo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
