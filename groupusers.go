package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	groupMemberType = "AzureADGroupMember"
)

func (ad *AzureADPlugin) GetGroupUsers(ctx context.Context, data interface{}) error {
			var groupMembers *models.GroupMembers
			groupMembersJson, err := json.Marshal(data) 
			if err != nil {
				return err
			}

			err = json.Unmarshal(groupMembersJson, &groupMembers)
			if err != nil {
				return err
			}
			for _, groupMember := range(groupMembers.Members) {
				var user *models.User
				userJson, err := json.Marshal(groupMember.Member) 
				if err != nil {
					return err
				}

				err = json.Unmarshal(userJson, &user)
				if err != nil {
					return err
				}
				groupMemberMapInterface, err := plugin.StructToMap(user)
				if err != nil {
					fmt.Errorf("failed to marshal: %v", err)
				}

				if user.DisplayName != "" {
					splitName := strings.Split(user.DisplayName, " ")
					if len(splitName) >= 2 {
						groupMemberMapInterface["first_name"]  = splitName[0]
						groupMemberMapInterface["last_name"] = splitName[1]
					}
					if len(splitName) == 1 {
						groupMemberMapInterface["first_name"]  = splitName[0]
					}
				}

				groupMemberMapInterface["type"]    	 = "employee"
				// groupMemberMapInterface["service"] = "dsc_service_policy_manager"
				groupMemberMapInterface["personnel"]    = "personnel"
				// groupMemberMapInterface["personnel_id"] = fmt.Sprintf("%s_%s", ad.plugin.Name, user.Mail)

				labels := []string{groupMemberType}
				graph := ad.Plugin.AddOrFindGraph(groupMemberType, plugin.NewSchema(nil))

				
				newNode, err := graph.NewNode(plugin.Personnel, groupMemberType, user.Id, user.DisplayName, labels, groupMemberMapInterface)
				if err != nil {
					fmt.Errorf("unable to create groupmember node: %v", err)
				}

				// relation to the user 
				_, err = newNode.NewRelation(user.Mail, plugin.BELONGS_TO, nil)
				if err != nil {
					return fmt.Errorf("unable to create user to groupmember relation: %v", err)
				}

				// relation to the group 
				_, err = newNode.NewRelation(groupMembers.GroupId, plugin.BELONGS_TO, nil)
				if err != nil {
					return fmt.Errorf("unable to create group to user relation: %v", err)
				}
			}

	return nil
}