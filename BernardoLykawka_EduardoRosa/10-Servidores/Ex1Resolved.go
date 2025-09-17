package main

import (
	"fmt"
	"math/rand"
)

const NClients = 10

type Request struct {
	value   int
	replyCh chan int
}

// client: gera pedidos continuamente
func client(id int, srv chan<- Request) {
	myReply := make(chan int)
	for {
		v := rand.Intn(1000)
		srv <- Request{value: v, replyCh: myReply}
		r := <-myReply
		fmt.Println("client:", id, " req:", v, " resp:", r)
	}
}

// handler: trata um pedido e responde no canal do cliente
func handle(id int, req Request) {
	// (cálculo qualquer; aqui só dobra)
	req.replyCh <- req.value * 2
	fmt.Println("                       handled by goroutine", id)
}

// concurrent server: cria uma goroutine por request
func serverConcurrent(in <-chan Request) {
	reqID := 0
	for req := range in {
		reqID++
		go handle(reqID, req)
	}
}

func main() {
	fmt.Println("------ Concurrent Server (1 goroutine/request) -------")
	srvCh := make(chan Request)
	go serverConcurrent(srvCh)
	for i := 0; i < NClients; i++ {
		go client(i, srvCh)
	}
	select {} // bloqueia main
}