package main

import (
	"fmt"
	"bufio"
	"log"
	"os"
	"encoding/json"
	"context"
	"errors"

    cmd "github.com/defensestation/azurehound/cmd"
    enums "github.com/defensestation/azurehound/enums"
    plugin "github.com/defensestation/pluginutils"
    // models "github.com/defensestation/azurehound/models"
    client "github.com/defensestation/azurehound/client"
    client_config "github.com/defensestation/azurehound/client/config"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"   
)

const (
	pluginType = "azuread"
	usersType = "user"
)

type Options struct {
	ApplicationId string 
	ClientSecret string 
	Tenant string
	Users bool
	Groups bool
	// Apps bool
} 

func (o *Options) Validate () error {
	if o.ApplicationId == "" || o.ClientSecret == "" || o.Tenant == "" {
		return errors.New("not proper creds given")
	}
	return nil
}

type AzureADPlugin struct {
	client *client.AzureClient
	Response []map[string]string
	plugin *plugin.Plugin
}

func NewAzureADPlugin(ctx context.Context, options *Options, name, accountId string) *AzureADPlugin {
	// ctx := context.Background()
	config := client_config.Config{
		ApplicationId:  options.ApplicationId,
		Authority:      "",
		ClientSecret:   options.ClientSecret,
		ClientCert:     "",
		ClientKey:      "",
		ClientKeyPass:  "",
		Graph:          "",
		JWT:            "",
		Management:     "",
		MgmtGroupId:    []string{""},
		Password:       "",
		ProxyUrl:       "",
		RefreshToken:   "",
		Region:         "cloud",
		SubscriptionId: []string{""},
		Tenant:         options.Tenant,
		Username:       "",
	}

	log.Println("creating client...")
	azClient, err := client.NewClient(config)
	if err != nil {
		log.Println(err)
		return nil
	}

	// start plugin
	newPlugin := plugin.New(name, accountId)

	return &AzureADPlugin{
		client: &azClient,
		Response: []map[string]string{},
		plugin: newPlugin,
	}
}

func startPlugin(ctx context.Context, mainEvent events.CloudWatchEvent) ([]byte, error) { // mainEvent events.CloudWatchEvent
	
	event := &plugin.Event{}
	json.Unmarshal(mainEvent.Detail, event)
		
	options := &Options{}
	err := event.MarshalOptions(options)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("----")
	fmt.Println(options)
	fmt.Println("========================")
	fmt.Println(event)
	fmt.Println("----")

	// validate
	err = options.Validate()
	if err != nil {
		return nil, err
	}
	
	azuread := NewAzureADPlugin(ctx, options, event.Name, event.AccountId)
	fmt.Println(options)
	// // get user graph
	// if options.Users {
	// 	err := azuread.getUsers("users")
	// 	fmt.Println(err)
	// }

	log.Println("client created, listing all data...")
	stream := cmd.ListAll(ctx, *azuread.client)
	file, err := os.Create("output.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// Iterate over the stream and write each item to the file
	for item := range stream {
		// Assuming the data in the stream is structured, you may need to adjust this part
		if err := encoder.Encode(item); err != nil {
			panic(err)
		}
	}

	// Open the file
	file, err = os.Open("output.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

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

		if azureWrapper.Kind == enums.KindAZUser {
			azuread.CreateUserNode(azureWrapper)
		}
		// if azureWrapper.Kind == enums.KindAZGroup {
		// 	CreateGroupNode(azureWrapper)
		// }
		// if azureWrapper.Kind == enums.AZGroupMember {
		// 	CreateGroupRelation(azureWrapper)
		// }
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	response, _ := json.Marshal(azuread.Response)
	return response, nil
}

func main() {
	// startPlugin(context.Background(), event)
	lambda.Start(startPlugin)
}