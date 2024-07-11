package main

import (
	"fmt"
	"strings"
	// "encoding/json"
	"context"

	// cmd "github.com/defensestation/azurehound/cmd"
    // enums "github.com/defensestation/azurehound/enums"
    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	userType = "AzureADUser"
	noMailFound = "NoMailFound"
)

func (ad *AzureADPlugin) GetUsers(ctx context.Context, data interface{}) error {
	user := data.(models.User)

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
	// log to debug
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Printf("adding user node %s\n", user.Mail)
	labels := []string{userType}
	graph := ad.Plugin.AddOrFindGraph(userType, plugin.NewSchema(nil))

	nodeId := user.Mail
	if user.Mail == "" {
		nodeId = user.Id
		labels = append(labels, noMailFound)
	}
	_, err = graph.NewNode(plugin.Personnel, userType, nodeId, user.DisplayName, labels, props)
	if err != nil {
		fmt.Errorf("unable to create user node: %v", err)
	
	}

	return nil
}