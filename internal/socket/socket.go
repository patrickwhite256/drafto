package socket

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type SocketServer struct {
	upgrader websocket.Upgrader

	topicChannels map[string][]chan struct{}
	topicLock     sync.RWMutex
}

func NewServer() *SocketServer {
	return &SocketServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		topicChannels: map[string][]chan struct{}{},
	}
}

func (s *SocketServer) Handler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 2 || parts[len(parts)-2] != "ws" {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		topic := parts[len(parts)-1]
		log.Printf("opening websocket on topic %s\n", topic)

		conn, err := s.upgrader.Upgrade(rw, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		topicCh := make(chan struct{})
		doneCh := make(chan struct{})

		s.topicLock.Lock()
		s.topicChannels[topic] = append(s.topicChannels[topic], topicCh)
		s.topicLock.Unlock()

		conn.SetCloseHandler(func(int, string) error {
			close(doneCh)
			err := conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				log.Printf("received error writing close message: %v\n", err)
			}

			return nil
		})

		// read loop so we know when to close
		go func() {
			for {
				_, _, err := conn.NextReader()
				if err != nil {
					return
				}
			}
		}()

	outer:
		for {
			select {
			case <-topicCh:
				conn.WriteMessage(websocket.TextMessage, []byte("r"))
			case <-doneCh:
				break outer
			}
		}

		log.Printf("closing websocket on topic %s\n", topic)

		s.topicLock.Lock()
		defer s.topicLock.Unlock()

		close(topicCh)
		for i, ch := range s.topicChannels[topic] {
			if ch == topicCh {
				s.topicChannels[topic][i] = s.topicChannels[topic][len(s.topicChannels[topic])-1]
				s.topicChannels[topic] = s.topicChannels[topic][:len(s.topicChannels[topic])-1]
				break
			}
		}
	})
}

// TODO: this should maybe take a context
func (s *SocketServer) NotifyTopic(topic string) {
	s.topicLock.RLock()
	defer s.topicLock.RUnlock()

	for _, ch := range s.topicChannels[topic] {
		ch <- struct{}{}
	}
	log.Printf("notified websockets on topic %s\n", topic)
}
