# Introduction

Welcome to **Laju Go** — a high-performance SaaS boilerplate that combines the speed of Go with the developer experience of modern frontend frameworks.

## What is Laju Go?

Laju Go is a production-ready web application starter kit built with:

- **Go Fiber** - Fasthttp-based web framework for blazing-fast backend performance
- **Svelte 5** - Reactive frontend framework with compile-time optimization
- **SQLite** - Embedded database with production-grade optimizations
- **Inertia.js** - Bridge between backend and frontend without building a separate API

This stack gives you the performance of a traditional server-rendered application with the rich interactivity of a single-page app — all without the complexity of building and maintaining a separate API.

## Key Features

### 🔐 Authentication & Security
- Email/password authentication with bcrypt hashing
- Google OAuth 2.0 integration
- Password reset via email with secure tokens
- Database-backed persistent sessions
- CSRF protection middleware
- Rate limiting for sensitive endpoints

### 👥 User Management
- Role-based access control (Admin/User)
- User profile management
- Avatar upload with validation
- Session persistence across server restarts

### 🎨 Frontend Experience
- Svelte 5 with reactive components
- Inertia.js for SPA-like navigation
- Tailwind CSS for styling
- Dark mode support
- Responsive design out of the box
- Toast notifications
- Form validation with error display

### 🚀 Development Tools
- Hot Module Replacement (Vite)
- Go hot reload (Air)
- TypeScript support
- Component-based architecture
- Automatic asset versioning

### 📦 Production Ready
- SQLite with WAL mode and connection pooling
- Database migrations with Goose
- Docker multi-stage builds
- Systemd service configuration
- Nginx reverse proxy ready
- SSL/TLS support

## When to Use Laju Go

Laju Go is perfect for:

✅ **SaaS Applications** - Multi-tenant software as a service  
✅ **Internal Tools** - Admin panels, dashboards, CRUD apps  
✅ **MVPs** - Quick prototyping with production-ready code  
✅ **Small to Medium Apps** - User management, authentication, file uploads  
✅ **Solo Developers** - Full-stack development without context switching  

## When NOT to Use Laju Go

Consider other solutions if you need:

❌ **Microservices Architecture** - Laju Go is a monolithic application  
❌ **Real-time Features** - No WebSocket support out of the box  
❌ **High-Scale Databases** - SQLite may not suit high-concurrency needs (though it's surprisingly capable)  
❌ **GraphQL API** - Laju Go uses traditional HTTP routing  
❌ **Mobile Backend** - While possible, it's optimized for server-rendered web apps  

## Architecture Overview

Laju Go follows a **layered architecture** pattern:

```
┌─────────────────┐
│     Routes      │ → Route definitions and middleware setup
└────────┬────────┘
         │
┌────────▼────────┐
│   Middleware    │ → Auth, CSRF, rate limiting
└────────┬────────┘
         │
┌────────▼────────┐
│    Handlers     │ → HTTP request/response handling
└────────┬────────┘
         │
┌────────▼────────┐
│    Services     │ → Business logic
└────────┬────────┘
         │
┌────────▼────────┐
│  Repositories   │ → Database queries
└────────┬────────┘
         │
┌────────▼────────┐
│    Database     │ → SQLite storage
└─────────────────┘
```

This separation of concerns makes the codebase:
- **Maintainable** - Clear boundaries between layers
- **Testable** - Each layer can be tested independently
- **Scalable** - Easy to add new features without breaking existing code

## Technology Choices

### Why Go Fiber?

Fiber is inspired by Express.js and offers:
- **Performance** - Built on fasthttp, one of the fastest Go HTTP libraries
- **Familiar API** - Express-like syntax for easy onboarding
- **Middleware Support** - Rich ecosystem of middleware
- **Low Memory Footprint** - Efficient memory management

### Why Svelte 5?

Svelte compiles components at build time, resulting in:
- **Smaller Bundles** - No runtime framework overhead
- **Faster Runtime** - Direct DOM manipulation without virtual DOM
- **Reactive by Default** - Simple, intuitive state management
- **TypeScript Support** - Full type safety

### Why SQLite?

Despite being "just a file", SQLite offers:
- **Zero Configuration** - No database server to manage
- **Portability** - Single file database
- **Performance** - Faster than PostgreSQL for read-heavy workloads
- **Reliability** - Used in production by companies like Dropbox and Trello
- **WAL Mode** - Write-Ahead Logging for better concurrency

### Why Inertia.js?

Inertia bridgesges backend and frontend:
- **No API Required** - Return data directly from controllers
- **Server-Driven** - Backend controls routing and business logic
- **SPA Experience** - Client-side navigation without page reloads
- **Simple Mental Model** - Traditional HTTP requests with modern UX

## Next Steps

- [Installation](installation.md) - Quick start with CLI or manual setup
- [Configuration](configuration.md) - Configure environment variables
- [Architecture Guide](../guide/architecture.md) - Deep dive into the codebase

## Community & Support

- **GitHub**: [maulanashalihin/laju-go](https://github.com/maulanashalihin/laju-go)
- **Issues**: [Report a bug](https://github.com/maulanashalihin/laju-go/issues)
- **Discussions**: [Ask questions](https://github.com/maulanashalihin/laju-go/discussions)
- **CLI Tool**: [create-laju-go](https://www.npmjs.com/package/create-laju-go) - Quick project scaffolding
