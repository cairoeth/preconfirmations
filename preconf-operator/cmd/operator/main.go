package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/cairoeth/preconfirmations-avs/preconf-operator"
	"github.com/cairoeth/preconfirmations-avs/preconf-operator/core/config"
	"github.com/cairoeth/preconfirmations-avs/preconf-operator/types"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{config.ConfigFileFlag}
	app.Name = "preconf-operator"
	app.Usage = "Preconfirmations Operator"
	app.Description = "Service that subcribes to Preconf-Share and sends preconfirmations."

	app.Action = operatorMain
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Application failed. Message:", err)
	}
}

func operatorMain(ctx *cli.Context) error {
	configPath := ctx.GlobalString(config.ConfigFileFlag.Name)
	nodeConfig := types.NodeConfig{}
	err := sdkutils.ReadYamlConfig(configPath, &nodeConfig)
	if err != nil {
		return err
	}
	configJson, err := json.MarshalIndent(nodeConfig, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println("Config:", string(configJson))

	log.Println("Initializing operator")
	operator, err := operator.NewOperatorFromConfig(nodeConfig)
	if err != nil {
		return err
	}
	log.Println("Initialized operator")

	log.Println("Starting operator")
	err = operator.Start(context.Background())
	if err != nil {
		return err
	}
	log.Println("Started operator")

	return nil
}
