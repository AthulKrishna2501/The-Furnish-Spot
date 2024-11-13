# The-Furniture-Spot

An eCommerce platform for furniture, built with Golang, GORM, and Gin, featuring PostgreSQL for data storage. The platform includes a wide range of functionalities from product management to payment processing.

## Table of Contents
- [About the Project](#about-the-project)
- [Features](#features)
  - [User Side](#user-side)
  - [Admin Side](#admin-side)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Project Structure](#project-structure)

## About the Project
The-Furniture-Spot is a backend-focused eCommerce platform that offers a range of features tailored for an online furniture store. It includes inventory management, order tracking, payment integration with PayPal, and a comprehensive admin dashboard to handle store operations.

### User Side
- **Signup and Login**: Users can create an account or log into the platform using their email and password.
- **OTP Verification**: After login or signup, an OTP is sent to the user's registered phone number for verification.
- **Captcha Verification**: Users are required to complete a captcha challenge during signup or login to prevent bot activity.
- **Product Listings**: View and browse available furniture products with filtering and sorting options.
- **Cart and Wishlist**: Add items to the cart or wishlist for future purchase.
- **Order Management**: Track orders, view status, and cancel orders if necessary.
- **Coupons and Offers**: Apply discounts and offers during checkout.
- **Invoice Generation**: After completing a purchase, users can generate invoices that detail the order summary, payment method, and any applied discounts. The invoice is available for download in PDF format.

### Admin Side
- **Product and Inventory Management**: Add, edit, and delete products, monitor stock levels.
- **Order Management**: Track orders, update statuses, cancel orders.
- **Sales Reports**: Generate reports with daily, weekly, or custom date ranges, showing discounts, order totals, and more.
- **Analytics**: View overall sales, order count, and discount summaries.

## Technologies Used
- **Backend**: Golang, Gin (framework), GORM (ORM)
- **Database**: PostgreSQL
- **Payments**: PayPal integration
- **PDF Generation**: GoPDF for generating sales reports

## Getting Started
### Prerequisites
- Go (version 1.20+)
- PostgreSQL (version 13+)
- Docker (optional, for containerized setup)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/AthulKrishna2501/The-Furniture-Spot.git
   cd The-Furniture-Spot


### Project Structure
- ** The-Furniture-Spot/
- **│
- **├── controllers/        # API endpoint handlers
- **├── models/             # Data models for the database
- **├── routes/             # API route definitions
- **├── services/           # Business logic
- **├── utils/              # Helper functions (logging, currency conversion, etc.)
- **├── main.go             # Application entry point
- **├── README.md           # Project README file
- **└── .env.example        # Environment variables example file


