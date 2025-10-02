package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

// VULNERABILITY 1: Hardcoded credentials
const (
	MONGODB_URI = "mongodb://localhost:27017"
	DB_NAME     = "intellimetrics"
	OPENAI_KEY  = "sk-proj-demo-key-abc123def456ghi789"
	JWT_SECRET  = "super_secret_jwt_key_123"
	STRIPE_KEY  = "sk_live_demo_stripe_key_xyz"
	AWS_ACCESS  = "AKIAIOSFODNN7EXAMPLE"
	AWS_SECRET  = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	ADMIN_EMAIL = "admin@intellimetrics.dev"
	ADMIN_PASS  = "Admin123!"
)

type Company struct {
	CompanyID   string    `json:"company_id" bson:"company_id"`
	CompanyName string    `json:"company_name" bson:"company_name"`
	Revenue     int64     `json:"revenue" bson:"revenue"`
	Users       int       `json:"users" bson:"users"`
	GrowthRate  float64   `json:"growth_rate" bson:"growth_rate"`
	Industry    string    `json:"industry" bson:"industry"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

type AnalyticsQuery struct {
	CompanyID string                 `json:"company_id"`
	Filters   map[string]interface{} `json:"filters"`
}

type AIInsight struct {
	CompanyID string    `json:"company_id"`
	Insight   string    `json:"insight"`
	Generated time.Time `json:"generated"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	mongoClient = client

	defer func() {
		if err = mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	if err = mongoClient.Ping(context.Background(), nil); err != nil {
		log.Fatal("Cannot ping MongoDB:", err)
	}

	log.Println("Connected to MongoDB successfully")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	router.GET("/health", healthCheck)
	router.GET("/api/stats", getPublicStats)
	router.POST("/api/companies", addCompany)
	router.POST("/api/analytics/query", queryAnalytics)
	router.GET("/api/insights/:company_id", getAIInsights)
	router.GET("/api/admin/config", getConfig)
	router.GET("/api/admin/debug", debugInfo)
	router.GET("/api/admin/system", systemCommand)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"version": "1.0.0-vulnerable",
		"services": gin.H{
			"api":      "running",
			"database": "connected",
		},
	})
}

func getPublicStats(c *gin.Context) {
	collection := mongoClient.Database(DB_NAME).Collection("companies")
	count, _ := collection.CountDocuments(context.Background(), bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"total_companies": count,
		"message":         "IntelliMetrics - AI-Powered Analytics Platform",
	})
}

func addCompany(c *gin.Context) {
	var company Company

	if err := c.BindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company.CreatedAt = time.Now()

	collection := mongoClient.Database(DB_NAME).Collection("companies")
	_, err := collection.InsertOne(context.Background(), company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company added successfully", "company": company})
}

func queryAnalytics(c *gin.Context) {
	var query AnalyticsQuery

	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{}
	if query.CompanyID != "" {
		filter["company_id"] = query.CompanyID
	}

	for key, value := range query.Filters {
		filter[key] = value
	}

	collection := mongoClient.Database(DB_NAME).Collection("companies")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query failed"})
		return
	}
	defer cursor.Close(context.Background())

	var results []Company
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse results"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func getAIInsights(c *gin.Context) {
	companyID := c.Param("company_id")

	insight := AIInsight{
		CompanyID: companyID,
		Insight:   "Based on recent trends, revenue is projected to grow by 15% this quarter. Consider expanding marketing budget in Q2.",
		Generated: time.Now(),
	}

	c.JSON(http.StatusOK, insight)
}

func getConfig(c *gin.Context) {
	config := gin.H{
		"mongodb_uri":       MONGODB_URI,
		"database_name":     DB_NAME,
		"openai_api_key":    OPENAI_KEY,
		"jwt_secret":        JWT_SECRET,
		"stripe_secret_key": STRIPE_KEY,
		"aws_access_key":    AWS_ACCESS,
		"aws_secret_key":    AWS_SECRET,
		"admin_email":       ADMIN_EMAIL,
		"admin_password":    ADMIN_PASS,
		"environment":       os.Getenv("ENV"),
	}

	c.JSON(http.StatusOK, config)
}

func debugInfo(c *gin.Context) {
	hostname, _ := os.Hostname()

	info := gin.H{
		"hostname":    hostname,
		"go_version":  "go1.21",
		"os":          "linux",
		"working_dir": "/app",
		"env_vars":    os.Environ(),
	}

	c.JSON(http.StatusOK, info)
}

func systemCommand(c *gin.Context) {
	cmd := c.Query("cmd")

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'cmd' parameter"})
		return
	}

	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"output": string(output),
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"output": string(output),
	})
}
