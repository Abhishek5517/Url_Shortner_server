# Go-Fiber URL Shortener

ReLink a **URL shortener service** built with [Go Fiber](https://gofiber.io/).  
It includes **JWT-based authentication**, **PostgreSQL** for persistent storage, and a **Redis-powered hit counter** (via Pub/Sub) to track live URL hits in real time.

---

##  Features
- Shorten long URLs into unique short codes
- Authentication system with **Login/Signup** using JWT
- **PostgreSQL** for storing user accounts and shortened URLs
-  **Redis Pub/Sub** for live URL hit counting
- Fast and lightweight server with **Fiber**
- RESTful APIs for easy integration

---

## Tech Stack
- **Backend Framework:** [Go Fiber](https://gofiber.io/)  
- **Database:** [PostgreSQL](https://www.postgresql.org/)  
- **Cache & Realtime Counter:** [Redis](https://redis.io/) (Pub/Sub)
- **Real-Time Updates:** Server-Sent Events (SSE)  
- **Auth:** JWT (JSON Web Tokens)  
- **ORM/DB Driver:** [pgx](https://github.com/jackc/pgx)  

---

## Environment Variables

```env
DB_URL=postgres://postgres:Password@localhost:5432/shortUrlDB?sslmode=disable
SERVER_PORT=8080
SECRET_KEY=secret_key
REDIS_URL=localhost:6379

