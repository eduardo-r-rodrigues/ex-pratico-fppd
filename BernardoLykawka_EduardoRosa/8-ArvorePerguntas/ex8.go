// por Fernando Dotti - PUCRS
// dado abaixo um exemplo de estrutura em arvore, uma arvore inicializada
// e uma operação de caminhamento, pede-se fazer:
//   1.a) a operação que soma todos elementos da arvore.
//        func soma(r *Nodo) int {...}
//   1.b) uma operação concorrente que soma todos elementos da arvore
//   [OS ACIMA ESTAO RESOLVIDOS]
//   2.a) a operação de busca de um elemento v, dizendo true se encontrou v na árvore, ou falso
//        func busca(r* Nodo, v int) bool {}...}
//   2.b) a operação de busca concorrente de um elemento, que informa imediatamente
//        por um canal se encontrou o elemento (sem acabar a busca), ou informa
//        que nao encontrou ao final da busca
//   3.a) a operação que escreve todos pares em um canal de saidaPares e
//        todos impares em um canal saidaImpares, e ao final avisa que acabou em um canal fin
//        func retornaParImpar(r *Nodo, saidaP chan int, saidaI chan int, fin chan struct{}){...}
//   3.b) a versao concorrente da operação acima, ou seja, os varios nodos sao testados
//        concorrentemente se pares ou impares, escrevendo o valor no canal adequado
//
//  ABAIXO: RESPOSTAS A QUESTOES 1a e b
//  APRESENTE A SOLUÇÃO PARA AS DEMAIS QUESTÕES

package main

import (
	"fmt"
	"sync"
)

type Node struct {
	value int
	left  *Node
	right *Node
}

func inorder(root *Node) {
	if root != nil {
		inorder(root.left)
		fmt.Print(root.value, ", ")
		inorder(root.right)
	}
}

func sum(root *Node) int {
	if root != nil {
		return root.value + sum(root.left) + sum(root.right)
	}
	return 0
}

func sumConcurrent(root *Node) int {
	result := make(chan int)
	go sumConcurrentCh(root, result)
	return <-result
}

func sumConcurrentCh(root *Node, result chan int) {
	if root != nil {
		child := make(chan int)
		go sumConcurrentCh(root.left, child)
		go sumConcurrentCh(root.right, child)
		result <- (root.value + <-child + <-child)
	} else {
		result <- 0
	}
}

// Sequential search
func search(root *Node, target int) bool {
	for root != nil {
		if target == root.value {
			return true
		}
		if target < root.value {
			root = root.left
		} else {
			root = root.right
		}
	}
	return false
}

// Concurrent search
func searchConcurrent(root *Node, target int) bool {
	result := make(chan bool, 1)
	var wg sync.WaitGroup
	var once sync.Once

	var walk func(*Node)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		if n.value == target {
			once.Do(func() { result <- true })
			return
		}
		wg.Add(2)
		go func() { defer wg.Done(); walk(n.left) }()
		go func() { defer wg.Done(); walk(n.right) }()
	}

	wg.Add(1)
	go func() { defer wg.Done(); walk(root) }()

	go func() {
		wg.Wait()
		once.Do(func() { result <- false })
	}()

	return <-result
}

// Even/Odd concurrent output
func splitEvenOdd(root *Node, outEven chan int, outOdd chan int, done chan struct{}) {
	var wg sync.WaitGroup

	var traverse func(*Node)
	traverse = func(n *Node) {
		if n == nil {
			return
		}
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			if val%2 == 0 {
				outEven <- val
			} else {
				outOdd <- val
			}
		}(n.value)

		traverse(n.left)
		traverse(n.right)
	}

	traverse(root)

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
}

func main() {
	root := &Node{value: 10,
		left: &Node{value: 5,
			left: &Node{value: 3,
				left:  &Node{value: 1},
				right: &Node{value: 4}},
			right: &Node{value: 7,
				left:  &Node{value: 6},
				right: &Node{value: 8}}},
		right: &Node{value: 15,
			left: &Node{value: 13,
				left:  &Node{value: 12},
				right: &Node{value: 14}},
			right: &Node{value: 18,
				left:  &Node{value: 17},
				right: &Node{value: 19}}}}

	outEven := make(chan int)
	outOdd := make(chan int)
	done := make(chan struct{})

	fmt.Println("\nEven/Odd values:")
	go splitEvenOdd(root, outEven, outOdd, done)

	finished := false
	for count := 0; count < 20 && !finished; {
		select {
		case even := <-outEven:
			fmt.Println("Even:", even)
			count++
		case odd := <-outOdd:
			fmt.Println("Odd:", odd)
			count++
		case <-done:
			finished = true
		}
	}

	fmt.Print("\nInorder: ")
	inorder(root)
	fmt.Println("\n")

	fmt.Println("Sum:      ", sum(root))
	fmt.Println("SumConc:  ", sumConcurrent(root))
	fmt.Println()
	fmt.Println("Search 17:", search(root, 17))
	fmt.Println("Search 99:", search(root, 99))
	fmt.Println()
	fmt.Println("SearchC 17:", searchConcurrent(root, 17))
	fmt.Println("SearchC 99:", searchConcurrent(root, 99))
}
