<div align="center">

<img src="./assets/img (2).png" alt="ProductBid" width="100%"/>

<br/>
<br/>

## Bid Higher. Get Seen.

<br/>

An open source platform where startups compete for visibility
and users discover the products worth watching.


</div>

---

## 🚀 About

ProductBid is a simple bidding platform built around one idea:

> **The highest bid gets the top spot.**

No accounts. No unnecessary complexity. Just products, bids, rankings, and
a transparent marketplace for visibility.

This repository contains the backend API responsible for products, categories,
bidding, rankings, payments, and database management.

---

## ✨ Features

- 🚀 Product submission and management
- 💰 Product visibility through bidding
- 🏆 Real-time leaderboard rankings
- 📊 Category-based product discovery
- 💳 Payment integration


</div>

---

<br/>


# ⚡ Getting Started

Follow these steps to set up and run the ProductBid backend locally.

## Prerequisites

Make sure you have the following installed:

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- Git

Verify your installation:

```bash
go version
docker --version
git --version



## Clone the Repository

```bash
git clone https://github.com/YOUR_USERNAME/productbid.git
cd productbid/Backend
```

## Configure Environment Variables

Create a `.env` file in the `Backend` directory:

```bash
cp .env.example .env
```

Then update the environment variables according to your local setup:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=productbid
```

> **Note:** Never commit your `.env` file. Keep credentials and secrets private.

## Start PostgreSQL

ProductBid uses PostgreSQL for data storage and Docker for the local database environment.

From the `Backend` directory, start PostgreSQL:

```bash
docker compose up -d
```

Verify that the PostgreSQL container is running:

```bash
docker ps
```

To stop the database:

```bash
docker compose down
```

## Install Dependencies

Download the required Go dependencies:

```bash
go mod download
```

## Run the Backend

Start the backend server:

```bash
go run main.go
```

If everything is configured correctly, you should see:

```text
Database connected successfully
Migration completed successfully!!
Server starting on port 8080
```

The API will be available at:

```text
http://localhost:8080
```