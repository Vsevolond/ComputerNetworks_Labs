package main

import (
	"ComputerNetworks/Peer_to_peer/proto"
	"bufio"
	"encoding/json"
	"fmt"
	_ "github.com/mgutz/logxi/v1"
	"net"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	Connections map[string]bool
	Posts       map[string]string
	Address     Address
}

type Address struct {
	IPv4 string
	Port string
}

type Package struct {
	From string
	Data string
}

// ./main :8080
func init() {
	if len(os.Args) != 2 {
		panic("len args != 2")
	}
}

func main() {
	NewNode(os.Args[1]).Run(handleServer, handleClient)
}

//ipv4:port
//127.0.0.1:8080
func NewNode(address string) *Node {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil
	}
	return &Node{
		Connections: make(map[string]bool),
		Posts:       make(map[string]string),
		Address: Address{
			IPv4: splited[0],
			Port: ":" + splited[1],
		},
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(node)
	handleClient(node)
}

func handleConnection(node *Node, conn net.Conn) {
	//defer conn.Close()
	/*var (
		buffer  = make([]byte, 512)
		message string
		pack    Package
		addr    string
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}*/
	//принимаем запрос и отправляем соглашение на подключение
	decoder := json.NewDecoder(conn)
	var req proto.Request
	if err := decoder.Decode(&req); err != nil {
		return
	} else {
		if !node.handleRequest(&req) {
			fmt.Println("ERROR: smth went wrong")
			os.Exit(0)
		}
	}
	//fmt.Println(pack.Data)
}

func handleServer(node *Node) {
	listen, err := net.Listen("tcp", "127.0.0.1"+node.Address.Port)
	if err != nil {
		fmt.Println("ERROR: this port has been occupied")
		os.Exit(0)
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		go handleConnection(node, conn)
	}
}

func handleClient(node *Node) {
	node.ConnectToAll()
	for {
		message := InputString()
		splited := strings.Split(message, " ")
		switch splited[0] {
		case "/exit":
			node.DisconnectFromAll()
			os.Exit(0)
		case "/network":
			node.PrintNetwork()
		case "/post":
			fmt.Print("#Post: ")
			post := InputString()
			node.AddPost(post)
		case "/takeoff":
			node.DeletePost()
		case "/board":
			node.PrintAllPosts()
		default:
			fmt.Println("WTF?")
		}
	}
}

func (node *Node) AddPost(post string) {
	if len(node.Posts[node.Address.Port]) == 0 {
		fmt.Println("OK: Post added")
	} else {
		fmt.Println("OK: Post updated")
	}
	node.Posts[node.Address.Port] = post
	for addr := range node.Connections {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		encoder := json.NewEncoder(conn)
		var new_pack = Package{
			From: node.Address.Port,
			Data: post,
		}
		node.send_request(encoder, "post", &new_pack)
		conn.Close()
	}
}

func (node *Node) DeletePost() {
	if len(node.Posts[node.Address.Port]) == 0 {
		fmt.Println("ERROR: You don't have any post")
		return
	}
	delete(node.Posts, node.Address.Port)
	fmt.Println("OK: Post deleted")
	for addr := range node.Connections {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			continue
		}
		encoder := json.NewEncoder(conn)
		var new_pack = Package{
			From: node.Address.Port,
			Data: "remove",
		}
		node.send_request(encoder, "takeoff", &new_pack)
		conn.Close()
	}
}

func (node *Node) PrintAllPosts() {
	f := false
	for addr, post := range node.Posts {
		if len(post) > 0 {
			fmt.Println("#Post from ", addr, " - ", post)
		}
		f = f || (len(post) > 0)
	}
	if !f {
		fmt.Println("THERE ARE NO POSTS")
	}
}

func (node *Node) PrintNetwork() {
	for addr := range node.Connections {
		fmt.Println("|", addr)
	}
}

func (node *Node) ConnectTo(address string) {
	if !node.Connections[address] {
		//fmt.Println("OK: connected to address ", address)
		node.Connections[address] = true
	} /*else {
		fmt.Println("OK: Already been connected to ", address)
	}*/
}

func (node *Node) ConnectToAll() {
	for i := 1000; i < 10000; i++ {
		addr := ":" + strconv.Itoa(i)
		conn, err := net.Dial("tcp", addr)
		if err != nil || addr == node.Address.Port {
			continue
		}

		//отправляем запрос на соединение и получаем ответ
		encoder := json.NewEncoder(conn)
		var new_pack = Package{
			From: node.Address.Port,
			Data: "hello",
		}
		node.send_request(encoder, "connect", new_pack)
		conn.Close()
	}
	fmt.Println("READY")
}

func (node *Node) DisconnectFromAll() {
	for addr := range node.Connections {
		conn, err := net.Dial("tcp", addr)
		if err != nil || addr == node.Address.Port {
			continue
		}
		encoder := json.NewEncoder(conn)
		var new_pack = Package{
			From: node.Address.Port,
			Data: "goodbye",
		}
		node.send_request(encoder, "disconnect", new_pack)
		conn.Close()
	}
}

func (node *Node) SendMessageToAll(message string) {
	var new_pack = Package{
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	for addr := range node.Connections {
		node.Send(&new_pack, addr)
	}
}

func (node *Node) Send(pack *Package, address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return
	}
	defer conn.Close()
	json_pack, _ := json.Marshal(pack)
	conn.Write(json_pack)
}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}

func (node *Node) send_request(enc *json.Encoder, command string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	enc.Encode(&proto.Request{command, &raw})
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (node *Node) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "connect":
		var pack Package
		if err := json.Unmarshal(*req.Data, &pack); err != nil {
			return false
		}
		conn, err := net.Dial("tcp", pack.From)
		if err != nil {
			return false
		}
		encoder := json.NewEncoder(conn)
		var new_pack = Package{
			From: node.Address.Port,
			Data: node.Posts[node.Address.Port],
		}
		node.send_request(encoder, "agree", &new_pack)
		node.ConnectTo(pack.From)
		conn.Close()
		return true
	case "agree":
		var pack Package
		if err := json.Unmarshal(*req.Data, &pack); err != nil {
			return false
		}
		node.Posts[pack.From] = pack.Data
		node.ConnectTo(pack.From)
		return true
	case "disconnect":
		var pack Package
		if err := json.Unmarshal(*req.Data, &pack); err != nil {
			return false
		}
		delete(node.Connections, pack.From)
		delete(node.Posts, pack.From)
		return true
	case "post":
		var pack Package
		if err := json.Unmarshal(*req.Data, &pack); err != nil {
			return false
		}
		node.Posts[pack.From] = pack.Data
		return true
	case "takeoff":
		var pack Package
		if err := json.Unmarshal(*req.Data, &pack); err != nil {
			return false
		}
		delete(node.Posts, pack.From)
		return true
	}
	return false
}
