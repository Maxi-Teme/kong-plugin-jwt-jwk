package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const PluginName = "jwt-jwk"
const Version = "0.0.1"
const Priority = 1450

type Config struct {
	JwksUrl             string `json:"jwks_url"`
	AuthorizationHeader string `json:"authorization_header"`
	BearerPrefix        string `json:"bearer_prefix"`
	SubHeader           string `json:"sub_header"`
	_jwk_cache          *jwk.Cache
	_ctx                context.Context
}

func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{
		AuthorizationHeader: "authorization",
		BearerPrefix:        "Bearer",
		SubHeader:           "x-verified-sub",
		_jwk_cache:          nil,
		_ctx:                context.Background(),
	}
}

func (conf *Config) Access(kong *pdk.PDK) {
	prefix := conf.BearerPrefix + " "
	auth_header, err := kong.Request.GetHeader(conf.AuthorizationHeader)
	if err != nil || !strings.HasPrefix(auth_header, prefix) {
		return
	}

	auth_header = strings.TrimPrefix(auth_header, prefix)

	jwks_url := conf.JwksUrl
	if jwks_url == "" {
		log.Fatalf("Error 'jwks_url' is empty")
		return
	}

	if conf._jwk_cache == nil {
		conf._jwk_cache = jwk.NewCache(conf._ctx)
		conf._jwk_cache.Register(jwks_url, jwk.WithMinRefreshInterval(5*time.Minute))
	}

	jwks, err := conf._jwk_cache.Get(conf._ctx, jwks_url)
	if err != nil {
		log.Fatalf("Error while fetching JWKS from %s. Error: %s", jwks_url, err.Error())
		return
	}

	// TODO: check if jwt.WithValidate(true) already validates 'exp' or if we use validator
	// validator := jwt.ValidatorFunc(func(_ context.Context, t jwt.Token) jwt.ValidationError {
	// 	return nil
	// })

	token, err := jwt.Parse(
		[]byte(auth_header),
		jwt.WithKeySet(jwks, jws.WithUseDefault(true)),
		jwt.WithValidate(true),
		// jwt.WithValidator(validator),
	)

	if err != nil {
		log.Printf("Error parsing 'jwt' from auth header. Error: %s", err.Error())
		return
	}

	kong.Response.SetHeader(conf.SubHeader, token.Subject())
}
