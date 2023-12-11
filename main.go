package main

import (
	"context"
	"encoding/json"


	plugin "github.com/defensestation/pluginutils"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"   
)

const (
	pluginName = "azuread"
	employeeType = "employee"
)


func startPlugin(ctx context.Context, mainEvent events.CloudWatchEvent) ([]byte, error) {
	event := &plugin.Event{}
	json.Unmarshal(mainEvent.Detail, event)

	newPlugin, err := plugin.New(event)
	if err != nil {
		return nil, err
	}

	azureAD := NewAzureADPlugin(newPlugin)

	err = azureAD.Run(ctx)	
	if err != nil {
		return nil, err
	}

	return newPlugin.Complete() 
}


func main() {
	lambda.Start(startPlugin)
}