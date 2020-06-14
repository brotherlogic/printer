package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pbp "github.com/brotherlogic/printer/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("PrintCLI", "printer")
	defer cancel()

	host, port, _ := utils.Resolve("printer", "printer-cli")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	client := pbp.NewPrintServiceClient(conn)

	if os.Args[1] == "clear" {
		client.Clear(ctx, &pbp.ClearRequest{})
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
