

# Chirpy Webserver


## Overview

Chirpy is a webserver that provides an API for a social media platform. It allows users to create chirps (short messages), retrieve chirps, delete chirps, and perform various user-related operations. The webserver is built using the Go programming language and the Chi router.


## Getting Started

To run the Chirpy webserver, follow these steps:

1.  Clone the repository: `git clone https://github.com/tcluri/chirpy.git`
2.  Install the dependencies: `go mod download`
3.  Create a `.env` file in the root directory and set the required environment variables (e.g., `JWT_SECRET` and `POLKA_KEY`).
4.  Build the application: `go build`
5.  Run the webserver: `./chirpy`


## API Endpoints

The Chirpy webserver provides the following API endpoints:

-   `GET /api/healthz`: Health check endpoint to verify the server&rsquo;s availability.

-   `POST /api/chirps`: Create a new chirp.
-   `GET /api/chirps`: Retrieve all chirps.
-   `GET /api/chirps/{chirpID}`: Retrieve a specific chirp by ID.
-   `DELETE /api/chirps/{chirpID}`: Delete a specific chirp by ID.

-   `PUT /api/users`: Update a user&rsquo;s information.
-   `POST /api/users`: Create a new user.
-   `POST /api/login`: User login.

-   `POST /api/refresh`: Refresh an authentication token.
-   `POST /api/revoke`: Revoke an authentication token.

-   `POST /api/polka/webhooks`: Handle Polka webhooks for user upgrades().

-   `GET /metrics`: Retrieve server metrics.

For examples, see [request examples](EXAMPLES.md).


## Configuration

The Chirpy webserver supports the following configuration options:

-   `JWT_SECRET`: Secret key for JWT token generation and validation.
-   `POLKA_KEY`: Secret key for handling Polka webhooks.


## Development Mode

To run the webserver in debug mode, use the `-debug` flag: `./chirpy -debug`. This will reset the database and enable debug logging.


## Dependencies

The Chirpy webserver uses the following external dependencies:

-   `github.com/go-chi/chi/v5`: Lightweight and expressive HTTP router for Go.
-   `github.com/joho/godotenv`: Go library for loading environment variables from a `.env` file.
-   `github.com/tcluri/chirpy/internal/database`: Internal package for managing the database.

