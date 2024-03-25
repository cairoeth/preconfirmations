package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/cairoeth/preconfirmations/preconf-operator"
	"github.com/cairoeth/preconfirmations/preconf-operator/core/config"
	"github.com/cairoeth/preconfirmations/preconf-operator/types"

	sdkutils "github.com/Layr-Labs/eigensdk-go/utils"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{config.ConfigFileFlag}
	app.Name = "preconf-operator"
	app.Usage = "Preconfirmations Operator"
	app.Description = "Service that subcribes to preconf-share and sends preconfirmations."

	app.Action = operatorMain
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("Application failed. Message:", err)
	}
}

func operatorMain(ctxCli *cli.Context) error {
	configPath := ctxCli.GlobalString(config.ConfigFileFlag.Name)
	nodeConfig := types.NodeConfig{}
	err := sdkutils.ReadYamlConfig(configPath, &nodeConfig)
	if err != nil {
		return err
	}
	configJSON, err := json.MarshalIndent(nodeConfig, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Println("Config:", string(configJSON))

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
	// log.Println("Started operator")

	// jsonRPCServer, err := jsonrpcserver.NewHandler(jsonrpcserver.Methods{
	// 	"receive": operator.Receive,
	// })
	// if err != nil {
	// 	log.Fatal("Failed to create jsonrpc server", err)
	// }

	// ctx, ctxCancel := context.WithCancel(context.Background())

	// http.Handle("/", jsonRPCServer)
	// server := &http.Server{
	// 	Addr:              fmt.Sprintf(":%s", "8000"),
	// 	ReadHeaderTimeout: 5 * time.Second,
	// }

	// connectionsClosed := make(chan struct{})
	// go func() {
	// 	notifier := make(chan os.Signal, 1)
	// 	signal.Notify(notifier, os.Interrupt, syscall.SIGTERM)
	// 	<-notifier
	// 	log.Println("Shutting down...")
	// 	ctxCancel()
	// 	if err := server.Shutdown(context.Background()); err != nil {
	// 		log.Println("Failed to shutdown server", err)
	// 	}
	// 	close(connectionsClosed)
	// }()

	// err = server.ListenAndServe()
	// if err != nil && !errors.Is(err, http.ErrServerClosed) {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	// <-ctx.Done()
	// <-connectionsClosed

	return nil
}
