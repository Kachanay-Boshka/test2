/*
Программа реализует:
1. Демультиплексирование (разделение одного канала на несколько)
2. Мультиплексирование (объединение нескольких каналов в один)
*/
package main

import (
	"fmt"
	"sync"
)

// Количество сообщений для генерации
const messageCount = 5

// Demultiplexer разделяет входной канал на несколько выходных каналов
// Принимает: inputChan: канал-источник данных; count: количество выходных каналов
// Возвращает: Слайс выходных каналов
func Demultiplexer(inputChan chan int, count int) []chan int {
	// Создаем слайс выходных каналов
	outputChannels := make([]chan int, count)

	// Канал для сигнала завершения работы
	doneSignal := make(chan struct{})

	// Инициализируем выходные каналы
	for i := range outputChannels {
		outputChannels[i] = make(chan int)

		// Для каждого канала запускаем горутину, которая закроет канал по сигналу
		go func(outChan chan int) {
			<-doneSignal
			close(outChan)
		}(outputChannels[i])
	}

	go func() {
		defer close(doneSignal) // По завершении посылаем сигнал закрытия

		// Читаем данные из входного канала и рассылаем по всем выходным
		for data := range inputChan {
			for _, outChan := range outputChannels {
				outChan <- data
			}
		}
	}()
	fmt.Print("asaas")
	return outputChannels
}

// Multiplexer объединяет несколько каналов в один
// Принимает: channels: каналы для объединения
// Возвращает: Канал с объединенными данными
func Multiplexer(channels ...chan int) chan int {
	multiplexingChan := make(chan int)
	var wg sync.WaitGroup

	// Функция для чтения из одного канала и записи в объединенный
	multipelx := func(inputChan <-chan int) {
		defer wg.Done()
		for data := range inputChan {
			multiplexingChan <- data
		}
	}

	// Запускаем горутины для каждого входного канала
	wg.Add(len(channels))
	for _, inputChan := range channels {
		go multipelx(inputChan)
	}

	// Горутина для закрытия объединенного канала после завершения всех читателей
	go func() {
		wg.Wait()
		close(multiplexingChan)
	}()

	return multiplexingChan
}

func main() {
	dataGenerator := func() chan int {
		dataChan := make(chan int)

		go func() {
			defer close(dataChan)
			for i := 1; i <= messageCount; i++ {
				dataChan <- i
			}
		}()

		return dataChan
	}

	demuxOutputs := Demultiplexer(dataGenerator(), 5)

	resultChan := Multiplexer(demuxOutputs...)
	for data := range resultChan {
		fmt.Println(data)
	}
}
