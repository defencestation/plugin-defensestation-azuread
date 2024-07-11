package main

import (
	"fmt"
	// "strings"
	"encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	groupMemberType = "AzureADGroupMember"
)

func (ad *AzureADPlugin) GetGroupUsers(ctx context.Context, data interface{}) error {
			groupMembers := data.(models.GroupMembers)
			
			for _, groupMember := range(groupMembers.Members) {
				var user *models.User
				// userJson, err := json.Marshal(groupMember.Member) 
				// if err != nil {
				// 	return err
				// }

				err := json.Unmarshal(groupMember.Member, &user)
				if err != nil {
					return err
				}
				// groupMemberMapInterface, err := plugin.StructToMap(user)
				// if err != nil {
				// 	fmt.Errorf("failed to marshal: %v", err)
				// }

				// if user.DisplayName != "" {
				// 	splitName := strings.Split(user.DisplayName, " ")
				// 	if len(splitName) >= 2 {
				// 		groupMemberMapInterface["first_name"]  = splitName[0]
				// 		groupMemberMapInterface["last_name"] = splitName[1]
				// 	}
				// 	if len(splitName) == 1 {
				// 		groupMemberMapInterface["first_name"]  = splitName[0]
				// 	}
				// }

				// groupMemberMapInterface["type"]    	 = "employee"
				// // groupMemberMapInterface["service"] = "dsc_service_policy_manager"
				// groupMemberMapInterface["personnel"]    = "personnel"
				// // groupMemberMapInterface["personnel_id"] = fmt.Sprintf("%s_%s", ad.plugin.Name, user.Mail)

				// labels := []string{groupMemberType}
				graph := ad.Plugin.AddOrFindGraph(userType, plugin.NewSchema(nil))

				// fmt.Println("000000000")
				// fmt.Printf("adding group member node %s\n", user.DisplayName)
				// newNode, err := graph.NewNode(plugin.Role, groupMemberType, user.Id, user.DisplayName, labels, groupMemberMapInterface)
				// if err != nil {
				// 	fmt.Errorf("unable to create groupmember node: %v", err)
				// }

				// fmt.Println("000000000")
				fmt.Printf("adding group member relation to node %s\n", user.Mail)
				startNodeId := user.Mail
				if user.Mail == "" {
					startNodeId = user.Id
				}
				// relation to the user 
				_, err = graph.DirectRelation(startNodeId, groupMembers.GroupId, plugin.BELONGS_TO, nil)
				if err != nil {
					return fmt.Errorf("unable to create user to groupmember relation: %v", err)
				}

				// fmt.Println("000000000")
				// fmt.Printf("adding group member relation to node %s\n", groupMembers.GroupId)
				// // relation to the group 
				// _, err = newNode.NewRelation(groupMembers.GroupId, plugin.BELONGS_TO, nil)
				// if err != nil {
				// 	return fmt.Errorf("unable to create group to user relation: %v", err)
				// }
			}

	return nil
}