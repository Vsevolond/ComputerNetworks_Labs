package main

import (
	"ComputerNetworks/CustomProtocol/src/proto"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"net"
	"strings"
)

// Client - состояние клиента.
type Client struct {
	logger  log.Logger    // Объект для печати логов
	conn    *net.TCPConn  // Объект TCP-соединения
	enc     *json.Encoder // Объект для кодирования и отправки сообщений
	str     string        // Текущая строка для поиска в ней
	substrs []string      // Массив подстрок, которые надо найти в текущей строке
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger:  log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:    conn,
		enc:     json.NewEncoder(conn),
		str:     "",
		substrs: make([]string, 0),
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "set":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var str string
			if err := json.Unmarshal(*req.Data, &str); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("performing addition", "value", str)
				client.str = str
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "find":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var substr string
			if err := json.Unmarshal(*req.Data, &substr); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("performing addition", "value", substr)
				client.substrs = append(client.substrs, substr)
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("addition failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "search":
		if len(client.substrs) == 0 {
			client.logger.Error("searching failed", "reason", "there are no substrings")
			client.respond("failed", "there are no substrings")
		} else {
			res := "\n#### searching in " + client.str + " ####\n"
			for _, substr := range client.substrs {
				positions := proto.Substring(client.str, substr)
				if len(positions) > 0 {
					res += substr + " on position: " + strings.Join(positions, ",") + "\n"
				} else {
					res += substr + " not exist\n"
				}
			}
			client.substrs = []string{}
			client.respond("result", &res)
		}
	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{status, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
}
