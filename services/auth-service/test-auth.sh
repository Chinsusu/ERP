#!/bin/bash

# Auth Service - Quick Test Script
# This script tests all auth endpoints

set -e

BASE_URL="${BASE_URL:-http://localhost:8081}"
EMAIL="admin@company.vn"
PASSWORD="Admin@123"

echo "========================================="
echo "   Auth Service - Quick Test Suite"
echo "========================================="
echo "Base URL: $BASE_URL"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
PASSED=0
FAILED=0

test_endpoint() {
    local name="$1"
    local expected_code="$2"
    shift 2
    
    echo -n "Testing: $name... "
    
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$@")
    
    if [ "$HTTP_CODE" == "$expected_code" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $HTTP_CODE)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected $expected_code, got $HTTP_CODE)"
        ((FAILED++))
        return 1
    fi
}

# Test 1: Health Check
test_endpoint "Health Check" "200" \
    -X GET "$BASE_URL/health"

# Test 2: Readiness Check
test_endpoint "Readiness Check" "200" \
    -X GET "$BASE_URL/ready"

# Test 3: Login with valid credentials
echo ""
echo -n "Testing: Login (valid credentials)... "
RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")

if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo -e "${GREEN}✓ PASS${NC}"
    ((PASSED++))
    
    # Extract tokens
    ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
    REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r '.data.refresh_token')
    
    echo "  Access Token: ${ACCESS_TOKEN:0:20}..."
    echo "  Refresh Token: ${REFRESH_TOKEN:0:20}..."
else
    echo -e "${RED}✗ FAIL${NC}"
    ((FAILED++))
    echo "Response: $RESPONSE"
fi

# Test 4: Login with invalid credentials
test_endpoint "Login (invalid credentials)" "401" \
    -X POST "$BASE_URL/api/v1/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@company.vn","password":"WrongPassword"}'

# Test 5: Refresh Token
if [ -n "$REFRESH_TOKEN" ]; then
    echo ""
    echo -n "Testing: Refresh Token... "
    REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
    
    if echo "$REFRESH_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
        
        # Update tokens
        ACCESS_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.data.access_token')
        NEW_REFRESH_TOKEN=$(echo "$REFRESH_RESPONSE" | jq -r '.data.refresh_token')
        echo "  New Access Token: ${ACCESS_TOKEN:0:20}..."
    else
        echo -e "${RED}✗ FAIL${NC}"
        ((FAILED++))
    fi
fi

# Test 6: Logout
if [ -n "$ACCESS_TOKEN" ] && [ -n "$NEW_REFRESH_TOKEN" ]; then
    echo ""
    echo -n "Testing: Logout... "
    LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"refresh_token\":\"$NEW_REFRESH_TOKEN\"}")
    
    if echo "$LOGOUT_RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
        echo -e "${GREEN}✓ PASS${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}"
        ((FAILED++))
    fi
fi

# Summary
echo ""
echo "========================================="
echo "   Test Results"
echo "========================================="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"
echo "Total:  $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some tests failed${NC}"
    exit 1
fi
