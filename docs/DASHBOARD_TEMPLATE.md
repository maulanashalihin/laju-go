# VeloStack Go - Dashboard Template Features

This document describes the standard dashboard template features that have been added to VeloStack Go (laju-go), matching the capabilities available in laju-ts.

## 🎨 Frontend Components

### 1. **Header Component** (`frontend/src/components/Header.svelte`)
- Desktop sidebar navigation with logo
- Mobile responsive hamburger menu with drawer
- User dropdown menu
- Active state highlighting
- Dark mode toggle integration
- Logout functionality

### 2. **DarkModeToggle Component** (`frontend/src/components/DarkModeToggle.svelte`)
- System preference detection
- localStorage persistence
- Sun/Moon icons (lucide-svelte)
- Smooth transitions
- No layout shift on mount

### 3. **Toast Notification System** (`frontend/src/lib/utils/helpers.js`)
- Success, error, warning, and info variants
- Auto-dismiss with smooth animations
- Customizable duration
- Icon indicators for each type
- Queue management for multiple toasts

### 4. **Utility Helpers** (`frontend/src/lib/utils/helpers.js`)
- `clickOutside` - Action for detecting clicks outside elements
- `debounce` - Function debouncing
- `password_generator` - Random password generation
- `getCsrfToken` - CSRF token retrieval
- `fetchWithCsrf` - Fetch wrapper with CSRF token injection
- `Toast` - Toast notification function

### 5. **Translation System (i18n)** (`frontend/src/lib/i18n/`)
- English (`en.json`) and Indonesian (`id.json`) translations
- `translation.js` - Translation helper with:
  - `t(key, params)` - Translate function
  - `setLocale(locale)` - Set current locale
  - `getLocale()` - Get current locale
  - Nested key support (e.g., `auth.login`)

## 📄 Enhanced Pages

### 1. **Dashboard** (`frontend/src/pages/app/Dashboard.svelte`)
- Welcome card with gradient icon
- Stats grid (Performance, Latency, Uptime)
- Quick action cards (Documentation, Profile Settings)
- Animated blob backgrounds
- Fly-in transitions
- Dark mode support

### 2. **Profile** (`frontend/src/pages/app/Profile.svelte`)
- Avatar upload with live preview
- Profile information form (name, email)
- Password change form with validation
- Dark mode toggle
- Breadcrumb navigation
- Flash message display
- Show/hide password toggle
- Responsive design

### 3. **Forgot Password** (`frontend/src/pages/auth/ForgotPassword.svelte`)
- Email input form
- Modern gradient background
- Loading states
- Success/error message display
- Back to login link

### 4. **Reset Password** (`frontend/src/pages/auth/ResetPassword.svelte`)
- Password input with show/hide toggle
- Password confirmation
- Password strength indicator
- Token validation
- Loading states

## 🔧 Backend Features

### 1. **Password Reset Handler** (`app/handlers/password-reset.go`)
- `ShowForgotPasswordForm` - Display forgot password page
- `SendResetLink` - Send password reset email
- `ShowResetPasswordForm` - Display reset password page
- `ResetPassword` - Process password reset

### 2. **Email Service** (`app/services/mailer.go`)
- SMTP email sending
- HTML email templates
- Password reset token generation
- Token validation and expiration
- Automatic token cleanup

### 3. **Rate Limiting Middleware** (`app/middlewares/rate-limit.go`)
- Configurable rate limits
- In-memory storage with cleanup
- Predefined limiters:
  - `AuthRateLimit` - 5 requests/15 minutes
  - `PasswordResetRateLimit` - 3 requests/hour
  - `APIRateLimit` - 100 requests/15 minutes
  - `UploadRateLimit` - 50 requests/hour

### 4. **CSRF Protection Middleware** (`app/middlewares/csrf.go`)
- Token generation and validation
- Session-based token storage
- Configurable expiry
- Constant-time comparison
- Skip paths and methods configuration

### 5. **Enhanced User Service** (`app/services/user.go`)
- `GetProfileByEmail` - Get user by email
- `UpdatePassword` - Update user password
- `ChangePassword` - Change password with verification

### 6. **Enhanced App Handler** (`app/handlers/app.go`)
- `UpdatePassword` - Handle password change requests

## 🛣️ Routes

### New Routes Added:
```
GET  /forgot-password          - Forgot password page
POST /forgot-password          - Send reset link (rate limited)
GET  /reset-password/:token    - Reset password page
POST /reset-password/:token    - Process password reset

PUT  /app/profile/password     - Change password (protected)
```

## ⚙️ Configuration

### Environment Variables (`.env.example`)
```env
# Email Configuration (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-smtp-username
SMTP_PASS=your-smtp-password
FROM_EMAIL=noreply@example.com
FROM_NAME=VeloStack

# Application URL
APP_URL=http://localhost:8080
```

### Config Structure (`app/config/config.go`)
- Added email configuration fields
- `getEnvAsInt` helper function
- Default values for all email settings

## 📁 File Structure

```
laju-go/
├── app/
│   ├── handlers/
│   │   └── password-reset.go         # NEW
│   ├── middlewares/
│   │   ├── rate-limit.go             # NEW
│   │   └── csrf.go                   # NEW
│   └── services/
│       └── mailer.go                 # NEW
├── frontend/src/
│   ├── components/
│   │   ├── Header.svelte             # NEW
│   │   └── DarkModeToggle.svelte     # NEW
│   ├── lib/
│   │   ├── utils/
│   │   │   └── helpers.js            # NEW
│   │   └── i18n/
│   │       ├── en.json               # NEW
│   │       ├── id.json               # NEW
│   │       └── translation.js        # NEW
│   └── pages/
│       ├── app/
│       │   ├── Dashboard.svelte      # ENHANCED
│       │   └── Profile.svelte        # ENHANCED
│       └── auth/
│           ├── ForgotPassword.svelte # NEW
│           └── ResetPassword.svelte  # NEW
├── routes/
│   └── web.go                        # UPDATED
├── .env.example                      # UPDATED
└── main.go                           # UPDATED
```

## 🚀 Quick Start

1. **Configure Email Settings**
   ```bash
   cp .env.example .env
   # Edit .env with your SMTP credentials
   ```

2. **Install Dependencies**
   ```bash
   npm install
   go mod tidy
   ```

3. **Run Development Server**
   ```bash
   npm run dev
   ```

4. **Build for Production**
   ```bash
   npm run build
   ```

## 🎯 Key Features Summary

| Feature | Status |
|---------|--------|
| Dark Mode Toggle | ✅ |
| Responsive Sidebar Navigation | ✅ |
| Toast Notifications | ✅ |
| Password Reset Flow | ✅ |
| Rate Limiting | ✅ |
| CSRF Protection | ✅ |
| Avatar Upload | ✅ |
| Password Change | ✅ |
| i18n Support (EN/ID) | ✅ |
| Modern Dashboard UI | ✅ |
| Enhanced Profile Page | ✅ |

## 🔐 Security Features

- **CSRF Protection**: All state-changing requests require valid CSRF token
- **Rate Limiting**: Prevents brute force attacks on auth endpoints
- **Password Reset Tokens**: Secure random tokens with 1-hour expiration
- **Constant-Time Comparison**: Prevents timing attacks on token validation

## 📝 Notes

- The mailer service uses in-memory token storage. For production, consider using Redis or a database.
- CSRF cookies are set with `HTTPOnly: false` to allow JavaScript access for Inertia.js requests.
- Rate limiters clean up expired entries automatically every minute.
- All new pages support dark mode out of the box.
