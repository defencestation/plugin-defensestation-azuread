package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"context"

	// cmd "github.com/defensestation/azurehound/cmd"
    // enums "github.com/defensestation/azurehound/enums"
    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	userType = "AzureADUser"
)

func (ad *AzureADPlugin) GetUsers(ctx context.Context, data interface{}) error {
	var user *models.User
	userJson, err := json.Marshal(data) 
	if err != nil {
		return err
	}

	err = json.Unmarshal(userJson, &user)
	if err != nil {
		return err
	}

	props, err := plugin.StructToMap(user)
	if err != nil {
		fmt.Errorf("failed to marshal: %v", err)
	}
	if user.DisplayName != "" {
		splitName := strings.Split(user.DisplayName, " ")
		if len(splitName) >= 2 {
			props["first_name"]  = splitName[0]
			props["last_name"] = splitName[1]
		}
		if len(splitName) == 1 {
			props["first_name"]  = splitName[0]
		}
	}

	props["type"]    	 = "employee"
	// userMapInterface["service"] = "dsc_service_policy_manager"
	props["personnel"]    = "personnel"
	// props["personnel_id"] = fmt.Sprintf("%s_%s", ad.plugin.Name, user.Mail)

	labels := []string{userType}
	graph := ad.Plugin.AddOrFindGraph(userType, plugin.NewSchema(nil))

	
	_, err = graph.NewNode(plugin.Personnel, userType, user.Mail, user.DisplayName, labels, props)
	if err != nil {
		fmt.Errorf("unable to create user node: %v", err)
	
	}

	return nil
}