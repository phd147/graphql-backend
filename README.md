# graphql-backend

A GraphQL backend service for managing users, products, and orders. Built with Go, gqlgen, and local JSON data persistence.

---

## Local Setup

### Prerequisites
- Go 1.20+
- Docker & Docker Compose (optional, for containerized setup)

### Install Go Dependencies
```bash
go mod download
```

### Run Locally (Development)
```bash
make run
```

#### Run Tests Locally
```bash
make test
```
This will run all integration tests against your locally running server. Make sure the server is running before executing tests.

### Using Docker
```bash
make docker-run
```

#### Run Tests in Docker
```bash
make docker-test
```
This will run all integration tests in a fresh Go container, using the same Docker network as the API server (see Docker instructions above).

### Data Persistence
- On startup, the app loads data from these files. On changes, it writes back to them.

---
### GraphQL Endpoint
- The GraphQL API is available at `http://localhost:8080/query`.
- GraphQL Playground is available at http://localhost:8080 in browser.
## GraphQL API Reference

### Queries

#### 1. Get Products
```graphql
query {
  products(limit: 10, offset: 0, category: "Books") {
    id
    name
    price
    inStock
    description
    category
  }
}
```

#### 2. Get Single Product
```graphql
query {
  product(id: "PRODUCT_ID") {
    id
    name
    price
    inStock
    description
    category
  }
}
```

#### 3. Get Orders (for current user)
```graphql
query {
  orders(limit: 10, offset: 0) {
    id
    total
    status
    createdAt
    products {
      id
      name
    }
    user {
      id
      name
    }
  }
}
```

#### 4. Get Single Order
```graphql
query {
  order(id: "ORDER_ID") {
    id
    total
    status
    createdAt
    products {
      id
      name
    }
    user {
      id
      name
    }
  }
}
```

#### 5. Get Current User
```graphql
query {
  me {
    id
    name
    email
  }
}
```

---

### Mutations

#### 1. Create Product (Admin only)
```graphql
mutation {
  createProduct(input: {
    name: "Sample Product"
    price: 29.99
    inStock: 100
    description: "A sample product"
    category: "Books"
  }) {
    id
    name
    price
    inStock
    description
    category
  }
}
```

#### 2. Update Product (Admin only)
```graphql
mutation {
  updateProduct(input: {
    id: "PRODUCT_ID"
    name: "Updated Name"
    price: 39.99
    inStock: 80
    description: "Updated description"
    category: "Updated Category"
  }) {
    id
    name
    price
    inStock
    description
    category
  }
}
```

#### 3. Place Order (Authenticated user)
```graphql
mutation {
  placeOrder(productIds: ["PRODUCT_ID_1", "PRODUCT_ID_2"]) {
    id
    total
    status
    createdAt
    products {
      id
      name
    }
    user {
      id
      name
    }
  }
}
```

#### 4. Login
```graphql
mutation {
  login(input: { email: "user@example.com", password: "yourpassword" }) {
    accessToken
    refreshToken
    user {
      id
      name
      email
    }
  }
}
```

---

## Default User Credentials

The initial `users.json` file contains two default users for testing:

- **Admin**
  - Email: `admin@example.com`
  - Password: `secret`
- **Customer**
  - Email: `customer@example.com`
  - Password: `secret`

You can use these credentials to log in and test the API with different roles.

---

## Project Structure
- `cmd/` - Application entrypoint
- `app/` - Core business logic
- `entity/` - Data models (User, Product, Order)
- `store/` - Data persistence (repo, JSON files)
- `graph/` - GraphQL schema, resolvers
- `pkg/` - HTTP transport, JWT utilities
- `data-loader/` - DataLoader utilities to batch and cache requests, reducing the N+1 query problem in GraphQL resolvers
- `tests/` - Integration tests 

---
