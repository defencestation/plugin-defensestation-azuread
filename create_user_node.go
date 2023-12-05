package main

import (
	"encoding/json"
	"fmt"
	"strings"
	cmd "github.com/defensestation/azurehound/cmd"
	models "github.com/defensestation/azurehound/models"
	plugin "github.com/defensestation/pluginutils"
	// grapher "github.com/defensestation/grapher"
)

func (a *AzureADPlugin) CreateUserNode(azureWrapper []*cmd.AzureWrapper) error {
	// fmt.Printf("Kind: %s, Data: %s\n", azureWrapper.Kind, azureWrapper.Data)
	// unmarshal
	var mapValues []interface{}
	for _, value := range azureWrapper {
		var user *models.User
		userJson, err := json.Marshal( value.Data) 
		if err != nil {
			fmt.Println("Error:", err)
		}


		err = json.Unmarshal(userJson, &user)
		if err != nil {
			fmt.Println("Error:", err)
		}

		mapValue, _ := StructToInterface(user.User)
		mapValues = append(mapValues, mapValue)
	}
	// fmt.Println(user.User.Id)
	// create schema
	
	if len(mapValues) == 0 {
		return nil
	}

	opts := &plugin.MapOptions{
		Ignore: []string{
			"password",
			"id",
		},
	}

	propsMap := plugin.NewPropsTypeMap(mapValues[0].(map[string]interface{}), opts)
	propsMap.AddType("user_type","string")
	propsMap.AddType("email","string")
	propsMap.AddType("last_name","string")
	propsMap.AddType("first_name","string")
	propsMap.AddType("name","string")
	propsMap.AddType("personnel_id","string")

	schema := &plugin.Schema{
		ID: []string{"email"},
		Labels: []string{"email", "personnel"},
		Props: propsMap.GetMap(),
		// RelationsKey: "relations",
		// Relation: &plugin.Relation{},
		ReplaceKeys: map[string]string{
			"mail": "email",
			"displayName": "name",
		},
	}

	newGraph, err := a.plugin.NewGraph(usersType, schema)
	if err != nil {
		return err
	}

	for _, mapValue := range mapValues {
		// newGraph.SetReverseRelation(true
		itemType := mapValue.(map[string]interface{})
		if itemType["displayName"] != nil {
			splitName := strings.Split(itemType["displayName"].(string), " ")
			if len(splitName) >= 2 {
				itemType["first_name"]  = splitName[0]
				itemType["last_name"] = splitName[1]
			}
			if len(splitName) == 1 {
				itemType["first_name"]  = splitName[0]
			}
		}
		itemType["personnel"] = "personnel"
		itemType["personnel_id"] = fmt.Sprintf("%s_%s",a.plugin.Name,itemType["mail"])

		err = newGraph.AddNode(itemType)
		if err != nil {
			return err
		}
	}

	response, err := newGraph.Save()
	if err != nil {
		return err
	}


	// g.Response = append(g.Response, alaisResponse)
	a.Response = append(a.Response, response)
	
	// save to graph
	return nil
}