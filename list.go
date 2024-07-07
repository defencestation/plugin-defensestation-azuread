package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "reflect"

	cmd "github.com/defensestation/azurehound/cmd"
	enums "github.com/defensestation/azurehound/enums"
)

func (ad *AzureADPlugin) List(ctx context.Context) error {
	// List all
	stream := cmd.ListAll(ctx, *ad.Client)

	for item := range stream {
		// Print the type of item
		// fmt.Printf("Type of item: %s\n", reflect.TypeOf(item))

		// Assuming the data in the stream is structured
		var azureWrapper *cmd.AzureWrapper

		// Convert the item to JSON if necessary
		data, err := json.Marshal(item)
		if err != nil {
			fmt.Println("Error marshalling item:", err)
			continue
		}

		// Unmarshal the JSON data from the stream item into the AzureWrapper struct
		err = json.Unmarshal(data, &azureWrapper)
		if err != nil {
			fmt.Println("Error unmarshalling item:", err)
			continue
		}

		switch azureWrapper.Kind {
		case enums.KindAZUser:
			// fmt.Println("getting users")
			err := ad.GetUsers(ctx, azureWrapper.Data)
			if err != nil {
				fmt.Println(err)
				return err
			}

		case enums.KindAZApp:
			err := ad.GetApps(ctx, azureWrapper.Data)
			if err != nil {
				fmt.Println(err)
				return err
			}

		case enums.KindAZGroup:
			err := ad.GetGroups(ctx, azureWrapper.Data)
			if err != nil {
				fmt.Println(err)
				return err
			}

		case enums.KindAZGroupMember:
			err := ad.GetGroupUsers(ctx, azureWrapper.Data)
			if err != nil {
				fmt.Println(err)
				return err
			}

		case enums.KindAZServicePrincipal:
			err := ad.GetServicePrincipal(ctx, azureWrapper.Data)
			if err != nil {
				fmt.Println(err)
				return err
			}

		default:
			// fmt.Println("not handled by plugin ", azureWrapper.Kind)
		}
	}

	return nil
}
