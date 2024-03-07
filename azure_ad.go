package main

import (
	"fmt"
	"context"

	client_config "github.com/defensestation/azurehound/client/config"
	client "github.com/defensestation/azurehound/client"

	plugin "github.com/defensestation/pluginutils"
)

type AzureAd struct {
	plugin *plugin.Plugin
	client *client.AzureClient
}

func NewAzureADPlugin(plugin *plugin.Plugin) *AzureAd {
	return &AzureAd{
		plugin: plugin,
	}
}

func (ad *AzureAd) Run(ctx context.Context) (error) {
	err := ad.plugin.ValidateOptions("application_id", "client_secret", "tenant", "group_name")
	if err != nil {
		return err
	}

	applicationId, _ := ad.plugin.GetOption("application_id");
    clientSecret, _ := ad.plugin.GetOption("client_secret");
    Tenant, _ := ad.plugin.GetOption("tenant");
    groupName, _ := ad.plugin.GetOption("group_name");

   err = ad.setAzureADClient(ctx,  applicationId.(string), clientSecret.(string), Tenant.(string))
   if err != nil {
   		return fmt.Errorf("Unable setup azuread client: %v", err)
   }

   // get all userons
   err = ad.GetUsers(ctx)
   if err != nil {
   	fmt.Println(err)
   	return err
   }

   if groupName.(string) != "" {
   	err = ad.GetGroupUsers(ctx, groupName.(string))
	   if err != nil {
	   	fmt.Println(err)
	   	return err
	   }
   }
   
   return nil
}

func (ad *AzureAd) setAzureADClient(ctx context.Context, applicationId, clientSecret, Tenant string) error {
	// ctx := context.Background()
	config := client_config.Config{
		ApplicationId:  applicationId,
		Authority:      "",
		ClientSecret:   clientSecret,
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
		SubscriptionId: []string{""},
		Tenant:         Tenant,
		Username:       "",
	}

	azClient, err := client.NewClient(config)
	if err != nil {
		return err
	}

	ad.client = &azClient
	return nil
}
