package dispatcher

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	. "github.com/usetheplatform/ci-system/pkg/common"
)

// @TODO: https://stackoverflow.com/questions/54290124/sending-a-websocket-message-to-a-specific-client-in-go-using-gorilla
type Server struct {
	config     *Options
	dispatcher *Dispatcher
}

func NewServer(config *Options) Server {
	return Server{
		config: config,
	}
}

func (s *Server) Serve() {
	http.HandleFunc("/observer", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		CheckIfError(err)

		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Errorf("read failed: %v", err)
				break
			}
			input := string(message)
			cmd := getCmd(input)

			if cmd == Status {
				message := []byte("status:ok")
				err := conn.WriteMessage(mt, message)
				if err != nil {
					fmt.Printf("write failed: %v", err)
					break
				}
			}

			if cmd == Dispatch {
				// TODO: Parse commit from the payload
				s.dispatcher.Enqueue("13")
			}
		}
	})

	http.HandleFunc("/runners", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		CheckIfError(err)

		defer conn.Close()

		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Errorf("read failed: %v", err)
				break
			}
			input := string(message)
			cmd := getCmd(input)

			if cmd == AddRunner {
				// TODO: runner commit from the payload
				s.dispatcher.AddRunner(Runner{})
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}

type Command = int

const (
	Status Command = iota
	Dispatch
	AddRunner
	RemoveRunner
	UpdateRunnerStatus
	Results
)

// @TODO: AUTH https://stackoverflow.com/questions/55536439/how-can-i-upgrade-a-client-http-connection-to-websockets-in-golang-after-sending

func getCmd(input string) Command {
	inputArr := strings.Split(input, " ")
	cmd := inputArr[0]

	if cmd == "dispatch" {
		return Dispatch
	}

	return Status
}
