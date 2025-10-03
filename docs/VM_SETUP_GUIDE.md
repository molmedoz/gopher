# Virtual Machine Setup Guide for Gopher Testing
**Complete guide for setting up Windows, Linux, and macOS VMs**

---

## 🎯 Overview

This guide covers multiple VM solutions, from simple to advanced, to help you test Gopher on all platforms.

### **Quick Recommendation by Platform**

| Your Host OS | Best for Windows VM | Best for Linux VM | Best for macOS VM |
|--------------|-------------------|------------------|------------------|
| **macOS** | UTM (free) or Parallels | UTM or Multipass | Not needed (native) |
| **Windows** | Native testing | WSL2 or VirtualBox | Not legally possible |
| **Linux** | VirtualBox + Wine | Native testing | Not legally possible |

---

## 📋 Table of Contents

1. [VM Software Options](#vm-software-options)
2. [Windows VM Setup](#windows-vm-setup)
3. [Linux VM Setup](#linux-vm-setup)
4. [macOS VM Setup](#macos-vm-setup)
5. [Cloud VM Alternatives](#cloud-vm-alternatives)
6. [Quick Testing with Docker](#quick-testing-with-docker)
7. [Troubleshooting](#troubleshooting)

---

## 🔧 VM Software Options

### **Free Options**

#### **1. UTM (macOS - RECOMMENDED)**
- **Best for:** macOS users wanting Windows/Linux VMs
- **Cost:** Free & Open Source
- **Pros:** 
  - Native Apple Silicon support
  - Easy to use GUI
  - Good performance
  - Built on QEMU
- **Cons:** macOS only
- **Download:** https://mac.getutm.app/

#### **2. VirtualBox (All Platforms)**
- **Best for:** Cross-platform testing
- **Cost:** Free & Open Source
- **Pros:**
  - Works on Windows, Linux, macOS
  - Well documented
  - Large community
- **Cons:**
  - Apple Silicon support limited
  - Slower than native solutions
- **Download:** https://www.virtualbox.org/

#### **3. Multipass (All Platforms)**
- **Best for:** Quick Ubuntu VMs
- **Cost:** Free
- **Pros:**
  - Extremely fast setup
  - Cloud-init support
  - CLI focused
  - Canonical supported
- **Cons:**
  - Ubuntu only
  - Limited customization
- **Download:** https://multipass.run/

#### **4. WSL2 (Windows)**
- **Best for:** Linux testing on Windows
- **Cost:** Free (built into Windows)
- **Pros:**
  - Native integration
  - Fast
  - Easy to use
- **Cons:**
  - Windows only
  - Not a full VM
- **Setup:** Built into Windows 10/11

#### **5. QEMU (All Platforms)**
- **Best for:** Advanced users
- **Cost:** Free & Open Source
- **Pros:**
  - Very flexible
  - Good performance
  - Works everywhere
- **Cons:**
  - Command-line only
  - Steep learning curve
- **Download:** https://www.qemu.org/

### **Paid Options**

#### **1. Parallels Desktop (macOS)**
- **Best for:** Best macOS VM experience
- **Cost:** $99/year (or $129 perpetual)
- **Pros:**
  - Excellent performance
  - Great integration
  - Easy to use
  - Official Windows support
- **Cons:**
  - Expensive
  - macOS only
- **Download:** https://www.parallels.com/

#### **2. VMware Workstation/Fusion**
- **Best for:** Professional use
- **Cost:** ~$200 (or free personal license)
- **Pros:**
  - Excellent performance
  - Professional features
  - Good support
- **Cons:**
  - Expensive (unless free license)
- **Download:** https://www.vmware.com/

---

## 🪟 Windows VM Setup

### **Option 1: UTM on macOS (Recommended for Mac users)**

#### **Step 1: Install UTM**
```bash
# Install via Homebrew
brew install --cask utm

# Or download from https://mac.getutm.app/
```

#### **Step 2: Download Windows 11 ISO**
1. Visit: https://www.microsoft.com/software-download/windows11
2. Select "Download Windows 11 Disk Image (ISO)"
3. Choose "Windows 11 (multi-edition ISO)"
4. Select language and download (6+ GB)

Or use the free Windows 11 Development Environment:
- Visit: https://developer.microsoft.com/en-us/windows/downloads/virtual-machines/
- Download the UTM/QEMU image
- This includes Visual Studio and expires after 90 days

#### **Step 3: Create VM in UTM**

1. **Open UTM** and click "Create a New Virtual Machine"

2. **Select Virtualize** (not Emulate)

3. **Operating System:**
   - Select "Windows"
   - Click "Browse" and select your Windows 11 ISO

4. **Hardware Configuration:**
   - **Memory:** 4096 MB (4 GB minimum, 8 GB recommended)
   - **CPU Cores:** 2-4 cores
   - **Storage:** 64 GB minimum

5. **Summary:**
   - Name: "Windows 11 - Gopher Testing"
   - Click "Save"

6. **Install Windows:**
   - Click "Play" to start the VM
   - Follow Windows installation wizard
   - Choose "I don't have a product key" (can use trial version)
   - Select "Windows 11 Pro"
   - Choose "Custom: Install Windows only"
   - Wait for installation (15-30 minutes)

7. **Post-Installation:**
   ```powershell
   # In Windows VM, open PowerShell as Administrator
   
   # Install Chocolatey (package manager)
   Set-ExecutionPolicy Bypass -Scope Process -Force
   [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
   iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
   
   # Install Git
   choco install git -y
   
   # Install Go (for testing)
   choco install golang -y
   
   # Refresh environment
   refreshenv
   ```

8. **Transfer Gopher to VM:**
   ```bash
   # On your Mac, in gopher directory
   make build-all
   
   # The Windows binary is at: build/gopher-windows-amd64.exe
   # Copy it to a shared folder or use GitHub
   
   # Option A: Via GitHub
   git push
   # Then clone in Windows VM
   
   # Option B: Via UTM shared folder
   # In UTM: Settings → Sharing → Enable Directory Sharing
   ```

---

### **Option 2: VirtualBox (Cross-Platform)**

#### **Step 1: Install VirtualBox**
```bash
# On macOS
brew install --cask virtualbox

# On Linux (Ubuntu/Debian)
sudo apt update
sudo apt install virtualbox

# On Linux (Fedora)
sudo dnf install VirtualBox

# Or download from: https://www.virtualbox.org/
```

#### **Step 2: Download Windows ISO**
Same as UTM Option 1, Step 2

#### **Step 3: Create VM**

1. **Open VirtualBox** → Click "New"

2. **Name and Operating System:**
   - Name: "Windows 11 - Gopher Testing"
   - Type: Microsoft Windows
   - Version: Windows 11 (64-bit)
   - Click "Next"

3. **Memory Size:**
   - Set to 4096 MB (4 GB) or more
   - Click "Next"

4. **Hard Disk:**
   - Select "Create a virtual hard disk now"
   - Click "Create"

5. **Hard Disk File Type:**
   - Select "VDI (VirtualBox Disk Image)"
   - Click "Next"

6. **Storage:**
   - Select "Dynamically allocated"
   - Click "Next"

7. **Size:**
   - Set to 64 GB minimum
   - Click "Create"

8. **Settings (before starting):**
   - Click "Settings"
   - **System → Processor:** Set to 2-4 cores
   - **Display → Video Memory:** Set to 128 MB
   - **Storage:** Click "Empty" → Click disk icon → "Choose a disk file" → Select Windows ISO
   - Click "OK"

9. **Start VM** and follow Windows installation

10. **Install Guest Additions** (after Windows installed):
    - In VirtualBox menu: Devices → Insert Guest Additions CD image
    - In Windows: Run the installer from the CD
    - Reboot after installation
    - This enables shared folders and better integration

---

### **Option 3: Windows on Windows (Native Testing)**

If you're already on Windows, you can test directly! But you might want a clean VM for testing:

#### **Using Hyper-V (Built into Windows Pro/Enterprise)**

1. **Enable Hyper-V:**
   ```powershell
   # Run PowerShell as Administrator
   Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
   # Reboot after this
   ```

2. **Open Hyper-V Manager:**
   - Press Windows key
   - Type "Hyper-V Manager"
   - Open it

3. **Create VM:**
   - Click "Quick Create"
   - Select "Windows 11 dev environment"
   - Click "Create Virtual Machine"
   - Wait for download and setup

4. **Start and configure VM**

---

## 🐧 Linux VM Setup

### **Option 1: Multipass (FASTEST - Recommended)**

This is the easiest way to get Ubuntu VMs running!

#### **Step 1: Install Multipass**
```bash
# On macOS
brew install multipass

# On Windows
# Download from: https://multipass.run/download/windows
# Or: choco install multipass

# On Linux
sudo snap install multipass
```

#### **Step 2: Create Ubuntu VM**
```bash
# Create a VM named "gopher-test" with 2 CPUs, 4GB RAM, 20GB disk
multipass launch --name gopher-test --cpus 2 --memory 4G --disk 20G

# Or use defaults (1 CPU, 1GB RAM, 5GB disk)
multipass launch --name gopher-test

# Access the VM
multipass shell gopher-test

# Inside VM - Set up for testing
sudo apt update
sudo apt install -y git build-essential curl

# Install Go (optional, for testing system Go detection)
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Clone your project
git clone https://github.com/yourusername/gopher.git
cd gopher

# Build and test
go build -o gopher cmd/gopher/main.go
./gopher version
```

#### **Multipass Management Commands**
```bash
# List all VMs
multipass list

# Stop a VM
multipass stop gopher-test

# Start a VM
multipass start gopher-test

# Delete a VM
multipass delete gopher-test
multipass purge  # Actually removes deleted VMs

# Mount a folder from host to VM
multipass mount ~/gopher gopher-test:/home/ubuntu/gopher

# Copy files to VM
multipass transfer ./file.txt gopher-test:/home/ubuntu/

# Get VM info
multipass info gopher-test

# Execute command without entering VM
multipass exec gopher-test -- ls -la
```

---

### **Option 2: UTM on macOS**

#### **Step 1: Download Ubuntu ISO**
```bash
# Visit: https://ubuntu.com/download/desktop
# Download Ubuntu 22.04 LTS (3+ GB)
# Or use the minimal ISO for servers
```

#### **Step 2: Create VM in UTM**

1. **Open UTM** → "Create a New Virtual Machine"

2. **Select Virtualize**

3. **Operating System:**
   - Select "Linux"
   - Click "Browse" and select Ubuntu ISO

4. **Hardware:**
   - Memory: 2048 MB (2 GB minimum, 4 GB recommended)
   - CPU: 2 cores
   - Storage: 20 GB

5. **Summary:**
   - Name: "Ubuntu 22.04 - Gopher Testing"
   - Click "Save"

6. **Start and Install:**
   - Click "Play"
   - Follow Ubuntu installation
   - Choose "Minimal installation"
   - Create user account

7. **Post-Installation:**
   ```bash
   # Update system
   sudo apt update && sudo apt upgrade -y
   
   # Install essential tools
   sudo apt install -y git build-essential curl vim
   
   # Install Go (optional)
   curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
   source ~/.bashrc
   
   # Clone and test Gopher
   git clone https://github.com/yourusername/gopher.git
   cd gopher
   make build
   ```

---

### **Option 3: VirtualBox**

#### **Same as Windows setup, but:**
- Select "Linux" as type
- Select "Ubuntu (64-bit)" as version
- Use 2 GB RAM minimum
- Use 20 GB disk minimum
- Follow Ubuntu installation wizard

---

### **Option 4: WSL2 on Windows (Easiest for Windows users)**

#### **Step 1: Install WSL2**
```powershell
# Run PowerShell as Administrator

# Install WSL2 with Ubuntu
wsl --install

# Or specify Ubuntu version
wsl --install -d Ubuntu-22.04

# Reboot if prompted
```

#### **Step 2: Set up Ubuntu**
```bash
# After reboot, Ubuntu terminal will open
# Create username and password

# Update system
sudo apt update && sudo apt upgrade -y

# Install tools
sudo apt install -y git build-essential curl

# Install Go (optional)
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Your Windows files are accessible at /mnt/c/
# Clone your project
git clone https://github.com/yourusername/gopher.git
cd gopher
make build
```

#### **WSL Management Commands**
```powershell
# List installed distributions
wsl --list

# Set default distribution
wsl --set-default Ubuntu-22.04

# Shutdown WSL
wsl --shutdown

# Enter WSL
wsl

# Or enter specific distribution
wsl -d Ubuntu-22.04
```

---

## 🍎 macOS VM Setup

### **Important Note About macOS VMs**

⚠️ **Legal Restrictions:**
- macOS license only allows virtualization on Apple hardware
- Running macOS VM on non-Apple hardware violates Apple's EULA
- Running macOS VM on Apple hardware is allowed

### **Option 1: You're Already on macOS (No VM Needed!)**

If you're on macOS, you can test Gopher natively:

```bash
# Just build and test
make build
./build/gopher version

# Follow the macOS testing script from TESTING_GUIDE.md
```

### **Option 2: Nested VM on macOS (For Testing Different macOS Versions)**

If you want to test on a different macOS version:

#### **Using UTM (Free)**

1. **Download macOS Installer:**
   ```bash
   # For macOS 13 Ventura and later
   # Download from App Store
   
   # Or use softwareupdate command
   softwareupdate --list-full-installers
   softwareupdate --fetch-full-installer --full-installer-version 13.5.2
   ```

2. **Create Ventura ISO** (this is complex, consider using the simpler option below)

#### **Simpler Option: Use Docker for macOS Simulation**

Since setting up macOS VMs is complex, use Docker:

```bash
# Use the pre-configured Docker tests
make docker-test-macos-with-go
make docker-test-macos-no-go
```

These simulate macOS-like environments (not perfect, but good for basic testing).

---

## ☁️ Cloud VM Alternatives

If local VMs are problematic, use cloud VMs:

### **Option 1: GitHub Codespaces (EASIEST)**

Free for personal use with GitHub!

#### **Setup:**
1. Go to your GitHub repository
2. Click "Code" → "Codespaces" tab
3. Click "Create codespace on main"
4. Wait for environment to load
5. You get a full Ubuntu environment in browser!

#### **Test Gopher:**
```bash
# Already in your repository
make build
./build/gopher version

# Run all tests
make test
```

#### **Pricing:**
- Free: 60 hours/month for personal accounts
- Free: 120 core-hours/month

---

### **Option 2: AWS Free Tier**

Get free VMs for 12 months!

#### **Setup:**
1. Sign up at https://aws.amazon.com/free/
2. Get 750 hours/month of t2.micro instances (free for 12 months)

#### **Launch Ubuntu VM:**
```bash
# Using AWS CLI (after configuring credentials)
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t2.micro \
  --key-name YourKeyPair

# Or use AWS Console web interface
```

#### **Launch Windows VM:**
- Use same process but select Windows Server AMI
- Note: Windows Server, not Windows 11, but good for testing

---

### **Option 3: DigitalOcean**

Simple and cheap ($6/month, $200 credit for new users).

#### **Setup:**
1. Sign up at https://www.digitalocean.com/
2. Get $200 credit (60 days)
3. Create a "Droplet" (their term for VM)

#### **Create Ubuntu Droplet:**
1. Click "Create" → "Droplets"
2. Choose Ubuntu 22.04 LTS
3. Choose Basic plan ($6/month)
4. Choose datacenter region
5. Add SSH key
6. Click "Create Droplet"

```bash
# SSH into droplet
ssh root@your-droplet-ip

# Set up testing environment
apt update && apt upgrade -y
apt install -y git build-essential curl

# Install Go
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Clone and test
git clone https://github.com/yourusername/gopher.git
cd gopher
make build
```

---

### **Option 4: Microsoft Azure**

Free tier includes Windows VMs!

#### **Free Tier Includes:**
- 12 months free with $200 credit
- Windows Server VMs
- Linux VMs

#### **Setup:**
1. Sign up at https://azure.microsoft.com/free/
2. Go to "Virtual Machines"
3. Click "Create"
4. Choose Windows Server 2022 or Ubuntu
5. Select B1s size (free tier eligible)

---

## 🐳 Quick Testing with Docker (NO VM NEEDED!)

**The fastest way to test on all platforms:**

```bash
# Test on Linux (Ubuntu)
make docker-test-unix-no-go
make docker-test-unix-with-go

# Test on simulated macOS
make docker-test-macos-no-go
make docker-test-macos-with-go

# Test Windows build (simulation)
make docker-test-windows-simulated

# Or run all at once
make docker-test
```

**Why Docker is easier:**
- ✅ No VM software installation
- ✅ No OS ISOs to download
- ✅ No disk space for VM images
- ✅ Fast startup (seconds vs minutes)
- ✅ Automated testing
- ✅ Already configured in your Makefile

**Limitations:**
- ❌ Not "true" Windows/macOS (simulated)
- ❌ Can't test GUI features
- ❌ Can't test shell integration fully
- ✅ But perfect for command-line testing!

---

## 🔧 Troubleshooting

### **Common VM Issues**

#### **Issue: VM is too slow**
```
Solutions:
1. Increase RAM allocation (4-8 GB recommended)
2. Increase CPU cores (2-4 recommended)
3. Enable hardware virtualization in BIOS
   - Intel: Enable VT-x
   - AMD: Enable AMD-V
4. Use SSD instead of HDD for VM storage
5. Close other applications while running VM
```

#### **Issue: Cannot enable virtualization on Mac**
```
Solutions:
1. On Intel Macs: Check that VT-x is supported
2. On Apple Silicon: Use UTM (native support)
3. Ensure Parallels/VMware is not running simultaneously
4. Try UTM instead of VirtualBox
```

#### **Issue: Windows VM won't boot**
```
Solutions:
1. Check that TPM 2.0 is enabled in VM settings
2. Enable UEFI boot mode
3. Allocate at least 4 GB RAM
4. Try Windows 10 instead of Windows 11
5. Use the official Windows 11 Dev VM
```

#### **Issue: Linux VM has no network**
```bash
# Check network adapter in VM settings
# Should be set to "NAT" or "Bridged"

# Inside VM, check network:
ip addr show
ping google.com

# Restart networking
sudo systemctl restart NetworkManager
```

#### **Issue: Cannot access files between host and VM**
```
VirtualBox:
1. Install Guest Additions
2. Settings → Shared Folders → Add folder
3. In VM: mount shared folder
   sudo mount -t vboxsf ShareName /mnt/shared

UTM:
1. Settings → Sharing → Enable Directory Sharing
2. Files appear automatically in VM

Multipass:
multipass mount ~/gopher gopher-test:/home/ubuntu/gopher
```

#### **Issue: VM uses too much disk space**
```
Solutions:
1. Use "dynamically allocated" storage
2. Delete snapshots if not needed
3. Clean up inside VM:
   - Windows: Disk Cleanup tool
   - Linux: sudo apt clean && sudo apt autoremove
4. Compact VM disk:
   - VirtualBox: VBoxManage modifymedium --compact
   - UTM: Settings → Compress image
```

### **macOS-Specific Issues**

#### **Issue: "System Extension Blocked" when installing VirtualBox**
```
Solution:
1. Open System Preferences
2. Security & Privacy
3. Click "Allow" next to Oracle message
4. Restart VirtualBox
```

#### **Issue: UTM won't run on Apple Silicon**
```
Solution:
1. Make sure you downloaded UTM (not UTM SE)
2. Use "Virtualize" not "Emulate" for x86_64 guests
3. For ARM guests, use "Virtualize" for best performance
4. Update to latest macOS version
```

### **Windows-Specific Issues**

#### **Issue: Hyper-V conflicts with VirtualBox**
```powershell
# They can't run simultaneously
# Choose one:

# Disable Hyper-V (to use VirtualBox):
bcdedit /set hypervisorlaunchtype off
# Reboot

# Enable Hyper-V (to use Hyper-V):
bcdedit /set hypervisorlaunchtype auto
# Reboot
```

#### **Issue: WSL2 won't install**
```powershell
# Run as Administrator

# Enable required features
dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart
dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart

# Reboot

# Download and install WSL2 kernel update
# From: https://aka.ms/wsl2kernel

# Set WSL2 as default
wsl --set-default-version 2
```

---

## 📊 Comparison Matrix

| Solution | Setup Time | Cost | Platforms | Performance | Best For |
|----------|-----------|------|-----------|-------------|----------|
| **Multipass** | 5 min | Free | All | ⭐⭐⭐⭐⭐ | Quick Ubuntu testing |
| **UTM** | 30 min | Free | macOS | ⭐⭐⭐⭐ | Mac users |
| **WSL2** | 5 min | Free | Windows | ⭐⭐⭐⭐⭐ | Windows users + Linux |
| **VirtualBox** | 60 min | Free | All | ⭐⭐⭐ | Cross-platform |
| **Parallels** | 30 min | $99 | macOS | ⭐⭐⭐⭐⭐ | Best Mac experience |
| **Docker** | 5 min | Free | All | ⭐⭐⭐⭐ | Automated testing |
| **GitHub Codespaces** | 2 min | Free* | All | ⭐⭐⭐⭐ | Cloud-based |
| **Cloud VMs** | 15 min | $-$$ | All | ⭐⭐⭐⭐ | Production-like |

---

## 🚀 Recommended Setup Path

### **If you're on macOS:**
```
1. Use native macOS for Mac testing (no VM needed)
2. Install Multipass for Linux testing (5 minutes)
3. Install UTM for Windows testing (30 minutes + download)
4. Use Docker for quick automated testing
```

### **If you're on Windows:**
```
1. Use native Windows for Windows testing (no VM needed)
2. Install WSL2 for Linux testing (5 minutes)
3. Use Docker for automated testing
4. Cannot legally test macOS (use Docker simulation)
```

### **If you're on Linux:**
```
1. Use native Linux for Linux testing (no VM needed)
2. Install VirtualBox for Windows testing (60 minutes)
3. Use Docker for automated testing
4. Cannot legally test macOS (use Docker simulation)
```

---

## 📝 Quick Start Checklist

- [ ] Choose VM solution based on your host OS
- [ ] Install VM software
- [ ] Download OS images/ISOs
- [ ] Create VMs (1 for each platform you need)
- [ ] Install guest tools/additions
- [ ] Set up development environment in each VM
- [ ] Clone Gopher repository
- [ ] Run test scripts from TESTING_GUIDE.md
- [ ] Document results

---

## 📚 Additional Resources

### **Official Documentation**
- [UTM Documentation](https://docs.getutm.app/)
- [VirtualBox Manual](https://www.virtualbox.org/manual/)
- [Multipass Documentation](https://multipass.run/docs)
- [WSL2 Documentation](https://docs.microsoft.com/en-us/windows/wsl/)

### **Video Tutorials**
- Search YouTube for: "UTM macOS tutorial"
- Search YouTube for: "VirtualBox Windows 11 tutorial"
- Search YouTube for: "WSL2 setup tutorial"

### **Community Help**
- [r/virtualization](https://www.reddit.com/r/virtualization/)
- [r/homelab](https://www.reddit.com/r/homelab/)
- [Ask Ubuntu](https://askubuntu.com/)

---

**Need Help?** 

If you're still having issues:
1. Specify your host OS and what problems you're encountering
2. Share error messages
3. I can provide more specific guidance!

---

**Last Updated:** 2025-10-11
**Version:** 1.0

