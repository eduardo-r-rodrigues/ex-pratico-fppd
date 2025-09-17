package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	NJ = 5
	M  = 4
)
type Carta string

func jogador(
	id int,
	maoInicial []Carta,
	in <-chan Carta,
	out chan<- Carta,
	canalBatida chan<- int,
	inicioBatida <-chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	mao := make([]Carta, M)
	copy(mao, maoInicial)
	jaBati := false

	for !jaBati {
		if temJogo(mao) {
			fmt.Printf("Jogador %d formou um jogo com a mão: %v e vai bater!\n", id, mao)
			canalBatida <- id
			jaBati = true
			continue
		}

		select {
		case cartaRecebida, ok := <-in:
			if !ok {
				return
			}
			
			mao = append(mao, cartaRecebida)
			cartaParaPassar := escolherCartaParaPassar(mao)
			mao = removerCarta(mao, cartaParaPassar)
			out <- cartaParaPassar

		case <-inicioBatida:
			if !jaBati {
				fmt.Printf("Jogador %d viu que alguém bateu! Batendo também...\n", id)
				canalBatida <- id
				jaBati = true
			}
		}
	}

	fmt.Printf("Jogador %d bateu e agora está apenas repassando cartas.\n", id)
	for cartaRecebida := range in {
		out <- cartaRecebida
	}
}

func temJogo(mao []Carta) bool {
	if len(mao) != M {
		return false
	}
	primeiraCarta := mao[0]
	for i := 1; i < M; i++ {
		if mao[i] != primeiraCarta {
			return false
		}
	}
	return true
}

func escolherCartaParaPassar(mao []Carta) Carta {
	contagem := make(map[Carta]int)
	for _, c := range mao {
		contagem[c]++
	}

	maiorContagem := 0
	cartaMaisComum := mao[0]
	for c, n := range contagem {
		if n > maiorContagem {
			maiorContagem = n
			cartaMaisComum = c
		}
	}

	var cartaParaPassar Carta
	for _, c := range mao {
		if c != cartaMaisComum {
			if c != "@" {
				return c
			}
			cartaParaPassar = c
		}
	}
	if cartaParaPassar != "" {
		return cartaParaPassar
	}
	return mao[len(mao)-1]
}

func removerCarta(mao []Carta, cartaRemover Carta) []Carta {
	for i, c := range mao {
		if c == cartaRemover {
			return append(mao[:i], mao[i+1:]...)
		}
	}
	return mao
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("--- Iniciando o Jogo Dorminhoco ---")

	canaisCartas := make([]chan Carta, NJ)
	for i := 0; i < NJ; i++ {
		canaisCartas[i] = make(chan Carta, 1)
	}
	canalBatida := make(chan int, NJ) 
	inicioBatida := make(chan struct{})   
	var wg sync.WaitGroup

	baralho := make([]Carta, 0, NJ*M+1)
	tiposDeCarta := []Carta{"A", "B", "C", "D", "E", "F", "G"}[:NJ]
	for i := 0; i < NJ; i++ {
		for j := 0; j < M; j++ {
			baralho = append(baralho, tiposDeCarta[i])
		}
	}
	baralho = append(baralho, "@")
	rand.Shuffle(len(baralho), func(i, j int) {
		baralho[i], baralho[j] = baralho[j], baralho[i]
	})
	fmt.Printf("Baralho com %d cartas criado e embaralhado.\n", len(baralho))

	for i := 0; i < NJ; i++ {
		maoInicial := baralho[i*M : (i+1)*M]
		
		in := canaisCartas[i]
		out := canaisCartas[(i+1)%NJ]

		wg.Add(1)
		go jogador(i, maoInicial, in, out, canalBatida, inicioBatida, &wg)
		fmt.Printf("Jogador %d criado com a mão: %v\n", i, maoInicial)
	}

	cartaInicial := baralho[NJ*M]
	fmt.Printf("\nO jogo começa! A carta coringa '%s' foi enviada para o Jogador 0.\n\n", cartaInicial)
	canaisCartas[0] <- cartaInicial

	ordemDeBatida := make([]int, 0, NJ) // registra ordem de batida
	
	primeiroBatedor := <-canalBatida
	ordemDeBatida = append(ordemDeBatida, primeiroBatedor)
	fmt.Printf("\n--- BATIDA! O Jogador %d foi o primeiro a bater! ---\n", primeiroBatedor)
	fmt.Println("--- Todos os outros jogadores devem bater agora! ---")

	close(inicioBatida)

	for i := 0; i < NJ-1; i++ {
		batedor := <-canalBatida
		ordemDeBatida = append(ordemDeBatida, batedor)
	}

	fmt.Println("\n--- Fim de Jogo! ---")
	fmt.Println("A ordem de batida foi:")
	for i, jogadorID := range ordemDeBatida {
		fmt.Printf("%dº lugar: Jogador %d\n", i+1, jogadorID)
	}

	dorminhoco := ordemDeBatida[len(ordemDeBatida)-1]
	fmt.Printf("\nO DORMIONHOCO (último a bater) é o Jogador %d!\n", dorminhoco)

	for _, ch := range canaisCartas {
		close(ch)
	}
	wg.Wait()
}