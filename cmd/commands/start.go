package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/binacsgo/inject"
	"github.com/binacsgo/log"

	"github.com/OpenKikCoc/raftkv/config"
	"github.com/OpenKikCoc/raftkv/service"
)

func init() {
	rootStartFlags(RootCmd)
}

func rootStartFlags(cmd *cobra.Command) {
}

var (
	// StartCmd the root command
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start Command",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			fmt.Println(*cfg)
			node := initService(logger, cfg)
			fmt.Println("node = ", node)
			if err = node.OnStart(); err != nil {
				fmt.Println(err)
			}

			return nil
		},
	}
)

func initService(logger log.Logger, cfg *config.Config) *service.NodeServiceImpl {
	nodeSvc := service.NodeServiceImpl{}

	inject.Regist(Inject_Config, cfg)

	inject.Regist(Inject_LOGGER, logger)
	inject.Regist(Inject_ZAPLOGGER, zaplogger)
	inject.Regist(Inject_Node_LOGGER, logger.With("module", "node"))
	inject.Regist(Inject_Web_LOGGER, logger.With("module", "web"))
	inject.Regist(Inject_GRPC_LOGGER, logger.With("module", "grpc"))

	inject.Regist(Inject_Node_Service, &nodeSvc)

	inject.Regist(Inject_Web_Service, &service.WebServiceImpl{})

	inject.Regist(Inject_GRPC_Service, &service.GRPCServiceImpl{})

	err := inject.DoInject()
	if err != nil {
		panic(err.Error())
	}
	return &nodeSvc
}
