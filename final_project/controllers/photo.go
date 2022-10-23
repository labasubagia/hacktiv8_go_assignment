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
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	p := models.Photo{}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid payload"})
		return
	}

	validation := helpers.Validate(p)
	if len(validation) > 0 {
		c.JSON(http.StatusBadRequest, dto.Response{
			Error:           "Validation Error",
			ValidationError: validation,
		})
		return
	}

	p.UserID = uint(userID)

	if err := ctrl.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Error: "Cannot create photo",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.Photo{
			ID:        &p.ID,
			Title:     &p.Title,
			Caption:   &p.Caption,
			URL:       &p.URL,
			UserID:    &p.UserID,
			CreatedAt: &p.CreatedAt,
		},
		Message: "Create photo success",
	})
}

func (ctrl *photoCtrl) List(c *gin.Context) {

	photos := []models.Photo{}
	if err := ctrl.db.Preload("User").Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Error: "Cannot list photo",
		})
		return
	}

	result := []dto.Photo{}
	for _, photo := range photos {
		p := photo
		result = append(result, dto.Photo{
			ID:        &p.ID,
			Title:     &p.Title,
			URL:       &p.URL,
			UserID:    &p.UserID,
			CreatedAt: &p.CreatedAt,
			UpdatedAt: &p.UpdatedAt,
			Caption:   &p.Caption,
			User: &dto.User{
				Username: &p.User.Username,
				Email:    &p.User.Email,
			},
		})
	}
	c.JSON(http.StatusCreated, dto.Response{Data: result})
}

func (ctrl *photoCtrl) Update(c *gin.Context) {
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

	IDStr := c.Param("photoId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid photo id"})
		return
	}

	var payload struct {
		Title   string `json:"title" validate:"required"`
		Caption string `json:"caption"`
		URL     string `json:"photo_url" validate:"required,url"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid input"})
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

	p := models.Photo{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&p).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Error: "Cannot update photo"})
		return
	}

	p.Caption = payload.Caption
	p.Title = payload.Title
	p.URL = payload.URL

	if err := ctrl.db.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "Failed to update photo"})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.Photo{
			ID:        &p.ID,
			Title:     &p.Title,
			URL:       &p.URL,
			UserID:    &p.UserID,
			UpdatedAt: &p.UpdatedAt,
			Caption:   &p.Caption,
		},
		Message: "Update photo success",
	})
}

func (ctrl *photoCtrl) Delete(c *gin.Context) {
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

	IDStr := c.Param("photoId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid photo id"})
		return
	}

	if err := ctrl.db.Unscoped().Where("id=? AND user_id=?", ID, userID).Delete(&models.Photo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message: "Your photo has been successfully deleted",
	})
}
