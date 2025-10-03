# Gopher Docker Testing Environment

This directory contains Docker configurations for testing Gopher in isolated environments across different operating systems and Go installation states.

## 🐳 Quick Start

### Build All Images

```bash
# Using Makefile
make docker-build

# Or using build script directly
./docker/build.sh build
```

### Run All Tests

```bash
# Using Makefile
make docker-test

# Or using build script directly
./docker/build.sh test
```

### Run Specific Scenario

```bash
# Using Makefile
make docker-test-unix-no-go
make docker-test-unix-with-go
make docker-test-windows-no-go
make docker-test-windows-with-go
make docker-test-macos-no-go
make docker-test-macos-with-go

# Or using build script directly
./docker/build.sh test-unix-no-go
./docker/build.sh test-unix-with-go
./docker/build.sh test-windows-no-go
./docker/build.sh test-windows-with-go
./docker/build.sh test-macos-no-go
./docker/build.sh test-macos-with-go
```

## 📋 Available Scenarios

| Scenario | Image | Base OS | Go Status | Description |
|----------|-------|---------|-----------|-------------|
| **Unix No Go** | `molmedoz/gopher:unix-no-go` | Alpine Linux | Not installed | Clean Linux environment for testing fresh installation |
| **Unix With Go** | `molmedoz/gopher:unix-with-go` | Alpine Linux + Go 1.21 | System Go installed | Linux with system Go for testing system integration |
| **Windows No Go** | `molmedoz/gopher:windows-no-go` | Ubuntu 22.04 | Not installed | Windows-like environment for testing fresh installation |
| **Windows With Go** | `molmedoz/gopher:windows-with-go` | Ubuntu + Go 1.21 | System Go installed | Windows-like environment with system Go |
| **macOS No Go** | `molmedoz/gopher:macos-no-go` | Ubuntu 22.04 | Not installed | macOS-like environment for testing fresh installation |
| **macOS With Go** | `molmedoz/gopher:macos-with-go` | Ubuntu + Go 1.21 | System Go installed | macOS-like environment with system Go |

## 🏗️ Architecture

### Base Image (`molmedoz/gopher:base`)
- **Purpose**: Builds Gopher binary and provides base runtime
- **Base**: Alpine Linux
- **Size**: ~41MB
- **Contains**: Gopher binary, runtime dependencies

### Scenario Images
- **Purpose**: Test specific scenarios with pre-configured environments
- **Base**: Various Linux distributions
- **Size**: 44MB - 1.09GB (depending on Go installation)
- **Contains**: Gopher binary, test scripts, environment setup

## 🧪 Test Coverage

Each scenario tests the following functionality:

### Core Features
- ✅ Gopher version display
- ✅ Version listing (empty vs populated)
- ✅ System Go detection
- ✅ Go installation via Gopher
- ✅ Version switching
- ✅ Current version display
- ✅ JSON output support

### System Integration
- ✅ System Go information display
- ✅ Switching between system and Gopher-managed versions
- ✅ Cross-platform compatibility
- ✅ Environment-specific features

### Error Handling
- ✅ No system Go scenarios
- ✅ Installation failures
- ✅ Permission issues
- ✅ Network connectivity

## 🚀 Usage Examples

### Interactive Testing

```bash
# Start interactive container
docker run -it molmedoz/gopher:unix-no-go /bin/bash

# Run specific test
docker run --rm molmedoz/gopher:unix-no-go /home/gopher/test-gopher.sh
```

### Custom Testing

```bash
# Test with custom Go version
docker run --rm molmedoz/gopher:unix-no-go sh -c "
  gopher install 1.22.0
  gopher use 1.22.0
  go version
"

# Test with debug mode
docker run --rm -e GOPHER_DEBUG=1 molmedoz/gopher:unix-no-go
```

### Volume Mounting

```bash
# Persist Gopher data
docker run --rm -v gopher-data:/home/gopher/.gopher molmedoz/gopher:unix-no-go

# Mount custom config
docker run --rm -v $(pwd)/config.json:/home/gopher/.gopher/config.json molmedoz/gopher:unix-no-go
```

## 🔧 Development

### Building Images

```bash
# Build specific image
docker build -f docker/Dockerfile.unix-no-go -t molmedoz/gopher:unix-no-go .

# Build all images
./docker/build.sh build
```

### Customizing Tests

Edit the test scripts in each Dockerfile:

```dockerfile
# Add custom test
RUN echo 'echo "Custom test: $(gopher version)"' >> /home/gopher/test-gopher.sh
```

### Adding New Scenarios

1. Create new Dockerfile in `docker/` directory
2. Add scenario to `docker-compose.yml`
3. Add build target to `Makefile`
4. Add test target to `build.sh`

## 📊 Performance

### Image Sizes

| Image | Size | Description |
|-------|------|-------------|
| `base` | 41MB | Minimal base image |
| `unix-no-go` | 44MB | Alpine + Gopher |
| `unix-with-go` | 363MB | Alpine + Go + Gopher |
| `windows-no-go` | 132MB | Ubuntu + Gopher |
| `windows-with-go` | 1.09GB | Ubuntu + Go + Gopher |
| `macos-no-go` | 132MB | Ubuntu + Gopher |
| `macos-with-go` | 1.09GB | Ubuntu + Go + Gopher |

### Build Times

- **Base image**: ~30 seconds
- **Unix images**: ~5 seconds
- **Windows images**: ~60 seconds
- **macOS images**: ~10 seconds

## 🐛 Troubleshooting

### Common Issues

#### Build Failures

```bash
# Clean up and rebuild
docker system prune -f
./docker/build.sh clean
./docker/build.sh build
```

#### Permission Issues

```bash
# Fix script permissions
chmod +x docker/build.sh

# Run with sudo if needed
sudo ./docker/build.sh build
```

#### Container Issues

```bash
# Check container logs
docker logs gopher-unix-no-go

# Debug container
docker run -it molmedoz/gopher:unix-no-go /bin/bash
```

### Debug Mode

Enable debug output:

```bash
docker run --rm -e GOPHER_DEBUG=1 molmedoz/gopher:unix-no-go
```

## 🔄 CI/CD Integration

### GitHub Actions

```yaml
name: Docker Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Build Docker images
      run: make docker-build
    
    - name: Run tests
      run: make docker-test
```

### GitLab CI

```yaml
test:
  stage: test
  script:
    - make docker-build
    - make docker-test
```

## 📚 Documentation

- **[Test Scenarios](test-scenarios.md)** - Detailed testing documentation
- **[User Guide](../docs/USER_GUIDE.md)** - Gopher user documentation
- **[API Reference](../docs/API_REFERENCE.md)** - API documentation
- **[Examples](../docs/EXAMPLES.md)** - Usage examples

## 🧹 Cleanup

### Remove All Images

```bash
make docker-clean
# or
./docker/build.sh clean
```

### Remove Specific Images

```bash
docker rmi molmedoz/gopher:unix-no-go
docker rmi molmedoz/gopher:unix-with-go
```

### Complete Cleanup

```bash
docker system prune -a -f
```

## 🎯 Best Practices

1. **Always test in isolated environments** before deploying
2. **Use specific image tags** for reproducible builds
3. **Clean up after testing** to save disk space
4. **Test both success and failure scenarios**
5. **Verify cross-platform compatibility**
6. **Test with different Go versions**
7. **Monitor resource usage during tests**

## 🤝 Contributing

When adding new scenarios or modifying existing ones:

1. Update this README
2. Update `docker-compose.yml`
3. Update `Makefile` targets
4. Update `build.sh` script
5. Test all scenarios
6. Update documentation

This Docker testing environment provides comprehensive coverage for testing Gopher across different platforms and scenarios without affecting your local system.
