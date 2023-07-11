package main

import (
	"context"
	"log"
	"strings"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

const PluginName = "jwt-jwk"
const Version = "0.0.1"
const Priority = 1450

type Config struct {
	JwksUrl             string `json:"jwks_url"`
	AuthorizationHeader string `json:"authorization_header"`
	BearerPrefix        string `json:"bearer_prefix"`
	SubHeader           string `json:"sub_header"`
}

func main() {
	server.StartServer(New, Version, Priority)
}

func New() interface{} {
	return &Config{
		AuthorizationHeader: "authorization",
		BearerPrefix:        "Bearer",
		SubHeader:           "x-verified-sub",
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

	jwks, err := jwk.Fetch(context.Background(), jwks_url)
	if err != nil {
		log.Fatalf("Error while fetching JWKS from %s. Error: %s", jwks_url, err.Error())
		return
	}

	token, err := jwt.Parse([]byte(auth_header), jwt.WithKeySet(jwks), jwt.UseDefaultKey(true))
	if err != nil {
		log.Printf("Error parsing 'jwt' from auth header. Error: %s", err.Error())
		return
	}

	kong.Response.SetHeader(conf.SubHeader, token.Subject())
}
