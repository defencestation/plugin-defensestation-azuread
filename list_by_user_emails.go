package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	plugin "github.com/defensestation/pluginutils"
)

const (
	roleType  = "AzureADRole"
	// groupType = "AzureADGroup"
	// userType  = "AzureADUser"
)

func (a *AzureADPlugin) ListByUserEmails(ctx context.Context) error {
	cred, err := azidentity.NewClientSecretCredential(a.Options.TenantID, a.Options.ClientID, a.Options.ClientSecret, nil)
	if err != nil {
		return fmt.Errorf("unable to create credential: %w", err)
	}
	a.Credential = cred

	if err := a.ListByUserGraph(ctx); err != nil {
		return err
	}

	return nil
}

func (a *AzureADPlugin) ListByUserGraph(ctx context.Context) error {
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(a.Credential, nil)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	// Create the filter query to find users by email
    var filters []string
    for _, email := range a.Options.UserEmailArray {
        filters = append(filters, fmt.Sprintf("mail eq '%s'", email))
    }
    filterQuery := strings.Join(filters, " or ")

    // Define the request configuration
    requestConfig := &users.UsersRequestBuilderGetQueryParameters{
        Filter: &filterQuery,
    }
    requestOptions := &users.UsersRequestBuilderGetRequestConfiguration{
        QueryParameters: requestConfig,
    }

	result, err := client.Users().Get(ctx, requestOptions)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.Userable](result, client.GetAdapter(), models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return fmt.Errorf("error creating page iterator: %w", err)
	}

	return pageIterator.Iterate(ctx, func(user models.Userable) bool {
		userProperties := GetUserProperties(user)
		props, err := plugin.StructToMap(userProperties)
		if err != nil {
			log.Printf("failed to marshal: %v", err)
			return true
		}

		populateUserProps(*userProperties, props)

		graph := a.Plugin.AddOrFindGraph(userType, plugin.NewSchema(nil))
		
		nodeId := ""
		if userProperties.Mail != nil{
			nodeId = *userProperties.Mail
		}
		nodeId = *userProperties.UserPrincipalName
		fmt.Println(nodeId)
		userNode, err := graph.NewNode(plugin.Personnel, userType, nodeId, *userProperties.DisplayName, []string{userType}, props)
		if err != nil {
			log.Printf("unable to create user node: %v", err)
			return true
		}

		if err := a.processUserMemberships(ctx, client, user, *userNode); err != nil {
			log.Printf("error processing user memberships: %v", err)
		}

		return true
	})
}

func populateUserProps(userProperties UserProperties, props map[string]interface{}) {
	if *userProperties.DisplayName != "" {
		splitName := strings.Split(*userProperties.DisplayName, " ")
		props["first_name"] = splitName[0]
		if len(splitName) >= 2 {
			props["last_name"] = splitName[1]
		}
	}
	props["type"] = "employee"
	props["personnel"] = "personnel"
}

func (a *AzureADPlugin) processUserMemberships(ctx context.Context, client *msgraphsdk.GraphServiceClient, user models.Userable, userNode plugin.Node) error {
	memberOf, err := client.Users().ByUserId(*user.GetUserPrincipalName()).MemberOf().Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("error getting user memberships: %w", err)
	}

	members := memberOf.GetValue()
	for _, member := range members {
		directoryObject, ok := member.(models.DirectoryObjectable)
		if !ok {
			continue
		}

		switch obj := directoryObject.(type) {
		case *models.Group:
			if err := a.addGroupNode(obj); err != nil {
				log.Printf("error adding group node: %v", err)
			}
		case *models.DirectoryRole:
			if err := a.addRoleNode(obj); err != nil {
				log.Printf("error adding role node: %v", err)
			}
		case *models.AdministrativeUnit:
			if err := a.addAdminUnitNode(obj); err != nil {
				log.Printf("error adding admin unit node: %v", err)
			}
		}

		if _, err := userNode.NewRelation(*directoryObject.GetId(), plugin.BELONGS_TO, nil); err != nil {
			log.Printf("error creating relation: %v", err)
		}
	}

	return nil
}

func (a *AzureADPlugin) addGroupNode(group *models.Group) error {
	groupProperties := GetGroupProperties(group)
	props, err := plugin.StructToMap(groupProperties)
	if err != nil {
		return fmt.Errorf("failed to marshal group properties: %w", err)
	}

	graph := a.Plugin.AddOrFindGraph(groupType, plugin.NewSchema(nil))
	_, err = graph.NewNode(plugin.Group, groupType, *group.GetId(), *groupProperties.DisplayName, []string{groupType}, props)
	if err != nil {
		return fmt.Errorf("unable to create group node: %w", err)
	}

	return nil
}

func (a *AzureADPlugin) addRoleNode(role *models.DirectoryRole) error {
	roleProperties := GetRoleProperties(role)
	props, err := plugin.StructToMap(roleProperties)
	if err != nil {
		return fmt.Errorf("failed to marshal role properties: %w", err)
	}

	graph := a.Plugin.AddOrFindGraph(roleType, plugin.NewSchema(nil))
	_, err = graph.NewNode(plugin.Role, roleType, *role.GetId(), *roleProperties.DisplayName, []string{roleType}, props)
	if err != nil {
		return fmt.Errorf("unable to create role node: %w", err)
	}

	return nil
}

func (a *AzureADPlugin) addAdminUnitNode(adminUnit *models.AdministrativeUnit) error {
	adminUnitProperties := GetAdminUnitProperties(adminUnit)
	props, err := plugin.StructToMap(adminUnitProperties)
	if err != nil {
		return fmt.Errorf("failed to marshal admin unit properties: %w", err)
	}

	graph := a.Plugin.AddOrFindGraph(roleType, plugin.NewSchema(nil))
	_, err = graph.NewNode(plugin.Role, roleType, *adminUnit.GetId(), *adminUnitProperties.DisplayName, []string{roleType}, props)
	if err != nil {
		return fmt.Errorf("unable to create admin unit node: %w", err)
	}

	return nil
}
