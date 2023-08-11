package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brotherlogic/goserver/utils"

	pbp "github.com/brotherlogic/printer/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("PrintCLI", "printer")
	defer cancel()

	conn, err := utils.LFDialServer(ctx, "printer")
	if err != nil {
		log.Fatalf("Bad dial: %v", err)
	}
	defer conn.Close()

	client := pbp.NewPrintServiceClient(conn)

	if os.Args[1] == "clear" {
		_, err := client.Clear(ctx, &pbp.ClearRequest{})
		fmt.Printf("CLEAR: %v\n", err)
	} else if os.Args[1] == "ping" {
		pong, err := client.Ping(ctx, &pbp.PingRequest{})
		fmt.Printf("PING: %v %v\n", pong, err)
	} else if os.Args[1] == "list" {
		re, err := client.List(ctx, &pbp.ListRequest{})
		if err != nil {
			log.Fatalf("Error on list: %v", err)
		}
		for _, elem := range re.GetQueue() {
			fmt.Printf("%v\n", elem)
		}
	} else {

		r, err := client.Print(ctx, &pbp.PrintRequest{Lines: os.Args, Origin: "recordprinter"})
		fmt.Printf("%v and %v -> %v\n", r, err, &pbp.PrintRequest{Lines: os.Args})

	}
}
