package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Task representa uma tarefa na lista de tarefas
type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var tasks []Task

func main() {
	// Carrega as tarefas do arquivo
	loadTasksFromFile("tasks.json")

	// Inicia o servidor HTTP em uma goroutine
	go startHTTPServer()

	// Loop principal do programa
	for {
		// Exibe as opções do menu
		fmt.Println("\n-- Menu --")
		fmt.Println("1. Adicionar tarefa")
		fmt.Println("2. Exibir tarefas")
		fmt.Println("3. Completar tarefa")
		fmt.Println("4. Sair")
		fmt.Print("Escolha uma opção: ")

		// Lê a escolha do usuário
		choice := readInput()

		// Executa a ação com base na escolha do usuário
		switch choice {
		case "1":
			addTask()
		case "2":
			displayTasks()
		case "3":
			completeTask()
		case "4":
			fmt.Println("Saindo...")
			saveTasksToFile("tasks.json", tasks)
			return
		default:
			fmt.Println("Opção inválida. Por favor, escolha novamente.")
		}
	}
}

// readInput lê a entrada do usuário
func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// addTask adiciona uma nova tarefa à lista
func addTask() {
	for {
		fmt.Print("Digite a tarefa: ")
		task := readInput()
		tasks = append(tasks, Task{ID: len(tasks) + 1, Text: task})
		fmt.Println("Tarefa adicionada com sucesso!")

		// Pergunta se deseja adicionar mais tarefas
		fmt.Print("Deseja adicionar mais uma tarefa? (s/n): ")
		choice := readInput()
		if choice != "s" {
			return
		}
	}
}

// displayTasks exibe as tarefas e oferece a opção de voltar ao menu
func displayTasks() {
	fmt.Println("\n--- Tarefas ---")
	for _, task := range tasks {
		fmt.Printf("%d. %s\n", task.ID, task.Text)
	}

	// Pergunta se deseja voltar ao menu
	fmt.Print("\nPressione Enter para voltar ao menu...")
	readInput()
}

// completeTask completa uma tarefa e remove da lista
func completeTask() {
	for {
		fmt.Println("\n--- Tarefas ---")
		for _, task := range tasks {
			fmt.Printf("%d. %s\n", task.ID, task.Text)
		}
		fmt.Print("Digite o número da tarefa concluída: ")
		taskID := readInput()
		id, err := strconv.Atoi(taskID)
		if err != nil || id < 1 || id > len(tasks) {
			fmt.Println("Número de tarefa inválido.")
			continue
		}
		tasks = removeTaskByID(tasks, id)
		reorganizeTaskIDs(&tasks)
		fmt.Println("Tarefa concluída e removida da lista.")

		// Pergunta se deseja completar mais tarefas
		fmt.Print("Deseja completar mais uma tarefa? (s/n): ")
		choice := readInput()
		if choice != "s" {
			return
		}
	}
}

// reorganizeTaskIDs reorganiza os IDs das tarefas após a exclusão de uma tarefa
func reorganizeTaskIDs(tasks *[]Task) {
	for i := range *tasks {
		(*tasks)[i].ID = i + 1
	}
}

// removeTaskByID remove a tarefa com o ID especificado da lista de tarefas
func removeTaskByID(tasks []Task, id int) []Task {
	index := -1
	for i, task := range tasks {
		if task.ID == id {
			index = i
			break
		}
	}
	if index != -1 {
		return append(tasks[:index], tasks[index+1:]...)
	}
	return tasks
}

// startHTTPServer inicia o servidor HTTP
func startHTTPServer() {
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/add", handleAddTask)
	http.HandleFunc("/tasks/delete", handleDeleteTask)

	fmt.Println("Servidor HTTP escutando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleTasks retorna todas as tarefas como JSON
func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// handleAddTask adiciona uma nova tarefa à lista
func handleAddTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks = append(tasks, newTask)
	w.WriteHeader(http.StatusCreated)
}

// handleDeleteTask exclui uma tarefa da lista
func handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	// Implemente a lógica de exclusão da tarefa aqui
}

// saveTasksToFile salva as tarefas em um arquivo JSON
func saveTasksToFile(filename string, tasks []Task) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Erro ao criar o arquivo %s: %v", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(tasks)
	if err != nil {
		log.Fatalf("Erro ao codificar as tarefas: %v", err)
	}
}

// loadTasksFromFile carrega as tarefas de um arquivo JSON
func loadTasksFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Arquivo %s não encontrado. Criando um novo...\n", filename)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatalf("Erro ao decodificar as tarefas: %v", err)
	}
}
