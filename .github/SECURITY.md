# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in Gopher, please report it to us as described below.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via one of the following methods:

1. **Email**: Send an email to [security@molmedoz.dev](mailto:security@molmedoz.dev)
2. **GitHub Security Advisory**: Use GitHub's private vulnerability reporting feature

### What to Include

When reporting a vulnerability, please include:

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if any)
- Your contact information

### Response Timeline

We will respond to security vulnerability reports within:

- **Initial response**: 24 hours
- **Status update**: 72 hours
- **Resolution**: 7 days (for critical vulnerabilities)

### Security Measures

Gopher implements several security measures:

- **Cryptographic verification**: All downloaded Go binaries are verified using SHA256 checksums
- **Secure downloads**: Downloads are performed over HTTPS
- **Input validation**: All user inputs are validated and sanitized
- **Path traversal protection**: File operations are protected against path traversal attacks
- **Permission checks**: Proper file permissions are enforced

### Security Best Practices

When using Gopher:

1. **Verify checksums**: Always verify the integrity of downloaded files
2. **Use HTTPS**: Ensure downloads are performed over secure connections
3. **Keep updated**: Use the latest version of Gopher
4. **Review permissions**: Check file permissions after installation
5. **Monitor logs**: Review logs for any suspicious activity

### Acknowledgments

We would like to thank the following security researchers who have helped improve Gopher's security:

- [Add acknowledgments here]

### Security Updates

Security updates will be released as soon as possible after a vulnerability is confirmed and a fix is available. Updates will be announced through:

- GitHub releases
- Security advisories
- Project documentation

### Contact

For security-related questions or concerns, please contact:

- Email: [security@molmedoz.dev](mailto:security@molmedoz.dev)
- GitHub: [@molmedoz](https://github.com/molmedoz)

Thank you for helping keep Gopher and its users safe!
