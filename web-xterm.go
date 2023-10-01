package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"github.com/olahol/melody"
)

//go:embed index.html node_modules/xterm/css/xterm.css node_modules/xterm/lib/xterm.js
var content embed.FS

type ConnectMessage struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
}

func connectSSH(username, ip string) (*os.File, error) {
	for {
		sshCommand := exec.Command("ssh", username+"@"+ip)
		f, err := pty.Start(sshCommand)
		if err != nil {
			log.Printf("SSH connection error: %v. Retrying...", err)
			time.Sleep(5 * time.Second) // 等待一段时间后尝试重新连接
			continue
		}
		return f, nil
	}
}

func main() {
	m := melody.New()
	m.Config.MaxMessageSize = 10000000
	connected := false
	var f *os.File

	m.HandleConnect(func(s *melody.Session) {
		connected = false
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if !connected {
			var connectMsg ConnectMessage
			err := json.Unmarshal(msg, &connectMsg)
			if err != nil {
				log.Printf("Failed to parse connect message: %v", err)
				return
			}

			f, err = connectSSH(connectMsg.Username, connectMsg.IP)
			if err != nil {
				log.Printf("Failed to connect SSH: %v", err)
				return
			}

			go func() {
				for {
					buf := make([]byte, 1024*1024*5)
					read, err := f.Read(buf)
					if err != nil {
						log.Printf("SSH connection closed: %v. Reconnecting...", err)
						f, err = connectSSH(connectMsg.Username, connectMsg.IP) // 连接中断时尝试重新连接
						if err != nil {
							log.Printf("SSH reconnection failed: %v. Retrying...", err)
							time.Sleep(5 * time.Second) // 等待一段时间后尝试重新连接
							continue
						}
						log.Println("SSH reconnected.")
						connected = true
						continue
					}
					m.Broadcast(buf[:read]) // 向所有客户端广播消息
				}
			}()

			connected = true
		} else {
			f.Write(msg)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	fs := http.FileServer(http.FS(content))
	http.Handle("/", http.StripPrefix("/", fs))

	err := http.ListenAndServe("0.0.0.0:22333", nil)
	if err != nil {
		log.Fatal(err)
	}
}
