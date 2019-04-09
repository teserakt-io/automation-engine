package main

import (
	"log"

	"google.golang.org/grpc"

	"gitlab.com/teserakt/c2se/internal/cli"
	"gitlab.com/teserakt/c2se/internal/pb"
)

func main() {
	cnx, err := grpc.Dial("127.0.0.1:5556", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer cnx.Close()
	client := pb.NewC2ScriptEngineClient(cnx)

	rootCmd := cli.NewRootCommand(client)
	rootCmd.Execute()
}
