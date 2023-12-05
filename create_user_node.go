package main

import (
	"encoding/json"
	"fmt"
	cmd "github.com/defensestation/azurehound/cmd"
	models "github.com/defensestation/azurehound/models"
	plugin "github.com/defensestation/pluginutils"
	// grapher "github.com/defensestation/grapher"
)

func (a *AzureADPlugin) CreateUserNode(azureWrapper *cmd.AzureWrapper) error {
	// fmt.Printf("Kind: %s, Data: %s\n", azureWrapper.Kind, azureWrapper.Data)
	// unmarshal
	var user *models.User
	userJson, err := json.Marshal( azureWrapper.Data) 
	if err != nil {
		fmt.Println("Error:", err)
	}


	err = json.Unmarshal(userJson, &user)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// fmt.Println(user.User.Id)
	// create schema
	mapValue, _ := StructToInterface(user.User)
	fmt.Println(mapValue)

	opts := &plugin.MapOptions{
		Ignore: []string{
			"password",
		},
	}

	propsMap := plugin.NewPropsTypeMap(mapValue.(map[string]interface{}), opts)
	propsMap.AddType("user_type","string")
	propsMap.AddType("email","string")
	propsMap.AddType("last_name","string")
	propsMap.AddType("first_name","string")
	propsMap.AddType("name","string")
	propsMap.AddType("personnel_id","string")

	schema := &plugin.Schema{
		ID: []string{"mail"},
		Labels: []string{"mail", "personnel", "alais"},
		Props: propsMap.GetMap(),
		RelationsKey: "relations",
		Relation: &plugin.Relation{
			ID: []string{"relations_id"},
			EndID: []string{"email"},
			Type: []string{"rtype"},
			Props: map[string]string{},
		},
		ReplaceKeys: map[string]string{
			"mail": "email",
			"name__familyName": "last_name",
			"name__givenName": "first_name",
			"name__fullName": "name",
		},
	}

	newGraph, err := a.plugin.NewGraph(usersType, schema)
	if err != nil {
		return err
	}
	// newGraph.SetReverseRelation(true
	itemType := mapValue.(map[string]interface{})

	itemType["personnel"] = "personnel"
	itemType["personnel_id"] = fmt.Sprintf("%s_%s",a.plugin.Name,itemType["mail"])

	err = newGraph.AddNode(itemType)
	if err != nil {
		return err
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