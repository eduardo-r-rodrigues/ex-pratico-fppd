// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
//@Author Eduardo Rosa Rodrigues
//@Author Bernardo Lykawka Medeiros Silva
// EXERCÍCIO:  dado o programa abaixo
//    1) quantos processos concorrentes são gerados ?
//		R: 40
//    2) execute e observe: que se pode supor sobre a velocidade relativa dos mesmos ?
// OBSERVACAO:o sleep no método main serve para este nao acabar, o que acabaria todos processos em execucao.
//     mais adiante veremos outras formas de sincronizar isto
//		R: A velocidade relativa é diferente, já que depende da ordem de execução das Threads e
// 		 dos recursos disponibilizados pelo hardware

package main

import (
	"fmt"
	"time"
)

var N int = 40

func funcaoA(id int, s string) {
	for {
		fmt.Println(s, id)
	}
}

func geraNespacos(n int) string {
	s := "  "
	for j := 0; j < n; j++ {
		s = s + "   "
	}
	return s
}

func main() {
	for i := 0; i < N; i++ {
		go funcaoA(i, geraNespacos(i))
	}
	for true {
		time.Sleep(100 * time.Millisecond)
	}
}
