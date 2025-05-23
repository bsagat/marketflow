package datafetcher

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
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
	Exchanges []*Exchange
}

func NewLiveModeFetcher() *LiveMode {
	return &LiveMode{Exchanges: make([]*Exchange, 0)}
}

var _ domain.DataFetcher = (*LiveMode)(nil)

func (m *LiveMode) CheckHealth() error {
	var unhealthy string
	for i := 0; i < len(m.Exchanges); i++ {
		select {
		case _, ok := <-m.Exchanges[i].messageChan:
			if !ok {
				unhealthy += m.Exchanges[i].number + " "
			}
		default:
			continue
		}

	}
	if len(unhealthy) != 0 {
		return errors.New("unhealthy exchanges: " + unhealthy)
	}
	return nil
}

func (m *LiveMode) Close() {
	for i := 0; i < len(m.Exchanges); i++ {
		if err := m.Exchanges[i].conn.Close(); err != nil {
			log.Println("Failed to close connection: ", err.Error())
		}
	}

}

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

		m.Exchanges = append(m.Exchanges, exch)
	}

	mergedCh := MergeFlows(dataFlows)

	aggregated := Aggregate(mergedCh)

	go func() {
		wg.Wait()
		for i := 0; i < len(m.Exchanges); i++ {
			m.Exchanges[i].conn.Close()
		}

		fmt.Println("All workers have finished processing.")
	}()
	return aggregated
}

func Aggregate(mergedCh chan []domain.Data) chan map[string]domain.ExchangeData {
	aggregatedCh := make(chan map[string]domain.ExchangeData)

	go func() {
		defer close(aggregatedCh)

		for dataBatch := range mergedCh {
			fmt.Println(len(dataBatch))

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
					ed.Timestamp = time.Now()
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
	ch := make(chan []domain.Data, 3)

	closedCount := 0
	var muClosed sync.Mutex

	go func() {
		defer close(mergedCh)

		for {
			select {
			case e1, ok := <-dataFlows[0]:
				if !ok {
					muClosed.Lock()
					closedCount++
					muClosed.Unlock()
					dataFlows[0] = nil
				} else {
					mergedCh <- e1
				}
			case e2, ok := <-dataFlows[1]:
				if !ok {
					muClosed.Lock()
					closedCount++
					muClosed.Unlock()
					dataFlows[1] = nil
				} else {
					mergedCh <- e2
				}
			case e3, ok := <-dataFlows[2]:
				if !ok {
					muClosed.Lock()
					closedCount++
					muClosed.Unlock()
					dataFlows[2] = nil
				} else {
					mergedCh <- e3
				}
			}

			muClosed.Lock()
			if closedCount == 3 {
				muClosed.Unlock()
				break
			}
			muClosed.Unlock()
		}
	}()

	t := time.NewTicker(time.Second)
	rawData := make([]domain.Data, 0)
	mu := sync.Mutex{}

	go func() {
		defer close(ch)

		for tick := range t.C {
			slog.Debug(tick.String())

			mu.Lock()
			if len(rawData) == 0 {
				mu.Unlock()
				select {
				case _, ok := <-mergedCh:
					if !ok {
						return
					}
				default:
				}
				continue
			}
			ch <- rawData
			rawData = make([]domain.Data, 0)
			mu.Unlock()
		}
	}()

	go func() {
		for data := range mergedCh {
			mu.Lock()
			rawData = append(rawData, data)
			mu.Unlock()
		}

		t.Stop()
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
	close(exch.messageChan)
}

// SetWorkers starts goroutine workers to process data
func (exch *Exchange) SetWorkers(globalWg *sync.WaitGroup, fan_in chan domain.Data) {
	workerWg := &sync.WaitGroup{}
	for w := 1; w <= 5; w++ {
		workerWg.Add(1)
		globalWg.Add(1)
		go func() {
			Worker(exch.number, exch.messageChan, fan_in, workerWg)
			globalWg.Done()
		}()
	}

	go func() {
		workerWg.Wait()
		close(fan_in)
	}()
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
