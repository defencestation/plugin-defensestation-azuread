package main

import (
	"fmt"
	// "encoding/json"
	"context"

    plugin "github.com/defensestation/pluginutils"
    models "github.com/defensestation/azurehound/models"
)

const (
	appType = "AzureADApp"
)

func (ad *AzureADPlugin) GetApps(ctx context.Context, data interface{}) error {
			app := data.(models.App)
			// appJson, err := json.Marshal(data) 
			// if err != nil {
			// 	return err
			// }

			// err = json.Unmarshal(appJson, &app)
			// if err != nil {
			// 	return err
			// }

			props, err := plugin.StructToMap(app)
			if err != nil {
				fmt.Errorf("failed to marshal: %v", err)
			}

			labels := []string{appType}
			graph := ad.Plugin.AddOrFindGraph(appType, plugin.NewSchema(nil))

			fmt.Println("##############")
			fmt.Printf("adding app  node %s\n", app.DisplayName)
			_, err = graph.NewNode(plugin.Application, appType, app.AppId, app.DisplayName, labels, props)
			if err != nil {
				fmt.Errorf("unable to create app node: %v", err)
			}
return nil
}