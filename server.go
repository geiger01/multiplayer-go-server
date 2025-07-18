package main

import (
    "encoding/json"
    "fmt"
    "net"
    "sync"
    "time"
)

type Player struct {
    ID string  `json:"id"`
    X  float64 `json:"x"`
    Y  float64 `json:"y"`
}

var (
    players   = make(map[string]Player)
    addresses = make(map[string]*net.UDPAddr)
    mu        sync.Mutex
)

func main() {
    addr, _ := net.ResolveUDPAddr("udp", ":9999")
    conn, _ := net.ListenUDP("udp", addr)
    defer conn.Close()

    fmt.Println("🚀 Server listening on UDP port 9999")

    go handleBroadcast(conn)

    buf := make([]byte, 1024)
    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            continue
        }

        var p Player
		fmt.Printf("Received: %s from %s\n", buf[:n], clientAddr.String())
        if err := json.Unmarshal(buf[:n], &p); err == nil {
            mu.Lock()
            players[p.ID] = p
            addresses[p.ID] = clientAddr
            mu.Unlock()
        }
    }
}

func handleBroadcast(conn *net.UDPConn) {
    ticker := time.NewTicker(50 * time.Millisecond)
    for range ticker.C {
        mu.Lock()
        for _, addr := range addresses {
            for _, p := range players {
                data, _ := json.Marshal(p)
                conn.WriteToUDP(data, addr)
            }
        }
        mu.Unlock()
    }
}
