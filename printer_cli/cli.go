package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pbgs "github.com/brotherlogic/goserver/proto"
	pbp "github.com/brotherlogic/printer/proto"
	pbt "github.com/brotherlogic/tracer/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	ctx, cancel := utils.BuildContext("PrintCLI", "printer", pbgs.ContextType_MEDIUM)
	defer cancel()

	host, port, _ := utils.Resolve("printer")
	conn, _ := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	client := pbp.NewPrintServiceClient(conn)
	r, err := client.Print(context.Background(), &pbp.PrintRequest{Lines: os.Args})
	fmt.Printf("%v and %v\n", r, err)

	utils.SendTrace(ctx, "PrintCLI", time.Now(), pbt.Milestone_END, "printer")
}
