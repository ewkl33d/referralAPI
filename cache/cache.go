package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"service/config"
	"service/db"
	"service/models"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.C.RedisHost, config.C.RedisPort),
	})
}

func SetReferralCode(code string, userID uint, email string, expiry time.Time) error {
	// Сохраняем пару "реферальный код - userID создавший"
	if userID != 0 {
		err := rdb.Set(ctx, code, userID, expiry.Sub(time.Now())).Err()
		if err != nil {
			return fmt.Errorf("failed to set referral code in cache: %w", err)
		}
	}

	if email != "" {
		err := rdb.Set(ctx, email, code, expiry.Sub(time.Now())).Err()
		if err != nil {
			return fmt.Errorf("failed to set email-code pair in cache: %w", err)
		}
	}

	return nil
}

func GetUIDByReferralCode(code string) (uint, time.Time, error) {
	val, err := rdb.Get(ctx, code).Result()
	if err == redis.Nil {
		return 0, time.Time{}, nil
	} else if err != nil {
		panic(err)
	}

	userID := uint(0)
	fmt.Sscanf(val, "%d", &userID)

	expiry, err := rdb.TTL(ctx, code).Result()
	if err != nil {
		panic(err)
	}

	return userID, time.Now().Add(expiry), nil
}

func GetReferralCodeByEmail(email string) (string, error) {
	val, err := rdb.Get(ctx, email).Result()
	if err == redis.Nil {
		return "", nil // Ключ не найден в кэше
	} else if err != nil {
		return "", fmt.Errorf("failed to get referral code from cache: %w", err)
	}

	return val, nil
}

func DeleteReferralCode(code string) error {
	// Получаем userID по реферальному коду из кеша
	userID, _, err := GetUIDByReferralCode(code)
	if err != nil {
		return fmt.Errorf("failed to get userID by referral code from cache: %w", err)
	}

	// Если userID не найден в кеше, ищем в базе данных
	if userID == 0 {
		var referral models.Referral
		result := db.DB.Where("code = ?", code).First(&referral)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("referral code not found in database")
			} else {
				return fmt.Errorf("failed to get referral code from database: %w", result.Error)
			}
		}

		userID = referral.UserID
	}

	// Получаем email пользователя по userID из базы данных
	var user models.User
	result := db.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found in database")
		} else {
			return fmt.Errorf("failed to get user from database: %w", result.Error)
		}
	}

	// Удаляем обе пары из кеша
	err = rdb.Del(ctx, code).Err()
	if err != nil {
		return fmt.Errorf("failed to delete referral code from cache: %w", err)
	}

	err = rdb.Del(ctx, user.Email).Err()
	if err != nil {
		return fmt.Errorf("failed to delete email-code pair from cache: %w", err)
	}

	return nil
}
