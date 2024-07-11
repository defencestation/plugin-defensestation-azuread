package main

import (
	"context"
	// "encoding/json"
	"fmt"
	// "sync"

	cmd "github.com/defensestation/azurehound/cmd"
	enums "github.com/defensestation/azurehound/enums"
)

func (ad *AzureADPlugin) List(ctx context.Context) error {
	// List all
	stream := cmd.ListAll(ctx, *ad.Client)

	fmt.Println("started getting data")


	for item := range stream {
		processItem(ctx, ad, item)
	}

	fmt.Println("done getting data")

	return nil
}

func processItem(ctx context.Context, ad *AzureADPlugin, item interface{}) {
	// Assuming the data in the stream is structured
	// Assuming the data in the stream is structured
	azureWrapper, ok := item.(cmd.AzureWrapper)
	if !ok {
		fmt.Println("Error: item is not of type cmd.AzureWrapper")
		return
	}

	switch azureWrapper.Kind {
	case enums.KindAZUser:
		err := ad.GetUsers(ctx, azureWrapper.Data)
		if err != nil {
			fmt.Println(err)
		}

	case enums.KindAZApp:
		err := ad.GetApps(ctx, azureWrapper.Data)
		if err != nil {
			fmt.Println(err)
		}

	case enums.KindAZGroup:
		err := ad.GetGroups(ctx, azureWrapper.Data)
		if err != nil {
			fmt.Println(err)
		}

	case enums.KindAZGroupMember:
		err := ad.GetGroupUsers(ctx, azureWrapper.Data)
		if err != nil {
			fmt.Println(err)
		}

	case enums.KindAZServicePrincipal:
		err := ad.GetServicePrincipal(ctx, azureWrapper.Data)
		if err != nil {
			fmt.Println(err)
		}

	default:
		fmt.Println("not handled by plugin", azureWrapper.Kind)
	}
}
