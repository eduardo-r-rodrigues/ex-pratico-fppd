// Ex2-Servidores2-SOLUTION.go
package main

import (
	"fmt"
	"math/rand"
)

const (
	NClients = 100
	Pool     = 10 // limite de concorrência
)

type Request struct {
	value   int
	replyCh chan int
}

func client(id int, srv chan<- Request) {
	myReply := make(chan int)
	for {
		v := rand.Intn(1000)
		srv <- Request{value: v, replyCh: myReply}
		r := <-myReply
		fmt.Println("client:", id, " req:", v, " resp:", r)
	}
}

// worker: lê pedidos e trata enquanto houver no canal
func worker(wid int, in <-chan Request) {
	for req := range in {
		req.replyCh <- req.value * 2
		// (log opcional para ver qual worker atendeu)
		// fmt.Println("                         worker", wid, "handled")
	}
}

// servidor com pool fixo de workers
func serverWithPool(in <-chan Request) {
	for i := 1; i <= Pool; i++ {
		go worker(i, in)
	}
}

func main() {
	fmt.Println("------ Server with fixed worker pool (max 10) -------")
	srvCh := make(chan Request)
	serverWithPool(srvCh)      // inicia 10 workers
	for i := 0; i < NClients; i++ {
		go client(i, srvCh)
	}
	select {} // bloqueia main
}
