package main

import (
    "sync"
    plugin "github.com/defensestation/pluginutils"
    // "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    client "github.com/defensestation/azurehound/client"
)

type AzureADPlugin struct {
    Options  *Options
    Client *client.AzureClient
    Plugin   *plugin.Plugin
    wg       sync.WaitGroup
}

type Options struct {
    SubscriptionID  string `json:"subscription_id"`
    TenantID  string `json:"tenant_id"`
    ClientID   string `json:"client_id"`
    ClientSecret   string `json:"client_secret"`
}