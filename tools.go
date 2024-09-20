//go:build tools

package Involvio

import (
	_ "github.com/go-openapi/errors"
	_ "github.com/go-openapi/runtime"
	_ "github.com/go-openapi/spec"
	_ "github.com/go-openapi/strfmt"
	_ "github.com/go-openapi/swag"
	_ "github.com/go-openapi/validate"
	_ "github.com/go-swagger/go-swagger"
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jessevdk/go-flags"
	_ "github.com/pkg/errors"
	_ "github.com/stretchr/testify"
)
