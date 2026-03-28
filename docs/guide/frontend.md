# Frontend Development

This guide covers frontend development with Svelte 5 and Inertia.js in Laju Go.

## Overview

Laju Go uses **Svelte 5** for the frontend with **Inertia.js** as the bridge between backend and frontend. This gives you:

- **SPA Experience** - Client-side navigation without page reloads
- **Server-Driven** - Backend controls routing and business logic
- **No API Required** - Return data directly from controllers
- **Simple Mental Model** - Traditional HTTP requests with modern UX

## Project Structure

```
frontend/
├── src/
│   ├── components/        # Reusable UI components
│   ├── pages/             # Page components
│   ├── lib/               # Utilities and helpers
│   ├── main.ts            # Inertia.js entry point
│   └── app.css            # Global styles
├── package.json
└── vite.config.js
```

## Svelte 5 Basics

### Component Structure

```svelte
<!-- frontend/src/components/Card.svelte -->
<script>
  // Props with default values
  export let title = '';
  export let subtitle = '';
  
  // Reactive state
  let count = $state(0);
  
  // Functions
  function increment() {
    count++;
  }
</script>

<div class="card">
  <h2>{title}</h2>
  {#if subtitle}
    <p>{subtitle}</p>
  {/if}
  
  <button onclick={increment}>
    Count: {count}
  </button>
  
  <slot />
</div>

<style>
  .card {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 16px;
    margin: 16px;
  }
</style>
```

### Reactive State

Svelte 5 uses `$state` for reactive variables:

```svelte
<script>
  // Reactive state
  let count = $state(0);
  let user = $state({ name: '', email: '' });
  
  // Derived values (automatically reactive)
  $: doubled = count * 2;
  
  // Side effects
  $effect(() => {
    console.log('Count changed:', count);
  });
</script>
```

### Using Components

```svelte
<!-- frontend/src/pages/app/Dashboard.svelte -->
<script>
  import Card from '../../components/Card.svelte';
</script>

<Card title="Welcome" subtitle="Dashboard">
  <p>Dashboard content here</p>
</Card>
```

## Inertia.js Integration

### Page Component

```svelte
<!-- frontend/src/pages/app/Dashboard.svelte -->
<script>
  import { page } from '@inertiajs/svelte';
  
  // Access props from server
  const props = $page.props;
  const user = props.user;
  const stats = props.stats;
</script>

<h1>Welcome, {user.name}!</h1>

<div class="stats">
  <div class="stat">
    <span class="value">{stats.totalUsers}</span>
    <span class="label">Total Users</span>
  </div>
</div>
```

### Server Props

From the backend:

```go
// app/handlers/app.go
func (h *AppHandler) Dashboard(c *fiber.Ctx) error {
    return h.inertiaService.Render(c, "Dashboard", fiber.Map{
        "user": c.Locals("user"),
        "stats": fiber.Map{
            "totalUsers": 100,
            "activeUsers": 50,
        },
    })
}
```

## Navigation

### Inertia Links

```svelte
<script>
  import { Link } from '@inertiajs/svelte';
</script>

<nav>
  <Link href="/">Home</Link>
  <Link href="/about">About</Link>
  <Link href="/app" class="active-class">Dashboard</Link>
</nav>
```

### Programmatic Navigation

```svelte
<script>
  import { router } from '@inertiajs/svelte';
  
  function goToProfile() {
    router.visit('/app/profile');
  }
</script>

<button onclick={goToProfile}>
  Go to Profile
</button>
```

## Form Handling

### Basic Form

```svelte
<!-- frontend/src/pages/auth/Login.svelte -->
<script>
  import { router } from '@inertiajs/svelte';
  
  let email = $state('');
  let password = $state('');
  let processing = $state(false);
  
  function submit() {
    processing = true;
    
    router.post('/login/login', {
      email,
      password,
    }, {
      onSuccess: () => {
        // Redirect happens automatically
      },
      onError: (errors) => {
        console.log('Validation errors:', errors);
        processing = false;
      },
    });
  }
</script>

<form onsubmit={(e) => { e.preventDefault(); submit(); }}>
  <div>
    <label for="email">Email</label>
    <input 
      type="email" 
      id="email" 
      bind:value={email}
      required
    />
  </div>
  
  <div>
    <label for="password">Password</label>
    <input 
      type="password" 
      id="password" 
      bind:value={password}
      required
    />
  </div>
  
  <button type="submit" disabled={processing}>
    {processing ? 'Logging in...' : 'Login'}
  </button>
</form>
```

### Handling Validation Errors

```svelte
<script>
  import { page } from '@inertiajs/svelte';
  
  let email = $state('');
  
  // Access errors from server
  $: errors = $page.props.errors || {};
</script>

<form>
  <div>
    <label for="email">Email</label>
    <input 
      type="email" 
      id="email" 
      bind:value={email}
      class:error={errors.email}
    />
    {#if errors.email}
      <span class="error">{errors.email}</span>
    {/if}
  </div>
</form>

<style>
  .error {
    color: red;
    font-size: 0.875rem;
  }
  
  input.error {
    border-color: red;
  }
</style>
```

### File Upload

```svelte
<script>
  import { router } from '@inertiajs/svelte';
  
  let avatar = $state(null);
  let preview = $state(null);
  
  function handleFileChange(event) {
    const file = event.target.files[0];
    if (file) {
      avatar = file;
      preview = URL.createObjectURL(file);
    }
  }
  
  function upload() {
    if (!avatar) return;
    
    const formData = new FormData();
    formData.append('avatar', avatar);
    
    router.post('/upload', formData, {
      onSuccess: () => {
        alert('Upload successful!');
      },
      onError: (errors) => {
        alert('Upload failed: ' + (errors.avatar || 'Unknown error'));
      },
    });
  }
</script>

<div>
  {#if preview}
    <img src={preview} alt="Preview" />
  {/if}
  
  <input 
    type="file" 
    accept="image/*" 
    onchange={handleFileChange}
  />
  
  <button onclick={upload}>
    Upload Avatar
  </button>
</div>
```

## Components

### Button Component

```svelte
<!-- frontend/src/components/Button.svelte -->
<script>
  export let variant = 'primary';
  export let size = 'md';
  export let disabled = false;
  export let loading = false;
  
  let variants = {
    primary: 'bg-blue-500 hover:bg-blue-600 text-white',
    secondary: 'bg-gray-500 hover:bg-gray-600 text-white',
    danger: 'bg-red-500 hover:bg-red-600 text-white',
  };
  
  let sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };
</script>

<button 
  class="{variants[variant]} {sizes[size]}"
  disabled={disabled || loading}
>
  {#if loading}
    <span class="spinner">Loading...</span>
  {:else}
    <slot />
  {/if}
</button>

<style>
  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>
```

### Input Component

```svelte
<!-- frontend/src/components/Input.svelte -->
<script>
  export let label = '';
  export let type = 'text';
  export let value = '';
  export let error = '';
  export let required = false;
  
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  
  function handleInput(event) {
    dispatch('input', event.target.value);
  }
</script>

<div class="input-group">
  {#if label}
    <label>
      {label}
      {#if required}<span class="required">*</span>{/if}
    </label>
  {/if}
  
  <input 
    type={type}
    value={value}
    oninput={handleInput}
    class:error={!!error}
    required={required}
  />
  
  {#if error}
    <span class="error-message">{error}</span>
  {/if}
</div>

<style>
  .input-group {
    margin-bottom: 1rem;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
  }
  
  .required {
    color: red;
  }
  
  input {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ddd;
    border-radius: 4px;
  }
  
  input.error {
    border-color: red;
  }
  
  .error-message {
    color: red;
    font-size: 0.875rem;
    margin-top: 0.25rem;
  }
</style>
```

### Header Component

```svelte
<!-- frontend/src/components/Header.svelte -->
<script>
  import { page, router } from '@inertiajs/svelte';
  import DarkModeToggle from './DarkModeToggle.svelte';
  
  const user = $page.props.user;
  
  function logout() {
    if (confirm('Are you sure you want to logout?')) {
      router.post('/logout');
    }
  }
</script>

<header class="header">
  <div class="container">
    <nav>
      <a href="/" class="logo">Laju Go</a>
      
      <div class="nav-links">
        {#if user}
          <a href="/app">Dashboard</a>
          <a href="/app/profile">Profile</a>
          <DarkModeToggle />
          <button onclick={logout}>Logout</button>
        {:else}
          <a href="/login">Login</a>
          <a href="/register">Register</a>
        {/if}
      </div>
    </nav>
  </div>
</header>

<style>
  .header {
    background: white;
    border-bottom: 1px solid #ddd;
    padding: 1rem 0;
  }
  
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
  }
  
  nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .logo {
    font-weight: bold;
    font-size: 1.5rem;
    color: #333;
    text-decoration: none;
  }
  
  .nav-links {
    display: flex;
    gap: 1rem;
    align-items: center;
  }
  
  .nav-links a {
    color: #333;
    text-decoration: none;
  }
  
  .nav-links a:hover {
    text-decoration: underline;
  }
</style>
```

## Dark Mode

### DarkModeToggle Component

```svelte
<!-- frontend/src/components/DarkModeToggle.svelte -->
<script>
  import { onMount } from 'svelte';
  import { Sun, Moon } from 'lucide-svelte';
  
  let isDark = $state(false);
  
  onMount(() => {
    // Check localStorage
    const stored = localStorage.getItem('theme');
    if (stored) {
      isDark = stored === 'dark';
    } else {
      // Check system preference
      isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    
    applyTheme();
  });
  
  function toggle() {
    isDark = !isDark;
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
    applyTheme();
  }
  
  function applyTheme() {
    if (isDark) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
</script>

<button onclick={toggle} class="theme-toggle" aria-label="Toggle theme">
  {#if isDark}
    <Sun size={20} />
  {:else}
    <Moon size={20} />
  {/if}
</button>

<style>
  .theme-toggle {
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    color: inherit;
  }
  
  .theme-toggle:hover {
    background: rgba(128, 128, 128, 0.1);
  }
</style>
```

### Global Styles with Dark Mode

```css
/* frontend/src/app.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --bg-primary: #ffffff;
  --text-primary: #1a1a1a;
}

.dark {
  --bg-primary: #1a1a1a;
  --text-primary: #ffffff;
}

body {
  background: var(--bg-primary);
  color: var(--text-primary);
}
```

## Utilities

### Helper Functions

```javascript
// frontend/src/lib/utils/helpers.js

// Click outside action
export function clickOutside(node, callback) {
  const handleClick = (event) => {
    if (node && !node.contains(event.target)) {
      callback();
    }
  };
  
  document.addEventListener('click', handleClick, true);
  
  return {
    destroy() {
      document.removeEventListener('click', handleClick, true);
    }
  };
}

// Debounce function
export function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}

// Get CSRF token from cookie
export function getCsrfToken() {
  const match = document.cookie.match(/csrf_token=([^;]+)/);
  return match ? match[1] : null;
}

// Fetch wrapper with CSRF token
export async function fetchWithCsrf(url, options = {}) {
  const csrfToken = getCsrfToken();
  
  const headers = {
    ...options.headers,
    'X-CSRF-Token': csrfToken,
  };
  
  if (options.body && !(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }
  
  const response = await fetch(url, {
    ...options,
    headers,
  });
  
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }
  
  return response;
}
```

### Using Actions

```svelte
<script>
  import { clickOutside } from '../lib/utils/helpers';
  
  let dropdownOpen = $state(false);
  
  function closeDropdown() {
    dropdownOpen = false;
  }
</script>

<div class="dropdown" use:clickOutside={closeDropdown}>
  <button onclick={() => dropdownOpen = !dropdownOpen}>
    Menu
  </button>
  
  {#if dropdownOpen}
    <div class="dropdown-content">
      <!-- Dropdown items -->
    </div>
  {/if}
</div>
```

## Toast Notifications

```svelte
<script>
  // frontend/src/lib/utils/toast.js
  
  let toasts = $state([]);
  
  export function toast(message, type = 'info', duration = 3000) {
    const id = Date.now();
    toasts = [...toasts, { id, message, type }];
    
    setTimeout(() => {
      toasts = toasts.filter(t => t.id !== id);
    }, duration);
  }
  
  export function getToasts() {
    return toasts;
  }
</script>
```

```svelte
<!-- Usage in component -->
<script>
  import { toast, getToasts } from '../lib/utils/toast';
  
  const toasts = getToasts();
  
  function handleSuccess() {
    toast('Operation successful!', 'success');
  }
  
  function handleError() {
    toast('Something went wrong', 'error');
  }
</script>

<div class="toast-container">
  {#each toasts as t (t.id)}
    <div class="toast toast-{t.type}">
      {t.message}
    </div>
  {/each}
</div>
```

## Layouts

### Base Layout

```svelte
<!-- frontend/src/layouts/AppLayout.svelte -->
<script>
  import Header from '../components/Header.svelte';
</script>

<div class="layout">
  <Header />
  <main class="main">
    <slot />
  </main>
</div>

<style>
  .layout {
    min-height: 100vh;
  }
  
  .main {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }
</style>
```

### Using Layouts

```svelte
<!-- frontend/src/pages/app/Dashboard.svelte -->
<script>
  import AppLayout from '../../layouts/AppLayout.svelte';
</script>

<AppLayout>
  <h1>Dashboard</h1>
  <p>Dashboard content here</p>
</AppLayout>
```

## Best Practices

### 1. Keep Components Small

```svelte
<!-- ✅ Good: Small, focused component -->
<script>
  export let user;
</script>

<div class="user-card">
  <img src={user.avatar} alt={user.name} />
  <h3>{user.name}</h3>
</div>

<!-- ❌ Bad: Large component with too much logic -->
```

### 2. Use TypeScript

```svelte
<!-- frontend/src/pages/app/Dashboard.svelte -->
<script lang="ts">
  interface User {
    id: number;
    name: string;
    email: string;
  }
  
  interface Props {
    user: User;
  }
  
  let { user }: Props = $props();
</script>
```

### 3. Extract Reusable Logic

```svelte
<!-- Use stores for shared state -->
<script>
  import { writable } from 'svelte/store';
  
  export const userStore = writable(null);
</script>
```

### 4. Handle Loading States

```svelte
<script>
  let loading = $state(false);
  
  async function loadData() {
    loading = true;
    try {
      // Fetch data
    } finally {
      loading = false;
    }
  }
</script>

{#if loading}
  <div>Loading...</div>
{:else}
  <!-- Content -->
{/if}
```

## Next Steps

- [Styling Guide](styling.md) - Tailwind CSS styling
- [Forms Guide](forms.md) - Form handling and validation
- [Inertia.js Guide](inertia.md) - Deep dive into Inertia.js
