package streaming

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"
)

func NewDecoderSSE(res *http.Response) Decoder[Event] {
	if res == nil || res.Body == nil {
		return nil
	}

	var decoder Decoder[Event]
	contentType := res.Header.Get("content-type")
	if t, ok := decoderTypes[contentType]; ok {
		decoder = t(res.Body)
	} else {
		scanner := bufio.NewScanner(res.Body)
		decoder = &eventStreamDecoder{rc: res.Body, scn: scanner}
	}
	return decoder
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder[Event]){}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder[Event]) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for s.scn.Scan() && s.err == nil {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
			}
			return true
		}

		// Split a string like "event: bar" into name="event" and value=" bar".
		name, value, _ := bytes.Cut(txt, []byte(":"))

		// Consume an optional space after the colon if it exists.
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			// An empty line in the for ": something" is a comment and should be ignored.
			continue
		case "event":
			event = string(value)
		case "data":
			_, s.err = data.Write(value)

			// Not sure why they are adding an new line here
			// _, s.err = data.WriteRune('\n')

		}
	}

	return false
}

func (s *eventStreamDecoder) Current() Event {
	return s.evt
}

func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

func (s *eventStreamDecoder) Err() error {
	return s.err
}
