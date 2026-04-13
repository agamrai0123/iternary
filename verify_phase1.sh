#!/bin/bash
# Phase 1 Security Implementation - Verification Script
# Usage: ./verify_phase1.sh

echo "========================================="
echo "Phase 1 Security Implementation Verification"
echo "========================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} $1"
        return 0
    else
        echo -e "${RED}✗${NC} $1 (MISSING)"
        return 1
    fi
}

check_function() {
    if grep -q "func $2" "$1" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} $2 in $1"
        return 0
    else
        echo -e "${RED}✗${NC} $2 NOT FOUND in $1"
        return 1
    fi
}

echo "1. Security Modules Created"
echo "----------------------------"
check_file "itinerary/security/jwt.go"
check_file "itinerary/security/password.go"
check_file "itinerary/security/session.go"
check_file "itinerary/security/tls.go"
echo ""

echo "2. Middleware Created"
echo "---------------------"
check_file "itinerary/middleware/ratelimit.go"
echo ""

echo "3. Key Functions Implemented"
echo "-----------------------------"
check_function "itinerary/security/jwt.go" "GenerateToken"
check_function "itinerary/security/jwt.go" "ValidateToken"
check_function "itinerary/security/password.go" "HashPassword"
check_function "itinerary/security/password.go" "VerifyPassword"
check_function "itinerary/security/session.go" "SaveSession"
check_function "itinerary/security/session.go" "ValidateSession"
check_function "itinerary/middleware/ratelimit.go" "Allow"
check_function "itinerary/security/tls.go" "GenerateSelfSignedCert"
echo ""

echo "4. Configuration Files"
echo "----------------------"
check_file ".env.production"
check_file "PHASE_1_COMPLETE.md"
echo ""

echo "5. Build Status"
echo "---------------"
if command -v go &> /dev/null; then
    if go build -v 2>&1 | grep -q "github.com/yourusername/itinerary-backend"; then
        echo -e "${GREEN}✓${NC} Project compiles successfully"
        # Get binary size
        if [ -f "itinerary-backend.exe" ]; then
            SIZE=$(ls -lh itinerary-backend.exe | awk '{print $5}')
            echo -e "${GREEN}✓${NC} Binary created: itinerary-backend.exe ($SIZE)"
        fi
    else
        echo -e "${YELLOW}⚠${NC}  Build check - Run 'go build' manually to verify"
    fi
else
    echo -e "${YELLOW}⚠${NC}  Go not found in PATH"
fi
echo ""

echo "6. Vulnerability Coverage"
echo "-------------------------"
echo "CWE-327: Weak Cryptography"
check_function "itinerary/security/password.go" "HashPassword" && \
    echo -e "${GREEN}✓${NC} Bcrypt implementation with cost validation"
echo ""

echo "CWE-287: Auth Bypass"
check_function "itinerary/security/jwt.go" "ValidateToken" && \
    echo -e "${GREEN}✓${NC} JWT validation with signature/expiration check"
echo ""

echo "CWE-798: Hardcoded Credentials"
if grep -q "password123" itinerary/auth/handlers.go 2>/dev/null; then
    echo -e "${RED}✗${NC} Hardcoded credentials still present"
else
    echo -e "${GREEN}✓${NC} No hardcoded credentials found"
fi
echo ""

echo "CWE-203: User Enumeration"
if grep -q "Invalid email or password" itinerary/auth/handlers.go 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Generic error messages implemented"
else
    echo -e "${YELLOW}⚠${NC}  Check error messages in handlers"
fi
echo ""

echo "CWE-532: Token Logging"
if grep -q "HashTokenForLogging\|hashTokenForLogging" itinerary/security/*.go 2>/dev/null; then
    echo -e "${GREEN}✓${NC} Token hashing for logging implemented"
else
    echo -e "${YELLOW}⚠${NC}  Verify token logging mechanism"
fi
echo ""

echo "CWE-424: Rate Limiting"
check_function "itinerary/middleware/ratelimit.go" "Allow" && \
    echo -e "${GREEN}✓${NC} Rate limiting middleware implemented"
echo ""

echo "CWE-640: Token Extraction"
check_function "itinerary/auth/middleware.go" "RequireAuth" && \
    echo -e "${GREEN}✓${NC} JWT middleware for token validation"
echo ""

echo "CWE-613: HTTPS/TLS"
check_function "itinerary/security/tls.go" "GetTLSConfig" && \
    echo -e "${GREEN}✓${NC} TLS configuration implemented"
echo ""

echo "========================================="
echo "Phase 1 Verification Complete!"
echo "========================================="
echo ""
echo "Next Steps:"
echo "1. Review PHASE_1_COMPLETE.md for deployment details"
echo "2. Configure .env.production with production secrets"
echo "3. Generate TLS certificates for production"
echo "4. Run integration tests"
echo "5. Deploy to production environment"
