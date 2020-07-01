package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/figassis/goinagbe/pkg/utl/model"
	"github.com/figassis/goinagbe/pkg/utl/zaplog"

	echo "github.com/labstack/echo/v4"

	jwt "github.com/dgrijalva/jwt-go"
)

// New generates new JWT service necessery for auth middleware
func New(secret, algo string, d, r int) *Service {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid jwt signing method")
	}
	return &Service{
		key:             []byte(secret),
		algo:            signingMethod,
		duration:        time.Duration(d) * time.Minute,
		refreshDuration: time.Duration(r) * time.Minute,
	}
}

// Service provides a Json-Web-Token authentication implementation
type Service struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	duration time.Duration

	// Duration for which the jwt refresh token is valid.
	refreshDuration time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod
}

// MWFunc makes JWT implement the Middleware interface.
func (j *Service) MWFunc() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := j.ParseToken(c)
			if err != nil || !token.Valid {
				zaplog.ZLog("Invalid auth token")
				return c.NoContent(http.StatusUnauthorized)
			}

			claims := token.Claims.(jwt.MapClaims)
			// zaplog.ZLog(claims)

			id, ok := claims["id"].(string)
			if !ok {
				zaplog.ZLog("Invalid user id in token")
				return c.NoContent(http.StatusUnauthorized)
			}

			organization, ok := claims["o"].(string)
			if !ok {
				zaplog.ZLog("Invalid organization id in token")
				return c.NoContent(http.StatusUnauthorized)
			}

			username, ok := claims["u"].(string)
			if !ok {
				zaplog.ZLog("Invalid username in token")
				return c.NoContent(http.StatusUnauthorized)
			}
			email, ok := claims["e"].(string)
			if !ok {
				zaplog.ZLog("Invalid email in token")
				return c.NoContent(http.StatusUnauthorized)
			}

			tempRole, ok := claims["r"].(float64)
			if !ok {
				zaplog.ZLog("Invalid role in token")
				return c.NoContent(http.StatusUnauthorized)
			}
			role := model.AccessRole(tempRole)

			c.Set("id", id)
			c.Set("organization", organization)
			c.Set("username", username)
			c.Set("email", email)
			c.Set("role", role)

			return next(c)
		}
	}
}

// ParseToken parses token from Authorization header
func (j *Service) ParseToken(c echo.Context) (*jwt.Token, error) {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return nil, model.ErrGeneric
	}
	parts := strings.SplitN(token, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, model.ErrGeneric
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if j.algo != token.Method {
			return nil, model.ErrGeneric
		}
		return j.key, nil
	})

}

// GenerateToken generates new JWT token and populates it with user data
func (j *Service) GenerateToken(u *model.User) (string, string, string, string, error) {
	expire := time.Now().Add(j.duration)
	refresh_expire := time.Now().Add(j.refreshDuration)

	token := jwt.NewWithClaims((j.algo), jwt.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   int64(u.Role),
		"o":   u.Organization,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", "", "", "", err
	}

	refreshToken := jwt.NewWithClaims((j.algo), jwt.MapClaims{
		"id":  u.ID,
		"u":   u.Username,
		"e":   u.Email,
		"r":   int64(u.Role),
		"o":   u.Organization,
		"exp": refresh_expire.Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(j.key)
	if err != nil {
		return "", "", "", "", err
	}

	return tokenString, expire.Format(time.RFC3339), refreshTokenString, refresh_expire.Format(time.RFC3339), err
}
