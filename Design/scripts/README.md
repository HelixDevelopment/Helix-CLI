# HelixCode Design Scripts

## Security Notice

**⚠️ IMPORTANT: Never commit access tokens or secrets to version control!**

This directory contains scripts for generating and managing HelixCode designs. Some scripts require API access tokens which should be stored securely.

## Setup Instructions

### 1. Environment Configuration

Copy the example environment file:
```bash
cp .env.example .env
```

Edit the `.env` file and add your actual tokens:
```bash
# .env
FIGMA_ACCESS_TOKEN=your_actual_figma_token_here
```

### 2. Running Scripts

Use environment variables when running scripts:
```bash
# Set token temporarily
FIGMA_ACCESS_TOKEN=your_token node figma_integration.js

# Or load from .env file
source .env && node figma_integration.js
```

### 3. Security Best Practices

- **Never commit** `.env` files
- **Use environment variables** for tokens
- **Rotate tokens** regularly
- **Use token scopes** with minimal permissions
- **Store tokens** in secure password managers

## Available Scripts

### `figma_integration.js`
Generates Figma design specifications and connects to Figma API.

**Usage:**
```bash
FIGMA_ACCESS_TOKEN=your_token node figma_integration.js
```

### `generate_exports.js`
Generates mock PNG and PDF exports for design review.

**Usage:**
```bash
node generate_exports.js
```

### `generate_figma_designs.js`
Generates Figma-compatible design specifications.

**Usage:**
```bash
node generate_figma_designs.js
```

### `generate_penpot_designs.js`
Generates PenPot-compatible design specifications.

**Usage:**
```bash
node generate_penpot_designs.js
```

## Token Management

### Figma Personal Access Token
1. Go to Figma Account Settings
2. Navigate to "Personal access tokens"
3. Create a new token with appropriate scopes
4. Copy the token to your `.env` file

### Required Scopes
- `file_read` - Read design files
- `file_write` - Create/update design files (if needed)

## Troubleshooting

### Common Issues

**"Token not found" error**
- Ensure `.env` file exists and contains the token
- Verify token is set in environment variables

**API connection failures**
- Check token validity and expiration
- Verify network connectivity
- Ensure token has correct scopes

### Security Issues

If you accidentally committed a token:
1. **Immediately revoke** the token in Figma settings
2. **Remove** the token from commit history
3. **Create a new token** with updated permissions
4. **Update** your local `.env` file

## File Structure

```
scripts/
├── .env.example          # Example configuration
├── .gitignore           # Git ignore rules
├── figma_integration.js # Figma API integration
├── generate_exports.js  # Export generation
├── generate_figma_designs.js
├── generate_penpot_designs.js
└── README.md           # This file
```

## Support

For security concerns or token issues, contact the development team immediately.