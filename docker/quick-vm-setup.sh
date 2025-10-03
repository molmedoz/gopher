#!/bin/bash

# Quick Windows VM Setup Script for Gopher Testing
# This script helps you get started with Windows VM setup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}================================${NC}"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
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

# Check if running on macOS
if [[ "$OSTYPE" != "darwin"* ]]; then
    print_error "This script is designed for macOS. For other platforms, see VM_SETUP_GUIDE.md"
    exit 1
fi

print_header "Windows VM Setup for Gopher Testing"

echo "This script will help you set up a Windows virtual machine for testing Gopher."
echo ""

# Step 1: Check for VM software
print_step "1. Checking for VM software..."

if command -v vmrun &> /dev/null; then
    print_success "VMware Fusion detected"
    VM_TYPE="vmware"
elif command -v VBoxManage &> /dev/null; then
    print_success "VirtualBox detected"
    VM_TYPE="virtualbox"
elif command -v prlctl &> /dev/null; then
    print_success "Parallels Desktop detected"
    VM_TYPE="parallels"
else
    print_warning "No VM software detected"
    echo ""
    echo "Please install one of the following:"
    echo "  - VMware Fusion: brew install --cask vmware-fusion"
    echo "  - VirtualBox: brew install --cask virtualbox"
    echo "  - Parallels Desktop: https://www.parallels.com/"
    echo ""
    read -p "Press Enter to continue after installing VM software..."
fi

# Step 2: Build Windows binary
print_step "2. Building Windows binary..."

if [ ! -f "build/gopher-windows-amd64.exe" ]; then
    print_warning "Windows binary not found. Building..."
    mkdir -p build
    GOOS=windows GOARCH=amd64 go build -o build/gopher-windows-amd64.exe cmd/gopher/main.go
    print_success "Windows binary built: build/gopher-windows-amd64.exe"
else
    print_success "Windows binary already exists: build/gopher-windows-amd64.exe"
fi

# Step 3: Create setup files
print_step "3. Preparing setup files..."

# Copy test scripts to build directory
cp docker/test-windows-native.bat build/
cp docker/setup-windows-vm.ps1 build/
cp docker/test-gopher.ps1 build/

print_success "Setup files prepared in build/ directory"

# Step 4: Instructions
print_step "4. Next steps:"

echo ""
print_header "Windows VM Setup Instructions"
echo ""
echo "1. Download Windows 11 ISO:"
echo "   https://www.microsoft.com/en-us/software-download/windows11"
echo ""
echo "2. Create Windows VM:"
if [ "$VM_TYPE" = "vmware" ]; then
    echo "   - Open VMware Fusion"
    echo "   - Create New Virtual Machine"
    echo "   - Select Windows ISO"
    echo "   - Configure: 8GB RAM, 100GB storage, 2 CPU cores"
elif [ "$VM_TYPE" = "virtualbox" ]; then
    echo "   - Open VirtualBox"
    echo "   - Click 'New'"
    echo "   - Configure: Windows 11, 4GB RAM, 60GB storage"
elif [ "$VM_TYPE" = "parallels" ]; then
    echo "   - Open Parallels Desktop"
    echo "   - Install Windows"
    echo "   - Configure: 8GB RAM, 100GB storage"
else
    echo "   - Follow instructions in VM_SETUP_GUIDE.md"
fi
echo ""
echo "3. Install Windows:"
echo "   - Follow standard Windows installation"
echo "   - Create local account (easier for testing)"
echo "   - Install Go from: https://golang.org/dl/"
echo ""
echo "4. Set up Gopher in VM:"
echo "   - Copy build/ directory to Windows VM"
echo "   - Run: setup-windows-vm.ps1 (as Administrator)"
echo "   - Or manually: copy gopher-windows-amd64.exe to C:\\gopher\\gopher.exe"
echo ""
echo "5. Test Gopher:"
echo "   - Run: test-windows-native.bat"
echo "   - Or run: test-gopher.ps1"
echo ""

# Step 5: Create quick start script
print_step "5. Creating quick start script..."

cat > build/quick-start-windows.bat << 'EOF'
@echo off
echo Quick Start for Gopher Testing on Windows
echo ========================================
echo.

echo 1. Copy gopher-windows-amd64.exe to C:\gopher\gopher.exe
echo 2. Add C:\gopher to PATH
echo 3. Run: gopher version
echo 4. Run: test-windows-native.bat
echo.

echo Files in this directory:
dir /b
echo.

echo Ready for Gopher testing!
pause
EOF

print_success "Quick start script created: build/quick-start-windows.bat"

# Step 6: Summary
print_header "Setup Complete!"
echo ""
print_success "Windows binary built: build/gopher-windows-amd64.exe"
print_success "Test scripts ready: build/test-windows-native.bat"
print_success "Setup script ready: build/setup-windows-vm.ps1"
print_success "Quick start guide: build/quick-start-windows.bat"
echo ""
echo "Next steps:"
echo "1. Set up Windows VM (see instructions above)"
echo "2. Copy build/ directory to Windows VM"
echo "3. Run setup-windows-vm.ps1 in Windows VM"
echo "4. Test with test-windows-native.bat"
echo ""
echo "For detailed instructions, see: docker/VM_SETUP_GUIDE.md"
echo ""
print_success "Ready for Windows VM testing! 🚀"
