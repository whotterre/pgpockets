FinPay Backend
A FinTech Solution built with Go, Fiber, and PostgreSQL

This repository contains the backend services for FinPay, a comprehensive FinTech application designed to empower users with robust financial management tools. The application will enable users to create invoices, generate virtual cards, set up virtual accounts, and process payments in multiple currencies.

🚀 Project Overview
The FinPay project aims to build a high-performance and secure backend system capable of handling various financial operations. This solution is developed using Go, leveraging the speed of the Fiber web framework and the reliability of PostgreSQL as its primary data store.

Key Features (Planned)
User Authentication & Authorization: Secure login, registration, password management, and access control.

Dashboard: Overview of balances, payment accounts, invoice summaries, exchange rates, and virtual assets.

Invoicing: Create, view, manage, and delete invoices with various statuses (draft, pending, due, overdue).

Cards: Generation and management of virtual cards.

Wallets: Multi-currency balance management, account statements, currency conversion, sending/receiving money, funding, and withdrawals.

Transactions: Comprehensive view, search, filter, and pagination of all financial transactions.

User Profile: Management of profile details, beneficiaries, and 2FA activation.

Notifications: Real-time alerts and management.

🛠️ Technologies Used
Go: The primary programming language for high-performance backend services.

Fiber: An Express-inspired web framework for Go, built on Fasthttp, known for its speed.

PostgreSQL: A powerful, open-source relational database for reliable data storage.

jackc/pgx: A high-performance PostgreSQL driver for Go.

jmoiron/sqlx: An extension to Go's database/sql for easier struct mapping to query results.

golang-migrate/migrate: For managing database schema migrations.

golang.org/x/crypto/bcrypt: For secure password hashing.

golang-jwt/jwt/v5: For JSON Web Token (JWT) based authentication.

go-playground/validator: For robust request input validation.

joho/godotenv: For loading environment variables from .env files.

Mono API / Open Banking: (Future Integration) External FinTech APIs for core functionalities like virtual cards and currency exchange.

💡 What You Will Learn (or are Learning!)
Building robust backend systems with Go.

Implementing secure authentication flows (registration, login, JWTs).

Database design and interaction with PostgreSQL.

Effective API design principles.

Error handling, data validation, and security best practices in a FinTech context.

(Future) Integrating with external FinTech APIs.

📦 Getting Started
Follow these steps to get your local development environment set up and running.

Prerequisites
Before you begin, ensure you have the following installed:

Go: Version 1.20 or newer.

PostgreSQL: Running locally (e.g., via Docker, Homebrew, or direct installation).

Git: For version control.

Postman/Insomnia/curl: For API testing.

1. Clone the Repository
git clone https://github.com/yourusername/finpay-backend.git
cd finpay-backend

2. Set Up Environment Variables
Create a file named .env (or app.env) in the root of your project directory and add the following:

# Database Configuration
DATABASE_URL="postgres://user:password@localhost:5432/finpay_db?sslmode=disable"
# Replace user, password, and finpay_db with your actual PostgreSQL credentials and database name.

# JWT Configuration
JWT_SECRET="your_very_secret_jwt_key_here" # Change this to a strong, random key in production

Important: Replace user, password, finpay_db, and your_very_secret_jwt_key_here with your actual credentials and a strong secret.

3. Database Setup
a. Create Database & User (if not already done)
Access your PostgreSQL client and run:

CREATE DATABASE finpay_db;
CREATE USER finpay_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE finpay_db TO finpay_user;

(Replace finpay_db, finpay_user, and your_secure_password as needed, matching your .env file.)

b. Run Migrations
Install the migration tool:

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

Generate your initial user table migration (this creates the file, you'll fill it in):

migrate create -ext sql -dir database/migrations -seq create_users_table

Then, edit the newly created .sql files in database/migrations (e.g., XXXXXXXXXXXXXX_create_users_table.up.sql and XXXXXXXXXXXXXX_create_users_table.down.sql) to define your users table:

XXXXXXXXXXXXXX_create_users_table.up.sql:

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

XXXXXXXXXXXXXX_create_users_table.down.sql:

DROP TABLE users;

Finally, run the migrations:

migrate -path database/migrations -database "$DATABASE_URL" up
# Or directly: migrate -path database/migrations -database "postgres://user:password@localhost:5432/finpay_db?sslmode=disable" up

4. Install Go Dependencies
go mod tidy

5. Run the Application
go run main.go

The server should start on http://localhost:3000 (or whatever port you configure in Fiber).

💡 Current Development Focus (This Week: June 23rd - June 29th)
This week, the primary focus is on establishing the foundational backend infrastructure and implementing the core user authentication functionalities:

Project Initialization: Setting up the Go module, Fiber app, and basic project structure.

Database Integration: Connecting to PostgreSQL using pgx and sqlx.

Database Migrations: Setting up and running initial migrations for the users table.

User Registration: Implementing the API endpoint (POST /api/v1/register) with input validation, password hashing (bcrypt), and saving user data to the database.

User Login: Implementing the API endpoint (POST /api/v1/login) with credential verification and JSON Web Token (JWT) generation.

📝 API Endpoints (Under Development)
Authentication
POST /api/v1/register

Body: {"email": "user@example.com", "password": "SecurePassword123!"}

Description: Registers a new user.

POST /api/v1/login

Body: {"email": "user@example.com", "password": "SecurePassword123!"}

Description: Authenticates a user and returns a JWT token.

🤝 Contributing
Contributions are welcome! Please feel free to fork the repository, make changes, and submit pull requests.

📄 License
[Specify your project's license here, e.g., MIT, Apache 2.0]