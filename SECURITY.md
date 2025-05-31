# Security Policy

## Supported Versions

We actively maintain and provide security updates for the following versions:

| Version  | Supported          |
| -------- | ------------------ |
| Latest   | :white_check_mark: |
| < Latest | :x:                |

We recommend always using the latest version to ensure you have the most recent security fixes.

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in NetBird API Exporter, please report it to us privately.

### How to Report

1. **Email**: Send details to [maintainer email - please update this]
2. **Subject Line**: Use "SECURITY: [brief description]"
3. **Include**:
   - Description of the vulnerability
   - Steps to reproduce (if applicable)
   - Affected versions
   - Potential impact assessment
   - Any suggested fixes (optional)

### What to Expect

- **Acknowledgment**: We'll acknowledge receipt within 48 hours
- **Initial Assessment**: We'll provide an initial assessment within 5 business days
- **Updates**: We'll keep you informed of our progress
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days
- **Credit**: We'll credit you in the fix announcement (unless you prefer to remain anonymous)

### Responsible Disclosure

Please allow us reasonable time to investigate and fix vulnerabilities before public disclosure. We commit to:

- Working with you to understand and resolve the issue
- Keeping you informed of our progress
- Providing credit for responsible disclosure
- Releasing security updates in a timely manner

## Security Considerations

### For Users

#### API Token Security

- **Never log API tokens**: Ensure `NETBIRD_API_TOKEN` is not logged or exposed
- **Rotate tokens regularly**: Use token rotation best practices
- **Limit token scope**: Use minimal required permissions for the API token
- **Secure storage**: Store tokens in secure credential management systems

#### Network Security

- **Use HTTPS**: Always use HTTPS for NetBird API connections
- **Network isolation**: Deploy the exporter in appropriate network segments
- **Firewall rules**: Restrict access to the metrics endpoint (`/metrics`)
- **Reverse proxy**: Consider using a reverse proxy with authentication

#### Container Security

- **Run as non-root**: The provided Docker image runs as `nobody` user
- **Image scanning**: Regularly scan container images for vulnerabilities
- **Resource limits**: Set appropriate CPU and memory limits
- **Security contexts**: Use appropriate Kubernetes security contexts

#### Configuration Security

- **Environment variables**: Use secrets management for sensitive configuration
- **File permissions**: Secure configuration files with appropriate permissions
- **Monitoring**: Monitor exporter logs for suspicious activity

### For Developers

#### Code Security

- **Input validation**: Validate all external inputs
- **Error handling**: Avoid exposing sensitive information in errors
- **Dependencies**: Keep dependencies updated and scan for vulnerabilities
- **Secrets**: Never commit secrets or API tokens to version control

#### Build Security

- **Reproducible builds**: Use pinned dependencies and container base images
- **Supply chain**: Verify integrity of dependencies
- **Signing**: Consider signing releases and container images

## Security Features

### Built-in Security Measures

1. **Structured Logging**: Sensitive data is filtered from logs
2. **Error Handling**: Errors don't expose system internals
3. **Minimal Privileges**: Container runs as non-root user
4. **Input Validation**: API responses are validated before processing
5. **Rate Limiting**: Built-in protection against API abuse

### Configuration Hardening

1. **Default Configurations**: Secure defaults are provided
2. **TLS Support**: HTTPS connections to NetBird API
3. **Metrics Filtering**: Only necessary metrics are exposed
4. **Health Checks**: Built-in health endpoints for monitoring

## Security Updates

### Update Notifications

- **GitHub Releases**: Security updates are clearly marked in release notes
- **Changelog**: Security fixes are documented in `CHANGELOG.md`
- **Container Images**: Updated images are published to container registries

### Emergency Updates

For critical security vulnerabilities:

- Emergency releases will be published within 24-48 hours
- Clear upgrade instructions will be provided
- Mitigation steps will be documented

## Compliance and Standards

### Security Standards

- We follow OWASP security guidelines
- Regular security reviews of code and dependencies
- Automated security scanning in CI/CD pipeline

### Privacy

- The exporter only collects metrics data from NetBird API
- No personal data is stored or transmitted beyond what's necessary for metrics
- Logs are structured to avoid capturing sensitive information

## Contact

For security-related questions or concerns:

- **Security Issues**: [maintainer email - please update this]
- **General Questions**: Open an issue on GitHub (for non-security matters)

## Acknowledgments

We appreciate the security research community and responsible disclosure of vulnerabilities. Contributors who report security issues will be acknowledged in our security advisories (unless they prefer to remain anonymous).

---

**Note**: This security policy is subject to updates. Please check back regularly for the latest version.
