# Microservices Project

This project is a microservices-based application built using Docker and Docker Compose. It consists of multiple services: `account`, `catalog`, `order`, `graphql`, and corresponding databases (`account_db`, `catalog_db`, `order_db`). Each service is containerized and can be orchestrated using Docker Compose.

## Project Structure

### Services

- **Account Service**: A service managing account data, connecting to a PostgreSQL database (`account_db`).
- **Catalog Service**: A service managing catalog data, connecting to an Elasticsearch database (`catalog_db`).
- **Order Service**: A service handling orders, connecting to a PostgreSQL database (`order_db`) and interacting with the `Account` and `Catalog` services.
- **GraphQL Gateway**: A GraphQL API service that connects to the `Account`, `Catalog`, and `Order` services.

### Databases

- **account_db**: PostgreSQL database for storing account information.
- **catalog_db**: Elasticsearch database for storing catalog data.
- **order_db**: PostgreSQL database for managing order information.

## Prerequisites

- Docker and Docker Compose must be installed on your system.

## Getting Started

### Setup

Clone the repository and navigate to the project directory:

```bash
git clone <repository_url>
cd <project_directory>