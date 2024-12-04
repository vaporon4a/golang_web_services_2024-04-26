package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

const goroutinesNum = 5

func startWorker(workerNum int, in <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for input := range in {
		fmt.Printf(formatWork(workerNum, input))
		runtime.Gosched() // попробуйте закомментировать
	}
	printFinishWork(workerNum)
}

func main() {
	wg := &sync.WaitGroup{}
	runtime.GOMAXPROCS(0)               // попробуйте с 0 (все доступные) и 1
	worketInput := make(chan string, 2) // попробуйте увеличить размер канала
	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, worketInput, wg)
	}

	months := []string{"Январь", "Февраль", "Март",
		"Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь",
		"Октябрь", "Ноябрь", "Декабрь",
	}

	for _, monthName := range months {
		worketInput <- monthName
	}
	close(worketInput) // попробуйте закомментировать

	// time.Sleep(time.Millisecond)
	wg.Wait()
}

func formatWork(in int, input string) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"th", in,
		"recieved", input)
}

func printFinishWork(in int) {
	fmt.Println(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum-in),
		"===", in,
		"finished")
}
