# Go-Auth

**This project is only for learning purposes.**
A simple microservice for storing user credentials and authenticating them.

The REST api is documented and generated with OpenApi spec.
The tool for generation is deepmap/oapi-codegen.

## Setup

### Prerequisites

- Golang 1.18 or higher

- running MySQL Database instance
	- use files in /SQL to setup the db
	- dsn for the database is setup in app.go at 24th line

### Actions to run

1. `git clone https://github.com/Nesquiko/go-auth.git`
2. `cd go-auth`
3. `go run .`

## Interact

To interact with running Go-Auth service, either go through the 
OpenApi spec and use a client of your choice. Or use Postman,
you can import the collection and environment in /postman dir.

### Interaction flow

1. signup new user
2. log in with credentials, and get unauthenticated token
3. setup 2FA
4. verify 2FA TOTP password, and get fully authenticated token
5. use test endpoint test endpoint if you are correctly authenticated

