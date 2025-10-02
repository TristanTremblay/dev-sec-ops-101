package main

import (
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/md4"
)

// VULNERABILITY: Storing passwords as LM and NTLM hashes (Windows-style)
// This allows pass-the-hash attacks and easier cracking
var users = map[string]User{
	"admin": {
		Username:     "admin",
		Email:        "admin@intellimetrics.dev",
		LMHash:       "E52CAC67419A9A224A3B108F3FA6CB6D", // Admin123!
		NTHash:       "209C6174DA490CAEB422F3FA5A7AE634", // Admin123!
		Role:         "admin",
		PasswordHint: "Company name + 123!",
	},
	"john": {
		Username:     "john",
		Email:        "john@intellimetrics.dev",
		LMHash:       "AAD3B435B51404EEAAD3B435B51404EE", // Password1
		NTHash:       "E19CCF75EE54E06B06A5907AF13CEF42", // Password1
		Role:         "user",
		PasswordHint: "Common password",
	},
	"sarah": {
		Username:     "sarah",
		Email:        "sarah@intellimetrics.dev",
		LMHash:       "44EFCE164AB921CAAAD3B435B51404EE", // Welcome2024
		NTHash:       "32ED87BDB5FDC5E9CBA88547376818D4", // Welcome2024
		Role:         "user",
		PasswordHint: "Greeting + current year",
	},
	"dbadmin": {
		Username:     "dbadmin",
		Email:        "dbadmin@intellimetrics.dev",
		LMHash:       "C23413A8A1E7665FAAD3B435B51404EE", // Database!
		NTHash:       "B7C899154197D6E04DE72F9DDA8F7E8C", // Database!
		Role:         "admin",
		PasswordHint: "Your role + !",
	},
	"test": {
		Username:     "test",
		Email:        "test@intellimetrics.dev",
		LMHash:       "AAD3B435B51404EEAAD3B435B51404EE", // test
		NTHash:       "0CB6948805F797BF2A82807973B89537", // test
		Role:         "user",
		PasswordHint: "Same as username",
	},
}

type User struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	LMHash       string `json:"lm_hash"`
	NTHash       string `json:"nt_hash"`
	Role         string `json:"role"`
	PasswordHint string `json:"password_hint,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PassTheHashRequest struct {
	Username string `json:"username"`
	NTHash   string `json:"nt_hash"`
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// VULNERABILITY 1: Exposes user database with hashes
	router.GET("/api/auth/users", listUsers)

	// VULNERABILITY 2: Exposes password hashes in debug endpoint
	router.GET("/api/auth/debug/hashes", dumpHashes)

	// VULNERABILITY 3: Password spray vulnerable - no rate limiting
	router.POST("/api/auth/login", login)

	// VULNERABILITY 4: Allows pass-the-hash attacks
	router.POST("/api/auth/pth", passTheHash)

	// VULNERABILITY 5: Reveals password hints
	router.GET("/api/auth/hint/:username", getPasswordHint)

	// VULNERABILITY 6: User enumeration via timing/response differences
	router.POST("/api/auth/reset-password", resetPassword)

	// VULNERABILITY 7: Exposes SAM-like database dump
	router.GET("/api/auth/sam-dump", samDump)

	log.Println("Auth Service starting on port 8081")
	router.Run(":8081")
}

func listUsers(c *gin.Context) {
	// VULNERABILITY: Returns all users without authentication
	userList := []string{}
	for username := range users {
		userList = append(userList, username)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": userList,
		"count": len(userList),
	})
}

func dumpHashes(c *gin.Context) {
	// CRITICAL VULNERABILITY: Dumps password hashes without authentication
	c.JSON(http.StatusOK, users)
}

func login(c *gin.Context) {
	var req LoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// VULNERABILITY: No rate limiting, perfect for password spraying
	user, exists := users[strings.ToLower(req.Username)]
	if !exists {
		// Timing attack: different response time for non-existent users
		time.Sleep(100 * time.Millisecond)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Calculate NTLM hash of provided password
	providedNTHash := calculateNTLMHash(req.Password)

	if strings.EqualFold(providedNTHash, user.NTHash) {
		// VULNERABILITY: Weak JWT secret (same as in main API)
		token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IiArIHJlcS5Vc2VybmFtZSArICIsInJvbGUiOiIgKyB1c2VyLlJvbGUgKyAifQ.signature"

		c.JSON(http.StatusOK, gin.H{
			"token":    token,
			"username": user.Username,
			"role":     user.Role,
			"message":  "Login successful",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func passTheHash(c *gin.Context) {
	// CRITICAL VULNERABILITY: Accepts NTLM hash directly for authentication
	var req PassTheHashRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, exists := users[strings.ToLower(req.Username)]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Accept the hash directly!
	if strings.EqualFold(req.NTHash, user.NTHash) {
		c.JSON(http.StatusOK, gin.H{
			"token":    "pth-token-" + req.Username,
			"username": user.Username,
			"role":     user.Role,
			"message":  "Pass-the-hash authentication successful",
			"warning":  "This is a critical vulnerability!",
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid hash"})
}

func getPasswordHint(c *gin.Context) {
	username := strings.ToLower(c.Param("username"))

	user, exists := users[username]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// VULNERABILITY: Exposes password hints
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"hint":     user.PasswordHint,
	})
}

func resetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// VULNERABILITY: User enumeration via response differences
	for _, user := range users {
		if user.Email == req.Email {
			c.JSON(http.StatusOK, gin.H{
				"message":  "Password reset link sent to " + req.Email,
				"username": user.Username, // VULNERABILITY: Leaks username
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Email not found in system"})
}

func samDump(c *gin.Context) {
	// CRITICAL VULNERABILITY: Dumps SAM-like format for john/hashcat
	samFormat := "# Windows-style SAM dump - IntelliMetrics Auth Service\n"
	samFormat += "# Format: username:uid:lmhash:nthash:::\n"
	samFormat += "# Use this with john the ripper or hashcat!\n\n"

	uid := 1000
	for username, user := range users {
		samFormat += username + ":" + string(rune(uid)) + ":" + user.LMHash + ":" + user.NTHash + ":::\n"
		uid++
	}

	c.JSON(http.StatusOK, gin.H{
		"format": "sam",
		"data":   samFormat,
		"hint":   "john --format=nt hash.txt or hashcat -m 1000 hash.txt wordlist.txt",
	})
}

// Calculate NTLM hash (simplified - for demo purposes)
func calculateNTLMHash(password string) string {
	// Convert password to UTF-16LE
	utf16 := make([]byte, len(password)*2)
	for i, r := range password {
		utf16[i*2] = byte(r)
		utf16[i*2+1] = byte(r >> 8)
	}

	// MD4 hash
	hasher := md4.New()
	hasher.Write(utf16)
	hash := hasher.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(hash))
}
