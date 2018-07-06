package main
import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
)
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		io.Copy(os.Stdout, conn) 
		log.Println("done")
	}()
	mustCopy(conn, os.Stdin)
	conn.Close()
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}