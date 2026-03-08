# Aureus - Real-Time Auction API 🔨

## 📖 Description

Aureus is a backend REST API system for an online auction platform. It handles user authentication, digital wallets, auction item management, and a real-time bidding system.

This system is designed to safely handle money transactions using an escrow system and provide real-time updates to users without overloading the server.

## 🎯 Why I Built This

I built this project mainly to **learn and practice advanced backend engineering**. I wanted to step out of my comfort zone of making simple CRUD applications and learn how to solve real-world problems.

Through this project, I learned how to:

- Build a solid **Clean Architecture** to make the code easy to maintain.
- Use **Dependency Injection** (Google Wire) to manage project dependencies.
- Prevent money loss and race conditions using **Database Transactions & Pessimistic Locking**.
- Build a fast, real-time leaderboard using **Redis Sorted Sets (ZSET)** and **Server-Sent Events (SSE)** instead of heavy WebSockets.
- Automate background tasks using **Redis Keyspace Notifications** and Go Background Workers.

## ✨ Key Features

- **User & Wallet System:** Users can register, log in, and have a wallet with `Active Balance` and `Held Balance` (Escrow) to make sure bids are safely locked.
- **Auction Management:** Sellers can create auctions, upload images, and set a start/end time and bid increments.
- **Safe Bidding System:** When a user places a bid, the system locks the database row, checks the wallet balance, deducts the money, and automatically refunds the previous highest bidder.
- **Real-Time Leaderboard:** Users can see the top 10 highest bidders live, powered by Redis and Server-Sent Events (SSE).
- **Automated Settlement:** When an auction time is up, a background worker automatically finalizes the auction and transfers the locked money to the seller.

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Web Framework:** Fiber v3
- **Database:** PostgreSQL
- **ORM:** GORM
- **Payment Gateway:** Stripe
- **Cache:** Redis (go-redis/v9)
- **File Storage:** SeaweedFS
- **Dependency Injection:** Google Wire
- **JSON Parser:** Sonic

## 🚀 How to Run (Local Development)

**1. Clone the repository**

```bash
git clone https://github.com/Hilaladiii/aureus.git
```

**2. Install The Dependency**

```bash
go mod tidy
```

**3. Setup Environment Variables**

Create a .env file in the root directory and add your configurations (Database, Redis, JWT Secret, etc.).

```bash
cp .env.example .env
```

**4. Generate Dependency Injection (Wire)**

```bash
wire gen ./di
```

**5. Run Docker (DB,S3,Redis)**

```bash
docker compose up -d
```

**6. Run The Application**

If you have Air installed for live-reloading:

```bash
air
```

Otherwise, use the standard Go command:

```bash
go run ./cmd/api main.go
```

Built with ❤️ for learning and exploring Backend Engineering.
