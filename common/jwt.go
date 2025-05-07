package common

import (
	"github.com/RizkiMufrizal/gofiber-clean-architecture/configuration"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/entity"
	"github.com/RizkiMufrizal/gofiber-clean-architecture/exception"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
	"time"
)

func GenerateToken(username string, roles []map[string]interface{}, user entity.User, config configuration.Config) string {
	jwtSecret := config.Get("JWT_SECRET_KEY")
	jwtExpired, err := strconv.Atoi(config.Get("JWT_EXPIRE_MINUTES_COUNT"))
	exception.PanicLogging(err)

	claims := jwt.MapClaims{
		"username":     username,
		"roles":        roles,
		"exp":          time.Now().Add(time.Hour * time.Duration(jwtExpired)).Unix(),
		"userId":       user.ID,
		"emailAddress": user.Email,
		"phoneNumber":  user.PhoneNumber,
		"refineryId":   user.RefineryId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(jwtSecret))
	exception.PanicLogging(err)

	return tokenSigned
}
