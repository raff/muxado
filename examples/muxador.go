package main

import (
    "flag"
    "log"
    "strings"

    "github.com/raff/muxado"
)

func streamInfo(s muxado.Stream, prefix string) {
    log.Printf("%slocal: %v, remote: %v, id: %v, related: %v\n", prefix, s.LocalAddr(), s.RemoteAddr(), s.Id(), s.RelatedStreamId())
}

func handleStream(stream muxado.Stream) {
}

func handleSession(sess muxado.Session) {
    for {
        stream, err := sess.Accept()
        if err != nil {
            log.Println(err)
            return
        } else {
            streamInfo(stream, "server ")
            go handleStream(stream)
        }
    }
}

func main() {
    server_mode := flag.Bool("server", false, "server mode (vs. client mode")
    port := flag.String("port", ":1111", "port to listen or dial to")

    flag.Parse()

    if !strings.Contains(*port, ":") { // make sure it looks like a port
        flag.Set("port", ":" + *port)
    }

    if !*server_mode && strings.HasPrefix(*port, ":") { // the client should point to a server address
        flag.Set("port", "localhost" + *port)
    }

    if *server_mode {
        log.Println("listening on port", *port)

        l, err := muxado.Listen("tcp", *port)
        if err != nil {
            log.Fatal(err)
        }

        for {
            sess, err := l.Accept()
            if err != nil {
                log.Fatal(err)
            }
            go handleSession(sess)
        }
    } else {
        sess, err := muxado.Dial("tcp", *port)
        if err != nil {
            log.Fatal(err)
        }

        stream, err := sess.Open()
        if err != nil {
            log.Fatal(err)
        }

        streamInfo(stream, "client ")

        stream.Write([]byte("hello there"))
        stream.Close()
    }
}
