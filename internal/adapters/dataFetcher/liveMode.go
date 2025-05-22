package datafetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"marketflow/internal/domain"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Exchange struct {
	number      string
	conn        net.Conn
	messageChan chan string
}

type LiveMode struct {
}

var _ domain.DataFetcher = (*LiveMode)(nil)

func (m *LiveMode) SetupDataFetcher() chan map[string]domain.ExchangeData {
	dataFlows := [3]chan domain.Data{make(chan domain.Data), make(chan domain.Data), make(chan domain.Data)}
	ports := []string{"40101", "40102", "40103"}

	wg := &sync.WaitGroup{}

	for i := 0; i < len(ports); i++ {
		wg.Add(1)
		exch, err := GenerateExchange(strconv.Itoa(i+1), "0.0.0.0:"+ports[i])
		if err != nil {
			log.Printf("Failed to connect exchange number: %d, error: %s", i, err.Error())
			return nil
		}

		// Receive data from the server
		go exch.FetchData(wg)

		// Start the vorker to process the received data
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

func Aggregate(mergedCh chan []domain.Data) chan map[string]domain.ExchangeData {
	aggregatedCh := make(chan map[string]domain.ExchangeData)

	go func() {
		defer close(aggregatedCh)

		for dataBatch := range mergedCh {
			exchangesData := make(map[string]domain.ExchangeData)
			counts := make(map[string]int)
			sums := make(map[string]float64)

			for _, data := range dataBatch {
				keys := []string{
					data.ExchangeName + " " + data.Symbol, // by exchange
					"All " + data.Symbol,                  // by all exchanges
				}

				for _, key := range keys {
					val, exists := exchangesData[key]
					if !exists {
						val = domain.ExchangeData{
							Exchange:  strings.Split(key, " ")[0],
							Pair_name: data.Symbol,
							Min_price: math.Inf(1),
							Max_price: math.Inf(-1),
						}
					}

					// обновление мин/макс
					if data.Price < val.Min_price {
						val.Min_price = data.Price
					}
					if data.Price > val.Max_price {
						val.Max_price = data.Price
					}

					sums[key] += data.Price
					counts[key]++

					exchangesData[key] = val
				}
			}

			// Counting avg price
			for key, ed := range exchangesData {
				if count, ok := counts[key]; ok && count > 0 {
					ed.Average_price = sums[key] / float64(count)
					exchangesData[key] = ed
				}
			}

			aggregatedCh <- exchangesData
		}
	}()

	return aggregatedCh
}

func MergeFlows(dataFlows [3]chan domain.Data) chan []domain.Data {
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

	ch := make(chan []domain.Data, 3)
	t := time.NewTicker(time.Second)
	rawData := make([]domain.Data, 0)
	mu := sync.Mutex{}

	go func() {
		for tick := range t.C {
			fmt.Println(tick)

			mu.Lock()
			ch <- rawData
			rawData = make([]domain.Data, 0)
			mu.Unlock()
		}
	}()
	go func() {
		for data := range mergedCh {
			rawData = append(rawData, data)
		}
	}()

	return ch
}

// GenerateExchange returns pointer to Exchange data with messageChan
func GenerateExchange(number string, address string) (*Exchange, error) {
	messageChan := make(chan string)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	exchangeServ := &Exchange{number: number, conn: conn, messageChan: messageChan}
	return exchangeServ, nil
}

// FetchData reads data from the connection and sends it to the channel
func (exch *Exchange) FetchData(wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(exch.conn)
	for scanner.Scan() {
		line := scanner.Text()
		exch.messageChan <- line
	}
	close(exch.messageChan) // Close the channel after reading is completed
}

// SetWorkers starts goroutine workers to process data
func (exch *Exchange) SetWorkers(wg *sync.WaitGroup, fan_in chan domain.Data) {
	resultChSlc := []chan domain.Data{}
	for w := 1; w <= 5; w++ {
		resultCh := make(chan domain.Data)
		wg.Add(1)
		go Worker(exch.number, exch.messageChan, resultCh, wg)
		resultChSlc = append(resultChSlc, resultCh)
	}

	// Wait for the Worckers to complete
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

// Worker processes tasks from the jobs channel and sends the results to the results channel
func Worker(number string, jobs chan string, results chan domain.Data, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		data := domain.Data{}
		err := json.Unmarshal([]byte(j), &data)
		if err != nil {
			log.Printf("Unmarshalling error in worker %s", err.Error())
			continue
		}

		// Assign the name of the exchange and send it to the results channel
		data.ExchangeName = number
		results <- data
	}
}
