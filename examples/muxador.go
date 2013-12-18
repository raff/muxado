package main

import (
	"flag"
        "fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/raff/muxado"
)

type Stream struct {
    muxado.Stream
}

func (s *Stream) log(fmts string, args ...interface{}) {
    id := string(s.StreamInfo()) + ": "
    log.Println(id + fmt.Sprintf(fmts, args...))
}

func handleStream(stream *Stream, id string) {
	defer func() {
		stream.log("closing stream")
		stream.Close()
	}()

	stream.log("handle stream")

	buffer := make([]byte, 1024)

	if n, err := stream.Read(buffer); err == nil {
		stream.Write(buffer[:n])
	} else {
                stream.log("read error: %s", err)
	}
}

func handleSession(sess muxado.Session, id string) {
	defer sess.Close()

	for {
		log.Println(id, "accept")

		stream, err := sess.Accept()
		if err != nil {
			log.Println(id, "accept", err)
			return
		} else {
			go handleStream(&Stream{stream}, id)
		}
	}
}

func client(sess muxado.Session, id string, wg *sync.WaitGroup) {
	s, err := sess.OpenEx(muxado.StreamInfo(id))
	if err != nil {
		log.Fatal(id, "open", err)
	}

        stream := &Stream{s}

	defer func() {
		stream.log("closing stream")
		stream.Close()
		if wg != nil {
			wg.Done()
		}
	}()

	stream.log("client loop")

	for i := 0; i < 10; i++ {
		stream.Write([]byte("hello there " + strconv.Itoa(i)))

		buffer := make([]byte, 1024)

		if n, err := stream.Read(buffer); err == nil {
			log.Println(id, "read", string(buffer[:n]))
		} else {
			log.Println(id, "read", err)
			break
		}
	}
}

func main() {
	server_mode := flag.Bool("server", false, "server mode (vs. client mode")
	port := flag.String("port", ":1111", "port to listen or dial to")

	flag.Parse()

	if !strings.Contains(*port, ":") { // make sure it looks like a port
		flag.Set("port", ":"+*port)
	}

	if !*server_mode && strings.HasPrefix(*port, ":") { // the client should point to a server address
		flag.Set("port", "localhost"+*port)
	}

	if *server_mode {
		//
		// the server
		//
		log.Println("listening on port", *port)

		l, err := muxado.Listen("tcp", *port)
		if err != nil {
			log.Fatal("server listen", err)
		}

		for {
			sess, err := l.Accept()
			if err != nil {
				log.Fatal("listen accept", err)
			}
			go handleSession(sess, "server")
			go client(sess, "server", nil)
		}
	} else {
		//
		// the client
		//
		sess, err := muxado.Dial("tcp", *port)
		if err != nil {
			log.Fatal(err)
		}

		defer sess.Close()

		var wg sync.WaitGroup

		for _, v := range []string{"1", "2", "3"} {
			wg.Add(1)

			go handleSession(sess, "client")
			go client(sess, "client "+v, &wg)
		}

		wg.Wait()
	}
}
