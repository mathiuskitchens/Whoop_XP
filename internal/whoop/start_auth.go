package whoop

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

// generateRandomState creates a cryptographically secure random string.
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// StartAuthHandler begins the OAuth flow by redirecting the user to Whoop.
func StartAuthHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
			return
		}

		// Optionally link to a logged-in user (set to nil for now)
		var userID *int = nil

		if err := StoreState(db, state, userID, 10*time.Minute); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store state"})
			return
		}

		clientID := os.Getenv("WHOOP_CLIENT_ID")
		redirectURI := os.Getenv("WHOOP_REDIRECT_URI") // e.g. http://localhost:8080/whoop/callback
		scope := "offline+read:workout+read:profile"   // choose scopes as needed

		authURL := fmt.Sprintf(
			"https://api.prod.whoop.com/oauth/oauth2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s&state=%s",
			clientID, scope, redirectURI, state,
		)
		fmt.Println("clientID:", clientID)
		fmt.Println("redirectURI:", redirectURI)
		c.Redirect(http.StatusFound, authURL)
	}
}
