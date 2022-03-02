package cluster

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

// NodeActions define what a node can do
type NodeActions interface {
	getId() string
	getStatus() string
	sendMessage(msg string, receiverAddress string)
	receiveMessage()
	PrintInfo()
	Run()
}

// Node is a node of the cluster
type Node struct {
	NodeID         string
	NodeStatus     int8
	NodeAddress    string
	PeersID        []string
	PeersAddresses []string
}

// NewNode creates a new Node object
func NewNode(address string) Node {
	id := uuid.New()
	n := Node{id.String(), 0, address, []string{}, []string{}}
	return n
}

func (n Node) sendMessage(msg string, receiverAddress string) {
	// use ResolveTCPAddr to create address to connect to
	raddr, err := net.ResolveTCPAddr("tcp", receiverAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use DialTCP to create a connection to the
	// remote address. Note that there is no need
	// to specify the local address.
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// send text to server
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read response
	buf := make([]byte, 1024)
	readResp, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed reading response:", err)
		os.Exit(1)
	}
	fmt.Println(string(buf[:readResp]))
}

func (n Node) receiveMessage() {
	listener := n.createLocalListener()
	defer listener.Close()

	fmt.Println("listening at (tcp)", n.NodeAddress)

	// req/response loop
	for {
		// use TCPListener to block and wait for TCP
		// connection request using AcceptTCP which creates a TCPConn
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}
		fmt.Println("connected to: ", conn.RemoteAddr())

		go handleConnection(n, conn)
	}
}

func (n Node) createLocalListener() *net.TCPListener {
	// create local addr for socket
	laddr, err := net.ResolveTCPAddr("tcp", n.NodeAddress)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// announce service using ListenTCP
	// which a TCPListener.
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return l
}

// handleConnection reads request from connection
// with conn.Read() then write response using
// conn.Write().  Then the connection is closed.
func handleConnection(node Node, conn *net.TCPConn) {
	defer conn.Close() // clean up when done

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// echo buffer
	w, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("failed to write to client:", err)
		return
	}
	if w != n { // was all data sent
		fmt.Println("warning: not all data sent to client")
		return
	}

	msg := string(buf)

	// Accepts raw messages from tcp
	if strings.Contains(msg, "|") {
		cmd := strings.Split(msg, "|")
		fmt.Printf("command: %v\n", cmd)
		switch cmd[0] {
		case "SEND":
			node.sendMessage(cmd[2], cmd[1])
		default:
			fmt.Println(msg)
		}
	} else {
		fmt.Printf("Printing Anyway: %v\n", msg)
	}
}

// PrintInfo prints information about the node.
func (n Node) PrintInfo() {
	fmt.Printf("Node {%v}\n", n.NodeID)
	fmt.Printf("running at %v\n", n.NodeAddress)
	fmt.Printf("has following peers: %v\n", n.PeersAddresses)
}

// Run starts the node process
func (n Node) Run() {
	n.receiveMessage()
}
