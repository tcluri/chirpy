

# API Endpoint Request Examples


## Creating a User

To create a new user in the Chirpy webserver, you can send a POST request to the `/api/users` endpoint. This endpoint expects a JSON payload containing the necessary user information.


### Request and Response

-   Method: POST
-   URL: `/api/users`
-   Headers:
    -   Content-Type: application/json
-   Request Body(JSON object):
    -   `email`: The email address of the user.
    -   `password`: The password for the user.

Request Body:

    {
      "username": "john_doe",
      "password": "secretpassword"
    }

Response Body:

    {
    "id": 1,
    "email": "example@example.com",
    "is_chirpy_red": false
    }


## Updating a User

To update an existing user in the Chirpy webserver, you can send a PUT request to the `/api/users` endpoint. The request should include the necessary parameters in the request body in JSON format, along with a valid JWT (JSON Web Token) in the authorization header.


### Request and Response

-   Method: PUT
-   Endpoint: `/api/users`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {JWT}
-   Request Body:
    -   `email`: The new email address for the user.
    -   `password`: The new password for the user.

Request Body:

    {
    "email": "newemail@example.com",
    "password": "newsecretpassword"
    }

Response Body:

    {
    "id": 1,
    "email": "newemail@example.com",
    "is_chirpy_red": false
    }


## User Login

To authenticate and log in as a user in the Chirpy webserver, you can send a POST request to the `/api/login` endpoint.


### Request

-   Method: POST
-   Endpoint: `/api/login`

The request body should contain the following parameters:

-   `email` (string): The user&rsquo;s email address.
-   `password` (string): The user&rsquo;s password.

Request Body:

    {
    "email": "user@example.com",
    "password": "password123"
    }


### Response

If the login is successful, the API will respond with a status code of 200 (OK) and a JSON object containing the user information and access tokens:

-   `User` (object): An object containing user details, such as the user&rsquo;s ID, email, and membership status.
-   `Token` (string): An access token (JWT) used for authentication in subsequent requests.
-   `RefreshToken` (string): A refresh token (JWT) used to obtain new access tokens when the current access token expires.

Response Body:

    {
    "User": {
    "ID": 123,
    "Email": "user@example.com",
    "IsChirpyRed": true
    },
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.alskdjqweuhl.s0meR4nd0mT0k3n",
    "RefreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.alskdjqweuhl.s0meR4nd0mR3fr35hT0k3n"
    }


## Creating a Chirp

To create a new chirp in the Chirpy webserver, you can send a POST request to the `/api/chirps` endpoint. The request should include the necessary parameters in the request body in JSON format, along with a valid JWT (JSON Web Token) in the authorization header.


### Request and Response

-   Method: POST
-   Endpoint: `/api/chirps`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {JWT}
-   Request Body:
    -   `body`: The content of the chirp.

Request Body:

    {
    "body": "Hello, world! This is my first chirp."
    }

Response Body:

    {
    "id": 1,
    "author_id": 123,
    "body": "Hello, world! This is my first chirp."
    }


## Retrieving Chirps

To retrieve chirps from the Chirpy webserver, you can send a GET request to the `/api/chirps` endpoint. The request can include optional query parameters to specify the sorting order and filter by author ID.


### Request

-   Method: GET
-   Endpoint: `/api/chirps?sort=desc&author_id=123`
-   Query Parameters:
    -   `sort` (optional): Specify the sorting order of chirps. Use `sort=desc` to retrieve chirps in descending order based on their ID. The default sorting order is ascending.
    -   `author_id` (optional): Filter chirps by the author&rsquo;s ID. Only chirps created by the specified author will be retrieved.

Response Body:

    [
    {
    "id": 2,
    "author_id": 123,
    "body": "I'm enjoying the Chirpy webserver."
    },
    {
    "id": 1,
    "author_id": 123,
    "body": "Hello, world! This is my first chirp."
    }
    ]


## Get Chirp by ID

To retrieve a specific chirp by its ID from the Chirpy webserver, you can send a GET request to the `/api/chirps/{chirpID}` endpoint, where `{chirpID}` is the ID of the chirp you want to retrieve.


### Request

-   Method: GET
-   Endpoint: `/api/chirps/{chirpID}` (replace `{chirpID}` with the actual ID of the chirp)

No Request Body needed

Response Body:

    {
    "id": 123,
    "author_id": 456,
    "body": "This is a chirp about something interesting."
    }


## Deleting Chirps by ID

To delete a Chirp by its ID from the Chirpy webserver, you can send a DELETE request to the `/api/chirps/{chirpID}` endpoint, where `{chirpID}` should be replaced with the actual ID of the Chirp to be deleted.


### Request

-   Method: DELETE
-   Endpoint: `/api/chirps/{chirpID}`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {JWT}

Make sure to include a valid JWT (JSON Web Token) in the request header to authenticate and authorize the deletion operation.

No Request Body and no Response Body for the delete endpoint.


## Refresh Access Token

To refresh the access token for a user in the Chirpy webserver, you can send a POST request to the `/api/refresh` endpoint.


### Request

-   Method: POST
-   Endpoint: `/api/refresh`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {JWT}

The request should include the refresh token in the `Authorization` header as a Bearer token.


### Response

If the refresh token is valid and not revoked, the API will respond with a status code of 200 (OK) and a JSON object containing the new access token:

-   `Token` (string): A new access token (JWT) used for authentication in subsequent requests.

Response Body:

    {
    "Token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.asdqdfqwh.n3w4cc3sst0k3n"
    }


## Revoke Refresh Token

To revoke a refresh token in the Chirpy webserver, you can send a POST request to the `/api/revoke` endpoint.


### Request

-   Method: POST
-   Endpoint: `/api/revoke`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {JWT}

The request should include the refresh token in the `Authorization` header as a Bearer token.


### Response

If the refresh token is successfully revoked, the API will respond with a status code of 200 (OK) and an empty JSON object.


## User Upgrade Webhook

The user upgrade webhook endpoint in the Chirpy webserver allows the Polka service to send upgrade events and upgrade the user status in the Chirpy system. To utilize this endpoint, you can send a POST request to the `/api/polka/webhooks` endpoint.


### Request

-   Method: POST
-   Endpoint: `/api/polka/webhooks`
-   Headers:
    -   Content-Type: application/json
    -   Authorization: Bearer {POLKA<sub>KEY</sub>}

The request should include the Polka API key in the `Authorization` header.
The `event` field should be set to `"user.upgraded"`, and the `data` field should contain the `user_id` of the user to be upgraded.

Request Body:

    {
    "event": "user.upgraded",
    "data": {
    "user_id": 123
    }
    }


### Response

If the user upgrade event is successfully processed, the API will respond with a status code of 200 (OK) and an empty JSON object.

