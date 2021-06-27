package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

var (
	//Backlog - the print queue
	Backlog = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "printer_backlog",
		Help: "The size of the print queue",
	})
	printErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "printer_errors",
		Help: "The size of the print queue",
	}, []string{"error"})
)

func (s *Server) load(ctx context.Context) (*pb.Config, error) {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return nil, err
	}

	return data.(*pb.Config), nil
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	return s.KSclient.Save(ctx, KEY, config)
}

//Server main server type
type Server struct {
	*goserver.GoServer
	prints     int64
	printq     chan *pb.PrintRequest
	pretend    bool // Used for testing only
	pretendret error
	done       chan bool
	outOfPaper bool
}

func (s *Server) localPrint(text string, lines []string, ti time.Time) (time.Duration, error) {
	if s.pretend {
		if s.pretendret == nil {
			s.prints++
		}
		return time.Second, s.pretendret
	}

	if ti.Hour() < 9 || ti.Hour() > 20 { // || ((ti.Weekday() == time.Saturday || ti.Weekday() == time.Sunday) && (ti.Hour() != 10)) {
		return time.Minute, status.Errorf(codes.Unavailable, "Not the time to print right now")
	}

	s.prints++
	s.Log(fmt.Sprintf("NOW PRINTING: %v", lines))

	if len(text) != 0 {
		ioutil.WriteFile("/home/simon/print.txt", []byte(text), 0644)
	} else {
		os.Create("home/simon/print.txt")
		handle, err := os.Open("/home/simon/print.txt")
		if err != nil {
			return time.Second, err
		}
		for _, line := range lines {
			handle.WriteString(fmt.Sprintf("%v\n", line))
		}
		//Feed
		for i := 0; i < 5; i++ {
			handle.WriteString("\n")
		}
		serr := handle.Sync()
		cerr := handle.Close()
		s.Log(fmt.Sprintf("CLOSE %v and %v", serr, cerr))
	}

	cmd := exec.Command("lp", "/home/simon/print.txt")
	output := ""
	out, err := cmd.StdoutPipe()

	if err != nil {
		s.Log(fmt.Sprintf("Error in the actual stdout: %v", err))
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

	s.Log(fmt.Sprintf("OUTPUT = [%v] %v", err, output))
	return time.Second * 5, err
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer:   &goserver.GoServer{},
		prints:     int64(0),
		pretend:    false, // Prod version doesn't pretend to print
		pretendret: nil,
		printq:     make(chan *pb.PrintRequest, 200),
		done:       make(chan bool),
		outOfPaper: false,
	}
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
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func (s *Server) readyToPrint(ctx context.Context) error {
	config, err := s.load(ctx)
	if err != nil {
		return err
	}

	Backlog.Set(float64(len(config.GetRequests())))
	for _, r := range config.GetRequests() {
		s.printq <- r
	}

	go s.printQueue()

	return nil
}

func (s *Server) drainQueue() {
	close(s.printq)
	<-s.done
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	err := server.RegisterServerV2("printer", false, true)
	if err != nil {
		return
	}

	//Silent crash is we can't print
	ctx, cancel := utils.ManualContext("printerstart", time.Minute)
	err = server.readyToPrint(ctx)
	if err != nil {
		server.Log(fmt.Sprintf("Not ready to print: %v", err))
		time.Sleep(time.Minute)
		return
	}
	cancel()

	server.Log(fmt.Sprintf("PRETEND PRINTING %v", server.pretend))
	fmt.Printf("%v", server.Serve())
}
