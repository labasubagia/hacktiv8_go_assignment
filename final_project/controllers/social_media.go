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

type socialMediaCtrl struct {
	db *gorm.DB
}

func NewSocialMediaController(db *gorm.DB) *socialMediaCtrl {
	return &socialMediaCtrl{db: db}
}

func (ctrl *socialMediaCtrl) Create(c *gin.Context) {
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

	sm := models.SocialMedia{}
	if err := c.ShouldBindJSON(&sm); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	sm.UserID = uint(userID)

	if err := ctrl.db.Create(&sm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, socialMedia{
		ID:        &sm.ID,
		Name:      &sm.Name,
		URL:       &sm.URL,
		UserID:    &sm.UserID,
		CreatedAt: &sm.CreatedAt,
	})
}

func (ctrl *socialMediaCtrl) List(c *gin.Context) {
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

	socialMedias := []models.SocialMedia{}
	if err := ctrl.db.Preload("User").Where("user_id=?", userID).Find(&socialMedias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := []socialMedia{}
	for _, sm := range socialMedias {
		result = append(result, socialMedia{
			ID:        &sm.ID,
			Name:      &sm.Name,
			URL:       &sm.URL,
			UserID:    &sm.UserID,
			CreatedAt: &sm.CreatedAt,
			UpdatedAt: &sm.UpdatedAt,
			User: &user{
				ID:       &sm.User.ID,
				Username: &sm.User.Username,
			},
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"social_medias": result,
	})
}

func (ctrl *socialMediaCtrl) Update(c *gin.Context) {
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

	IDStr := c.Param("socialMediaId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("social media id invalid"))
		return
	}

	var payload struct {
		Name string `json:"name"`
		URL  string `json:"social_media_url"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	sm := models.SocialMedia{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&sm).Error; err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	sm.Name = payload.Name
	sm.URL = payload.URL

	if err := ctrl.db.Save(&sm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, socialMedia{
		ID:        &sm.ID,
		Name:      &sm.Name,
		URL:       &sm.URL,
		UserID:    &sm.UserID,
		UpdatedAt: &sm.UpdatedAt,
	})
}

func (ctrl *socialMediaCtrl) Delete(c *gin.Context) {
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

	IDStr := c.Param("socialMediaId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("social media id invalid"))
		return
	}

	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).Delete(&models.SocialMedia{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
