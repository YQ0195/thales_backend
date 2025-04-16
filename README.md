#  Setup Instructions for Thales Backend

This guide will help you run the Thales backend project locally using Docker.

---

## Prerequisites

- Install [Docker](https://www.docker.com/products/docker-desktop/)
- Clone this repository

```bash
git clone https://github.com/your-username/thales_backend.git
cd thales_backend

```

1. Build and start the containers

```bash
docker-compose up -d --build
```

2. (Important) Might need to run migration again after DB is ready
   ⚠️ On the first run, the migration container may start before the database is ready.
