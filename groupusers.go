package main

import (
	"fmt"
	// "strings"
	// "encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
    gjson "github.com/tidwall/gjson"
)

const (
	groupMemberType = "AzureADGroupMember"
)

func (ad *AzureADPlugin) GetGroupUsers(ctx context.Context, data interface{}) error {
			groupMembers := data.(models.GroupMembers)
			
			for _, groupMember := range(groupMembers.Members) {
				
				graph := ad.Plugin.AddOrFindGraph(userType, plugin.NewSchema(nil))

				
				// 
				startNodeId := gjson.GetBytes(groupMember.Member, "mail").String()
				fmt.Printf("adding group member relation to node %s\n", startNodeId)
				if startNodeId == "" {
					startNodeId = gjson.GetBytes(groupMember.Member, "id").String()
				}
				// relation to the user 
				_, err := graph.DirectRelation(startNodeId, groupMembers.GroupId, plugin.BELONGS_TO, nil)
				if err != nil {
					return fmt.Errorf("unable to create user to groupmember relation: %v", err)
				}
			}

	return nil
}