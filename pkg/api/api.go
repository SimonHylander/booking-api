// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"crypto/sha1"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/simonhylander/booker/pkg/api/auth"
	al "github.com/simonhylander/booker/pkg/api/auth/logging"
	at "github.com/simonhylander/booker/pkg/api/auth/transport"
	"github.com/simonhylander/booker/pkg/api/user"
	ul "github.com/simonhylander/booker/pkg/api/user/logging"
	ut "github.com/simonhylander/booker/pkg/api/user/transport"
	"github.com/simonhylander/booker/pkg/utl/config"
	"github.com/simonhylander/booker/pkg/utl/db"
	"github.com/simonhylander/booker/pkg/utl/jwt"
	authMw "github.com/simonhylander/booker/pkg/utl/middleware/auth"
	"github.com/simonhylander/booker/pkg/utl/rbac"
	"github.com/simonhylander/booker/pkg/utl/secure"
	"github.com/simonhylander/booker/pkg/utl/server"
	"github.com/simonhylander/booker/pkg/utl/zlog"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	dbSession, err := db.New()

	if err != nil {
		return err
	}

	securer := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.Service{}

	secret := "WBm2RTOfkjYNlOwdHuBedlvTsF2hyVMcFXnfD1Jbp8KPDNtQBJA5R1cnQdDE8i1p" //os.Getenv("JWT_SECRET")

	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, secret, cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	authMiddleware := authMw.Middleware(jwt)

	at.NewHTTP(al.New(auth.Initialize(jwt, securer, rbac), log), e, authMiddleware)

	v1 := e.Group("/v1")
	v1.Use(authMiddleware)

	ut.NewHTTP(ul.New(user.Initialize(rbac, securer), log), v1)
	//pt.NewHTTP(pl.New(password.Initialize(rbac, sec), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}