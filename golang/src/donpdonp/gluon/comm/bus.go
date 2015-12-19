package comm

import (
  "fmt"
  "net/url"
  "encoding/json"

  // message bus
  "github.com/gdamore/mangos"
  "github.com/gdamore/mangos/transport/tcp"
  "github.com/gdamore/mangos/protocol/bus"

)

type Bus struct {
  sock mangos.Socket
  Pipe chan map[string]interface{}
}

func Factory() (Bus, error) {
  new_bus := Bus{}
  bus_sock, err := bus.NewSocket()
  new_bus.sock = bus_sock
  new_bus.Pipe = make(chan map[string]interface{})
  return new_bus, err
}

func (comm *Bus) Start(url_str string) {
  u, _ := url.Parse(url_str)
  fmt.Printf("bus start %s\n", u)
  comm.addProtocol(u.Scheme)
  err := comm.sock.Listen(u.String())
  if err != nil {
    fmt.Println("can't listen on bus socket:", err.Error())
  }
}

func (comm *Bus) Loop() {
  var msg []byte
  var err error
  for {
    if msg, err = comm.sock.Recv(); err != nil {
      fmt.Println("Cannot recv: %s", err.Error())
    }
    jmsg := string(msg)
    fmt.Println("<-"+jmsg)
    var pkt map[string]interface{}
    json.Unmarshal(msg, &pkt)
    comm.Pipe <- pkt
  }
}

func (comm *Bus) Connect(url_str string) {
  u, _ := url.Parse(url_str)
  comm.addProtocol(u.Scheme)
  fmt.Println("bus connecting", u)
  err := comm.sock.Dial(u.String())
  if err != nil {
    fmt.Println("can't dial on bus socket:", err.Error())
  }
}

func (comm *Bus) Send(msg map[string]interface{}) {
  line, _ := json.Marshal(msg)
  err := comm.sock.Send(line)
  if err != nil {
    fmt.Println("Send err", err)
  } else{
    fmt.Println("->"+string(line))
  }
}

func (comm *Bus) addProtocol(scheme string) {
  switch scheme {
  case "tcp":
    comm.sock.AddTransport(tcp.NewTransport())
  }
}