package main

import (
	"fmt"
	"context"
	"encoding/json"
	"os"
	"bufio"

	cmd "github.com/defensestation/azurehound/cmd"
	enums "github.com/defensestation/azurehound/enums"
)

func (ad *AzureADPlugin) List(ctx context.Context) error {
	// list all
   stream := cmd.ListAll(ctx, *ad.Client)
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
			fmt.Println("not handled by plugin ", azureWrapper.Kind )

		}

	}

	return nil
}