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
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/printer/proto"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	print bool
	count int
}

func (s *Server) localPrint(text string, lines []string, ti time.Time) error {
	if ti.Hour() < 9 || ti.Hour() > 16 {
		return fmt.Errorf("Not the time to print right now")
	}

	s.count++
	if s.print {
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
	return nil
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		true,
		0,
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

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "count", Value: int64(s.count)},
		&pbg.State{Key: "enabled", Text: fmt.Sprintf("%v", s.print)},
	}
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
	server.RegisterServer("printer", false)
	server.Log("Starting!")
	fmt.Printf("%v", server.Serve())
}
