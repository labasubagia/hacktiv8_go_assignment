package controllers

import (
	"final_project/dto"
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
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	userID, ok := userData["id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid user id"})
		return
	}

	cm := models.Comment{}
	if err := c.ShouldBindJSON(&cm); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid payload"})
		return
	}

	cm.UserID = uint(userID)

	if err := ctrl.db.Create(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed create comment"})
		return
	}

	c.JSON(http.StatusCreated, dto.Comment{
		ID:        &cm.ID,
		Message:   &cm.Message,
		PhotoID:   &cm.PhotoID,
		UserID:    &cm.UserID,
		CreatedAt: &cm.CreatedAt,
	})
}

func (ctrl *commentCtrl) List(c *gin.Context) {

	comments := []models.Comment{}
	if err := ctrl.db.Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed get comment"})
		return
	}

	result := []dto.Comment{}
	for _, comment := range comments {
		cm := comment
		result = append(result, dto.Comment{
			ID:        &cm.ID,
			UserID:    &cm.UserID,
			CreatedAt: &cm.CreatedAt,
			UpdatedAt: &cm.UpdatedAt,
			Message:   &cm.Message,
			PhotoID:   &cm.PhotoID,
			User: &dto.User{
				ID:       &cm.User.ID,
				Username: &cm.User.Username,
				Email:    &cm.User.Email,
			},
			Photo: &dto.Photo{
				ID:      &cm.Photo.ID,
				Title:   &cm.Photo.Title,
				Caption: &cm.Photo.Caption,
				URL:     &cm.Photo.URL,
				UserID:  &cm.Photo.UserID,
			},
		})
	}

	c.JSON(http.StatusCreated, dto.Response{Data: result})
}

func (ctrl *commentCtrl) Update(c *gin.Context) {
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

	IDStr := c.Param("commentId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid comment id"})
		return
	}

	var payload struct {
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Error: "invalid payload"})
		return
	}

	cm := models.Comment{}
	if err := ctrl.db.Where("id=? AND user_id=?", ID, userID).First(&cm).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.Response{Error: "failed get comment"})
		return
	}

	cm.Message = payload.Message

	if err := ctrl.db.Save(&cm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed update comment"})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Data: dto.Comment{
			ID:        &cm.ID,
			UserID:    &cm.UserID,
			Message:   &cm.Message,
			UpdatedAt: &cm.UpdatedAt,
		},
		Message: "Success update comment",
	})
}

func (ctrl *commentCtrl) Delete(c *gin.Context) {
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

	IDStr := c.Param("commentId")
	ID, err := strconv.Atoi(IDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Error: "invalid comment id"})
		return
	}

	if err := ctrl.db.Unscoped().Where("id=? AND user_id=?", ID, userID).Delete(&models.Comment{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Error: "failed delete comment"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Message: "Your comment has been successfully deleted",
	})
}
