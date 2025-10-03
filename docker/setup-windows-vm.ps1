# Gopher Windows VM Setup Script
# This script sets up a Windows VM for testing Gopher

Write-Host "🪟 Setting up Windows VM for Gopher Testing" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan
Write-Host ""

# Check if running as Administrator
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")

if (-not $isAdmin) {
    Write-Host "⚠️  WARNING: Not running as Administrator" -ForegroundColor Yellow
    Write-Host "Some features may not work properly (like symlinks)" -ForegroundColor Yellow
    Write-Host "Consider running PowerShell as Administrator for full functionality" -ForegroundColor Yellow
    Write-Host ""
}

# Step 1: Install Go (if not already installed)
Write-Host "1. Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✅ Go is already installed: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go not found. Installing Go..." -ForegroundColor Red
    
    # Download Go
    $goUrl = "https://go.dev/dl/go1.21.5.windows-amd64.msi"
    $goInstaller = "$env:TEMP\go-installer.msi"
    
    Write-Host "Downloading Go installer..." -ForegroundColor Yellow
    Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller
    
    Write-Host "Installing Go..." -ForegroundColor Yellow
    Start-Process msiexec.exe -Wait -ArgumentList "/i $goInstaller /quiet"
    
    # Refresh PATH
    $env:PATH = [System.Environment]::GetEnvironmentVariable("PATH", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("PATH", "User")
    
    Write-Host "✅ Go installed successfully" -ForegroundColor Green
}

# Step 2: Create Gopher directory
Write-Host "2. Setting up Gopher directory..." -ForegroundColor Yellow
$gopherDir = "C:\gopher"
if (-not (Test-Path $gopherDir)) {
    New-Item -ItemType Directory -Path $gopherDir -Force
    Write-Host "✅ Created Gopher directory: $gopherDir" -ForegroundColor Green
} else {
    Write-Host "✅ Gopher directory already exists: $gopherDir" -ForegroundColor Green
}

# Step 3: Download Gopher binary (if not present)
Write-Host "3. Checking Gopher binary..." -ForegroundColor Yellow
$gopherExe = "$gopherDir\gopher.exe"
if (-not (Test-Path $gopherExe)) {
    Write-Host "❌ Gopher binary not found at: $gopherExe" -ForegroundColor Red
    Write-Host "Please copy gopher-windows-amd64.exe to $gopherExe" -ForegroundColor Yellow
    Write-Host "You can build it with: GOOS=windows GOARCH=amd64 go build -o gopher-windows-amd64.exe cmd/gopher/main.go" -ForegroundColor Yellow
} else {
    Write-Host "✅ Gopher binary found: $gopherExe" -ForegroundColor Green
}

# Step 4: Add Gopher to PATH
Write-Host "4. Adding Gopher to PATH..." -ForegroundColor Yellow
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($currentPath -notlike "*$gopherDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$gopherDir", "User")
    Write-Host "✅ Added Gopher to user PATH" -ForegroundColor Green
} else {
    Write-Host "✅ Gopher already in PATH" -ForegroundColor Green
}

# Step 5: Test Gopher
Write-Host "5. Testing Gopher..." -ForegroundColor Yellow
if (Test-Path $gopherExe) {
    try {
        $gopherVersion = & $gopherExe version
        Write-Host "✅ Gopher is working: $gopherVersion" -ForegroundColor Green
    } catch {
        Write-Host "❌ Gopher test failed: $_" -ForegroundColor Red
    }
} else {
    Write-Host "⚠️  Cannot test Gopher - binary not found" -ForegroundColor Yellow
}

# Step 6: Create test script
Write-Host "6. Creating test script..." -ForegroundColor Yellow
$testScript = "$gopherDir\test-gopher.bat"
@"
@echo off
echo Testing Gopher on Windows VM...
echo.

echo 1. Gopher version:
gopher version
echo.

echo 2. System Go detection:
gopher system
echo.

echo 3. List versions:
gopher list
echo.

echo 4. Current version:
gopher current
echo.

echo 5. Go version:
go version
echo.

echo 6. Install Go 1.20.7:
gopher install 1.20.7
echo.

echo 7. List after installation:
gopher list
echo.

echo 8. Switch to 1.20.7:
gopher use 1.20.7
echo.

echo 9. Current after switch:
gopher current
echo.

echo 10. Go version after switch:
go version
echo.

echo 11. Switch back to system:
gopher use system
echo.

echo 12. Current after switch back:
gopher current
echo.

echo 13. Go version after switch back:
go version
echo.

echo Testing completed!
pause
"@ | Out-File -FilePath $testScript -Encoding ASCII

Write-Host "✅ Created test script: $testScript" -ForegroundColor Green

# Step 7: Create PowerShell test script
Write-Host "7. Creating PowerShell test script..." -ForegroundColor Yellow
$psTestScript = "$gopherDir\test-gopher.ps1"
Copy-Item "test-gopher.ps1" -Destination $psTestScript -ErrorAction SilentlyContinue
if (Test-Path $psTestScript) {
    Write-Host "✅ PowerShell test script ready: $psTestScript" -ForegroundColor Green
} else {
    Write-Host "⚠️  PowerShell test script not found - using basic test" -ForegroundColor Yellow
}

# Step 8: Summary
Write-Host ""
Write-Host "🎉 Windows VM Setup Complete!" -ForegroundColor Green
Write-Host "=============================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "1. Copy gopher-windows-amd64.exe to: $gopherExe" -ForegroundColor White
Write-Host "2. Restart PowerShell to refresh PATH" -ForegroundColor White
Write-Host "3. Run test: $testScript" -ForegroundColor White
Write-Host "4. Or run PowerShell test: $psTestScript" -ForegroundColor White
Write-Host ""
Write-Host "Build command for Gopher binary:" -ForegroundColor Cyan
Write-Host "GOOS=windows GOARCH=amd64 go build -o gopher-windows-amd64.exe cmd/gopher/main.go" -ForegroundColor White
Write-Host ""
Write-Host "VM is ready for Gopher testing! 🚀" -ForegroundColor Green
