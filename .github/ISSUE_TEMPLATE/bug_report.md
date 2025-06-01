---
name: 🐛 Bug Report
about: Report a bug to help us improve the NetBird API Exporter
title: '[BUG] '
labels: 'bug'
assignees: ''

---

## 🐛 Bug Description
A clear and concise description of what the bug is.

## 🔄 Steps to Reproduce
Steps to reproduce the behavior:
1. Configure the exporter with...
2. Start the exporter with...
3. Query the metrics endpoint...
4. Observe the error/unexpected behavior

## ✅ Expected Behavior
A clear and concise description of what you expected to happen.

## ❌ Actual Behavior
A clear and concise description of what actually happened.

## 🌍 Environment Information
**NetBird API Exporter:**
- Version: [e.g., v1.2.3]
- Installation method: [e.g., Docker, Helm, binary]

**NetBird API:**
- API URL: [e.g., https://api.netbird.io or self-hosted]
- API Version: [if known]

**System:**
- OS: [e.g., Ubuntu 22.04, macOS 13.0, Windows 11]
- Architecture: [e.g., amd64, arm64]
- Container Runtime: [e.g., Docker 24.0.5, Podman 4.6.1] (if applicable)

**Prometheus:**
- Version: [e.g., 2.45.0]
- Configuration: [relevant scrape config if applicable]

## 📝 Configuration
**Environment Variables:**
```bash
NETBIRD_API_URL=...
NETBIRD_API_TOKEN=***
LISTEN_ADDRESS=...
METRICS_PATH=...
LOG_LEVEL=...
# Add other relevant env vars
```

**Command Line Arguments:**
```bash
# Include the command used to start the exporter
```

## 📊 Logs
**Exporter Logs:**
```
# Include relevant log output from the exporter
# Remember to redact sensitive information like API tokens
```

**Prometheus Logs (if applicable):**
```
# Include relevant Prometheus scrape errors or warnings
```

## 📸 Screenshots/Output
If applicable, add screenshots or text output to help explain your problem.

## 🔍 Additional Context
- Any additional context about the problem
- Related issues or discussions
- Workarounds you've tried
- Impact on your monitoring setup
