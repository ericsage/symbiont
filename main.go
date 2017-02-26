package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ericsage/symbiont/cx"
	"github.com/ericsage/symbiont/cyservice"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/grpc"
)

var (
	listeningAddress = getenv("LISTENING_ADDRESS", "0.0.0.0")
	listeningPort    = getenv("LISTENING_PORT", "80")
	serverAddress    = getenv("SERVICE_ADDRESS", "127.0.0.1")
	serverPort       = getenv("SERVICE_PORT", "8080")
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func streamNetwork(network io.ReadCloser, out http.ResponseWriter) {
	address := serverAddress + ":" + serverPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err == nil {
		client := cyservice.NewCyServiceClient(conn)
		stream, err := client.StreamFragments(context.Background())
		if err == nil {
			waitc := make(chan struct{})
			go func() {
				decoder := cx.NewDecoder()
				node := &cyservice.Node{}
				edge := &cyservice.Edge{}
				nodeAttr := &cyservice.NodeAttribute{}
				edgeAttr := &cyservice.EdgeAttribute{}
				networkAttr := &cyservice.NetworkAttribute{}
				decoder.RegisterAspectHandler("nodes", func(d *json.Decoder) {
					jsonpb.UnmarshalNext(d, node)
					frag := &cyservice.Fragment{Element: &cyservice.Fragment_Node{node}}
					stream.Send(frag)
				})
				decoder.RegisterAspectHandler("edges", func(d *json.Decoder) {
					jsonpb.UnmarshalNext(d, edge)
					frag := &cyservice.Fragment{Element: &cyservice.Fragment_Edge{edge}}
					stream.Send(frag)
				})
				decoder.RegisterAspectHandler("nodeAttributes", func(d *json.Decoder) {
					jsonpb.UnmarshalNext(d, nodeAttr)
					frag := &cyservice.Fragment{Element: &cyservice.Fragment_NodeAttribute{nodeAttr}}
					stream.Send(frag)
				})
				decoder.RegisterAspectHandler("edgeAttributes", func(d *json.Decoder) {
					jsonpb.UnmarshalNext(d, edgeAttr)
					frag := &cyservice.Fragment{Element: &cyservice.Fragment_EdgeAttribute{edgeAttr}}
					stream.Send(frag)
				})
				decoder.RegisterAspectHandler("networkAttributes", func(d *json.Decoder) {
					jsonpb.UnmarshalNext(d, networkAttr)
					frag := &cyservice.Fragment{Element: &cyservice.Fragment_NetworkAttribute{networkAttr}}
					stream.Send(frag)
				})
				decoder.Decode(network)
				fmt.Println("Done")
				stream.CloseSend()
				close(waitc)
			}()
			waitc2 := make(chan struct{})
			go func() {
				encoder := json.NewEncoder(out)
				for {
					in, err := stream.Recv()
					if err == io.EOF {
						close(waitc2)
						return
					} else {
					  encoder.Encode(in.Element)
				  }
				}
			}()
			<-waitc
			<-waitc2
		}
	}
	defer conn.Close()
}

func requestHandler(res http.ResponseWriter, req *http.Request) {
	streamNetwork(req.Body, res)
}

func main() {
	http.HandleFunc("/", requestHandler)
	address := listeningAddress + ":" + listeningPort
	log.Fatal(http.ListenAndServe(address, nil))
}
