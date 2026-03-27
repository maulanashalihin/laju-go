package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/velostack/velostack-go/app/models"
	"github.com/velostack/velostack-go/app/services"
	"github.com/velostack/velostack-go/app/session"
)

type AuthHandler struct {
	authService    *services.AuthService
	userService    *services.UserService
	store          *session.Store
	inertiaService *services.InertiaService
}

func NewAuthHandler(authService *services.AuthService, userService *services.UserService, store *session.Store, inertiaService *services.InertiaService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		userService:    userService,
		store:          store,
		inertiaService: inertiaService,
	}
}

// ShowLoginForm displays the login page
func (h *AuthHandler) ShowLoginForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Login", fiber.Map{
		"Title": "Login",
	})
}

// ShowRegisterForm displays the register page
func (h *AuthHandler) ShowRegisterForm(c *fiber.Ctx) error {
	return h.inertiaService.Render(c, "auth/Register", fiber.Map{
		"Title": "Register",
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "All fields are required",
		})
	}

	// Register user
	user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already registered",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	// Create session
	sess, _ := h.store.Get(c)
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))
	sess.Save()

	return c.JSON(fiber.Map{
		"message": "Registration successful",
		"user":    user.ToResponse(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Authenticate user
	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	// Create session
	sess, _ := h.store.Get(c)
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))
	sess.Save()

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user.ToResponse(),
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	sess.Destroy()

	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// GoogleLogin initiates Google OAuth login
func (h *AuthHandler) GoogleLogin(c *fiber.Ctx) error {
	state := generateState()
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HTTPOnly: true,
		SameSite: "Lax",
	})

	url := h.authService.GetOAuthURL(state)
	return c.Redirect(url)
}

// GoogleCallback handles Google OAuth callback
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	code := c.Query("code")

	// Validate state
	storedState := c.Cookies("oauth_state")
	if state != storedState {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid OAuth state",
		})
	}

	// Clear the state cookie
	c.ClearCookie("oauth_state")

	// Process the token
	user, err := h.authService.ProcessGoogleToken(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate with Google",
		})
	}

	// Create session
	sess, _ := h.store.Get(c)
	sess.Set("user_id", user.ID)
	sess.Set("email", user.Email)
	sess.Set("role", string(user.Role))
	sess.Save()

	// Redirect to app or return JSON for API calls
	if c.Get("Accept") == "application/json" {
		return c.JSON(fiber.Map{
			"message": "Login successful",
			"user":    user.ToResponse(),
		})
	}

	return c.Redirect("/app")
}

// Me returns the current authenticated user
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	sess, _ := h.store.Get(c)
	userID := sess.Get("user_id")

	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	user, err := h.authService.GetUserByID(userID.(int64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToResponse(),
	})
}

// generateState generates a random state string for OAuth
func generateState() string {
	// In production, use a proper random generator
	return "random-state-placeholder"
}
