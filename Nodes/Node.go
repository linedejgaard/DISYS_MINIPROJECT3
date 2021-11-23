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

	// write your port
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please write your port")
	input, _ := reader.ReadString('\n')

	// trim input
	input = strings.Replace(input, "\n", "", -1)
	input = strings.Replace(input, "\r", "", -1)

	var n Node
	port = input
	ports = make([]string, 0)

	// Create listener tcp on port *input*
	go n.startListen(port)

	leaderPort = "1111"

	n.sendJoinRequest()

	for {

		input, _ := reader.ReadString('\n')
		// trim input
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)

		if strings.Compare("/b", input) == 0 {
			if isAuction {
				n.sendGetStateRequest()
				fmt.Println("Give your bid:")
				input, _ := reader.ReadString('\n')
				// trim input
				input = strings.Replace(input, "\n", "", -1)
				input = strings.Replace(input, "\r", "", -1)
				n.sendBidRequest(input)
			} else {
				fmt.Println("No current auction, you can make an auction by using the /n command")
			}

		} else if strings.Compare("/n", input) == 0 {
			fmt.Println("Give your starting price:")
			input, _ := reader.ReadString('\n')
			// trim input
			input = strings.Replace(input, "\n", "", -1)
			input = strings.Replace(input, "\r", "", -1)
			n.sendMakeNewAuctionRequest(input)
		} else if strings.Compare("/l", input) == 0 {
			n.sendLeaveRequest()
		} else {
			fmt.Println("Invalid message")
		}
	}

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
		log.Fatalf("Could not connect: %s", err)
		fmt.Println("Please type new leader:")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		// trim input
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)
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
			log.Fatalf("Error when calling join: %s", err)
		}

		ports = strings.Split(response.Ports, " ")

		log.Printf("Join response: %s\n", response.Reply)
	}
}

func (n *Node) Join(ctx context.Context, in *Auction.JoinRequest) (*Auction.JoinReply, error) {
	portsString := ""
	if len(ports) == 0 {
		portsString = in.Port
	} else {
		portsString = sliceToString(ports) + " " + in.Port
	}

	ports = append(ports, in.Port)

	n.sendUpdatePortsRequest(portsString)

	return &Auction.JoinReply{
		Ports: portsString,
		Reply: "Join succeeded",
	}, nil
}

func (n *Node) sendUpdatePortsRequest(portsString string) { //called on leader
	fmt.Println(portsString)
	for _, p := range stringToSlice(portsString) {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			ports = removeByPort(stringToSlice(portsString), p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Fatalf("Could not connect: %s", err)
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
				log.Fatalf("Error when calling send update ports request: %s", err)

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

func (n *Node) UpdatePorts(ctx context.Context, in *Auction.UpdatePortsRequest) (*Auction.UpdatePortsReply, error) {
	ports = strings.Split(in.Ports, " ")

	return &Auction.UpdatePortsReply{
		Reply: "OK",
	}, nil
}

//GET STATE
func (n *Node) sendGetStateRequest() { //to leader
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {

		n.startElection()
		n.sendGetStateRequest()
		log.Fatalf("Could not connect: %s", err)
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
			//n.startElection()
			//n.sendGetStateRequest()
			log.Fatalf("Error when calling Get state: %s", err)
		}

		if response.State == 0 {
			fmt.Println("There is no current auction")
		} else {
			fmt.Printf("Current highest bid is: %v\n", response.State)
		}

		log.Printf("Get State response: %s\n", response.Reply)
	}
}

func (n *Node) GetState(context.Context, *Auction.GetStateRequest) (*Auction.GetStateReply, error) {

	return &Auction.GetStateReply{
		State: highestBid,
		Reply: "Get State succeded",
	}, nil
}

//BID
func (n *Node) sendBidRequest(input string) { //to leader
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendBidRequest(input)
		log.Fatalf("Could not connect: %s", err)
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
			input, _ := reader.ReadString('\n')
			// trim input
			input = strings.Replace(input, "\n", "", -1)
			input = strings.Replace(input, "\r", "", -1)
			n.sendBidRequest(input)
		}

		message := Auction.BidRequest{
			Bidder: port,
			Amount: int32(amount),
		}

		response, err := c.Bid(context.Background(), &message)

		if err != nil {
			//n.startElection()
			//n.sendBidRequest(input)
			log.Fatalf("Error when calling bid: %s", err)
		}

		log.Printf("Bid response: %s\n", response.Reply)
	}
}

func (n *Node) Bid(ctx context.Context, in *Auction.BidRequest) (*Auction.BidReply, error) {
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
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {
		n.startElection()
		n.sendLeaveRequest()
		log.Fatalf("Could not connect: %s", err)
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
		//n.startElection()
		//n.sendLeaveRequest()
		log.Fatalf("Error when calling leave: %s", err)
	}

	log.Printf("Leave response: %s\n", response.Reply)
}

func (n *Node) Leave(ctx context.Context, in *Auction.LeaveRequest) (*Auction.LeaveReply, error) {

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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeByPort(ports []string, port string) []string {
	fmt.Printf("ports before: %v\n", ports)

	newPorts := make([]string, 0)

	for _, p := range ports {
		if p != port {
			newPorts = append(newPorts, p)
		}
	}

	fmt.Printf("ports after: %v\n", newPorts)

	return newPorts
}

func (n *Node) sendPublishResultRequest() { //from leader
	for _, p := range ports {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			ports = removeByPort(ports, p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Fatalf("Could not connect: %s", err)
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
			ports = removeByPort(ports, p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Fatalf("Error when calling publish result ports request: %s", err)
		}

		log.Printf("publish result request response: %s\n", response.Reply)
	}
}

func (n *Node) PublishResult(ctx context.Context, in *Auction.PublishResultRequest) (*Auction.PublishResultReply, error) {

	fmt.Printf("Auction finished!\nWinner: %v\nFinal price:%v\n", highestBidder, highestBid)

	isAuction = false
	highestBidder = ""
	highestBid = 0

	return &Auction.PublishResultReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendMakeNewAuctionRequest(input string) { //to leader
	// Creat a virtual RPC Client Connection on leaderPort
	var conn *grpc.ClientConn
	portString := ":" + leaderPort
	conn, err := grpc.Dial(portString, grpc.WithInsecure())
	if err != nil {

		n.startElection()
		n.sendMakeNewAuctionRequest(input)
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	//  Create new Client from generated gRPC code from proto
	c := Auction.NewAuctionServiceClient(conn)

	amount, err := strconv.Atoi(input)

	if err != nil {

		// handle error
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Your input is not a valid number, please try again")
		input, _ := reader.ReadString('\n')
		// trim input
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Replace(input, "\r", "", -1)
		n.sendBidRequest(input)
	}

	// Send make new auction request
	message := Auction.MakeNewAuctionRequest{
		StartAmount: int32(amount),
		Port:        port,
	}

	response, err := c.MakeNewAuction(context.Background(), &message)
	if err != nil {
		//n.startElection()
		//n.sendMakeNewAuctionRequest(input)
		log.Fatalf("Error when calling Make New Auction: %s", err)

	}

	log.Printf("Make New Auction response: %s\n", response.Reply)

}

func (n *Node) MakeNewAuction(ctx context.Context, in *Auction.MakeNewAuctionRequest) (*Auction.MakeNewAuctionReply, error) {
	fmt.Println("HER 1 <---")
	reply := ""

	if !isAuction {
		fmt.Println("HER NOT 1 <---")
		reply = "Succeeded"
		isAuction = true
		highestBid = in.StartAmount
		highestBidder = in.Port
		fmt.Println("HER NOT 2 <---")
		go n.setTimer()
		fmt.Println("HER NOT 3 <---")

	} else {
		fmt.Println("HER ELSE 1 <---")
		reply = "There is already an auction, please wait until it is finished"
	}

	fmt.Println("HER 2 <---")

	n.sendUpdateAuctionStatusRequest()
	fmt.Println("HER 3 <---")

	return &Auction.MakeNewAuctionReply{
		Reply: reply,
	}, nil
}

//ELECTION!!
func (n *Node) startElection() {
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
				log.Fatalf("Could not connect: %s", err)
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
	message := Auction.ElectionRequest{
		Message: "Election",
	}
	response, err := c.Election(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling Elction: %s", err)
		return false
	}
	if response == nil {
		log.Printf("Response was nil")
		return false
	}

	log.Println("Election reply: ", response.Reply)
	return true
}

func (n *Node) Election(ctx context.Context, in *Auction.ElectionRequest) (*Auction.ElectionReply, error) {
	n.startElection()
	return &Auction.ElectionReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendLeaderRequest() {
	for _, p := range ports {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}

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
			log.Fatalf("Error when calling LeaderDeclaration: %s", err)
		}
		log.Printf("Leader Declaration response: %v", response.Reply)

	}
	if isAuction {
		go n.setTimer()
	}

}

func (n *Node) LeaderDeclaration(ctx context.Context, in *Auction.LeaderDeclarationRequest) (*Auction.LeaderDeclarationReply, error) {
	leaderPort = in.Port

	fmt.Printf("New leader: %v\n", port)

	return &Auction.LeaderDeclarationReply{
		Reply: "OK",
	}, nil
}

func (n *Node) sendUpdateAuctionStatusRequest() {
	for _, p := range ports {

		// Creat a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		portString := ":" + p
		conn, err := grpc.Dial(portString, grpc.WithInsecure())
		if err != nil {
			ports = removeByPort(ports, p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Fatalf("Could not connect: %s", err)
		}

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

		response, err := c.UpdateActionStatus(context.Background(), &message)
		if err != nil {
			ports = removeByPort(ports, p)
			n.sendUpdatePortsRequest(sliceToString(ports))
			log.Printf("Error when calling update auction ports request: %s", err)
		}

		log.Printf("update auction request response: %s\n", response.Reply)
	}
}

func (n *Node) UpdateActionStatus(ctx context.Context, in *Auction.UpdateAuctionStatusRequest) (*Auction.UpdateActionStatusReply, error) {
	highestBid = in.Bid
	highestBidder = in.Bidder
	isAuction = in.IsAuction

	return &Auction.UpdateActionStatusReply{
		Reply: "OK",
	}, nil
}

func (n *Node) setTimer() {
	time.Sleep(15 * time.Second)
	n.sendUpdateAuctionStatusRequest()
	n.sendPublishResultRequest()
}
