package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/raff/muxado"
)

func streamInfo(s muxado.Stream, prefix string) {
	log.Printf("%slocal: %v, remote: %v, id: %v, related: %v, info: %v\n", prefix, s.LocalAddr(), s.RemoteAddr(), s.Id(), s.RelatedStreamId(), string(s.StreamInfo()))
}

func handleStream(stream muxado.Stream, id string) {
	defer func() {
		streamInfo(stream, "closing "+id)
		stream.Close()
	}()

	streamInfo(stream, id+" ")

	buffer := make([]byte, 1024)

	if n, err := stream.Read(buffer); err == nil {
		stream.Write(buffer[:n])
	} else {
		log.Println(id, "read", err)
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
			go handleStream(stream, id)
		}
	}
}

func client(sess muxado.Session, id string, wg *sync.WaitGroup) {
	//stream, err := sess.OpenEx(muxado.StreamInfo("client" + v))
	stream, err := sess.Open()
	if err != nil {
		log.Fatal(id, "open", err)
	}

	defer func() {
		streamInfo(stream, "closing "+id+" ")
		stream.Close()
		if wg != nil {
			wg.Done()
		}
	}()

	streamInfo(stream, id+" ")

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
