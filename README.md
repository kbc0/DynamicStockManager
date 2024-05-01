# DynamicStockManager

DynamicStockManager is a comprehensive management system designed for managing dynamic stock data and forms. It supports user authentication and allows for detailed manipulation of forms, fields, and stock data through a RESTful API. This system is particularly tailored for the needs of the Perwatch case study but can be adapted for broader usage.

## Features

- User Registration and Authentication
- CRUD operations on forms
- CRUD operations on fields within forms
- CRUD operations on stock data linked to forms

## Getting Started

### Prerequisites

- MongoDB
- Go (at least version 1.15)
- Fiber v2 for the backend framework

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/DynamicStockManager.git
   cd DynamicStockManager
	
2. Set up your MongoDB database and ensure it is running.

3. Configure your environment variables (I didn't implement them because of simplicity, but the best practice is of course setting up an env var, for now you can change the URL from main.go):
   - `MONGO_URI`: Your MongoDB connection string.
   - `PORT`: Port number for the server to listen on (default `8080`).

4. Install Go dependencies:
   ```bash
   go mod tidy

5. Start the server:
   ```bash
   go run main.go

## API Documentation

### User Related APIs

- **Register User**
  - `POST /api/v1/register`
- **Login User**
  - `POST /api/v1/login`
- **Get User Account Info**
  - `GET /api/v1/account`

### Form Related APIs

- **Create Form**
  - `POST /api/v1/form/create`
- **List All Forms**
  - `GET /api/v1/form`
- **Get Specific Form**
  - `GET /api/v1/form/:_id`
- **Update Specific Form**
  - `PUT /api/v1/form/:_id`
- **Delete Specific Form**
  - `DELETE /api/v1/form/:_id`

### Field Related APIs

- **Add Field to Form**
  - `POST /api/v1/form/:_id/field`
- **List All Fields in Form**
  - `GET /api/v1/form/:_id/field`
- **Get Specific Field**
  - `GET /api/v1/form/:_id/field/:field_id`
- **Update Specific Field**
  - `PUT /api/v1/form/:_id/field/:field_id`
- **Delete Specific Field**
  - `DELETE /api/v1/form/:_id/field/:field_id`

### Stock Related APIs

- **Add Stock to Form**
  - `POST /api/v1/form/:_id/stock`
- **List All Stocks in Form**
  - `GET /api/v1/form/:_id/stock`
- **Get Specific Stock**
  - `GET /api/v1/form/:_id/stock/:stock_id`
- **Update Specific Stock**
  - `PUT /api/v1/form/:_id/stock/:stock_id`
- **Delete Specific Stock**
  - `DELETE /api/v1/form/:_id/stock/:stock_id`


