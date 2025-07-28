#!/bin/bash

# 🧪 SCD Backend API Testing Script
# Complete guide for testing database and API functionality

echo "🔥 SCD Backend API Testing & Database Guide"
echo "============================================"

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# API Base URL
API_URL="http://localhost:8080"

echo -e "\n${BLUE}📋 Available API Endpoints:${NC}"
echo "=============================="
echo "🏢 JOBS:"
echo "  GET    /jobs/:uid                     - Get job by UID"
echo "  POST   /jobs                         - Create new job"
echo "  PUT    /jobs/:uid                    - Update job (creates new version)"
echo "  PUT    /jobs/:uid/status?status=X    - Update status only"
echo "  GET    /companies/:id/jobs           - Get active jobs for company"
echo ""
echo "⏰ TIMELOGS:"
echo "  GET    /timelogs/:uid                - Get timelog by UID"
echo "  POST   /timelogs                     - Create new timelog"
echo "  PUT    /timelogs/:uid                - Update timelog (creates new version)"
echo "  GET    /jobs/:job_uid/timelogs       - Get timelogs for a job"
echo ""
echo "💰 PAYMENT LINE ITEMS:"
echo "  GET    /payment-line-items/:uid              - Get payment by UID"
echo "  POST   /payment-line-items                   - Create new payment"
echo "  PUT    /payment-line-items/:uid              - Update payment (creates new version)"
echo "  GET    /timelogs/:uid/payment-line-items     - Get payments for timelog"
echo "  GET    /jobs/:uid/payment-history            - Get payment history for job"

echo -e "\n${YELLOW}🗄️  DATABASE OPERATIONS:${NC}"
echo "======================="

check_api() {
    if curl -s $API_URL/jobs/00000000-0000-0000-0000-000000000003 > /dev/null; then
        echo -e "${GREEN}✅ API is running on $API_URL${NC}"
        return 0
    else
        echo -e "${RED}❌ API is not running. Start with: go run cmd/main.go${NC}"
        return 1
    fi
}

test_jobs() {
    echo -e "\n${BLUE}🏢 Testing Job Operations:${NC}"
    echo "=========================="
    
    echo "📖 Getting existing job (v3):"
    curl -s $API_URL/jobs/00000000-0000-0000-0000-000000000003 | jq '.'

    echo -e "\n➕ Creating new job:"
    NEW_JOB=$(curl -s -X POST $API_URL/jobs \
        -H "Content-Type: application/json" \
        -d '{
            "title": "API Test Job",
            "status": "active",
            "rate": 55.5,
            "companyId": "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
            "contractorId": "cccccccc-cccc-cccc-cccc-cccccccccccc"
        }')
    echo $NEW_JOB | jq '.'

    JOB_UID=$(echo $NEW_JOB | jq -r '.UID')
    echo -e "${GREEN}💾 Created job with UID: $JOB_UID${NC}"

    echo -e "\n🔄 Updating job status:"
    curl -s -X PUT "$API_URL/jobs/$JOB_UID/status?status=extended" | jq '.'

    echo -e "\n🔄 Full job update (creates new version):"
    curl -s -X PUT $API_URL/jobs/$JOB_UID \
        -H "Content-Type: application/json" \
        -d '{
            "title": "Updated API Test Job",
            "status": "active",
            "rate": 65.0
        }' | jq '.'
}

test_timelogs() {
    echo -e "\n${BLUE}⏰ Testing Timelog Operations:${NC}"
    echo "============================="

    echo "📖 Getting existing timelog:"
    curl -s $API_URL/timelogs/1c2e2ca7-a69d-421b-b278-f7f83a49e7e5 | jq '.'

    echo -e "\n➕ Creating new timelog:"
    START_TIME=$(date -u -d '1 hour ago' '+%Y-%m-%dT%H:%M:%SZ')
    END_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    NEW_TIMELOG=$(curl -s -X POST $API_URL/timelogs \
        -H "Content-Type: application/json" \
        -d "{
            \"startTime\": \"$START_TIME\",
            \"endTime\": \"$END_TIME\",
            \"contractorId\": \"cccccccc-cccc-cccc-cccc-cccccccccccc\"
        }")
    echo $NEW_TIMELOG | jq '.'

    TIMELOG_UID=$(echo $NEW_TIMELOG | jq -r '.UID')
    echo -e "${GREEN}💾 Created timelog with UID: $TIMELOG_UID${NC}"
}

test_payments() {
    echo -e "\n${BLUE}💰 Testing Payment Operations:${NC}"
    echo "=============================="

    echo "📖 Getting existing payment:"
    curl -s $API_URL/payment-line-items/de1dbf39-3e6c-4d3b-af19-4447e2c26571 | jq '.'

    echo -e "\n➕ Creating new payment:"
    ISSUED_AT=$(date -u '+%Y-%m-%dT%H:%M:%SZ')

    NEW_PAYMENT=$(curl -s -X POST $API_URL/payment-line-items \
        -H "Content-Type: application/json" \
        -d "{
            \"contractorId\": \"cccccccc-cccc-cccc-cccc-cccccccccccc\",
            \"amount\": 125.50,
            \"issuedAt\": \"$ISSUED_AT\"
        }")
    echo $NEW_PAYMENT | jq '.'

    PAYMENT_UID=$(echo $NEW_PAYMENT | jq -r '.UID')
    echo -e "${GREEN}💾 Created payment with UID: $PAYMENT_UID${NC}"
}

test_scd_versioning() {
    echo -e "\n${YELLOW}🔄 Testing SCD Versioning:${NC}"
    echo "========================="

    echo "📊 Getting all versions of the same job (ID: aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa):"
    echo "Version 1 (extended, rate: 20.0):"
    curl -s $API_URL/jobs/00000000-0000-0000-0000-000000000001 | jq '{Version, Status, Rate, UID}'

    echo -e "\nVersion 2 (active, rate: 20.0):"
    curl -s $API_URL/jobs/00000000-0000-0000-0000-000000000002 | jq '{Version, Status, Rate, UID}'

    echo -e "\nVersion 3 (active, rate: 15.5):"
    curl -s $API_URL/jobs/00000000-0000-0000-0000-000000000003 | jq '{Version, Status, Rate, UID}'
}

show_database_operations() {
    echo -e "\n${YELLOW}🗄️  Database Operations Guide:${NC}"
    echo "=============================="
    echo "1. 🔗 Connect to PostgreSQL:"
    echo "   sudo -u postgres psql -d mercor"
    echo ""
    echo "2. 📋 View table schemas:"
    echo "   \\d jobs"
    echo "   \\d timelogs" 
    echo "   \\d payment_line_items"
    echo ""
    echo "3. 📊 Query data:"
    echo "   SELECT id, uid, version, status, rate FROM jobs ORDER BY version;"
    echo "   SELECT id, uid, version, contractor_id FROM timelogs;"
    echo "   SELECT id, uid, version, amount FROM payment_line_items;"
    echo ""
    echo "4. 🔍 Understanding SCD:"
    echo "   - Same ID = Same logical entity"
    echo "   - Different UID = Different version"
    echo "   - Version number increments with each update"
    echo ""
    echo "5. ➕ Manual data insertion example:"
    echo "   INSERT INTO jobs (id, uid, version, status, rate, title, company_id, contractor_id)"
    echo "   VALUES ("
    echo "     gen_random_uuid(),"
    echo "     gen_random_uuid(),"
    echo "     1,"
    echo "     'active',"
    echo "     45.0,"
    echo "     'Manual Test Job',"
    echo "     'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',"
    echo "     'cccccccc-cccc-cccc-cccc-cccccccccccc'"
    echo "   );"
}

if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Usage: $0 [test-type]"
    echo "Available test types:"
    echo "  jobs       - Test job operations"
    echo "  timelogs   - Test timelog operations" 
    echo "  payments   - Test payment operations"
    echo "  scd        - Test SCD versioning"
    echo "  db         - Show database operations"
    echo "  all        - Run all tests (default)"
    exit 0
fi

if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}⚠️  Installing jq for JSON formatting...${NC}"
    sudo apt update && sudo apt install -y jq
fi

if ! check_api; then
    exit 1
fi

case "$1" in
    "jobs")
        test_jobs
        ;;
    "timelogs")
        test_timelogs
        ;;
    "payments") 
        test_payments
        ;;
    "scd")
        test_scd_versioning
        ;;
    "db")
        show_database_operations
        ;;
    *)
        test_jobs
        test_timelogs  
        test_payments
        test_scd_versioning
        show_database_operations
        ;;
esac

echo -e "\n${GREEN}🎉 Testing Complete!${NC}"
echo -e "${BLUE}💡 Pro Tips:${NC}"
echo "- Use 'curl -X POST' for creating data"
echo "- Use 'curl -X PUT' for updating (creates new version)"
echo "- Check database with: sudo -u postgres psql -d mercor"
echo "- View logs: tail -f /var/lib/postgresql/17/main/postgresql.log"
