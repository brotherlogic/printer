package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/printer/proto"
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/printer/config"
)

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

//Server main server type
type Server struct {
	*goserver.GoServer
	whitelist  []string
	prints     int64
	pretend    bool // Used for testing only
	pretendret error
	config     *pb.Config
}

func (s *Server) localPrint(text string, lines []string, ti time.Time) error {
	if s.pretend {
		s.prints++
		return s.pretendret
	}

	s.Log(fmt.Sprintf("Assessing print at %v", ti))
	if ti.Hour() < 9 || ti.Hour() > 17 || ((ti.Weekday() == time.Saturday || ti.Weekday() == time.Sunday) && (ti.Hour() != 10)) {
		return status.Errorf(codes.Unavailable, "Not the time to print right now")
	}

	s.prints++
	s.config.TotalPrints++
	cmd := exec.Command("sudo", "python", "/home/simon/gobuild/src/github.com/brotherlogic/printer/printText.py", text)
	if len(text) == 0 {
		all := []string{"sudo", "python", "/home/simon/gobuild/src/github.com/brotherlogic/printer/printText.py"}
		all = append(all, lines...)
		cmd = exec.Command("sudo", all...)
	}

	output := ""
	out, err := cmd.StdoutPipe()

	if err != nil {
		s.Log(fmt.Sprintf("Error stdout: %v", err))
	}

	if out != nil {
		scanner := bufio.NewScanner(out)
		go func() {
			for scanner != nil && scanner.Scan() {
				output += scanner.Text()
			}
			out.Close()
		}()
	}

	cmd.Start()
	err = cmd.Wait()

	s.Log(fmt.Sprintf("OUTPUT = %v", output))
	return err
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		[]string{
			"beerserver",
			"recordprinter",
		},
		int64(0),
		false, // Prod version doesn't pretend to print
		nil,
		&pb.Config{},
	}
	s.GoServer.KSclient = *keystoreclient.GetClient(s.DialMaster)
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterPrintServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.save(ctx)
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		return s.load(ctx)
	}

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "prints", Value: s.prints},
		&pbg.State{Key: "whitelisted", Value: int64(len(s.whitelist))},
		&pbg.State{Key: "backlog", Value: int64(len(s.config.Requests))},
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var init = flag.Bool("init", false, "Init the config")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	err := server.RegisterServerV2("printer", false, false)
	if err != nil {
		return
	}

	if *init {
		ctx, cancel := utils.BuildContext("printer", "printer")
		defer cancel()
		server.config.TotalPrints = 1
		server.save(ctx)
		return
	}

	server.RegisterRepeatingTask(server.processPrints, "process_prints", time.Minute)

	fmt.Printf("%v", server.Serve())
}
