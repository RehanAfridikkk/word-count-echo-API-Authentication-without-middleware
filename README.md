# API Endpoints

This project provides a set of multiple API endpoints to perform various operations. The API is secured using JSON Web Tokens (JWT) for authentication.

## Public Routes

### 1. Signup
- **Method:** POST
- **URL:** `/signup`
- **Description:** Endpoint to register a new user.

### 2. Login
- **Method:** POST
- **URL:** `/login`
- **Description:** Endpoint to authenticate and generate a JWT token for a user.

### 3. Upload
- **Method:** POST
- **URL:** `/upload`
- **Description:** Endpoint to upload files.

## Protected Routes (Require JWT Authentication)

### 1. My Processes
- **Method:** POST
- **URL:** `/my/processes`
- **Description:** Endpoint to retrieve processes for the authenticated user.

### 2. My Statistics
- **Method:** POST
- **URL:** `/my/statistics`
- **Description:** Endpoint to retrieve statistics for the authenticated user.

## Admin Routes (Require Admin Authentication)

### 1. Process by Username (Admin)
- **Method:** POST
- **URL:** `/Admin/process_by_username`
- **Description:** Endpoint to retrieve processes for a specific user (admin-only).

### 2. Admin Statistics
- **Method:** GET
- **URL:** `/Admin/statistics`
- **Description:** Endpoint to retrieve overall statistics (admin-only).

## Token Refresh Endpoint

### 1. Refresh Token
- **Method:** POST
- **URL:** `/refreshtoken`
- **Description:** Endpoint to refresh an expired JWT token.

### 2. Admin Login
- **Method:** POST
- **URL:** `/admin/login`
- **Description:** Endpoint for admin authentication.

## Database

This project uses PostgreSQL as its database. Ensure you have a PostgreSQL database running and configured appropriately. The connection details can be specified in the configuration files.

Make sure to set up the necessary tables for user information and any other required entities. You can find SQL scripts for table creation in the `sql` folder.

## Dependencies

- [PostgreSQL](https://www.postgresql.org/)
- [JWT Middleware](https://github.com/auth0/go-jwt-middleware) for JWT authentication.

Feel free to explore the Swagger UI for detailed API documentation: [Swagger UI](http://localhost:1303/swagger-ui/)
