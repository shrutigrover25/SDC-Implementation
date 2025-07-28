#!/bin/bash

# 🧪 SCD Backend API Testing Script (No jq)
echo "🔥 SCD Backend API Testing & Database Guide"
echo "============================================"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080"

echo -e "\n${BLUE}📋 Available API Endpoints:${NC}"
echo "🏢 JOBS: /jobs/:uid, POST /jobs, PUT /jobs/:uid, etc."
echo "⏰ TIMELOGS: /timelogs/:uid, POST /timelogs, etc."
echo "💰 PAYMENTS: /payment-line-items/:uid, POST /payment-line-items, etc."

check_api() {
    curl -s "$API_URL/jobs/00000000-0000-0000-0000-000000000003" > /dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ API is running on $API_URL${NC}"
        return 0
    else
        echo -e "${RED}❌ API is not running. Start with: go run cmd/main.go${NC}"
        return 1
    fi
}

test_jobs() {
    echo -e "\n${BLUE}🏢 Testing Job Operations:${NC}"

    echo -e "\n📖 Getting existing job:"
    curl -s "$API_URL/jobs/00000000-0000-0000-0000-000000000003"

    echo -e "\n➕ Creating new job:"
    NEW_JOB=$(curl -s -X POST "$API_URL/jobs" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "API Test Job",
            "status": "active",
            "rate": 55.5,
            "companyId": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
            "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc"
        }')
    echo "$NEW_JOB"

    JOB_UID=$(echo "$NEW_JOB" | grep -oE '"UID":"[^"]+' | cut -d':' -f2 | tr -d '"')
    echo -e "${GREEN}💾 Created job with UID: $JOB_UID${NC}"

    echo -e "\n🔄 Updating job status:"
    curl -s -X PUT "$API_URL/jobs/$JOB_UID/status?status=extended"

    echo -e "\n🔄 Full job update:"
    curl -s -X PUT "$API_URL/jobs/$JOB_UID" \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Updated API Test Job",
            "status": "active",
            "rate": 65.0
        }'
}

test_timelogs() {
    echo -e "\n${BLUE}⏰ Testing Timelog Operations:${NC}"

    echo -e "\n📖 Getting existing timelog:"
    curl -s "$API_URL/timelogs/1c2e2ca7-a69d-421b-b278-f7f83a49e7e5"

    echo -e "\n➕ Creating new timelog:"
    START_TIME=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
    END_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    NEW_TIMELOG=$(curl -s -X POST "$API_URL/timelogs" \
        -H "Content-Type: application/json" \
        -d "{
            \"startTime\": \"$START_TIME\",
            \"endTime\": \"$END_TIME\",
            \"contractorId\": \"cccccccc-cccc-cccc-cccc-cccccccccccc\"
        }")
    echo "$NEW_TIMELOG"

    TIMELOG_UID=$(echo "$NEW_TIMELOG" | grep -oE '"UID":"[^"]+' | cut -d':' -f2 | tr -d '"')
    echo -e "${GREEN}💾 Created timelog with UID: $TIMELOG_UID${NC}"
}

test_payments() {
    echo -e "\n${BLUE}💰 Testing Payment Operations:${NC}"

    echo -e "\n📖 Getting existing payment:"
    curl -s "$API_URL/payment-line-items/de1dbf39-3e6c-4d3b-af19-4447e2c26571"

    echo -e "\n➕ Creating new payment:"
    ISSUED_AT=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    NEW_PAYMENT=$(curl -s -X POST "$API_URL/payment-line-items" \
        -H "Content-Type: application/json" \
        -d "{
            \"contractorId\": \"cccccccc-cccc-cccc-cccc-cccccccccccc\",
            \"amount\": 125.50,
            \"issuedAt\": \"$ISSUED_AT\"
        }")
    echo "$NEW_PAYMENT"

    PAYMENT_UID=$(echo "$NEW_PAYMENT" | grep -oE '"UID":"[^"]+' | cut -d':' -f2 | tr -d '"')
    echo -e "${GREEN}💾 Created payment with UID: $PAYMENT_UID${NC}"
}

test_scd_versioning() {
    echo -e "\n${YELLOW}🔄 Testing SCD Versioning:${NC}"

    echo -e "\nVersion 1:"
    curl -s "$API_URL/jobs/00000000-0000-0000-0000-000000000001"

    echo -e "\nVersion 2:"
    curl -s "$API_URL/jobs/00000000-0000-0000-0000-000000000002"

    echo -e "\nVersion 3:"
    curl -s "$API_URL/jobs/00000000-0000-0000-0000-000000000003"
}

show_database_operations() {
    echo -e "\n${YELLOW}🗄️  Database Operations Guide:${NC}"
    echo "1. Connect to PostgreSQL: psql -U postgres -d mercor"
    echo "2. View schema: \\d jobs; \\d timelogs; \\d payment_line_items"
    echo "3. View data:"
    echo "   SELECT * FROM jobs;"
    echo "   SELECT * FROM timelogs;"
    echo "   SELECT * FROM payment_line_items;"
}

# Help
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Usage: $0 [test-type]"
    echo "  jobs       - Test job APIs"
    echo "  timelogs   - Test timelog APIs"
    echo "  payments   - Test payment APIs"
    echo "  scd        - Test version history"
    echo "  db         - Show DB commands"
    echo "  all        - Run all (default)"
    exit 0
fi

# Check API up
if ! check_api; then
    exit 1
fi

# Run the right tests
case "$1" in
    "jobs") test_jobs ;;
    "timelogs") test_timelogs ;;
    "payments") test_payments ;;
    "scd") test_scd_versioning ;;
    "db") show_database_operations ;;
    *) test_jobs; test_timelogs; test_payments; test_scd_versioning; show_database_operations ;;
esac

echo -e "\n${GREEN}🎉 All tests complete!${NC}"
