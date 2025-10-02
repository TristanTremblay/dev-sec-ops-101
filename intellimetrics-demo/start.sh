#!/bin/bash

echo "================================================"
echo "  IntelliMetrics CTF - Starting Services"
echo "================================================"
echo ""

echo "Building and starting all services..."
echo ""

docker-compose down -v 2>/dev/null
docker-compose build
docker-compose up -d

echo ""
echo "Waiting for services to start..."
sleep 10

echo ""
echo "================================================"
echo "  SERVICE STATUS"
echo "================================================"
echo ""

docker-compose ps

echo ""
echo "================================================"
echo "  AVAILABLE ENDPOINTS"
echo "================================================"
echo ""
echo "Frontend:         http://localhost:4321"
echo "Admin Panel:      http://localhost:4321/admin"
echo ""
echo "API Gateway:      http://localhost:8080"
echo "  Health:         http://localhost:8080/health"
echo "  Config:         http://localhost:8080/api/admin/config"
echo "  Debug:          http://localhost:8080/api/admin/debug"
echo "  Command Exec:   http://localhost:8080/api/admin/system?cmd=whoami"
echo ""
echo "Auth Service:     http://localhost:8081"
echo "  Users:          http://localhost:8081/api/auth/users"
echo "  Hashes:         http://localhost:8081/api/auth/debug/hashes"
echo "  SAM Dump:       http://localhost:8081/api/auth/sam-dump"
echo ""
echo "File Service:     http://localhost:8082"
echo "  List Files:     http://localhost:8082/api/files/list"
echo "  Read File (LFI): http://localhost:8082/api/files/read?file=/etc/passwd"
echo ""
echo "MongoDB:          mongodb://localhost:27017"
echo "SSH Server:       ssh admin@localhost -p 2222  (password: Admin123!)"
echo ""
echo "================================================"
echo "  Quick Tests"
echo "================================================"
echo ""

echo "Testing API Gateway..."
curl -s http://localhost:8080/health | jq '.' 2>/dev/null || echo "  Not ready yet"
echo ""

echo "Testing Auth Service..."
curl -s http://localhost:8081/api/auth/users | jq '.' 2>/dev/null || echo "  Not ready yet"
echo ""

echo "Testing File Service..."
curl -s http://localhost:8082/api/files/list | jq '.' 2>/dev/null || echo "  Not ready yet"
echo ""

echo "================================================"
echo ""
echo "All services should be running!"
echo "Open http://localhost:4321 in your browser"
echo ""
echo "To stop: docker-compose down"
echo "To view logs: docker-compose logs -f [service-name]"
echo ""
