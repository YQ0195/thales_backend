#  Setup Instructions for Thales Backend

This guide will help you run the Thales backend project locally using Docker.

---

## Prerequisites

- Install [Docker](https://www.docker.com/products/docker-desktop/)
- Clone this repository

```bash
git clone this repository
```
cd thales_backend

```

1. Build containers

```bash
docker-compose build
```
2. Start the database and run migrations

```bash
  
  
docker-compose up -d db
docker-compose run --rm migrate
```
3. Start the backend

```bash
  
docker-compose up -d backend

```
