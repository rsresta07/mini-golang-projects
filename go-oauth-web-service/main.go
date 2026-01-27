package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"                // Session management (cookies)
	"github.com/joho/godotenv"                   // Load .env file into environment
	"github.com/markbates/goth"                  // OAuth abstraction library
	"github.com/markbates/goth/gothic"           // Helpers for OAuth flow (handlers, sessions)
	"github.com/markbates/goth/providers/google" // Google OAuth provider
)

func main() {
	// Create Gin router with logger + recovery middleware
	r := gin.Default()

	// Create cookie-based session store
	store := sessions.NewCookieStore([]byte("super-secret-key-32-bytes"))

	// Session configuration
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400, // 1 day in seconds
		HttpOnly: true,  // Prevent JS access to cookie
	}

	// Tell Goth to use this session store instead of default
	gothic.Store = store

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file failed to load!")
	}

	// Read OAuth credentials
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallbackURL := os.Getenv("CLIENT_CALLBACK_URL")

	if clientID == "" || clientSecret == "" || clientCallbackURL == "" {
		log.Fatal("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	// Register the Google OAuth provider with required scopes
	goth.UseProviders(
		google.New(
			clientID,
			clientSecret,
			clientCallbackURL,
			"email",   // Access user's email
			"profile", // Access basic profile info
		),
	)

	// Route definitions
	r.GET("/", home)                                   // Landing page
	r.GET("/auth/:provider", signInWithProvider)       // Start OAuth flow
	r.GET("/auth/:provider/callback", callbackHandler) // OAuth callback
	r.GET("/success", Success)                         // Post-login page

	r.Run(":5000")
}

// Renders the home page
func home(c *gin.Context) {
	// Parse HTML template from disk
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Execute template and write to response
	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// Initiates OAuth login with provider (google)
func signInWithProvider(c *gin.Context) {
	// Extract provider name from URL (e.g., "google")
	provider := c.Param("provider")

	// Goth determines provider from query param
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Redirects user to provider's consent page
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// Handles OAuth callback after provider authentication
func callbackHandler(c *gin.Context) {
	// Extract provider again for Goth
	provider := c.Param("provider")

	// Attach provider to request query
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	// Completes OAuth flow and fetches user data
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Retrieve the session named "user-session"
	session, _ := gothic.Store.Get(c.Request, "user-session")

	// Store user details in session
	session.Values["email"] = user.Email
	session.Values["name"] = user.Name
	session.Values["avatar"] = user.AvatarURL

	// Persist session to cookie
	session.Save(c.Request, c.Writer)

	// Redirect user to the success page
	c.Redirect(http.StatusTemporaryRedirect, "/success")
}

// Success Displays logged-in user info
func Success(c *gin.Context) {
	// Load session
	session, err := gothic.Store.Get(c.Request, "user-session")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Type assertions from session values
	name, _ := session.Values["name"].(string)
	email, _ := session.Values["email"].(string)
	avatar, _ := session.Values["avatar"].(string)

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(fmt.Sprintf(`
        <div style="padding:40px;text-align:center;">
            <img src="%s" style="border-radius:50%%;width:80px;"><br><br>
            <h2>Signed in successfully</h2>
            <p>Name: %s</p>
            <p>Email: %s</p>
        </div>
    `, avatar, name, email)))
}
