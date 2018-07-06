package main 
import (
	"log"
	"net"
	"fmt"
	"bufio"
)
type client chan<- string
var (
	entering=make(chan client)
	leaving=make(chan client)
	messages=make(chan string)
)
func main() {
	listener,err:=net.Listen("tcp","localhost:8000")
	if err!=nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn,err:=listener.Accept()
		if err!=nil {
			log.Print(err)
		}
		go handleConn(conn)
	}
}
func broadcaster() {
	clients:=make(map[client]bool)
	for {
		select {
		case msg :=<-messages:
			/*
			Send each message received to the 
			client that has already connected the server
			*/
			for cli:= range clients {
				cli<-msg
			}
			//Log in to the incoming client
		case cli:=<-entering:
			clients[cli]=true
		case cli:=<-leaving:
			delete(clients,cli)
			close(cli)
		}
	}
}
func handleConn(conn net.Conn) {
	ch:=make(chan string)
	go clientWriter(conn,ch)
	who:=conn.RemoteAddr().String()
	ch<-"You:"+who
	messages<-who+" arrived"
	entering<-ch
	input:=bufio.NewScanner(conn)
	for input.Scan() {
		//This is used to wrap when outputting on the client side.
		messages<-who+":"+input.Text()+"\n" 
	}
	leaving<-ch
	messages<-who+" left"
	conn.Close()
}
func clientWriter(conn net.Conn,ch <-chan string) {
	for msg:=range ch {
		fmt.Fprintf(conn,msg)
	}
}
