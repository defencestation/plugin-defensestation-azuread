package main

import (
	"fmt"
	"encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	groupType = "AzureADGroup"
)

func (ad *AzureADPlugin) GetGroups(ctx context.Context, data interface{}) error {
			var group *models.Group
			groupJson, err := json.Marshal(data) 
			if err != nil {
				return err
			}

			err = json.Unmarshal(groupJson, &group)
			if err != nil {
				return err
			}

			props, err := plugin.StructToMap(group)
			if err != nil {
				fmt.Errorf("failed to marshal: %v", err)
			}

			labels := []string{groupType}
			graph := ad.Plugin.AddOrFindGraph(groupType, plugin.NewSchema(nil))

			fmt.Println("=================")
			fmt.Printf("adding group  node %s\n", group.Id)
			_, err = graph.NewNode(plugin.Group, groupType, group.Id, group.DisplayName, labels, props)
			if err != nil {
				fmt.Errorf("unable to create subscription node: %v", err)
			}
	
return nil
}