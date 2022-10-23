package controllers

import (
	"final_project/dto"
	"final_project/helpers"
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
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	sm := models.SocialMedia{}
	if err := c.ShouldBindJSON(&sm); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid input"})
		return
	}

	validation := helpers.Validate(sm)
	if len(validation) > 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:           "Validation Error",
			ValidationError: validation,
		})
		return
	}

	sm.UserID = uint(userID)

	if err := ctrl.db.Create(&sm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "cannot create social media"})
		return
	}

	c.JSON(http.StatusCreated, dto.SocialMedia{
		ID:        &sm.ID,
		Name:      &sm.Name,
		URL:       &sm.URL,
		UserID:    &sm.UserID,
		CreatedAt: &sm.CreatedAt,
	})
}

func (ctrl *socialMediaCtrl) List(c *gin.Context) {

	socialMedias := []models.SocialMedia{}
	if err := ctrl.db.Preload("User").Find(&socialMedias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Error: "cannot find user",
		})
		return
	}

	result := []dto.SocialMedia{}
	for _, social := range socialMedias {
		sm := social
		result = append(result, dto.SocialMedia{
			ID:        &sm.ID,
			Name:      &sm.Name,
			URL:       &sm.URL,
			UserID:    &sm.UserID,
			CreatedAt: &sm.CreatedAt,
			UpdatedAt: &sm.UpdatedAt,
			User: &dto.User{
				ID:       &sm.User.ID,
				Username: &sm.User.Username,
			},
		})
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: gin.H{
			"social_medias": result,
		},
	})
}

func (ctrl *socialMediaCtrl) Update(c *gin.Context) {
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

	IDStr := c.Param("socialMediaId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid social media id"})
		return
	}

	var payload struct {
		Name string `json:"name" validate:"required"`
		URL  string `json:"social_media_url" validate:"required,url"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid payload"})
		return
	}

	validation := helpers.Validate(payload)
	if len(validation) > 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:           "Validation Error",
			ValidationError: validation,
		})
		return
	}

	sm := models.SocialMedia{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&sm).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Error: "failed to get social media"})
		return
	}

	sm.Name = payload.Name
	sm.URL = payload.URL

	if err := ctrl.db.Save(&sm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed to save social media"})
		return
	}

	c.JSON(http.StatusCreated, dto.SocialMedia{
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
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	IDStr := c.Param("socialMediaId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid social media id"})
		return
	}

	if err := ctrl.db.Unscoped().Where("id=? AND user_id=?", ID, userID).Delete(&models.SocialMedia{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed to delete social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
