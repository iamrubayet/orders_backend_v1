# Orders Backend v1

This project is a backend service for managing orders. Follow the steps below to set up the project locally.

---

## Steps to Set Up

### 1. Cloning the Repository
```bash
git clone https://github.com/iamrubayet/orders_backend_v1.git
```

### 2. Install Dependencies
```bash
go mod tidy #or 
go get ./...
```

### 3. copy environment variables
```bash
cp .env.example .env
```
### 4. setup database called ordersdb in postgres and make adjustments in .env according to you

### 5. Download golang-migrate tool
```bash
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add 

curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash    
sudo apt-get update   
sudo apt-get install migrate
```

### 6. Migrate Database
```bash
migrate -database "postgres://yourusername:yourpassword@localhost:5432/yourdbname?sslmode=disable" -path migrations up
```

### 7. Run the Project
```bash
make run #or
go run cmd/main.go #or
go run ./...
```
