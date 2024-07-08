package main

import (
	"fmt"
	"encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	servicePrinicipalType = "AzureADServicePrinicipal"
)

func (ad *AzureADPlugin) GetServicePrincipal(ctx context.Context, data interface{}) error {
			var servicePrincipal *models.ServicePrincipal
			spJson, err := json.Marshal(data) 
			if err != nil {
				return err
			}

			err = json.Unmarshal(spJson, &servicePrincipal)
			if err != nil {
				return err
			}

			props, err := plugin.StructToMap(servicePrincipal)
			if err != nil {
				fmt.Errorf("failed to marshal: %v", err)
			}

			labels := []string{servicePrinicipalType}
			graph := ad.Plugin.AddOrFindGraph(servicePrinicipalType, plugin.NewSchema(nil))

			fmt.Println("@@@@@@@@@@@@@@@@@")
			fmt.Printf("adding servicePrincipal  node %s\n", servicePrincipal.Id)
			newNode, err := graph.NewNode(plugin.Role, servicePrinicipalType, servicePrincipal.Id, servicePrincipal.Id, labels, props)
			if err != nil {
				fmt.Errorf("unable to create servicePrincipal node: %v", err)
			}

			fmt.Println("@@@@@@@@@@@@@@@@@")
			fmt.Printf("adding servicePrincipal relation  node %s\n", servicePrincipal.AppId)
			// relation to the user 
			_, err = newNode.NewRelation(servicePrincipal.AppId, plugin.BELONGS_TO, nil)
			if err != nil {
				return fmt.Errorf("unable to create user to servicePrincipal to app relation: %v", err)
			}
	
	return nil
}