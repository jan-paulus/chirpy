# Chripy

Chripy is a social network similar to Twitter. This repository contains Chripys backend server written in Go.

## Context

Chirpy was built as an exercise to learn how to build webservers with Go.

## System Endpoints

### Health

#### GET /api/healthz

Returns the current status of the system

### Metrics

#### GET /admin/metrics

Returns metrics of the system

#### GET /api/reset

Resets the metrics of the system

### Webhooks

#### POST /api/polka/webhooks

Webhook used to upgrade users to Chirpy Red.

### Auth

#### POST /api/login

User login

Request Body:

```json
{
  "email": "jan@chirpy.com",
  "password": "superstrongpassword"
}
```

Response:

```json
{
  "id": 1,
  "email": "jan@chirpy.com",
  "password": "hash",
  "is_chirpy_red": false,
  "token": "jwt access token",
  "refresh_token": "jwt refresh token"
}
```

#### POST /api/refresh

Issue a new access token using the refresh token

#### POST /api/revoke

Revoke a refresh token

## Resources

### Users

```json
{
  "id": 1,
  "email": "jan@chirpy.com",
  "password": "hash",
  "is_chirpy_red": true
}
```

#### POST /api/users

Create a user

Request Body:

```json
{
  "email": "jan@chirpy.com",
  "password": "superstrongpassword"
}
```

Response:

```json
{
  "id": 1,
  "email": "jan@chirpy.com",
  "password": "hash",
  "is_chirpy_red": false
}
```

#### PUT /api/users

Update a user

Request Body:

```json
{
  "email": "jan-paulus@chirpy.com",
  "password": "adifferentpassword"
}
```

Response:

```json
{
  "id": 1,
  "email": "jan-paulus@chirpy.com",
  "password": "hash",
  "is_chirpy_red": false
}
```

### Chirps

```json
{
  "id": 1,
  "body": "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
  "author_id": 1
}
```

#### POST /api/chirps

Update a user

Request Body:

```json
{
  "body": "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat."
}
```

Response:

```json
{
  "id": 1,
  "body": "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
  "author_id": 1
}
```

#### GET /api/chirps

Response:

```json
[
  {
    "id": 1,
    "body": "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
    "author_id": 1
  },
  {
    "id": 2,
    "body": "Lorem ipsum dolor sit amet, laborsint consectetur cupidatat.",
    "author_id": 1
  }
]
```

#### GET /api/chirps/{chirpId}

Response:

```json
{
  "id": 1,
  "body": "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
  "author_id": 1
}
```

#### DELETE /api/chirps/{chirpId}

Returns HTTP Status 200 on successful deletion.
