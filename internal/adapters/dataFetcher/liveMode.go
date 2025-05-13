package datafetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"marketflow/internal/domain"
	"net"
	"strconv"
	"sync"
)

type Exchange struct {
	name        string
	conn        net.Conn
	messageChan chan string
}

type LiveMode struct {
}

func (m *LiveMode) SetupDataFetcher() chan domain.Data {
	dataFlows := [3]chan domain.Data{make(chan domain.Data), make(chan domain.Data), make(chan domain.Data)}
	ports := []string{"40101", "40102", "40103"}

	wg := &sync.WaitGroup{}

	for i := 0; i < len(ports); i++ {
		wg.Add(1)
		exch, err := GenerateExchange("Exchange"+strconv.Itoa(i+1), "0.0.0.0:"+ports[i])
		if err != nil {
			log.Printf("Failed to connect exchange number: %d, error: %s", i, err.Error())
			return nil
		}

		// Получаем данные с сервера
		go exch.FetchData(wg)

		// Запускаем воркеров для обработки полученных данных
		go exch.SetWorkers(wg, dataFlows[i])
	}

	mergedCh := MergeFlows(dataFlows)

	aggregated := Aggregate(mergedCh)

	go func() {
		wg.Wait()
		fmt.Println("All workers have finished processing.")
	}()
	return aggregated
}

func Aggregate(mergedCh chan domain.Data) chan domain.Data {

	aggregatedCh := make(chan domain.Data, 15)

	// функция аггрегирования
	go func() {
		for {
			data := <-mergedCh

			// Обработка данных

			// Отправляем данные
			aggregatedCh <- data
		}
	}()

	return aggregatedCh
}

func MergeFlows(dataFlows [3]chan domain.Data) chan domain.Data {
	mergedCh := make(chan domain.Data, 15)
	go func() {
		for {
			select {
			case e1 := <-dataFlows[0]:
				mergedCh <- e1
			case e2 := <-dataFlows[1]:
				mergedCh <- e2
			case e3 := <-dataFlows[2]:
				mergedCh <- e3
			}
		}
	}()
	return mergedCh
}

// GenerateExchange returns pointer to Exchange data with messageChan
func GenerateExchange(name, address string) (*Exchange, error) {
	messageChan := make(chan string)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	exchangeServ := &Exchange{name: name, conn: conn, messageChan: messageChan}
	return exchangeServ, nil
}

// FetchData читает данные из соединения и отправляет их в канал
func (exch *Exchange) FetchData(wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(exch.conn)
	for scanner.Scan() {
		line := scanner.Text()
		exch.messageChan <- line
	}
	close(exch.messageChan) // Закрываем канал после завершения чтения
}

// SetWorkers запускает рабочих горутин для обработки данных
func (exch *Exchange) SetWorkers(wg *sync.WaitGroup, fan_in chan domain.Data) {
	resultChSlc := []chan domain.Data{}
	for w := 1; w <= 5; w++ {
		resultCh := make(chan domain.Data)
		wg.Add(1)
		go Worker(exch.name, exch.messageChan, resultCh, wg)
		resultChSlc = append(resultChSlc, resultCh)
	}

	// Ожидаем завершения работы воркеров
	go func() {
		wg.Wait()
		for i := 0; i < len(resultChSlc); i++ {
			close(resultChSlc[i])
		}
	}()

	for {
		select {
		case d1 := <-resultChSlc[0]:
			fan_in <- d1
		case d2 := <-resultChSlc[1]:
			fan_in <- d2
		case d3 := <-resultChSlc[2]:
			fan_in <- d3
		case d4 := <-resultChSlc[3]:
			fan_in <- d4
		case d5 := <-resultChSlc[4]:
			fan_in <- d5
		}
	}
}

// Worker обрабатывает задачи из канала jobs и отправляет результаты в канал results
func Worker(exchName string, jobs chan string, results chan domain.Data, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		data := domain.Data{}
		err := json.Unmarshal([]byte(j), &data)
		if err != nil {
			log.Printf("Unmarshalling error in worker %s", err.Error())
			continue
		}
		// Присваиваем имя биржи и отправляем в канал результатов
		data.ExchangeName = exchName
		results <- data
	}
}
