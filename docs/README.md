# NetBird API Exporter Documentation

This directory contains the complete documentation for the NetBird API Exporter, built with Jekyll and the Just the Docs theme.

## 🌐 Live Documentation

The documentation is automatically deployed to GitHub Pages and available at:
**https://matanbaruch.github.io/netbird-api-exporter**

## 🚀 Setting Up GitHub Pages

To enable GitHub Pages for this repository:

### 1. Enable GitHub Pages

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Pages**
3. Under **Source**, select **GitHub Actions**
4. The documentation will automatically build and deploy when you push to the `main` branch

### 2. Verify GitHub Actions

The `.github/workflows/docs.yml` workflow will:
- Build the Jekyll site when docs are updated
- Deploy to GitHub Pages automatically
- Only deploy from the `main` branch

### 3. Access Your Documentation

Once deployed, your documentation will be available at:
- `https://YOUR_USERNAME.github.io/netbird-api-exporter`

## 📁 Documentation Structure

```
docs/
├── _config.yml                    # Jekyll configuration
├── Gemfile                        # Ruby dependencies
├── index.md                       # Homepage
├── getting-started.md             # Getting started guide
├── getting-started/
│   └── authentication.md         # API token setup
├── installation.md                # Installation overview
├── installation/
│   ├── docker-compose.md         # Docker Compose setup
│   ├── docker.md                 # Docker setup
│   ├── helm.md                   # Kubernetes/Helm setup
│   ├── systemd.md                # Linux systemd setup
│   └── binary.md                 # Binary/source setup
└── README.md                     # This file
```

## 🛠️ Local Development

To run the documentation locally:

### Prerequisites

- Ruby 3.1+
- Bundler

### Setup

```bash
# Navigate to docs directory
cd docs

# Install dependencies
bundle install

# Serve locally
bundle exec jekyll serve

# Open browser to http://localhost:4000
```

### Making Changes

1. Edit markdown files in the `docs/` directory
2. The site will auto-reload when you save changes
3. Test locally before committing
4. Push to `main` branch to deploy to GitHub Pages

## 🎨 Theme and Features

This documentation uses the [Just the Docs](https://just-the-docs.github.io/just-the-docs/) theme with:

- **Search functionality** - Full-text search across all pages
- **Navigation sidebar** - Organized hierarchical navigation
- **Mobile responsive** - Works on all devices
- **Code highlighting** - Syntax highlighting for all code blocks
- **Dark/light themes** - Automatic theme switching
- **Edit on GitHub** - Direct links to edit pages

## 📝 Writing Documentation

### Page Structure

Each page should have front matter at the top:

```yaml
---
layout: default
title: Page Title
parent: Parent Page (optional)
nav_order: 1
has_children: true (optional)
---
```

### Navigation

- `nav_order`: Controls the order in navigation (lower numbers first)
- `parent`: Creates a parent-child relationship
- `has_children`: Enables dropdown for child pages

### Content Guidelines

1. **Use clear headings** - H1 for page title, H2-H6 for sections
2. **Add table of contents** - Use the TOC snippet for long pages
3. **Include code examples** - Always provide working code samples
4. **Link between pages** - Use relative links like `[text](../other-page)`
5. **Use callouts** - Important notes, warnings, etc.

### Callouts

Use callouts for important information:

```markdown
{: .important }
> This is an important note

{: .warning }  
> This is a warning

{: .note }
> This is a note
```

## 🔧 Configuration

### Site Configuration

Key settings in `_config.yml`:

- `title`: Site title
- `description`: Site description
- `url`: Production URL
- `baseurl`: Path prefix (usually `/repo-name`)

### Theme Configuration

- `color_scheme`: light or dark
- `search_enabled`: Enable/disable search
- `nav_external_links`: External navigation links

## 📦 Dependencies

Ruby gems used:

- `jekyll`: Static site generator
- `just-the-docs`: Documentation theme
- `jekyll-feed`: RSS feed generation
- `jekyll-sitemap`: XML sitemap
- `jekyll-seo-tag`: SEO meta tags

## 🚀 Deployment

### Automatic Deployment

The GitHub Actions workflow automatically:

1. **Triggers on**:
   - Push to `main` branch with docs changes
   - Manual workflow dispatch

2. **Build process**:
   - Sets up Ruby environment
   - Installs dependencies
   - Builds Jekyll site
   - Uploads artifacts

3. **Deployment**:
   - Deploys to GitHub Pages
   - Available at the configured URL

### Manual Deployment

To deploy manually:

```bash
# Build the site
cd docs
bundle exec jekyll build

# The built site is in _site/
# Deploy _site/ contents to your web server
```

## 🆘 Troubleshooting

### Common Issues

1. **Pages not updating**: Check GitHub Actions for build errors
2. **404 errors**: Verify `baseurl` in `_config.yml`
3. **Missing styles**: Ensure all assets are properly linked
4. **Search not working**: Check that `search_enabled: true` in config

### Build Errors

Check the GitHub Actions logs:
1. Go to **Actions** tab in your repository
2. Click on the failed workflow
3. Expand the build step to see error details

### Local Development Issues

```bash
# Clear Jekyll cache
bundle exec jekyll clean

# Reinstall dependencies
bundle clean --force
bundle install

# Update dependencies
bundle update
```

## 📄 License

This documentation is part of the NetBird API Exporter project and follows the same license terms. 
