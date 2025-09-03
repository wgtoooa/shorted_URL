# URL Shortener

Welcome to my first Golang project - a URL shortening service!

<img src="/static/URLSHORT.png" alt="Логотип" style="width: 400px; height: auto;">

##  Tech Stack

*   **Backend:** Golang
*   **Frontend:** HTML, CSS, JavaScript
*   **Database:** PostgreSQL
*   **Cache:** Redis

##  Features
* **User Profiles:** Register and create your personal account
* **URL Shortening:** Convert long URLs into short, manageable links
* **Custom Short URLs:** Choose your own custom names for shortened links
* **Link Management:** Delete unwanted links, edit short URLs and copy full URLs
* **Quick Access:** Your profile stores up to 5 frequently used links for instant access



## 📦 Installation

## Prerequisites
* Docker and Docker Compose
* Git

## Setup Instructions

1.  Clone the repository:
    ```bash
    git clone https://github.com/wgtoooa/shorted_URL.git
    cd shorted_URL
    ```
2.  Configure environment variables:

#### Create a .env file in the project root with the following variables:

```env
# Database
POSTGRES_USER=user
POSTGRES_PASSWORD=12345
POSTGRES_HOST=postgres
POSTGRES_NAME=myDB
POSTGRES_PORT=5432

# Redis
REDIS_ADDRESS=redis:6379
REDIS_PASSWORD=12345
REDIS_DB=0

# Server
SERVER_PORT=80
SERVER_HOST=0.0.0.0

# API & Security
API_KEY=abc123xyz
JWT_SECRET=your-super-secret-jwt-key

# Environment
PRODUCTION=true 
```
3. Build and run with Docker:
```bash
# Build the Docker images
docker compose build

# Start the services in detached mode
docker compose up -d
```
4. Verify the installation:
```bash
curl http://localhost:80/health
# Expected response: status OK
```

## 🎯 Usage
1. Register a new account
2. Log in to your personal dashboard
3. Shorten URLs using the provided form
4. Manage your shortened links from your profile
5. Use your custom short URLs to redirect to original websites

## 🛠️ Development
To contribute to this project:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## ❓ Support
If you encounter any issues or have questions, please contact me on Telegram: **@wgtoooa**

## 📄 License
This project is open source and available under the MIT License.
