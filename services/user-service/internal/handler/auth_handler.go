package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/config"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/dto"
	"github.com/squ1ky/avigo-c2c-marketplace/services/user-service/internal/service"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

type AuthHandler struct {
	authService  *service.AuthService
	cookieConfig config.CookieConfig
	jwtConfig    config.JWTConfig
}

func NewAuthHandler(
	authService *service.AuthService,
	cookieConfig config.CookieConfig,
	jwtConfig config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		cookieConfig: cookieConfig,
		jwtConfig:    jwtConfig,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		InvalidJSON(c, err)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "registration successful, please check your email for confirmation code",
		"data":    resp,
	})
}

func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	var req dto.ConfirmEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		InvalidJSON(c, err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	if err := h.authService.ConfirmEmail(c.Request.Context(), userID, req.Code); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "email confirmed successfully, you can now login",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		InvalidJSON(c, err)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	h.setAuthCookies(c, resp.AccessToken, resp.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user":    resp.User,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "user not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	h.clearAuthCookies(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "refresh token not found",
			"code":  "MISSING_REFRESH_TOKEN",
		})
		return
	}

	resp, err := h.authService.RefreshTokens(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearAuthCookies(c)
		_ = c.Error(err)
		return
	}

	h.setAuthCookies(c, resp.AccessToken, resp.RefreshToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "tokens refreshed successfully",
		"user":    resp.User,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDHeader := c.GetHeader("X-User-ID")
	if userIDHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
			"code":  "NOT_AUTHENTICATED",
		})
		return
	}

	userID, err := uuid.Parse(userIDHeader)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID format",
			"code":  "INVALID_USER_ID",
		})
		return
	}

	userInfo, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	sameSite := h.parseSameSite()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    accessToken,
		MaxAge:   int(h.jwtConfig.AccessTokenDuration.Seconds()),
		Path:     "/",
		Domain:   h.cookieConfig.Domain,
		Secure:   h.cookieConfig.Secure,
		HttpOnly: h.cookieConfig.HTTPOnly,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    refreshToken,
		MaxAge:   int(h.jwtConfig.RefreshTokenDuration.Seconds()),
		Path:     "/",
		Domain:   h.cookieConfig.Domain,
		Secure:   h.cookieConfig.Secure,
		HttpOnly: h.cookieConfig.HTTPOnly,
		SameSite: sameSite,
	})
}

func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	sameSite := h.parseSameSite()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   h.cookieConfig.Domain,
		Secure:   h.cookieConfig.Secure,
		HttpOnly: h.cookieConfig.HTTPOnly,
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Domain:   h.cookieConfig.Domain,
		Secure:   h.cookieConfig.Secure,
		HttpOnly: h.cookieConfig.HTTPOnly,
		SameSite: sameSite,
	})
}

func (h *AuthHandler) parseSameSite() http.SameSite {
	switch h.cookieConfig.SameSite {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}
