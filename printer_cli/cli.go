package main

import (
	"context"
	"fmt"
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

	host, port, _ := utils.Resolve("printer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	client := pbp.NewPrintServiceClient(conn)

	if os.Args[1] == "clear" {
		client.Clear(context.Background(), &pbp.ClearRequest{})
	} else {

		r, err := client.Print(context.Background(), &pbp.PrintRequest{Lines: os.Args})
		fmt.Printf("%v and %v -> %v\n", r, err, &pbp.PrintRequest{Lines: os.Args})

	}
}
