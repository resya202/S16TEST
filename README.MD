# Go Delegation Tracker

This project fetches Cosmos staking delegations for a set of validators, stores hourly snapshots in PostgreSQL, and exposes a simple HTTP API (via Gin) to query both hourly and daily aggregated data.

---

Repository Structure

```
S16TEST/
├── cmd/                            # Application entrypoints
│   ├── collector/                  # Data collector service
│   │   └── main.go                 # Entrypoint for fetching & storing delegation data
│   └── api/                        # API server service
│       └── main.go                 # Entrypoint for HTTP API server
├── internal/                       # Core application logic
│   ├── config/                     # Configuration loader
│   │   └── config.go               # Loads .env variables using Viper (POSTGRES_URL, COSMOS_API_BASE)
│   ├── db/                         # Database utilities
│   │   └── db.go                   # GORM DB connection helper
│   ├── collector/                  # Delegation data fetch logic
│   │   └── collector.go            # `FetchAndStoreAll`, `fetchAndStore` (specify validator addresses here)
│   └── api/                        # API route logic
│       └── api.go                  # Registers routes and handlers for hourly & daily data
├── internal/model/                 # GORM model definitions
│   └── delegation.go              # `HourlyDelegation` model
├── migrations/                     # Database schema migrations
│   ├── 20250527022316_init_schema.up.sql
│   └── ...                         # Other migration files
├── .env                            # Environment variables file (not committed to version control)
├── go.mod, go.sum                  # Go module dependencies
└── README.md                       # Project documentation (this file)
```


---

Environment Setup

Create a `.env` file at the project root with the following keys:


POSTGRES_URL=postgres://<db_user>:<db_pass>@localhost:5432/delegations?sslmode=disable
COSMOS_API_BASE=https://cosmos-api.polkachu.com


POSTGRES_URL**: full connection string to your `delegations` database.
COSMOS_API_BASE**: base URL of the Cosmos REST API provider.

Load them into your shell:

bash
export $(grep -v '^#' .env | xargs)


---

Database Migrations

Ensure you have `migrate` CLI installed. Then:

bash
# Start Postgres (e.g. via Homebrew)
brew services start postgresql

# Create database
psql -U $(whoami) -d postgres -c "CREATE DATABASE delegations;"

# Apply migrations
migrate -path migrations -database "$POSTGRES_URL" up


This creates one table:

sql
-- delegation_hourly:
--  id SERIAL PK, validator_addr, delegator_addr, timestamp (with tz), amount_uatom, change_uatom

---

 Running the Collector

The collector fetches each validator’s delegations hourly, stores a snapshot in `delegation_hourly`, and can be extended to compute daily totals.

bash
go run cmd/collector/main.go


Check the `delegation_hourly` table to confirm rows.

---

 Starting the API Server

In another terminal:

bash
go run cmd/api/main.go


By default, Gin listens on `:8080`.

Endpoints

Hourly Delegations

  bash
  curl http://localhost:8080/api/v1/validators/{validator_address}/delegations/hourly
  

Daily Aggregation**

  bash
  curl http://localhost:8080/api/v1/validators/{validator_address}/delegations/daily
  

  * Groups `delegation_hourly` by date and sums `change_uatom`.

---

 Testing with Postman

1. Create a **new collection** (e.g. *Go Delegation Tracker*).
2. Add two GET requests:

   * **Get Hourly**

     * URL: `http://localhost:8080/api/v1/validators/{{validator_address}}/delegations/hourly`
   * **Get Daily**

     * URL: `http://localhost:8080/api/v1/validators/{{validator_address}}/delegations/daily`
3. Define an **Environment** variable: `validator_address` = your test val-oper.
4. Send and inspect JSON responses.

---

 What To Expect

* **Hourly**: an array of `HourlyDelegation` objects with raw snapshot rows.
* **Daily**: an array of `{ validator_addr, delegator_addr, date, total_change_uatom }` showing the net change per day per delegator.

If `total_change_uatom` remains `0` for today, it means no change since snapshots span a single date—so you may need to run multiple hourly snapshots across days to see a non-zero sum.

---

###CURL

```
{
  "info": {
    "name": "Go Delegation Tracker",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Get Hourly Delegations",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/api/v1/validators/{{validator_address}}/delegations/hourly",
          "host": ["http://localhost:8080"],
          "path": ["api","v1","validators","{{validator_address}}","delegations","hourly"]
        }
      }
    },
    {
      "name": "Get Daily Delegations",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/api/v1/validators/{{validator_address}}/delegations/daily",
          "host": ["http://localhost:8080"],
          "path": ["api","v1","validators","{{validator_address}}","delegations","daily"]
        }
      }
    },
    {
      "name": "Get Delegator History",
      "request": {
        "method": "GET",
        "header": [],
        "url": {
          "raw": "http://localhost:8080/api/v1/validators/{{validator_address}}/delegator/{{delegator_address}}/history",
          "host": ["http://localhost:8080"],
          "path": ["api","v1","validators","{{validator_address}}","delegator","{{delegator_address}}","history"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "validator_address",
      "value": "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn"
    },
    {
      "key": "delegator_address",
      "value": "cosmos1qqy2q4mtmncqh6yw3dy5qp5e37vytadckzpw2j"
    }
  ]
}

```

