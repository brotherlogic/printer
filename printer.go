package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
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
	backlog = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "printer_backlog",
		Help: "The size of the print queue",
	})
	goqueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "printer_queue",
		Help: "The size of the print queue",
	})
	printErrors = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "printer_errors",
		Help: "The size of the print queue",
	}, []string{"error"})
)

func (s *Server) metrics(config *pb.Config) {
	backlog.Set(float64(len(config.GetRequests())))
}

func (s *Server) load(ctx context.Context) (*pb.Config, error) {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return nil, err
	}

	s.metrics(data.(*pb.Config))
	return data.(*pb.Config), nil
}

func (s *Server) save(ctx context.Context, config *pb.Config) error {
	s.metrics(config)
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
	printlock  *sync.Mutex
}

func (s *Server) localPrint(ctx context.Context, text string, lines []string, ti time.Time, override bool) (time.Duration, error) {
	s.CtxLog(ctx, fmt.Sprintf("Trying to print %v/%v -> %v", text, lines, override))
	if s.pretend {
		if s.pretendret == nil {
			s.prints++
		}
		return time.Second, s.pretendret
	}

	if ti.Hour() < 7 || ti.Hour() > 19 { // || ((ti.Weekday() == time.Saturday || ti.Weekday() == time.Sunday) && (ti.Hour() != 10)) {
		return time.Minute, status.Errorf(codes.Unavailable, "Not the time to print right now")
	}

	//Only print if it's five to the hour
	if !override && ti.Minute() < 55 {
		return time.Minute, status.Errorf(codes.Unavailable, "Only print at five to the hour")
	}

	s.prints++

	if len(text) != 0 {
		ioutil.WriteFile("/home/simon/print.txt", []byte(text), 0644)
	} else {
		//os.Create("home/simon/print.txt")
		handle, err := os.Create("/home/simon/print.txt")
		if err != nil {
			return time.Second, err
		}
		for _, line := range lines {
			handle.WriteString(fmt.Sprintf("%v\n", line))
		}
		serr := handle.Sync()
		cerr := handle.Close()
		s.CtxLog(ctx, fmt.Sprintf("WHAT NOW LINES %v -> close errors are %v and %v", lines, serr, cerr))
	}

	cmd := exec.Command("lp", "/home/simon/print.txt")
	output := ""
	out, err := cmd.StdoutPipe()

	if err != nil {
		s.CtxLog(ctx, fmt.Sprintf("Error in the now resolved actual stdout: %v", err))
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

	s.CtxLog(ctx, fmt.Sprintf("OUTPUT = [%v] %v", err, output))
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
		printlock:  &sync.Mutex{},
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
	server.PrepServer("printer")
	server.Register = server
	err := server.RegisterServerV2(false)
	if err != nil {
		return
	}

	//Silent crash is we can't print
	ctx, cancel := utils.ManualContext("printerstart", time.Minute)
	err = server.readyToPrint(ctx)
	if err != nil {
		server.CtxLog(ctx, fmt.Sprintf("Not ready to print: %v", err))
		time.Sleep(time.Minute)
		return
	}

	server.CtxLog(ctx, fmt.Sprintf("PRETEND PRINTING %v", server.pretend))
	cancel()
	fmt.Printf("%v", server.Serve())
}
