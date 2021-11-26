package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/linedejgaard/DISYS_MINIPROJECT3/Auction"
	"google.golang.org/grpc"
)

type Node struct {
	Auction.UnimplementedAuctionServiceServer
}

var highestBidder string //leader
var highestBid int32     //leader

var isAuction bool //node

var leaderPort string //node
var ports []string    //node
var port string       //node

func main() {

	file, err := os.OpenFile("../info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)

	defer recoverError()
	// write your port
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please write your port")
	input := GetInputFromTerminal(reader)

	var n Node
	port = input
	ports = make([]string, 0)

	// Create listener tcp on port *input*
	go n.startListen(port)

	leaderPort = "1111"

	n.sendJoinRequest()
	for {
		QueryUserForInputViaTerminal(reader, n)
	}

}

func QueryUserForInputViaTerminal(reader *bufio.Reader, n Node) {
	input := GetInputFromTerminal(reader)

	if strings.Compare("/bid", input) == 0 {
		if isAuction {
			n.sendGetStateRequest()
			fmt.Println("Give your bid:")
			input = GetInputFromTerminal(reader)
			n.sendBidRequest(input)
		} else {
			fmt.Println("No current auction, you can make an auction by using the /new command")
		}

	} else if strings.Compare("/new", input) == 0 {
		fmt.Println("Give your starting price:")
		input := GetInputFromTerminal(reader)
		n.sendMakeNewAuctionRequest(input)
	} else if strings.Compare("/l", input) == 0 {
		n.sendLeaveRequest()
	} else if strings.Compare("/result", input) == 0 {
		if isAuction {
			n.sendGetStateRequest()
		} else {
			n.sendPublishResultRequest()
		}

	} else {
		fmt.Println("Invalid message")
	}
}

func GetInputFromTerminal(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')

	// trim input
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Replace(input, "\r", "", -1)
	return input
}

func (n *Node) startListen(port string) {
	//SERVER
	portString := ":" + port
	// Create listener tcp on portString
	list, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", port, err)
	}
	grpcServer := grpc.NewServer()
	Auction.RegisterAuctionServiceServer(grpcServer, &Node{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

//JOIN
func (n *Node) sendJoinRequest() {
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect: %s\n", err)
		fmt.Println("Please type new leader:")
		reader := bufio.NewReader(os.Stdin)
		input := GetInputFromTerminal(reader)
		leaderPort = input
		n.sendJoinRequest()

	} else {

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()
		//  Create new Client from generated gRPC code from proto
		c := Auction.NewAuctionServiceClient(conn)

		// Send leave request
		message := Auction.JoinRequest{
			Port: port,
		}

		response, err := c.Join(context.Background(), &message)
		if err != nil {
			log.Printf("Error when calling join: %s\n", err)
			fmt.Println("Please type new leader:")
			reader := bufio.NewReader(os.Stdin)
			input := GetInputFromTerminal(reader)
			leaderPort = input
			n.sendJoinRequest()
		} else {
			ports = strings.Split(response.Ports, " ")
			isAuction = response.IsAuction
			highestBidder = response.HighestBidder
			highestBid = response.HighestBid

			log.Printf("Join response: %s\n", response.Reply)
		}

	}
}

func (n *Node) Join(_ context.Context, in *Auction.JoinRequest) (*Auction.JoinReply, error) {
	portsString := ""
	if len(ports) == 0 {
		portsString = in.Port
	} else {
		portsString = sliceToString(ports) + " " + in.Port
	}

	ports = append(ports, in.Port)

	n.sendUpdatePortsRequest(portsString)

	return &Auction.JoinReply{
		Ports:         portsString,
		IsAuction:     isAuction,
		HighestBidder: highestBidder,
		HighestBid:    highestBid,
		Reply:         "Join succeeded",
	}, nil
}

func (n *Node) sendUpdatePortsRequest(portsString string) { //called on leader
	//fmt.Printf("---> SEND UPDATE PORTS REQUEST PORTSTRING:%v\n", portsString)

	for _, p := range stringToSlice(portsString) {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			ports = removeByPort(stringToSlice(portsString), p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Printf("Could not connect: %s", err)
		} else {

			// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
			defer conn.Close()
			//  Create new Client from generated gRPC code from proto
			c := Auction.NewAuctionServiceClient(conn)

			// Send leader request
			message := Auction.UpdatePortsRequest{
				Ports: portsString,
			}

			response, err := c.UpdatePorts(context.Background(), &message)
			if err != nil {
				ports = removeByPort(stringToSlice(portsString), p)
				n.sendUpdatePortsRequest(sliceToString(ports))
				log.Printf("Error when calling send update ports request: %s\n", err)

			} else {
				log.Printf("update ports request response: %s\n", response.Reply)
			}
		}

	}

}

func sliceToString(slice []string) (string string) {
	portsstring := ""

	for i, p := range slice {
		portsstring = portsstring + p
		if i != len(slice)-1 {
			portsstring = portsstring + " "
		}
	}

	return portsstring

}

func stringToSlice(string string) (slice []string) {
	return strings.Split(string, " ")
}

func (n *Node) UpdatePorts(_ context.Context, in *Auction.UpdatePortsRequest) (*Auction.UpdatePortsReply, error) {
	ports = strings.Split(in.Ports, " ")

	return &Auction.UpdatePortsReply{
		Reply: "OK",
	}, nil
}

//GET STATE
func (n *Node) sendGetStateRequest() { //to leader
	//fmt.Println("SEND GET STATE REQUEST")
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendGetStateRequest()
		log.Printf("Could not connect: %s\n", err)
	} else {

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()
		//  Create new Client from generated gRPC code from proto
		c := Auction.NewAuctionServiceClient(conn)
		// Send get state request
		message := Auction.GetStateRequest{
			Port: port,
		}

		response, err := c.GetState(context.Background(), &message)
		if err != nil {
			n.startElection()
			n.sendGetStateRequest()
			log.Printf("Error when calling Get state: %s\n", err)
		} else {
			if response.State == 0 {
				fmt.Println("There is no current auction")
			} else {
				fmt.Printf("Current highest bid is: %v\n", response.State)
			}

			log.Printf("Get State response: %s\n", response.Reply)
		}

	}
}

func (n *Node) GetState(context.Context, *Auction.GetStateRequest) (*Auction.GetStateReply, error) {
	//fmt.Println("GET STATE")
	return &Auction.GetStateReply{
		State: highestBid,
		Reply: "Get State succeded",
	}, nil
}

//BID
func (n *Node) sendBidRequest(input string) { //to leader
	//fmt.Println("SEND BID REQUEST")
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendBidRequest(input)
		log.Printf("Could not connect: %s\n", err)
	} else {

		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()
		//  Create new Client from generated gRPC code from proto
		c := Auction.NewAuctionServiceClient(conn)

		amount, err := strconv.Atoi(input)
		if err != nil {
			// handle error
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Your input is not a valid number, please try again")
			input := GetInputFromTerminal(reader)
			n.sendBidRequest(input)
		}

		message := Auction.BidRequest{
			Bidder: port,
			Amount: int32(amount),
		}

		response, err := c.Bid(context.Background(), &message)

		if err != nil {
			n.startElection()
			n.sendBidRequest(input)
			log.Printf("Error when calling bid: %s\n", err)
		} else {
			log.Printf("Bid response: %s\n", response.Reply)
		}

	}
}

func (n *Node) Bid(_ context.Context, in *Auction.BidRequest) (*Auction.BidReply, error) {
	//fmt.Println("BID")
	succeeded := in.Amount > highestBid
	reply := ""

	if succeeded {
		highestBid = in.Amount
		highestBidder = in.Bidder
		reply = "Your bid is now the highest"
	} else {
		reply = "Another bid was higher than yours"
	}
	n.sendUpdateAuctionStatusRequest()

	return &Auction.BidReply{
		Succeeded: succeeded,
		Reply:     reply,
	}, nil
}

//LEAVE
func (n *Node) sendLeaveRequest() { //to leader
	//fmt.Println("SEND LEAVE REQUEST")
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendLeaveRequest()
		log.Printf("Could not connect: %s\n", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	//  Create new Client from generated gRPC code from proto
	c := Auction.NewAuctionServiceClient(conn)

	// Send leave request
	message := Auction.LeaveRequest{
		Port: port,
	}

	response, err := c.Leave(context.Background(), &message)
	if err != nil {
		n.startElection()
		n.sendLeaveRequest()
		log.Printf("Error when calling leave: %s\n", err)
	} else {
		log.Printf("Leave response: %s\n", response.Reply)
		os.Exit(0)
	}

}

func (n *Node) Leave(_ context.Context, in *Auction.LeaveRequest) (*Auction.LeaveReply, error) {
	//fmt.Println("LEAVE")
	if contains(ports, port) {

		ports = removeByPort(ports, in.Port)

		reply := port + " left the auction house"
		n.sendUpdatePortsRequest(sliceToString(ports))
		return &Auction.LeaveReply{
			Reply: reply,
		}, nil
	} else {
		log.Println("Unknown leave reply")
		return &Auction.LeaveReply{
			Reply: "Unknown user tried to leave the server",
		}, nil
	}
}

func (n *Node) sendPublishResultRequest() { //from leader
	//fmt.Println("SEND PUBLISH RESULT REQUEST")
	canConnect := true
	for _, p := range ports {
		if canConnect {
			// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
			var conn *grpc.ClientConn
			portString := ":" + p
			conn, err := grpc.Dial(portString, grpc.WithInsecure())
			if err != nil {
				ports = removeByPort(ports, p)
				n.sendUpdatePortsRequest(sliceToString(ports))
				log.Printf("Could not connect: %s\n", err)
			}

			// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
			defer conn.Close()
			//  Create new Client from generated gRPC code from proto
			c := Auction.NewAuctionServiceClient(conn)

			// Send leader request
			message := Auction.PublishResultRequest{
				Highestbid: highestBid,
				Bidder:     highestBidder,
			}

			response, err := c.PublishResult(context.Background(), &message)
			if err != nil {
				canConnect = false
				ports = removeByPort(ports, p)
				n.sendUpdatePortsRequest(sliceToString(ports))
				n.sendPublishResultRequest()
				log.Printf("Error when calling publish result ports request: %s\n", err)
			} else {
				log.Printf("publish result request response: %s\n", response.Reply)
			}

		}

	}
}

func (n *Node) PublishResult(context.Context, *Auction.PublishResultRequest) (*Auction.PublishResultReply, error) {
	//fmt.Println("PUBLISH RESULT")
	fmt.Printf("Auction finished!\nWinner: %v\nFinal price:%v\n", highestBidder, highestBid)

	isAuction = false

	return &Auction.PublishResultReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendMakeNewAuctionRequest(input string) { //to leader
	//fmt.Println("SEND MAKE NEW AUCTION REQUEST")
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendMakeNewAuctionRequest(input)
		log.Printf("Could not connect: %s\n", err)
	}

	amount := n.GetBiddingAmountFromUser(input)

	response := n.SendMakeNewAuctionRequestReceiveReply(conn, amount, input)

	log.Printf("Make New Auction response: %s\n", response.Reply)

}

func (n *Node) SendMakeNewAuctionRequestReceiveReply(conn *grpc.ClientConn, amount int, input string) *Auction.MakeNewAuctionReply {
	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	//  Create new Client from generated gRPC code from proto
	c := Auction.NewAuctionServiceClient(conn)

	// Send make new auction request
	message := Auction.MakeNewAuctionRequest{
		StartAmount: int32(amount),
		Port:        port,
	}

	response, err := c.MakeNewAuction(context.Background(), &message)
	if err != nil {
		log.Printf("Error when calling Make New Auction: %s\n", err)
		n.startElection()
		n.sendMakeNewAuctionRequest(input)

	} else {
		return response
	}
	return &Auction.MakeNewAuctionReply{
		Reply: "NOT SUCCEEDED",
	}

}

func (n *Node) GetBiddingAmountFromUser(input string) int {
	amount, err := strconv.Atoi(input)
	if err != nil {
		// handle error
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Your input is not a valid number, please try again")
		input := GetInputFromTerminal(reader)
		n.sendBidRequest(input)
	}
	return amount
}

func (n *Node) MakeNewAuction(_ context.Context, in *Auction.MakeNewAuctionRequest) (*Auction.MakeNewAuctionReply, error) {
	//fmt.Println("MAKE NEW AUCTION")
	reply := ""

	if !isAuction {
		reply = "Succeeded"
		isAuction = true
		highestBid = in.StartAmount
		highestBidder = in.Port
		go n.setTimer()

	} else {
		reply = "There is already an auction, please wait until it is finished"
	}

	n.sendUpdateAuctionStatusRequest()

	return &Auction.MakeNewAuctionReply{
		Reply: reply,
	}, nil
}

//ELECTION!!
func (n *Node) startElection() {
	//fmt.Println("START ELECTION")
	ports = removeByPort(ports, leaderPort)
	n.sendUpdatePortsRequest(sliceToString(ports))
	leaderPort = ""
	sent := false
	recievedResponse := false
	for _, p := range ports {
		if p > port {
			sent = true
			// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
			var conn *grpc.ClientConn
			portString := ":" + p
			conn, err := grpc.Dial(portString, grpc.WithInsecure())
			if err != nil {
				ports = removeByPort(ports, p)
				n.sendUpdatePortsRequest(sliceToString(ports))
				log.Printf("Could not connect: %s\n", err)
			}

			// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
			defer conn.Close()
			//  Create new Client from generated gRPC code from proto
			c := Auction.NewAuctionServiceClient(conn)

			// Send election request
			if sendElectionRequest(c) {
				recievedResponse = true
			}
		}
	}
	if !sent || !recievedResponse {
		n.sendLeaderRequest()
	}
}

func sendElectionRequest(c Auction.AuctionServiceClient) bool {
	//fmt.Println("SEND ELECTION REQUEST")
	message := Auction.ElectionRequest{
		Message: "Election",
	}
	response, err := c.Election(context.Background(), &message)
	if err != nil {
		log.Printf("Error when calling Elction: %s\n", err)
		return false
	}
	if response == nil {
		log.Printf("Response was nil")
		return false
	}

	log.Println("Election reply: ", response.Reply)
	return true
}

func (n *Node) Election(context.Context, *Auction.ElectionRequest) (*Auction.ElectionReply, error) {
	//fmt.Println("ELECTION")
	n.startElection()
	return &Auction.ElectionReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendLeaderRequest() {
	//fmt.Println("SEND LEADER REQUEST")
	for _, p := range ports {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Printf("Could not connect: %s\n", err)
		}

		n.SendLeaderRequestReceiveReply(conn)

	}
	if isAuction {
		go n.setTimer()
	}

}

func (n *Node) SendLeaderRequestReceiveReply(conn *grpc.ClientConn) {
	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	//  Create new Client from generated gRPC code from proto
	c := Auction.NewAuctionServiceClient(conn)

	// Send leader request
	message := Auction.LeaderDeclarationRequest{
		Port: port,
	}

	response, err := c.LeaderDeclaration(context.Background(), &message)
	if err != nil {
		log.Printf("Error when calling LeaderDeclaration: %s\n", err)
	}
	log.Printf("Leader Declaration response: %v", response.Reply)
}

func (n *Node) LeaderDeclaration(_ context.Context, in *Auction.LeaderDeclarationRequest) (*Auction.LeaderDeclarationReply, error) {
	//fmt.Println("LEADER DECLARATION")
	leaderPort = in.Port

	fmt.Printf("New leader: %v\n", in.Port)

	return &Auction.LeaderDeclarationReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendUpdateAuctionStatusRequest() {
	canConnect := true
	//fmt.Printf("SEND UPDATE AUCTION STATUS REQUEST, PORT: %v, PORTS: %v\n", port, ports)
	for _, p := range ports {
		if canConnect {

			// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
			var conn *grpc.ClientConn
			portString := ":" + p
			conn, err := grpc.Dial(portString, grpc.WithInsecure())
			if err != nil {
				canConnect = false
				ports = removeByPort(ports, p)
				n.sendUpdatePortsRequest(sliceToString(ports))
				log.Printf("Could not connect: %s\n", err)
			} else {
				response := n.SendUpdateAuctionStatusRequestReceiveReply(conn, p)

				log.Printf("update auction request response: %s\n", response.Reply)
			}

		}

	}
}

func recoverError() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func (n *Node) SendUpdateAuctionStatusRequestReceiveReply(conn *grpc.ClientConn, p string) *Auction.UpdateAuctionStatusReply {

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()
	//  Create new Client from generated gRPC code from proto
	c := Auction.NewAuctionServiceClient(conn)

	// Send leader request
	message := Auction.UpdateAuctionStatusRequest{
		Bid:       highestBid,
		Bidder:    highestBidder,
		IsAuction: isAuction,
	}

	response, err := c.UpdateAuctionStatus(context.Background(), &message) //THIS IS THE ERROR --> RESPONSE MÅ VÆRE EN FEJL
	if err != nil {
		ports = removeByPort(ports, p)
		n.sendUpdatePortsRequest(sliceToString(ports))
		n.sendUpdateAuctionStatusRequest()
		log.Printf("Error when calling update auction status request: %s\n", err)
	} else {
		return response
	}

	return &Auction.UpdateAuctionStatusReply{
		Reply: "NOT OK",
	}
}

func (n *Node) UpdateAuctionStatus(_ context.Context, in *Auction.UpdateAuctionStatusRequest) (*Auction.UpdateAuctionStatusReply, error) {
	//fmt.Println("UPDATE AUCTION STATUS")
	highestBid = in.Bid
	highestBidder = in.Bidder
	isAuction = in.IsAuction

	return &Auction.UpdateAuctionStatusReply{
		Reply: "OK",
	}, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeByPort(ports []string, port string) []string {

	newPorts := make([]string, 0)

	for _, p := range ports {
		if p != port {
			newPorts = append(newPorts, p)
		}
	}
	return newPorts
}

func (n *Node) setTimer() {
	time.Sleep(15 * time.Second)
	n.sendUpdateAuctionStatusRequest()
	n.sendPublishResultRequest()
}
