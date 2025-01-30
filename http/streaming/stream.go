package streaming

// A Streamer is same as decoder
//
//go:generate go run go.uber.org/mock/mockgen@latest -destination=mock.go -package=streaming .  Streamer,Decoder
type Streamer[E any] Decoder[E]

type Decoder[E any] interface {
	Current() E
	Next() bool
	Close() error
	Err() error
}

// StreamHandler defines the interface for provider-specific stream handling
type StreamHandler[TypeOut any, TypeIn any] interface {
	HandleEvent(TypeIn) (TypeOut, error)
	ShouldContinue(TypeIn) bool
}

// Stream provides core streaming functionality that can be reused across providers
type Stream[TypeOut any, TypeIn any] struct {
	decoder Decoder[TypeIn] // Same as Stream
	handler StreamHandler[TypeOut, TypeIn]
	current TypeOut
	err     error
	done    bool
}

func NewStream[TypeOut any, TypeIn any](
	decoder Decoder[TypeIn],
	handler StreamHandler[TypeOut, TypeIn],
) *Stream[TypeOut, TypeIn] {
	return &Stream[TypeOut, TypeIn]{
		decoder: decoder,
		handler: handler,
		done:    false,
	}
}

func (s *Stream[T, E]) Next() bool {
	if s.err != nil || s.done {
		return false
	}

	for s.decoder.Next() {
		if !s.handler.ShouldContinue(s.decoder.Current()) {
			s.done = true
			return false
		}

		current, err := s.handler.HandleEvent(s.decoder.Current())
		if err != nil {
			s.err = err
			return false
		}

		s.current = current
		return true
	}

	s.err = s.decoder.Err()
	return false
}

func (s *Stream[T, E]) Current() T {
	return s.current
}

func (s *Stream[T, E]) Err() error {
	return s.err
}

func (s *Stream[T, E]) Close() error {
	return s.decoder.Close()
}
