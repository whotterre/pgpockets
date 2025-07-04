package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

const (
	InvoiceStatusDraft         string = "draft"
	InvoiceStatusSent          string = "sent"
	InvoiceStatusPending       string = "pending"
	InvoiceStatusPartiallyPaid string = "partially_paid"
	InvoiceStatusPaid          string = "paid"
	InvoiceStatusOverdue       string = "overdue"
	InvoiceStatusCancelled     string = "cancelled"
)

const (
	TransactionTypeDeposit    string = "deposit"
	TransactionTypeWithdrawal string = "withdrawal"
	TransactionTypeTransfer   string = "transfer"
	TransactionTypePayment    string = "payment"
	TransactionTypeFee        string = "fee"
	TransactionTypeRefund     string = "refund"
	TransactionTypeChargeback string = "chargeback"
)

const (
	TransactionStatusPending   string = "pending"
	TransactionStatusCompleted string = "completed"
	TransactionStatusFailed    string = "failed"
	TransactionStatusReversed  string = "reversed"
	TransactionStatusRefunded  string = "refunded"
	TransactionStatusCancelled string = "cancelled"
)

const (
	CurrencyUSD string = "USD"
	CurrencyNGN string = "NGN"
	CurrencyEUR string = "EUR"
	CurrencyGBP string = "GBP"
	CurrencyCAD string = "CAD"
	CurrencyAUD string = "AUD"
	CurrencyJPY string = "JPY"
)

const (
	CardTypeMastercard string = "mastercard"
	CardTypeVisa       string = "visa"
	CardTypeAmex       string = "amex"
	CardTypeDiscover   string = "discover"
	CardTypeVerve      string = "verve"
)

type Beneficiary struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null references" json:"user_id"`
	Description string    `gorm:"type:text" json:"description"`
	WalletID    uuid.UUID `gorm:"type:uuid;not null" json:"wallet_id"`
	Wallet      Wallet    `gorm:"foreignKey:WalletID;references:ID"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}

// User represents the users table in the database.
type User struct {
	ID              uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email           string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash    string    `gorm:"type:varchar(255);not null" json:"-"`
	IsEmailVerified bool      `gorm:"default:false" json:"is_email_verified"`
	CreatedAt       time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt       time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// Profile represents the profiles table in the database.
type Profile struct {
	ID          uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;unique;not null" json:"user_id"`
	FirstName   string     `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName    string     `gorm:"type:varchar(255);not null" json:"last_name"`
	DateOfBirth *time.Time `gorm:"type:date" json:"date_of_birth"`
	PhoneNumber string     `gorm:"type:varchar(20);unique" json:"phone_number"`
	Address     string     `gorm:"type:text" json:"address"`
	Gender      string     `gorm:"type:varchar(10)" json:"gender"`
	CreatedAt   time.Time  `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"not null;default:now()" json:"updated_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Wallet represents the wallets table in the database.
type Wallet struct {
	ID        uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null" json:"user_id"`
	Currency  string          `gorm:"type:varchar(3);not null;default:NGN" json:"currency"`
	Balance   decimal.Decimal `gorm:"type:numeric(18,2);not null;default:0.00" json:"balance"`
	Name      string          `gorm:"type:varchar(255);default:Naira Wallet" json:"name"`
	IsActive  bool            `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time       `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time       `gorm:"not null;default:now()" json:"updated_at"`

	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Invoice represents the invoices table in the database.
type Invoice struct {
	ID            uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	InvoiceNumber string    `gorm:"type:varchar(50);unique;not null" json:"invoice_number"`
	SenderID      uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID    uuid.UUID `gorm:"type:uuid;not null" json:"receiver_id"`
	IssueDate     time.Time `gorm:"type:date;not null" json:"issue_date"`
	DueDate       time.Time `gorm:"type:date;not null" json:"due_date"`
	TotalAmount   string    `gorm:"type:decimal(18,2);not null" json:"total_amount"`         // Use string for DECIMAL
	Currency      string    `gorm:"type:varchar(3);not null" json:"currency"`                // Maps to CurrencyCode enum
	Status        string    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // Maps to InvoiceStatus enum
	Description   string    `gorm:"type:text" json:"description"`
	PaymentTerms  string    `gorm:"type:varchar(255)" json:"payment_terms"`
	CreatedAt     time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	Sender   User          `gorm:"foreignKey:SenderID;references:ID"`
	Receiver User          `gorm:"foreignKey:ReceiverID;references:ID"`
	Items    []InvoiceItem `gorm:"foreignKey:InvoiceID;references:ID"`
}

// InvoiceItem represents the invoice_items table in the database.
type InvoiceItem struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	InvoiceID   uuid.UUID `gorm:"type:uuid;not null" json:"invoice_id"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Quantity    string    `gorm:"type:decimal(18,4);not null" json:"quantity"`       // Use string for DECIMAL
	UnitPrice   string    `gorm:"type:decimal(18,2);not null" json:"unit_price"`     // Use string for DECIMAL
	LineTotal   string    `gorm:"type:decimal(18,2);not null" json:"line_total"`     // Use string for DECIMAL
	TaxRate     string    `gorm:"type:decimal(5,4);default:0.00" json:"tax_rate"`    // Use string for DECIMAL
	TaxAmount   string    `gorm:"type:decimal(18,2);default:0.00" json:"tax_amount"` // Use string for DECIMAL
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	Invoice Invoice `gorm:"foreignKey:InvoiceID;references:ID"`
}

// Card represents the cards table in the database.
type Card struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CardToken      string    `gorm:"type:varchar(255);unique;not null" json:"card_token"`
	LastFourDigits string    `gorm:"type:varchar(4);not null" json:"last_four_digits"`
	CardType       string    `gorm:"type:varchar(20);not null" json:"card_type"`   // Maps to CardType enum
	ExpiryMonth    string    `gorm:"type:varchar(2);not null" json:"expiry_month"` // MM
	ExpiryYear     string    `gorm:"type:varchar(4);not null" json:"expiry_year"`  // YYYY
	CardBrand      string    `gorm:"type:varchar(50)" json:"card_brand"`
	BankName       string    `gorm:"type:varchar(255)" json:"bank_name"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Transaction represents the transactions table in the database.
type Transaction struct {
	ID               uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SenderWalletID   *uuid.UUID `gorm:"type:uuid" json:"sender_wallet_id"`
	ReceiverWalletID *uuid.UUID `gorm:"type:uuid" json:"receiver_wallet_id"`

	Amount          string `gorm:"type:decimal(18,2);not null" json:"amount"`
	Currency        string `gorm:"type:varchar(3);not null" json:"currency"`
	TransactionType string `gorm:"type:varchar(50);not null" json:"transaction_type"`
	Status          string `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	Description string `gorm:"type:text" json:"description"`
	ReferenceID string `gorm:"type:varchar(255);unique" json:"reference_id"`

	MadeAt    time.Time `gorm:"not null;default:now()" json:"made_at"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`

	SenderWallet   *Wallet `gorm:"foreignKey:SenderWalletID;references:ID"`
	ReceiverWallet *Wallet `gorm:"foreignKey:ReceiverWalletID;references:ID"`
}

// Session represents the sessions table in the database.
type Session struct {
	ID                    uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID                uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ClientIP              string    `gorm:"type:varchar(45)" json:"client_ip"`
	UserAgent             string    `gorm:"type:text" json:"user_agent"`
	AccessToken           string    `gorm:"type:text;not null" json:"access_token"`
	AccessTokenExpiresAt  time.Time `gorm:"not null" json:"access_token_expires_at"`
	RefreshToken          string    `gorm:"type:text;not null" json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `gorm:"not null" json:"refresh_token_expires_at"`
	IsActive              bool      `gorm:"default:true" json:"is_active"`
	CreatedAt             time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt             time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// Associations
	User User `gorm:"foreignKey:UserID;references:ID"`
}

type Notification struct {
	Title       string    `gorm:"type:text;not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	RecipientID uuid.UUID `gorm:"type:uuid;not null" json:"recipient_id"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
}
