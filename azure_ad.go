package main

import (
	"fmt"
	"context"
	// "encoding/json"
	// "os"
	// "bufio"

	// cmd "github.com/defensestation/azurehound/cmd"
	client_config "github.com/defensestation/azurehound/client/config"
	client "github.com/defensestation/azurehound/client"
)

func (ad *AzureADPlugin) Run(ctx context.Context) error {
    options := &Options{}

    err := ad.Plugin.MarshalOptions(options)
    if err != nil {
        fmt.Printf("unable to MarshalOptions")
        return err
    }

    ad.Options = options

    err = ad.setAzureADClient(ctx)
   	if err != nil {
   		return fmt.Errorf("Unable setup azuread client: %v", err)
   	}

    err = ad.List(ctx)
    if err != nil {
   		return fmt.Errorf("Failed to list: %v", err)
    }
   
   return nil
}

func (ad *AzureADPlugin) setAzureADClient(ctx context.Context) error {
	// ctx := context.Background()
	config := client_config.Config{
		ApplicationId:  ad.Options.ClientID,
		Authority:      "",
		ClientSecret:   ad.Options.ClientSecret,
		ClientCert:     "",
		ClientKey:      "",
		ClientKeyPass:  "",
		Graph:          "",
		JWT:            "",
		Management:     "",
		MgmtGroupId:    []string{""},
		Password:       "",
		ProxyUrl:       "",
		RefreshToken:   "",
		Region:         "cloud",
		SubscriptionId: []string{ad.Options.SubscriptionID},
		Tenant:         ad.Options.TenantID,
		Username:       "",
	}

	azClient, err := client.NewClient(config)
	if err != nil {
		return err
	}

	ad.Client = &azClient
	return nil
}