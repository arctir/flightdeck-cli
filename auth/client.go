package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/arctir/flightdeck-cli/client"
	"github.com/int128/oauth2cli"
	"github.com/int128/oauth2cli/oauth2params"
	"github.com/pkg/browser"
	"golang.org/x/sync/errgroup"

	arctirclient "github.com/arctir/go-flightdeck/pkg/client"
)

func Login(authEndpoint string, configPath string) error {
	pkce, err := oauth2params.NewPKCE()
	if err != nil {
		return err
	}
	ready := make(chan string, 1)
	defer close(ready)

	oauth2Config, err := arctirclient.NewOauthConfig(authEndpoint)
	if err != nil {
		return err
	}

	cfg := oauth2cli.Config{
		OAuth2Config:           *oauth2Config,
		AuthCodeOptions:        pkce.AuthCodeOptions(),
		TokenRequestOptions:    pkce.TokenRequestOptions(),
		LocalServerReadyChan:   ready,
		LocalServerBindAddress: []string{"localhost:30000"},
		Logf:                   log.Printf,
	}

	ctx := context.Background()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			log.Printf("Open %s", url)
			if err := browser.OpenURL(url); err != nil {
				log.Printf("could not open the browser: %s", err)
			}
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		token, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("could not get a token: %w", err)
		}
		log.Printf("You got a valid token until %s", token.Expiry)
		conf := client.Config{
			AuthEndpoint: authEndpoint,
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		}

		err = conf.Save(configPath)
		if err != nil {
			return err
		}

		return nil
	})
	if err := eg.Wait(); err != nil {
		log.Fatalf("authorization error: %s", err)
	}
	return nil
}
