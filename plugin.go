package main

import (
    // "context"
    // "fmt"
    // "log"

    "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
    // "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
    plugin "github.com/defensestation/pluginutils"
    // types "github.com/defencestation/plugin-defensestation-azuread/types"
)


func GetCredentials(clientID, clientSecret, tenantID string) (*azidentity.ClientSecretCredential, error) {
    cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
    if err != nil {
        return nil, err
    }
    return cred, nil
}


func NewAzurePlugin(plugin *plugin.Plugin) *AzureADPlugin {
    return &AzureADPlugin{
        Plugin: plugin,
    }
}





