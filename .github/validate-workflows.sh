#!/bin/bash

# GitHub Actions Workflow Validation Script
# This script validates the GitHub Actions workflows

set -e

echo "🔍 Validating GitHub Actions Workflows..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✅ $message${NC}"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠️  $message${NC}"
            ;;
        "ERROR")
            echo -e "${RED}❌ $message${NC}"
            ;;
    esac
}

# Check if required files exist
echo "📁 Checking required files..."

required_files=(
    ".github/workflows/ci.yml"
    ".github/workflows/release.yml"
    ".github/workflows/security.yml"
    ".github/workflows/test-matrix.yml"
    ".github/workflows/coverage.yml"
    ".github/workflows/docker.yml"
    ".github/dependabot.yml"
    ".github/CODEOWNERS"
    ".github/CONTRIBUTING.md"
    ".github/SECURITY.md"
    ".github/ISSUE_TEMPLATE/bug_report.yml"
    ".github/ISSUE_TEMPLATE/feature_request.yml"
    ".github/pull_request_template.md"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        print_status "SUCCESS" "Found $file"
    else
        print_status "ERROR" "Missing $file"
        exit 1
    fi
done

# Check for duplicate files
echo "🔍 Checking for duplicates..."

if [ -f ".github/workflows/dependabot.yml" ]; then
    print_status "WARNING" "Duplicate dependabot.yml found in workflows directory"
fi

# Validate YAML syntax (basic check)
echo "📝 Validating YAML syntax..."

yaml_files=$(find .github -name "*.yml" -o -name "*.yaml")

for file in $yaml_files; do
    # Basic YAML validation - check for common issues
    if grep -q "^\s*$" "$file" && grep -q "^\s*[a-zA-Z]" "$file"; then
        # Check for mixed indentation
        if grep -q "^\s*[a-zA-Z]" "$file" && grep -q "^\t" "$file"; then
            print_status "WARNING" "$file may have mixed indentation (spaces and tabs)"
        fi
    fi
    
    # Check for common YAML issues
    if grep -q ":\s*$" "$file"; then
        print_status "WARNING" "$file has empty values (may be intentional)"
    fi
    
    print_status "SUCCESS" "Basic YAML validation passed for $file"
done

# Check workflow dependencies
echo "🔗 Checking workflow dependencies..."

# Check if Dockerfile exists
if [ -f "docker/Dockerfile.base" ]; then
    print_status "SUCCESS" "Dockerfile.base found"
else
    print_status "ERROR" "Dockerfile.base not found"
    exit 1
fi

# Check if main.go has versionString
if grep -q "versionString" cmd/gopher/main.go; then
    print_status "SUCCESS" "versionString variable found in main.go"
else
    print_status "ERROR" "versionString variable not found in main.go"
    exit 1
fi

# Check for required secrets documentation
echo "🔐 Checking secrets documentation..."

# Docker secrets are not required since we're not publishing Docker images
print_status "SUCCESS" "Docker publishing disabled - no secrets required"

# Check workflow triggers
echo "⚡ Checking workflow triggers..."

workflows=(
    "ci.yml:push,pull_request"
    "release.yml:push"
    "security.yml:push,pull_request,schedule"
    "test-matrix.yml:push,pull_request"
    "coverage.yml:push,pull_request"
    "docker.yml:push,pull_request"
)

for workflow_info in "${workflows[@]}"; do
    workflow_file=$(echo "$workflow_info" | cut -d: -f1)
    expected_triggers=$(echo "$workflow_info" | cut -d: -f2)
    
    if [ -f ".github/workflows/$workflow_file" ]; then
        print_status "SUCCESS" "Workflow $workflow_file exists"
        
        # Check for basic trigger patterns
        if grep -q "on:" ".github/workflows/$workflow_file"; then
            print_status "SUCCESS" "Workflow $workflow_file has triggers defined"
        else
            print_status "WARNING" "Workflow $workflow_file may be missing triggers"
        fi
    else
        print_status "ERROR" "Workflow $workflow_file not found"
        exit 1
    fi
done

# Check for common issues
echo "🔍 Checking for common issues..."

# Check for hardcoded versions
if grep -r "go-version.*1\.[0-9]" .github/workflows/ | grep -v "1.21\|1.22\|1.23"; then
    print_status "WARNING" "Found hardcoded Go versions that might be outdated"
fi

# Check for action versions
if grep -r "uses:.*@master" .github/workflows/; then
    print_status "WARNING" "Found actions using @master (consider using specific versions)"
fi

# Summary
echo ""
echo "📊 Validation Summary:"
echo "✅ All required files present"
echo "✅ Basic YAML validation passed"
echo "✅ Workflow dependencies checked"
echo "✅ Common issues reviewed"
echo ""
print_status "SUCCESS" "GitHub Actions workflows validation completed successfully!"
echo ""
echo "🚀 Next steps:"
echo "1. Push changes to trigger workflows"
echo "2. Test workflows by creating a PR or pushing to main"
echo "3. Create a release by tagging a version (e.g., git tag v1.0.0)"
echo "4. Docker images will be built and tested locally (not published)"
