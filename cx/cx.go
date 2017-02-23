package cx

import "encoding/json"

type NumberVerification struct {
	LongNumber int64 `json:"longNumber"`
}

//MetaData holds aspect metadata
type MetaData struct {
	Name             string              `json:"name"`
	Version          string              `json:"version"`
	IDCounter        float64             `json:"idCounter"`
	Properties       []map[string]string `json:"properties"`
	ElementCount     float64             `json:"elementCount"`
	ConsistencyGroup float64             `json:"consistencyGroup"`
	Checksum         float64             `json:"checksum"`
}

type Aspect struct {
	Name      string             `json:"name"`
	Fragments []*json.RawMessage `json:"fragments"`
}

//KeyValue holds aspect properties
type KeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//Node represents a single node in a network
type Node struct {
	ID         int64  `json:"@id"`
	Name       string `json:"n"`
	Represents string `json:"r,omitempty"`
}

//Edge represents a single edge in a network
type Edge struct {
	ID          int64  `json:"@id"`
	SourceID    int64  `json:"s"`
	TargetID    int64  `json:"t"`
	Interaction string `json:"i"`
}

//NodeAttribute represents a single attribute on a single node in a network
type NodeAttribute struct {
	NodeID   int64  `json:"po"`
	Name     string `json:"n"`
	Value    string `json:"v"`
	Type     string `json:"d"`
	SubnetID int64  `json:"s"`
}

//EdgeAttribute represents a single attribute on a single edge in a network
type EdgeAttribute struct {
	EdgeID   int64  `json:"po"`
	Name     string `json:"n"`
	Value    string `json:"v"`
	Type     string `json:"d"`
	SubnetID int64  `json:"s,omitempty"`
}

//NetworkAttribute represents a single attribute on a network
type NetworkAttribute struct {
	Name     string      `json:"n"`
	Value    string      `json:"v"`
	Type     string      `json:"d"`
	SubnetID interface{} `json:"s"`
}

//MetaDataFragments contains all of the metaData aspect as a list of fragments
type MetaDataFragments struct {
	MetaData []MetaData `json:"metaData"`
}

//NodeFragments contains all of the node aspect as a list of fragments
type NodeFragments struct {
	Nodes []Node `json:"nodes"`
}

//EdgeFragments contains all of the edge aspect as a list of fragments
type EdgeFragments struct {
	Edges []Edge `json:"edges"`
}

//NodeAttributeFragments contains all of the nodeAttributes aspect as a list of fragments
type NodeAttributeFragments struct {
	NodeAttributes []NodeAttribute `json:"nodeAttributes"`
}

//EdgeAttributeFragments contains all of the edgeAttributes aspect as a list of fragments
type EdgeAttributeFragments struct {
	EdgeAttributes []EdgeAttribute `json:"edgeAttributes"`
}

//NetworkAttributeFragments contains all of the networkAttributes aspect as a list of fragments
type NetworkAttributeFragments struct {
	NetworkAttributes []NetworkAttribute `json:"networkAttributes"`
}
