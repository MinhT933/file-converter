package api

import (
	"strings"

	"github.com/MinhT933/file-converter/internal/domain"
	"github.com/MinhT933/file-converter/internal/infra/auth"
	"github.com/MinhT933/file-converter/internal/services"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	DisplayName string `json:"display_name" validate:"required"`
	PhotoURL    string `json:"photo_url"`
}

type LoginResponse struct {
	Token   string     `json:"token"`
	User    *auth.User `json:"user"`
	Message string     `json:"message"`
}

type AuthUser struct {
	UID   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LinkAccountRequest struct {
	ExistingUserID string              `json:"existing_user_id"`
	ProviderData   domain.ProviderData `json:"provider_data"`
}

type SocialLoginRequest struct {
	Provider    string `json:"provider" validate:"required"`     // "google", "facebook", "linkedin"
	AccessToken string `json:"access_token" validate:"required"` // Token từ provider
	IDToken     string `json:"id_token,omitempty"`               // ID token (cho Google)
}

func NewHandler(authProvider *auth.FirebaseProvider, authService *services.AuthService) *Handler {
	return &Handler{
		AuthProvider: authProvider,
		AuthService:  authService,
	}
}

// Login godoc
// @Summary     Login with email and password
// @Description Login user with email and password
// @Accept      json
// @Tags 	  auth
// @Param     request body LoginRequest true "Login request"
// @Success     200 {object} LoginResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, token, err := h.AuthProvider.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Login failed: "+err.Error())
	}

	return c.JSON(LoginResponse{
		Token:   token,
		User:    user,
		Message: "Login successful",
	})
}

// SocialLogin godoc
// @Summary      Login with social providers (Google, Facebook, LinkedIn)
// @Description Login user with social provider (Google, Facebook, etc.)
// @Tags 	   auth
// @Accept       json
// @produces    json
// @Param       request body SocialLoginRequest true "Social login request"
// @Success     200 {object} LoginResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /auth/social/login [post]
func (h *Handler) SocialLogin(c *fiber.Ctx) error {
	var req SocialLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Verify token dựa trên provider
	switch strings.ToLower(req.Provider) {
	case "google":
		return h.handleGoogleLogin(c, req)
	case "facebook":
		return h.handleFacebookLogin(c, req)
	// case "linkedin":
	//  return h.handleLinkedInLogin(c, req)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":     "Unsupported provider",
			"supported": []string{"google", "facebook", "linkedin"},
		})
	}
}

// VerifyToken godoc
// @Summary     Verify JWT token
// @Description Verify JWT token and return user information
// @Tags        auth
// @Accept      json
// @Param       token query string true "JWT token"
// @Success     200 {object} LoginResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /auth/verify [get]
// func (h *Handler) VerifyToken(c *fiber.Ctx) error {
// 	authHeader := c.Get("Authorization")
// 	if authHeader == "" {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Authorization header required",
// 		})
// 	}

// 	token := strings.TrimPrefix(authHeader, "Bearer ")
// 	if token == authHeader {
// 		token = authHeader // Fallback nếu không có "Bearer "
// 	}

// 	// Verify token với Auth Provider
// 	user, err := h.AuthProvider.VerifyToken(c.Context(), token)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "Invalid or expired token",
// 		})
// 	}

// 	return c.JSON(AuthUser{
// 		UID:   user.UID,
// 		Email: user.Email,
// 		Name:  user.DisplayName,
// 	})
// }

func (h *Handler) handleGoogleLogin(c *fiber.Ctx, req SocialLoginRequest) error {
	// Xử lý đăng nhập với Google
	FbUser, _, err := h.AuthProvider.GoogleLogin(c.Context(), req.AccessToken, req.IDToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Google login failed: "+err.Error())
	}

	providerData := domain.ProviderData{
		UID:         FbUser.UID,
		Email:       FbUser.Email,
		DisplayName: FbUser.DisplayName,
		PhotoURL:    FbUser.PhotoURL,
	}

	authResult, err := h.AuthService.HandleProviderLogin(c.Context(), providerData)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Google login failed: "+err.Error())
	}

	if authResult.RequiresLinking {
		// Link the social account with the existing user account
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"requires_linking": true,
			"message":          "Email already exists with different provider",
			"existing_user":    authResult.ExistingUser,
			"provider_data":    authResult.ProviderData,
		})
	}

	return c.JSON(LoginResponse{
		Token:   authResult.SessionToken,
		User:    convertDomainUserToAPIUser(authResult.User),
		Message: "Google login successful",
	})
}

func (h *Handler) handleFacebookLogin(c *fiber.Ctx, req SocialLoginRequest) error {
	// Similar implementation for Facebook
	// firebaseUser, err := h.AuthProvider.VerifyFacebookToken(c.Context(), req.AccessToken)
	// if err != nil {
	//     return fiber.NewError(fiber.StatusUnauthorized, "Invalid Facebook token: "+err.Error())
	// }

	// providerData := domain.ProviderData{
	//     UID:           firebaseUser.UID,
	//     Email:         firebaseUser.Email,
	//     DisplayName:   firebaseUser.DisplayName,
	//     PhotoURL:      firebaseUser.PhotoURL,
	//     ProviderID:    "facebook.com",
	//     EmailVerified: firebaseUser.EmailVerified,
	// }

	// authResult, err := h.AuthService.HandleProviderLogin(c.Context(), providerData)
	// if err != nil {
	//     return fiber.NewError(fiber.StatusInternalServerError, "Database error: "+err.Error())
	// }

	// if authResult.RequiresLinking {
	//     return c.Status(fiber.StatusConflict).JSON(fiber.Map{
	//         "requires_linking": true,
	//         "message":         "Email already exists with different provider",
	//         "existing_user":   authResult.ExistingUser,
	//         "provider_data":   authResult.ProviderData,
	//     })
	// }

	// return c.JSON(LoginResponse{
	//     Token:   authResult.SessionToken,
	//     User:    convertDomainUserToAPIUser(authResult.User),
	//     Message: "Facebook login successful",
	// })
	return nil
}

func convertDomainUserToAPIUser(user *domain.User) *auth.User {
	return &auth.User{
		UID:           user.UserID,
		Email:         user.Email,
		DisplayName:   user.Name,
		PhotoURL:      getStringValue(user.AvatarURL),
		EmailVerified: user.EmailVerified,
	}
}

func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// func (h *Handler) handleLinkedInLogin(c *fiber.Ctx, req SocialLoginRequest) error {
// 	// Xử lý đăng nhập với LinkedIn

// }
