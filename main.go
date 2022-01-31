package main

// This file is here just to make `go get` easier

import (
	_ "github.com/golang-jwt/jwt"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/google/uuid"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	_ "github.com/rs/zerolog"
	_ "github.com/urfave/cli/v2"
)

func main() {

}
