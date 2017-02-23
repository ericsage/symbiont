package cx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Decoder struct {
	Handlers map[string]func(*json.RawMessage)
}

func NewDecoder() *Decoder {
	return &Decoder{Handlers: map[string]func(*json.RawMessage){}}
}

func (d *Decoder) RegisterAspectHandler(aspectName string, handler func(*json.RawMessage)) {
	d.Handlers[aspectName] = handler
}

//Decode a cx binary into a list of aspects
func (d *Decoder) Decode(cx io.ReadCloser) error {
	dec := json.NewDecoder(cx)
	//Remove opening array [
	stripDelimiter(dec)
	//Parse every fragment in the array
	for dec.More() {
		frag := make(map[string]*json.RawMessage)
		err := dec.Decode(&frag)
		if err != nil {
			return err
		}
		//Convert raw frag to aspect type
		d.parseFragment(frag)
	}
	//Remove closing array ]
	stripDelimiter(dec)
	//Parse EOF
	stripEOF(dec)
	return nil
}

func (d *Decoder) parseFragment(frag map[string]*json.RawMessage) {
	for name, elements := range frag {
		if handler, ok := d.Handlers[name]; ok {
			parseElements(elements, handler)
		}
		fmt.Println("Finished parsing ", name)
	}
}

func parseElements(eles *json.RawMessage, handler func(*json.RawMessage)) {
	buf := bytes.NewBuffer(*eles)
	dec := json.NewDecoder(buf)
	//Remove opening array [
	stripDelimiter(dec)
	//Parse every fragment in the array
	for dec.More() {
		var ele *json.RawMessage
		err := dec.Decode(&ele)
		if err != nil {
			panic(err)
		}
		//Handle element
		handler(ele)
	}
	//Remove closing array ]
	stripDelimiter(dec)
}

func stripDelimiter(dec *json.Decoder) {
	token, err := dec.Token()
	if err != nil {
		panic(err)
	}
	switch token.(type) {
	case json.Delim:
		return
	default:
		panic("Should be delimiter")
	}
}

func stripEOF(dec *json.Decoder) {
	_, err := dec.Token()
	if err != io.EOF {
		panic("Should be EOF")
	}
}
