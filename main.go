package main

import (
    "fmt"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	plugin "github.com/defensestation/pluginutils"
)

func startPlugin(ctx context.Context, mainEvent events.CloudWatchEvent) ([]byte, error) {
    fmt.Println("running")
    event := &plugin.Event{}
    json.Unmarshal(mainEvent.Detail, event)
    fmt.Println(event)
    fmt.Println("dsgadsgsagsdf")

    newPlugin, err := plugin.New(ctx, event)
    if err != nil {
        return nil, err
    }

    azurePlugin := NewAzurePlugin(newPlugin)

    err = azurePlugin.Run(ctx)
    if err != nil {
        return nil, err
    }

    return newPlugin.Complete()
}

func main() {
    //test
    lambda.Start(startPlugin)
}
