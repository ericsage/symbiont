package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ericsage/symbiont/cx"
	"github.com/ericsage/symbiont/cxpb"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"sync"
)

var (
	requiredAspects = []string{
		"edges",
		"nodes",
		"nodeAttributes",
		"edgeAttributes",
		"networkAttributes",
	}
	knownAspects = map[string]proto.Message{
		"nodes":             &cxpb.Node{},
		"edges":             &cxpb.Edge{},
		"nodeAttributes":    &cxpb.NodeAttribute{},
		"edgeAttributes":    &cxpb.EdgeAttribute{},
		"networkAttributes": &cxpb.NetworkAttribute{},
	}
)

func wrapElement(element proto.Message) *cxpb.Fragment {
	switch e := element.(type) {
	case *cxpb.Node:
		fmt.Printf("Wrapping a node element %+v\n", element)
		return &cxpb.Fragment{&cxpb.Fragment_Node{e}}
	case *cxpb.Edge:
		fmt.Printf("Wrapping an edge element %+v\n", element)
		return &cxpb.Fragment{&cxpb.Fragment_Edge{e}}
	case *cxpb.NodeAttribute:
		return &cxpb.Fragment{&cxpb.Fragment_NodeAttribute{e}}
	case *cxpb.EdgeAttribute:
		return &cxpb.Fragment{&cxpb.Fragment_EdgeAttribute{e}}
	case *cxpb.NetworkAttribute:
		return &cxpb.Fragment{&cxpb.Fragment_NetworkAttribute{e}}
	default:
		fmt.Printf("Wrapping a default element %+v\n", e)
		return &cxpb.Fragment{&cxpb.Fragment_Aspect{}}
	}
}

func streamNetwork(network io.ReadCloser, out http.ResponseWriter) {
	address := serverAddress + ":" + serverPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		panic("Could not establish connection")
	}
	client := cxpb.NewCyServiceClient(conn)
	stream, err := client.StreamFragments(context.Background())
	if err != nil {
		panic("Could not open stream")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		decoder := cx.NewDecoder()
		for _, aspectName := range requiredAspects {
			fmt.Println("Registering handler for aspect ", aspectName)
			aspectName := aspectName
			decoder.RegisterAspectHandler(aspectName, func(d *json.Decoder) {
				if element, ok := knownAspects[aspectName]; ok {
					jsonpb.UnmarshalNext(d, element)
					frag := wrapElement(element)
					stream.Send(frag)
				}
			})
		}
		decoder.Decode(network)
		stream.CloseSend()
	}()
	go func() {
		defer wg.Done()
		encoder := json.NewEncoder(out)
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			} else {
				encoder.Encode(in.Element)
			}
		}
	}()
	wg.Wait()
}
