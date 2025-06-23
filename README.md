# Ecommerce Price Tracker – Backend

A modular, efficient backend in Go for tracking ecommerce product prices over time. Exposes RESTful APIs for secure user access, product tracking, historical pricing, and email alerts on price drops — all backed by scheduled scrapers and a PostgreSQL database.

## 📦 Features

- 🔐 User registration & login (`/api/register`, `/api/login`)
- 📦 Product tracking (`/api/product`)
- 🕒 Automated price checks every `N` hours
- 📨 Email alerts when prices drop
- 📈 View historical price trends (`/api/product/:id`)
- 📂 List tracked products (`/api/products`)

## 🚀 Getting Started

```bash
git clone https://github.com/Rudra-241/ecommerce-price-tracker.git
cd ecommerce-price-tracker
cp .env.example .env
# Fill in DB, SMTP, and config values

# Build and run the project using Make
make run
