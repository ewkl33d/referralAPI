package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"service/cache"
	"service/db"
	_ "service/docs"
	"service/models"
	"service/utils"
	"strconv"
	"time"
)

type CreateReferralCodeInput struct {
	Expiry time.Time `json:"expiry"`
}

type RegisterWithReferralCodeInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ReferralCode string `json:"referral_code"`
}

// @Summary Create a referral code
// @Description Create a new referral code for the authenticated user
// @Tags Referrals
// @Accept json
// @Produce json
// @Param input body CreateReferralCodeInput true "Expiry time for the referral code"
// @Success 200 {object} map[string]string "referral_code"
// @Failure 400 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security ApiKeyAuth
// @Router /referral/create [post]
func CreateReferralCode(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var input CreateReferralCodeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var referral models.Referral
	db.DB.Where("user_id = ?", userID).First(&referral)

	if referral.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Referral code already exists"})
		return
	}
	err := cache.SetReferralCode(referral.Code, userID, "", input.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set pair in cache"})
		return
	}
	referral = models.Referral{
		UserID: userID,
		Code:   utils.GenerateReferralCode(),
		Expiry: input.Expiry,
	}
	db.DB.Create(&referral)

	c.JSON(http.StatusOK, gin.H{"referral_code": referral.Code})
}

// @Summary Delete a referral code
// @Description Delete the referral code for the authenticated user
// @Tags Referrals
// @Produce json
// @Success 200 {object} map[string]string "message"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security ApiKeyAuth
// @Router /referral/delete [delete]
func DeleteReferralCode(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var referral models.Referral
	if err := db.DB.Where("user_id = ?", userID).First(&referral).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Referral code not found"})
		return
	}

	err := cache.DeleteReferralCode(referral.Code)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	db.DB.Delete(&referral)

	c.JSON(http.StatusOK, gin.H{"message": "Referral code deleted successfully"})
}

// @Summary Get referral code by email
// @Description Get the referral code for a user by their email
// @Tags Referrals
// @Param email path string true "User email"
// @Produce json
// @Success 200 {object} map[string]string "referral_code"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security ApiKeyAuth
// @Router /referral/get/{email} [get]
func GetReferralCodeByEmail(c *gin.Context) {
	email := c.Param("email")

	// Попытка получить реферальный код из кеша по email
	code, err := cache.GetReferralCodeByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral code from cache"})
		return
	}

	// Если код найден в кеше, возвращаем его
	if code != "" {
		c.JSON(http.StatusOK, gin.H{"referral_code": code})
		return
	}

	// Если код не найден в кеше, ищем в базе данных
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user from database"})
		}
		return
	}

	var referral models.Referral
	if err := db.DB.Where("user_id = ?", user.ID).First(&referral).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Referral code not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral code from database"})
		}
		return
	}

	expiry := referral.Expiry
	err = cache.SetReferralCode(referral.Code, user.ID, email, expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set pair in cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"referral_code": referral.Code})
}

// @Summary Register a new user with a referral code
// @Description Register a new user with a referral code
// @Tags Referrals
// @Accept json
// @Produce json
// @Param input body RegisterWithReferralCodeInput true "User registration data"
// @Success 200 {object} map[string]string "message"
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Router /referral/register [post]
func RegisterWithReferralCode(c *gin.Context) {
	var input RegisterWithReferralCodeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	referrerID, expiry, err := cache.GetUIDByReferralCode(input.ReferralCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral code from cache"})
		return
	}

	// Если код не найден в кеше, ищем в базе данных
	if referrerID == 0 {
		var referral models.Referral
		result := db.DB.Where("code = ?", input.ReferralCode).First(&referral)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invalid referral code"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get referral code from database"})
			}
			return
		}

		if referral.Expiry.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Referral code expired"})
			return
		}

		referrerID = referral.UserID
		expiry := referral.Expiry
		err = cache.SetReferralCode(referral.Code, referrerID, "", expiry)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set pair in cache"})
			return
		}
	}

	if expiry.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Referral code expired"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	user := models.User{
		Email:      input.Email,
		Password:   string(hashedPassword),
		ReferrerID: referrerID,
	}

	result := db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// @Summary Get referrals by referrer ID
// @Description Get the list of users referred by a specific referrer
// @Tags Referrals
// @Param id path int true "Referrer ID"
// @Produce json
// @Success 200 {object} map[string][]string "emails"
// @Failure 400 {object} map[string]string "error"
// @Failure 404 {object} map[string]string "error"
// @Failure 500 {object} map[string]string "error"
// @Security ApiKeyAuth
// @Router /referrals/referrer/{id} [get]
func GetReferralsByReferrerID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var users []models.User
	if err := db.DB.Where("referrer_id = ?", id).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		}
		return
	}

	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}
