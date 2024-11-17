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

### 5. Migrate Database
```bash
migrate -database "postgres://yourusername:yourpassword@localhost:5432/yourdbname?sslmode=disable" -path migrations up
```

### 6. Run the Project
```bash
make run #or
go run cmd/main.go #or
go run ./...
```
