# FinPay Backend
A FinTech Solution built with Go, Fiber, GORM, and PostgreSQL

This repository contains the backend services for FinPay, a comprehensive FinTech application designed to empower users with robust financial management tools. The application enables users to create invoices, generate virtual cards, set up virtual accounts, and process payments in multiple currencies.  
[Challenge](https://app.masteringbackend.com/projects/build-fin-pay-a-fin-tech-solution)

This project was initiated to explore the Go Fiber framework and GORM ORM. It also adopts a maintainable (Repository, Service, Handler) architecture while adhering to 12-Factor App principles.

🚀 **Project Overview**  
The FinPay project aims to build a high-performance and secure backend system capable of handling various financial operations. This solution is developed using Go, leveraging the speed of the Fiber web framework, the flexibility of GORM for ORM, and the reliability of PostgreSQL as its primary data store.

### Key Features (Planned)
- **User Authentication & Authorization**: Secure login, registration, password management, and access control.
- **Dashboard**: Overview of balances, payment accounts, invoice summaries, exchange rates, and virtual assets.
- **Invoicing**: Create, view, manage, and delete invoices with various statuses (draft, pending, due, overdue).
- **Cards**: Generation and management of virtual cards.
- **Wallets**: Multi-currency balance management, account statements, currency conversion, sending/receiving money, funding, and withdrawals.
- **Transactions**: Comprehensive view, search, filter, and pagination of all financial transactions.
- **User Profile**: Management of profile details, beneficiaries, and 2FA activation.
- **Notifications**: Real-time alerts and management.

🛠️ **Technologies Used**
- **Go**: The primary programming language for high-performance backend services.
- **Fiber**: An Express-inspired web framework for Go, built on Fasthttp, known for its speed.
- **GORM**: A powerful ORM library for Go, simplifying database interactions.
- **PostgreSQL**: A powerful, open-source relational database for reliable data storage.
- **golang.org/x/crypto/bcrypt**: For secure password hashing.
- **golang-jwt/jwt/v5**: For JSON Web Token (JWT) based authentication.
- **go-playground/validator**: For robust request input validation.
- **spf13/viper**: For loading environment variables from `.env` files and maintaining central config.
- **Mono API / Open Banking**: (Future Integration) External FinTech APIs for core functionalities like virtual cards and currency exchange.

💡 **What I Intend to Learn with This**
- Building robust backend systems with Go.
- Implementing secure authentication flows (registration, login, JWTs).
- Database design and interaction with PostgreSQL using GORM.
- Effective API design principles.
- Error handling, data validation, and security best practices in a FinTech context.
- (Future) Integrating with external FinTech APIs.

📦 **Getting Started**  
Follow these steps to set up your local development environment.

### Prerequisites
Ensure you have the following installed:
- **Go**: Version 1.20 or newer.
- **PostgreSQL**: Running locally (e.g., via Docker, Homebrew, or direct installation).
- **Git**: For version control.
- **Postman/Insomnia/curl**: For API testing.

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/pgpockets.git
cd pg-pockets
```

### 2. Set Up Environment Variables
Create a `.env` file in the root of your project directory and add the following:
```env
# Database Configuration
DATABASE_URL="postgres://user:password@localhost:5432/finpay_db?sslmode=disable"

# JWT Configuration
JWT_SECRET="your_very_secret_jwt_key_here"
```
Replace `user`, `password`, `finpay_db`, and `your_very_secret_jwt_key_here` with your actual credentials and a strong secret.

### 3. Database Setup
#### a. Create Database & User
Access your PostgreSQL client and run:
```sql
CREATE DATABASE finpay_db;
CREATE USER finpay_user WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE finpay_db TO finpay_user;
```
Ensure these match your `.env` file.

#### b. Run Migrations
Define your database schema using GORM's auto-migration feature:
```go
db.AutoMigrate(&User{})
```
Example `User` model:
```go
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Email        string    `gorm:"unique;not null"`
    PasswordHash string    `gorm:"not null"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}
```

### 4. Install Go Dependencies
```bash
go mod tidy
```

### 5. Run the Application
```bash
go run main.go
```
The server should start on `http://localhost:4000` (or the port configured in Fiber).

💡 **Current Development Focus**  
This week, the primary focus is on establishing the foundational backend infrastructure and implementing core user authentication functionalities:
- **Project Initialization**: Setting up the Go module, Fiber app, and basic project structure.
- **Database Integration**: Connecting to PostgreSQL using GORM.
- **User Registration**: Implementing the API endpoint (`POST /api/v1/register`) with input validation, password hashing (bcrypt), and saving user data to the database.
- **User Login**: Implementing the API endpoint (`POST /api/v1/login`) with credential verification and JWT generation.

📝 **API Endpoints (Under Development)**  
#### Authentication
- **POST /api/v1/register**  
  Body: `{"email": "user@example.com", "password": "SecurePassword123!"}`  
  Description: Registers a new user.

- **POST /api/v1/login**  
  Body: `{"email": "user@example.com", "password": "SecurePassword123!"}`  
  Description: Authenticates a user and returns a JWT token.

🤝 **Contributing**  
Contributions are welcome! Please feel free to fork the repository, make changes, and submit pull requests.
