package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Question struct {
	Text    string
	Options []string //slice de string
	Answer  int
}

// estado do jogo
type GameState struct {
	Name      string
	Points    int
	Questions []Question //slice do tipo question
}

// funcao init como metodo de atualizacao de gamestate (atualiza os dados) || recebe como ponteiro pq quero modificar o original
func (g *GameState) Init() {
	fmt.Println("Seja bem vindo(a) ao quiz!")
	fmt.Println("Escreva o seu nome:")
	//bufio package || chamdas IO (terminal) || Reader || ler em blocos como bits e guarda temporariamente || escuta entrada do terminal
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n') //o "\n" é até o usuario digitar enter no terminal || le tudo o que a pessoa coloca no terminal até um delimitador "Enter" || retorna um bit ou um erro

	if err != nil {
		panic("Erro ao ler a string")
	}

	g.Name = name

	fmt.Printf("Vamos ao jogo %s", g.Name) //o "%s" é como se colocasse uma variavel e falo que ela é do tipo string || vai ser substituido pelo valor que passar || "Printf" formatando
}

// lendo o csv
func (g *GameState) ProccessCSV() {
	f, err := os.Open("quizGo.csv") //ponteiro para a interface
	if err != nil {
		panic("Erro ao ler arquivo")
	}

	//fechar o arquivo pois Go usa muita memoria para ler o arquivo || funcao defer termina depois de executar || sempre sera executado ao fim do processamento
	defer f.Close()

	//lendo os dados
	reader := csv.NewReader(f)
	records, err := reader.ReadAll() //matriz de uma matriz
	if err != nil {
		panic("Erro ao ler csv")
	}

	for index, record := range records {
		//montando struct
		if index > 0 {
			correctAnswer, _ := toInt(record[5]) //convertendo para inteiro
			question := Question{
				Text:    record[0],
				Options: record[1:5],
				Answer:  correctAnswer,
			}

			g.Questions = append(g.Questions, question) //adicionar item no slice
		}
	}

}

// metodo RUN - valida todo o jogo
func (g *GameState) Run() {
	// Exibir a pergunta pro usuário
	for index, question := range g.Questions {
		fmt.Printf("\033[33m %d. %s \033[0m\n", index+1, question.Text) //pega o index (numero) + pergunta || o "%d" será substituido por um numero inteiro (index) || "%s" substitui por string

		//iterar sobre as opções que temos no game state e exivir no terminal para o usuário
		for j, option := range question.Options { //segundo for (for dentro do for) é usado "j"
			fmt.Printf("[%d] %s\n", j+1, option)
		}

		fmt.Println("Digite uma alternativa")

		//coletar entrada do usuario || validar caractere inserido || se for errado, usuario precisa inserir novamente
		var answer int
		var err error

		for { //valores a serem modificados || infinito || até o usuario não colocar o valor certo, não sai deste for
			reader := bufio.NewReader(os.Stdin)
			read, _ := reader.ReadString('\n') //precisa remover o "\n" da string

			answer, err = toInt(read[:len(read)-1]) //atribuindo o valor || precisa fazer uma operação de slice || remove o tamanho dela - 1
			if err != nil {
				fmt.Println(err.Error())
				continue //vai começar o for novamente
			}
			break //comando para sair do for infinito
		}

		//validar a resposta || exibir mensagem se foi correta ou não || calcular pontuação
		if answer == question.Answer {
			fmt.Println("Parabéns você acertou!")
			g.Points += 10
		} else {
			fmt.Println("Ops! Errou!")
			fmt.Println("-------------------------")
		}

	}
}

func main() {
	//chamando o init
	game := &GameState{Points: 0}
	go game.ProccessCSV()

	game.Init()
	game.Run()

	fmt.Printf("Fim de jogo, você fez %d pontos!\n", game.Points)
}

// funcao que converte em inteiro
func toInt(s string) (int, error) {
	i, err := strconv.Atoi(s) //pacote que converte texto em inteiro (string para int)
	if err != nil {
		return 0, errors.New("não é permitido caractere diferente de número, por favor insira um número")
	}
	return i, nil
}
