package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/arctir/flightdeck-cli/auth"
	"github.com/arctir/flightdeck-cli/commands/common"
	arctirclient "github.com/arctir/go-flightdeck/pkg/client"
	"github.com/golang-jwt/jwt"
	"golang.org/x/oauth2"
)

type AuthLoginCommand struct {
	AuthEndpoint string `name:"auth-endpoint" default:"${authEndpoint}"`
}

type AuthTokenCommand struct{}
type AuthWhoamiCommand struct{}

type AuthCommand struct {
	Login  AuthLoginCommand  `cmd:"login"`
	Token  AuthTokenCommand  `cmd:"token"`
	Whoami AuthWhoamiCommand `cmd:"whoami"`
}

func (c *AuthLoginCommand) BeforeApply(ctx *Context, globals *common.Globals) error {
	ctx.SkipConfigCheck = true
	return nil
}

func (a AuthLoginCommand) Run(globals *common.Globals) error {
	return auth.Login(a.AuthEndpoint, globals.ConfigPath)
}

func (a AuthTokenCommand) Run(ctx *Context) error {
	conf, err := arctirclient.NewOauthConfig(ctx.Config.AuthEndpoint)
	if err != nil {
		return err
	}

	exp, err := arctirclient.ExtractExpiry(ctx.Config.AccessToken, 10)
	if err != nil {
		return err
	}

	token := oauth2.Token{
		AccessToken:  ctx.Config.AccessToken,
		RefreshToken: ctx.Config.RefreshToken,
		Expiry:       *exp,
	}
	ts := conf.TokenSource(context.TODO(), &token)
	refreshedToken, err := ts.Token()
	if err != nil {
		return err
	}
	fmt.Print(refreshedToken.AccessToken)
	return nil
}

func (c *AuthWhoamiCommand) Run(ctx *Context, kc *kong.Context) error {
	claims := &auth.ArctirClaims{}
	_, _, err := new(jwt.Parser).ParseUnverified(ctx.Config.AccessToken, claims)
	if err != nil {
		return err
	}
	fmt.Printf("Issuer: %s\n", claims.Issuer)
	fmt.Printf("Username: %s\n", claims.Email)
	fmt.Printf("Groups: %s\n\n", strings.Join(claims.Groups, ", "))

	return nil
}
