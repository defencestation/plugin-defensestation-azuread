package main

import (
	"context"
	"encoding/json"
	"fmt"

	cmd "github.com/defensestation/azurehound/cmd"
	enums "github.com/defensestation/azurehound/enums"
)

func (ad *AzureADPlugin) List(ctx context.Context) error {
	// List all
	stream := cmd.ListAll(ctx, *ad.Client)

	for item := range stream {
		// Assert the type of item to []byte
		data, ok := item.([]byte)
		if !ok {
			fmt.Println("Error: item is not of type []byte")
			continue
		}

		// Assuming the data in the stream is structured
		var azureWrapper *cmd.AzureWrapper

		// Unmarshal the JSON data from the stream item into the AzureWrapper struct
		err := json.Unmarshal(data, &azureWrapper)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		switch azureWrapper.Kind {
		case enums.KindAZUser:
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
			fmt.Println("not handled by plugin ", azureWrapper.Kind)
		}
	}

	return nil
}
