package main

import (
	"fmt"
	"strings"
	"bufio"
	"os"
	"encoding/json"
	"context"

	cmd "github.com/defensestation/azurehound/cmd"
    enums "github.com/defensestation/azurehound/enums"
    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

func (ad *AzureAd) GetGroupUsers(ctx context.Context, groupName string) error {
	stream := cmd.ListAll(ctx, *ad.client)
	file, err := os.Create("/tmp/output.json")
	if err != nil {
		return err 
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// Iterate over the stream and write each item to the file
	for item := range stream {
		// Assuming the data in the stream is structured, you may need to adjust this part
		if err := encoder.Encode(item); err != nil {
			 return err
		}
	}

	// Open the file
	file, err = os.Open("/tmp/output.json")
	if err != nil {
		return  err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)



	graph, ok := ad.plugin.GetGraph(employeeType)
	if !ok {
		return fmt.Errorf("unable to find %s graph", employeeType)
	} 


	// Iterate through lines and parse each JSON object
	for scanner.Scan() {
		line := scanner.Text()

		// Create a new AppOwnerNode to store the parsed data
		var azureWrapper *cmd.AzureWrapper

		// Unmarshal the JSON data from the line into the AppOwnerNode struct
		err := json.Unmarshal([]byte(line), &azureWrapper)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if azureWrapper.Kind == enums.KindAZGroupMember {
			var groupMembers *models.GroupMembers
			groupMembersJson, err := json.Marshal(azureWrapper.Data) 
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
				groupMemberMapInterface := plugin.StructToMap(user)

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
				groupMemberMapInterface["service"] = "dsc_service_policy_manager"
				groupMemberMapInterface["personnel"]    = "personnel"
				groupMemberMapInterface["personnel_id"] = fmt.Sprintf("%s_%s", ad.plugin.Name, user.Mail)

				err = graph.AddNode(groupMemberMapInterface)
				if err != nil {
					return err
				}
			}
		}
}
return nil
}