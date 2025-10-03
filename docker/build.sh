#!/bin/bash

# Gopher Docker Build Script
# This script builds all Docker images for testing Gopher in different scenarios

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}================================${NC}"
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Build base image
build_base() {
    print_header "Building Base Image"
    print_status "Building molmedoz/gopher:base..."
    
    if docker build -f docker/Dockerfile.base -t molmedoz/gopher:base .; then
        print_success "Base image built successfully"
    else
        print_error "Failed to build base image"
        exit 1
    fi
}

# Build scenario images
build_scenarios() {
    local scenarios=(
        "unix-no-go:Unix OS, no system Go installed"
        "unix-with-go:Unix OS, Go already installed"
        "windows-no-go:Windows, no Go installed"
        "windows-with-go:Windows, Go already installed"
        "windows-simulated:Windows simulation container"
        "macos-no-go:macOS, no Go installed"
        "macos-with-go:macOS, Go already installed"
    )
    
    for scenario in "${scenarios[@]}"; do
        IFS=':' read -r name description <<< "$scenario"
        
        print_header "Building $description"
        print_status "Building molmedoz/gopher:$name..."
        
        if docker build -f "docker/Dockerfile.$name" -t "molmedoz/gopher:$name" .; then
            print_success "$description built successfully"
        else
            print_error "Failed to build $description"
            exit 1
        fi
    done
}

# List built images
list_images() {
    print_header "Built Images"
    docker images | grep molmedoz/gopher
}

# Run tests for a specific scenario
run_test() {
    local scenario=$1
    local description=$2
    
    print_header "Testing $description"
    print_status "Running container for $scenario..."
    
    if docker run --rm "molmedoz/gopher:$scenario"; then
        print_success "$description test completed successfully"
    else
        print_error "$description test failed"
        return 1
    fi
}

# Run all tests
run_all_tests() {
    local scenarios=(
        "unix-no-go:Unix OS, no system Go installed"
        "unix-with-go:Unix OS, Go already installed"
        "macos-no-go:macOS, no Go installed"
        "macos-with-go:macOS, Go already installed"
    )
    
    local failed_tests=()
    
    for scenario in "${scenarios[@]}"; do
        IFS=':' read -r name description <<< "$scenario"
        
        if ! run_test "$name" "$description"; then
            failed_tests+=("$description")
        fi
    done
    
    if [ ${#failed_tests[@]} -eq 0 ]; then
        print_success "All tests passed!"
    else
        print_error "The following tests failed:"
        for test in "${failed_tests[@]}"; do
            echo "  - $test"
        done
        exit 1
    fi
}

# Clean up images
cleanup() {
    print_header "Cleaning Up"
    print_status "Removing Gopher images..."
    
    docker images | grep molmedoz/gopher | awk '{print $3}' | xargs -r docker rmi -f
    
    print_success "Cleanup completed"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  build       Build all Docker images (default)"
    echo "  test        Run tests for all scenarios"
    echo "  test-all    Run tests for all scenarios"
    echo "  test-<name> Run test for specific scenario"
    echo "  list        List built images"
    echo "  clean       Clean up all images"
    echo "  help        Show this help message"
    echo ""
    echo "Scenarios:"
    echo "  unix-no-go      Unix OS, no system Go installed"
    echo "  unix-with-go    Unix OS, Go already installed"
    echo "  windows-no-go   Windows, no Go installed"
    echo "  windows-with-go Windows, Go already installed"
    echo "  windows-simulated Windows simulation container"
    echo "  macos-no-go     macOS, no Go installed"
    echo "  macos-with-go   macOS, Go already installed"
}

# Main function
main() {
    local command=${1:-build}
    
    case $command in
        "build")
            check_docker
            build_base
            build_scenarios
            list_images
            print_success "All images built successfully!"
            ;;
        "test")
            check_docker
            run_all_tests
            ;;
        "test-all")
            check_docker
            run_all_tests
            ;;
        "test-unix-no-go")
            check_docker
            run_test "unix-no-go" "Unix OS, no system Go installed"
            ;;
        "test-unix-with-go")
            check_docker
            run_test "unix-with-go" "Unix OS, Go already installed"
            ;;
        "test-windows-no-go")
            check_docker
            run_test "windows-no-go" "Windows, no Go installed"
            ;;
        "test-windows-with-go")
            check_docker
            run_test "windows-with-go" "Windows, Go already installed"
            ;;
        "test-windows-simulated")
            check_docker
            run_test "windows-simulated" "Windows simulation container"
            ;;
        "test-macos-no-go")
            check_docker
            run_test "macos-no-go" "macOS, no Go installed"
            ;;
        "test-macos-with-go")
            check_docker
            run_test "macos-with-go" "macOS, Go already installed"
            ;;
        "list")
            list_images
            ;;
        "clean")
            cleanup
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            print_error "Unknown command: $command"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
