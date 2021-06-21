package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	pb "golang-grpc-sample/proto"

	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"golang-grpc-sample/keycloak"
	"golang-grpc-sample/utils"
)

var logger = utils.NewLogger()

func initConfig() {
	var configPath string

	flag.StringVar(&configPath, "cf", ".", "Path to the config file")
	flag.Parse()

	viper.SetConfigFile("config/config.toml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("No valid config file is provided: %s", err.Error()))
	}

	viper.SetDefault("keycloak.uri", "http://localhost:8080")
	viper.SetDefault("keycloak.master_realm", "master")
}

func main() {
	initConfig()
	logger.Info("GRPC GO")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterKeycloakProvisionServer(
		grpcServer,
		keycloak.NewKcProvision(
			&keycloak.KcProvisionOpts{
				KcConfig: keycloak.KcConfig{
					MasterRealm:   viper.GetString("keycloak.master_realm"),
					AdminUsername: viper.GetString("keycloak.username"),
					AdminPassword: viper.GetString("keycloak.password"),
					KeycloakURI:   viper.GetString("keycloak.uri"),
				},
			},
		))

	reflection.Register(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
