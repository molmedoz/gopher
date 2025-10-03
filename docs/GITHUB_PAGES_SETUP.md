# GitHub Pages Setup Guide

This guide explains how to set up GitHub Pages for Gopher documentation.

## 📚 What is GitHub Pages?

GitHub Pages hosts your documentation directly from your repository, turning your `docs/` folder into a beautiful, searchable website.

**Benefits:**
- ✅ Uses your existing `docs/` folder (no duplication!)
- ✅ Automatically updates when you push to main
- ✅ Professional appearance with themes
- ✅ Free hosting from GitHub
- ✅ Custom domain support (optional)
- ✅ Search functionality (with theme)

**Your docs will be available at:**
```
https://molmedoz.github.io/gopher/
```

---

## 🚀 Setup Instructions

### **Step 1: Enable GitHub Pages**

1. Go to your repository on GitHub
2. Click **Settings** (top right)
3. Click **Pages** (left sidebar)
4. Under **Source**:
   - Select branch: `main`
   - Select folder: `/ (root)`
5. Click **Save**

**That's it!** Your documentation will be live in a few minutes.

---

### **Step 2: Choose Deployment Method**

#### **Option A: Basic (Simplest)**

No additional setup needed. GitHub will use the `_config.yml` we created.

**Access your docs at:**
- Main page: `https://molmedoz.github.io/gopher/`
- Docs index: `https://molmedoz.github.io/gopher/docs/`
- Specific doc: `https://molmedoz.github.io/gopher/docs/USER_GUIDE`

#### **Option B: Docs-Only (Recommended)**

Make GitHub Pages serve only the `docs/` folder:

1. In **Settings** → **Pages**
2. Under **Source**, select folder: `/docs`
3. Click **Save**

**Access your docs at:**
- Main page: `https://molmedoz.github.io/gopher/`
- User Guide: `https://molmedoz.github.io/gopher/USER_GUIDE`
- FAQ: `https://molmedoz.github.io/gopher/FAQ`

---

### **Step 3: Customize Theme (Optional)**

Edit `_config.yml` to change the theme:

```yaml
# Available GitHub themes:
theme: jekyll-theme-cayman      # Clean, modern (default)
# theme: jekyll-theme-minimal   # Simple, minimal
# theme: jekyll-theme-slate     # Dark theme
# theme: jekyll-theme-architect # Professional
# theme: jekyll-theme-hacker    # Terminal-style
```

**Preview themes:**
- [Cayman](https://pages-themes.github.io/cayman/)
- [Minimal](https://pages-themes.github.io/minimal/)
- [Slate](https://pages-themes.github.io/slate/)
- [Architect](https://pages-themes.github.io/architect/)

---

### **Step 4: Update README.md (Optional)**

Add a badge to your README.md:

```markdown
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://molmedoz.github.io/gopher/)
```

Add a documentation link:

```markdown
## Documentation

📚 **[Read the docs](https://molmedoz.github.io/gopher/)** - Complete documentation hosted on GitHub Pages
```

---

## 🔧 Advanced Configuration

### **Custom Domain**

To use a custom domain (e.g., `docs.gopher.dev`):

1. Create a `CNAME` file in your repo root:
   ```
   docs.gopher.dev
   ```

2. Configure DNS with your domain provider:
   ```
   Type: CNAME
   Name: docs
   Value: molmedoz.github.io
   ```

3. In **Settings** → **Pages** → **Custom domain**, enter: `docs.gopher.dev`

### **Automated Deployment (GitHub Actions)**

Already configured! GitHub automatically deploys when you push to `main`.

To customize deployment, create `.github/workflows/pages.yml`:

```yaml
name: Deploy GitHub Pages

on:
  push:
    branches: [ main ]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/configure-pages@v4
      - uses: actions/upload-pages-artifact@v3
        with:
          path: 'docs'
      - uses: actions/deploy-pages@v4
```

---

## 📊 Features You Get

### **With Basic Setup:**
- ✅ All your markdown files rendered as HTML
- ✅ Automatic navigation from links
- ✅ Themed appearance
- ✅ Mobile responsive
- ✅ Syntax highlighting for code blocks

### **With Jekyll Theme:**
- ✅ Professional design
- ✅ Automatic table of contents (some themes)
- ✅ SEO optimization
- ✅ Social media previews
- ✅ Google Analytics support (optional)

---

## 🎨 Customizing Appearance

### **Add Custom CSS**

Create `assets/css/style.scss`:

```scss
---
---

@import "{{ site.theme }}";

/* Custom styles */
.page-header {
  background: linear-gradient(120deg, #2563eb, #1d4ed8);
}

code {
  background-color: #f3f4f6;
  padding: 2px 6px;
  border-radius: 3px;
}
```

### **Add Custom Layouts**

Create `_layouts/default.html` to customize the page layout.

### **Add Search**

Add Algolia DocSearch (free for open source):

1. Apply at [docsearch.algolia.com](https://docsearch.algolia.com/)
2. Add the search widget to your pages

---

## 🔍 Verifying Your Setup

### **Check Build Status**

1. Go to your repository
2. Click **Actions** tab
3. Look for "pages build and deployment"
4. Green checkmark = successful deployment

### **Test Your Pages**

Visit these URLs to verify:

```
https://molmedoz.github.io/gopher/
https://molmedoz.github.io/gopher/docs/
https://molmedoz.github.io/gopher/docs/USER_GUIDE
https://molmedoz.github.io/gopher/QUICK_REFERENCE
```

### **Common Issues**

**404 Not Found:**
- Wait 5-10 minutes for first deployment
- Check **Settings** → **Pages** for build status
- Verify branch and folder are correct

**Broken Links:**
- Use relative links: `[Text](./FILE.md)` not `[Text](FILE.md)`
- Check `_config.yml` has correct `baseurl`

**No Styling:**
- Verify `_config.yml` exists in root
- Check theme name is correct
- Wait for deployment to complete

---

## 📝 Maintaining Your Documentation

### **When You Update Docs:**

1. Edit files in `docs/` folder
2. Commit and push to `main`
3. GitHub automatically rebuilds and deploys
4. Changes live in ~2-5 minutes

### **Best Practices:**

- ✅ Use relative links between docs
- ✅ Keep `_config.yml` in root
- ✅ Use frontmatter in markdown files (optional but recommended)
- ✅ Test links locally before pushing
- ✅ Keep docs in sync with code versions

---

## 🆚 GitHub Pages vs GitHub Wiki

### **GitHub Pages (Recommended for Gopher):**
- ✅ Uses your `docs/` folder (no duplication)
- ✅ Version-controlled with code
- ✅ Updated via pull requests (quality control)
- ✅ Professional appearance
- ✅ Custom themes and styling
- ✅ Stays in sync with code

### **GitHub Wiki:**
- ✅ Easy for community to edit
- ✅ Separate from code repository
- ❌ Can become outdated
- ❌ Separate git repo (`.wiki.git`)
- ❌ Harder to keep in sync
- ❌ Basic styling only

**For Gopher, GitHub Pages is better** because:
1. Your docs are already excellent and well-organized
2. You want version control and quality control
3. You want docs to stay in sync with code
4. Professional appearance matches your project quality

---

## 🎓 Summary

### **What We Created:**

1. ✅ `_config.yml` - Jekyll configuration for GitHub Pages
2. ✅ `docs/index.md` - Documentation landing page
3. ✅ This guide - Setup instructions

### **What You Need to Do:**

1. **Enable GitHub Pages** in repository settings
2. **Choose deployment folder** (root or `/docs`)
3. **Wait 5-10 minutes** for first deployment
4. **Visit your docs** at `https://molmedoz.github.io/gopher/`

### **Optional Enhancements:**

- Choose a different theme
- Add custom CSS
- Set up custom domain
- Add search functionality
- Add Google Analytics

---

## 🔗 Additional Resources

- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [Jekyll Themes](https://pages.github.com/themes/)
- [Jekyll Documentation](https://jekyllrb.com/docs/)
- [Markdown Guide](https://www.markdownguide.org/)

---

**Last Updated:** 2025-10-13  
**Version:** 1.0  
**Next Steps:** Enable GitHub Pages in your repository settings!

