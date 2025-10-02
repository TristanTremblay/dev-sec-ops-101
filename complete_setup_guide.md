# IntelliMetrics CTF - Complete Setup Guide

## Challenge Overview

**Target:** IntelliMetrics AI Analytics Platform  
**Difficulty:** Intermediate  
**Category:** Web Application, Cloud Infrastructure  
**Tags:** Command Injection, NoSQL Injection, Secrets Exposure, Post-Exploitation

### Attack Path
```
1. Reconnaissance → Port scanning, service enumeration
2. Initial Access → Command injection vulnerability
3. Exploitation → Metasploit payload delivery
4. Post-Exploitation → Credential harvesting, database access
5. Lateral Movement → Network pivoting
```

---

## Project Structure

```
intellimetrics-demo/
├── backend/
│   ├── cmd/api/
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/           # shadcn components
│   │   │   └── Dashboard.tsx
│   │   ├── layouts/
│   │   │   └── Layout.astro
│   │   ├── pages/
│   │   │   └── index.astro
│   │   └── styles/
│   │       └── globals.css
│   ├── package.json
│   ├── astro.config.mjs
│   ├── tailwind.config.mjs
│   └── tsconfig.json
├── exploits/
│   ├── 01-recon.sh
│   ├── 02-enumerate.sh
│   └── 03-exploit.sh
├── docker-compose.yml
└── README.md
```

---

## Part 1: Backend Setup

### File: `backend/go.mod`

```go
module intellimetrics-api

go 1.21

require (
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-gonic/gin v1.9.1
	go.mongodb.org/mongo-driver v1.13.1
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

### File: `backend/cmd/api/main.go`

```go
package main

import (
	"context"
	"encoding/json"
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
	MONGODB_URI   = "mongodb://mongo:27017"
	DB_NAME       = "intellimetrics"
	OPENAI_KEY    = "sk-proj-demo-key-abc123def456ghi789"
	JWT_SECRET    = "super_secret_jwt_key_123"
	STRIPE_KEY    = "sk_live_demo_stripe_key_xyz"
	AWS_ACCESS    = "AKIAIOSFODNN7EXAMPLE"
	AWS_SECRET    = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
	ADMIN_EMAIL   = "admin@intellimetrics.dev"
	ADMIN_PASS    = "Admin123!"
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
	// VULNERABILITY 2: No authentication on MongoDB
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

	// VULNERABILITY 3: Wide-open CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	// Public endpoints
	router.GET("/health", healthCheck)
	router.GET("/api/stats", getPublicStats)

	// VULNERABILITY 4: No authentication on sensitive endpoints
	router.POST("/api/companies", addCompany)
	router.POST("/api/analytics/query", queryAnalytics)
	router.GET("/api/insights/:company_id", getAIInsights)

	// VULNERABILITY 5: Dangerous admin endpoints
	router.GET("/api/admin/config", getConfig)
	router.GET("/api/admin/debug", debugInfo)
	router.GET("/api/admin/system", systemCommand) // CRITICAL: Command injection

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
		"message":        "IntelliMetrics - AI-Powered Analytics Platform",
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

// VULNERABILITY 6: NoSQL Injection
func queryAnalytics(c *gin.Context) {
	var query AnalyticsQuery

	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CRITICAL: Direct use of user-provided filters
	filter := bson.M{}
	if query.CompanyID != "" {
		filter["company_id"] = query.CompanyID
	}

	// DANGEROUS: Allows arbitrary MongoDB operators
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

// VULNERABILITY 7: Exposes all secrets
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

// VULNERABILITY 8: Leaks system information
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

// VULNERABILITY 9: CRITICAL - Command Injection
func systemCommand(c *gin.Context) {
	cmd := c.Query("cmd")

	if cmd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'cmd' parameter"})
		return
	}

	// CRITICAL: Direct command execution
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
```

### File: `backend/Dockerfile`

```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# VULNERABILITY: Running as root
# VULNERABILITY: No security updates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

EXPOSE 8080

# VULNERABILITY: No non-root user
CMD ["./main"]
```

---

## Part 2: Frontend Setup (Astro + shadcn/ui)

### File: `frontend/package.json`

```json
{
  "name": "intellimetrics-frontend",
  "type": "module",
  "version": "1.0.0",
  "scripts": {
    "dev": "astro dev",
    "start": "astro dev",
    "build": "astro check && astro build",
    "preview": "astro preview"
  },
  "dependencies": {
    "@astrojs/check": "^0.3.4",
    "@astrojs/node": "^8.2.0",
    "@astrojs/react": "^3.0.9",
    "@astrojs/tailwind": "^5.1.0",
    "@radix-ui/react-avatar": "^1.0.4",
    "@radix-ui/react-label": "^2.0.2",
    "@radix-ui/react-select": "^2.0.0",
    "@radix-ui/react-separator": "^1.0.3",
    "@radix-ui/react-slot": "^1.0.2",
    "@radix-ui/react-tabs": "^1.0.4",
    "astro": "^4.1.2",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.1.0",
    "lucide-react": "^0.294.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "recharts": "^2.10.3",
    "tailwind-merge": "^2.2.0",
    "tailwindcss": "^3.4.0",
    "tailwindcss-animate": "^1.0.7",
    "typescript": "^5.3.3"
  },
  "devDependencies": {
    "@types/react": "^18.2.48",
    "@types/react-dom": "^18.2.18"
  }
}
```

### File: `frontend/astro.config.mjs`

```javascript
import { defineConfig } from 'astro/config';
import react from '@astrojs/react';
import tailwind from '@astrojs/tailwind';
import node from '@astrojs/node';

export default defineConfig({
  output: 'server',
  adapter: node({
    mode: 'standalone'
  }),
  integrations: [
    react(),
    tailwind({
      applyBaseStyles: false,
    })
  ],
  vite: {
    ssr: {
      noExternal: ['react-icons']
    }
  }
});
```

### File: `frontend/tailwind.config.mjs`

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: ["class"],
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: { height: 0 },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: 0 },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
```

### File: `frontend/tsconfig.json`

```json
{
  "extends": "astro/tsconfigs/strict",
  "compilerOptions": {
    "jsx": "react-jsx",
    "jsxImportSource": "react",
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"]
    }
  }
}
```

### File: `frontend/src/styles/globals.css`

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
 
@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --card: 0 0% 100%;
    --card-foreground: 222.2 84% 4.9%;
    --popover: 0 0% 100%;
    --popover-foreground: 222.2 84% 4.9%;
    --primary: 222.2 47.4% 11.2%;
    --primary-foreground: 210 40% 98%;
    --secondary: 210 40% 96.1%;
    --secondary-foreground: 222.2 47.4% 11.2%;
    --muted: 210 40% 96.1%;
    --muted-foreground: 215.4 16.3% 46.9%;
    --accent: 210 40% 96.1%;
    --accent-foreground: 222.2 47.4% 11.2%;
    --destructive: 0 84.2% 60.2%;
    --destructive-foreground: 210 40% 98%;
    --border: 214.3 31.8% 91.4%;
    --input: 214.3 31.8% 91.4%;
    --ring: 222.2 84% 4.9%;
    --radius: 0.5rem;
  }
 
  .dark {
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --card: 222.2 84% 4.9%;
    --card-foreground: 210 40% 98%;
    --popover: 222.2 84% 4.9%;
    --popover-foreground: 210 40% 98%;
    --primary: 210 40% 98%;
    --primary-foreground: 222.2 47.4% 11.2%;
    --secondary: 217.2 32.6% 17.5%;
    --secondary-foreground: 210 40% 98%;
    --muted: 217.2 32.6% 17.5%;
    --muted-foreground: 215 20.2% 65.1%;
    --accent: 217.2 32.6% 17.5%;
    --accent-foreground: 210 40% 98%;
    --destructive: 0 62.8% 30.6%;
    --destructive-foreground: 210 40% 98%;
    --border: 217.2 32.6% 17.5%;
    --input: 217.2 32.6% 17.5%;
    --ring: 212.7 26.8% 83.9%;
  }
}
 
@layer base {
  * {
    @apply border-border;
  }
  body {
    @apply bg-background text-foreground;
  }
}
```

### File: `frontend/src/layouts/Layout.astro`

```astro
---
import '../styles/globals.css';

interface Props {
  title: string;
}

const { title } = Astro.props;
---

<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="description" content="IntelliMetrics AI Analytics Platform" />
    <meta name="viewport" content="width=device-width" />
    <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
    <title>{title}</title>
  </head>
  <body>
    <slot />
  </body>
</html>
```

### File: `frontend/src/pages/index.astro`

```astro
---
import Layout from '../layouts/Layout.astro';
import Dashboard from '../components/Dashboard';
---

<Layout title="IntelliMetrics - AI Analytics Platform">
  <Dashboard client:load />
</Layout>
```

### File: `frontend/src/lib/utils.ts`

```typescript
import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

---

## Part 3: shadcn/ui Components

You need to manually install shadcn/ui components. Run these commands in the `frontend/` directory:

```bash
npx shadcn-ui@latest init

# When prompted, choose:
# - TypeScript: Yes
# - Style: Default
# - Base color: Slate
# - CSS variables: Yes

# Then add the components:
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add input
npx shadcn-ui@latest add label
npx shadcn-ui@latest add tabs
```

This will create the component files in `src/components/ui/`

---

## Part 4: Docker Compose

### File: `docker-compose.yml`

```yaml
version: '3.8'

services:
  mongo:
    image: mongo:7.0
    container_name: intellimetrics-mongo
    ports:
      # VULNERABILITY: MongoDB exposed on host
      - "27017:27017"
    volumes:
      - mongo_data:/data/db
      - ./seed-database.js:/docker-entrypoint-initdb.d/seed.js:ro
    # VULNERABILITY: No authentication
    networks:
      - app_network

  api:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: intellimetrics-api
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - ENV=development
    depends_on:
      - mongo
    networks:
      - app_network

  frontend:
    build:
      context: ./frontend
    container_name: intellimetrics-frontend
    ports:
      - "4321:4321"
    environment:
      - PUBLIC_API_URL=http://localhost:8080
    depends_on:
      - api
    networks:
      - app_network

networks:
  app_network:
    driver: bridge

volumes:
  mongo_data:
```

### File: `frontend/Dockerfile`

```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

EXPOSE 4321

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
```

---

## Part 5: Seed Database Script

### File: `seed-database.js` (place in project root)

```javascript
db = db.getSiblingDB('intellimetrics');

db.companies.insertMany([
  {
    company_id: "COMP-10001",
    company_name: "TechStartup Inc",
    revenue: 1247000,
    users: 2847,
    growth_rate: 0.23,
    industry: "SaaS",
    created_at: new Date()
  },
  {
    company_id: "COMP-10002",
    company_name: "DataVision Corp",
    revenue: 3420000,
    users: 5120,
    growth_rate: 0.18,
    industry: "Analytics",
    created_at: new Date()
  },
  {
    company_id: "COMP-10003",
    company_name: "CloudScale Solutions",
    revenue: 890000,
    users: 1523,
    growth_rate: 0.45,
    industry: "Cloud Infrastructure",
    created_at: new Date()
  },
  {
    company_id: "COMP-10004",
    company_name: "AI Innovations Ltd",
    revenue: 2150000,
    users: 3890,
    growth_rate: 0.31,
    industry: "Artificial Intelligence",
    created_at: new Date()
  },
  {
    company_id: "COMP-10005",
    company_name: "SecureNet Systems",
    revenue: 1680000,
    users: 2340,
    growth_rate: 0.15,
    industry: "Cybersecurity",
    created_at: new Date()
  },
  {
    company_id: "COMP-10006",
    company_name: "FinTech Global",
    revenue: 4570000,
    users: 8920,
    growth_rate: 0.28,
    industry: "Financial Technology",
    created_at: new Date()
  },
  {
    company_id: "COMP-10007",
    company_name: "HealthTech Partners",
    revenue: 2980000,
    users: 4560,
    growth_rate: 0.22,
    industry: "Healthcare",
    created_at: new Date()
  },
  {
    company_id: "COMP-10008",
    company_name: "EduPlatform Co",
    revenue: 1120000,
    users: 6780,
    growth_rate: 0.35,
    industry: "EdTech",
    created_at: new Date()
  },
  {
    company_id: "COMP-10009",
    company_name: "RetailMetrics Inc",
    revenue: 3240000,
    users: 5450,
    growth_rate: 0.19,
    industry: "E-commerce",
    created_at: new Date()
  },
  {
    company_id: "COMP-10010",
    company_name: "GreenEnergy Analytics",
    revenue: 1890000,
    users: 2120,
    growth_rate: 0.41,
    industry: "Renewable Energy",
    created_at: new Date()
  }
]);

print("Database seeded with 10 companies");
```

---

## Part 6: Setup Instructions

### Prerequisites
- Docker & Docker Compose installed
- Git installed
- Basic understanding of Docker, Go, and Astro

### Step 1: Create Project Structure

```bash
mkdir intellimetrics-demo
cd intellimetrics-demo

mkdir -p backend/cmd/api
mkdir -p frontend/src/{components/ui,layouts,pages,styles}
mkdir -p exploits
```

### Step 2: Copy All Files

Copy all the code blocks above into their respective files according to the file paths shown.

### Step 3: Initialize Frontend

```bash
cd frontend
npm install

# Initialize shadcn/ui
npx shadcn-ui@latest init
# Choose: TypeScript, Default style, Slate color, CSS variables

# Add components
npx shadcn-ui@latest add button card input label tabs

cd ..
```

### Step 4: Initialize Backend

```bash
cd backend
go mod download
cd ..
```

### Step 5: Build and Run

```bash
# Start all services
docker-compose up