# Laju Go Documentation

Welcome to the Laju Go documentation. This folder contains comprehensive guides, references, and tutorials for building applications with Laju Go.

## 📚 Documentation Structure

### Getting Started

New to Laju Go? Start here.

| Document | Description |
|----------|-------------|
| [Introduction](getting-started/introduction.md) | What is Laju Go, key features, and when to use it |
| [Installation](getting-started/installation.md) | Prerequisites, setup steps, and verification |
| [Configuration](getting-started/configuration.md) | Environment variables and .env setup |

### Guide

In-depth guides for building applications.

| Document | Description |
|----------|-------------|
| [Architecture](guide/architecture.md) | Layered architecture, design patterns, and best practices |
| [Routing](guide/routing.md) | Route definitions, middleware, and request handling |
| [Handlers](guide/handlers.md) | Building HTTP handlers, request/response handling |
| [Database](guide/database.md) | SQLite setup, migrations, and query building with Squirrel |
| [Authentication](guide/authentication.md) | Auth flows, OAuth, sessions, and password reset |
| [Frontend](guide/frontend.md) | Svelte 5 components and Inertia.js integration |
| [File Upload](guide/file-upload.md) | File handling, validation, and storage |
| [Email](guide/email.md) | SMTP configuration and email sending |
| [Middleware](guide/middleware.md) | Creating custom middleware |
| [Validation](guide/validation.md) | Input validation techniques |
| [Svelte 5](guide/svelte.md) | Svelte 5 features and patterns |
| [Inertia.js](guide/inertia.md) | SPA bridge with Inertia.js |
| [Forms](guide/forms.md) | Form handling and validation |
| [Styling](guide/styling.md) | Tailwind CSS styling |
| [Storage](guide/storage.md) | File storage management |
| [Caching](guide/caching.md) | Caching strategies |
| [Testing](guide/testing.md) | Testing strategies and examples |
| [Building](guide/building.md) | Build process and optimization |
| [Production](guide/production.md) | Production deployment guide |
| [CI/CD](guide/cicd.md) | Continuous integration and deployment |

### Deployment

Deployment guides for various environments.

| Document | Description |
|----------|-------------|
| [Development Workflow](deployment/development.md) | Hot reload, scripts, and development best practices |
| [One-Click Deployment](deployment/one-click-deployment.md) | **NEW** Automated deployment via SSH with one command |
| [Production Deployment](deployment/production.md) | Ubuntu/Debian deployment with systemd and Nginx |
| [Docker Deployment](deployment/docker.md) | Containerized deployment with Docker |
| [GitHub Actions CI/CD](deployment/github-actions.md) | Automated CI/CD pipeline |
| [Performance Optimization](deployment/optimization.md) | SQLite optimization and performance tuning |

### Reference

Complete reference documentation.

| Document | Description |
|----------|-------------|
| [API Reference](reference/api-reference.md) | Complete HTTP endpoint documentation |
| [Project Structure](reference/project-structure.md) | Directory layout and file organization |
| [Troubleshooting](reference/troubleshooting.md) | Common issues and solutions |
| [Environment Variables](reference/environment.md) | Complete environment variable reference |

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/maulanashalihin/laju-go.git
cd laju-go

# Install dependencies
go mod download && npm install

# Configure environment
cp .env.example .env

# Start development
npm run dev:all
```

Visit `http://localhost:8080` to see your application.

## 📖 Reading Order

### For Beginners

1. [Introduction](getting-started/introduction.md) - Understand what Laju Go is
2. [Installation](getting-started/installation.md) - Set up your environment
3. [Configuration](getting-started/configuration.md) - Configure your application
4. [Architecture](guide/architecture.md) - Learn the codebase structure
5. [Routing](guide/routing.md) - Define your first routes
6. [Handlers](guide/handlers.md) - Build your first handlers
7. [Database](guide/database.md) - Connect to database
8. [Authentication](guide/authentication.md) - Add user authentication

### For Experienced Users

1. [Project Structure](reference/project-structure.md) - Quick reference
2. [API Reference](reference/api-reference.md) - Endpoint documentation
3. [Deployment](deployment/production.md) - Deploy to production
4. [Optimization](deployment/optimization.md) - Performance tuning

## 🎯 Common Tasks

### Development

| Task | Guide |
|------|-------|
| Set up development environment | [Installation](getting-started/installation.md) |
| Configure environment variables | [Configuration](getting-started/configuration.md) |
| Run development servers | [Development Workflow](deployment/development.md) |
| Create new route | [Routing](guide/routing.md) |
| Build handler | [Handlers](guide/handlers.md) |
| Add database model | [Database](guide/database.md) |
| Implement authentication | [Authentication](guide/authentication.md) |

### Deployment

| Task | Guide |
|------|-------|
| **One-click deploy to server** | [**One-Click Deployment**](deployment/one-click-deployment.md) |
| Build for production | [Building](guide/building.md) |
| Deploy to Ubuntu server | [Production Deployment](deployment/production.md) |
| Deploy with Docker | [Docker Deployment](deployment/docker.md) |
| Setup CI/CD pipeline | [GitHub Actions CI/CD](deployment/github-actions.md) |
| Optimize performance | [Optimization](deployment/optimization.md) |
| Set up SSL/TLS | [Production Deployment](deployment/production.md) |

### Troubleshooting

| Problem | Guide |
|---------|-------|
| Port already in use | [Troubleshooting](reference/troubleshooting.md) |
| Database locked | [Troubleshooting](reference/troubleshooting.md) |
| Session not persisting | [Troubleshooting](reference/troubleshooting.md) |
| OAuth not working | [Troubleshooting](reference/troubleshooting.md) |
| Email not sending | [Troubleshooting](reference/troubleshooting.md) |

## 🔧 Resources

### External Links

- [Go Fiber Documentation](https://docs.gofiber.io/)
- [Svelte Documentation](https://svelte.dev/docs)
- [Inertia.js Documentation](https://inertiajs.com/)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Squirrel Query Builder](https://github.com/Masterminds/squirrel)
- [Goose Migrations](https://github.com/pressly/goose)

### Community

- [GitHub Repository](https://github.com/maulanashalihin/laju-go)
- [GitHub Issues](https://github.com/maulanashalihin/laju-go/issues)
- [GitHub Discussions](https://github.com/maulanashalihin/laju-go/discussions)

## 📝 Documentation Conventions

### File Naming

- Lowercase with hyphens: `project-structure.md`
- Descriptive and consistent names

### Code Examples

- Language identifier always specified
- Complete, runnable examples when possible
- Comments for complex logic

### Cross-Referencing

- Relative links between documents
- Clear section anchors
- Descriptive link text

## 🤝 Contributing to Documentation

Found an error or want to improve the documentation?

1. Fork the repository
2. Make your changes
3. Test any code examples
4. Submit a pull request

### Documentation Guidelines

- Write in clear, concise English
- Use active voice
- Include examples for complex topics
- Add cross-references to related documents
- Keep formatting consistent

## 📄 License

Documentation is licensed under the MIT License - same as Laju Framework.

---

**Last Updated**: March 2026  
**Version**: 1.0.0
