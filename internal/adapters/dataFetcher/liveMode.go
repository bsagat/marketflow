package datafetcher

import (
	"bufio"
	"encoding/json"
	"errors"
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
	closeCh     chan bool
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
		if m.Exchanges[i] == nil || m.Exchanges[i].conn == nil {
			continue
		}

		if err := m.Exchanges[i].conn.Close(); err != nil {
			log.Println("Failed to close connection: ", err.Error())
			continue
		}

		m.Exchanges[i].closeCh <- true
	}

}

func (m *LiveMode) SetupDataFetcher() (chan map[string]domain.ExchangeData, chan []domain.Data, error) {
	dataFlows := [3]chan domain.Data{make(chan domain.Data), make(chan domain.Data), make(chan domain.Data)}
	ports := []string{"40101", "40102", "40103"}

	wg := &sync.WaitGroup{}

	for i := 0; i < len(ports); i++ {
		wg.Add(1)
		exch, err := GenerateExchange("Exchange"+strconv.Itoa(i+1), "0.0.0.0:"+ports[i])
		if err != nil {
			log.Printf("Failed to connect exchange number: %d, error: %s", i+1, err.Error())
			wg.Done()
			continue
		}

		// Receive data from the server
		go exch.FetchData(wg)

		// Start the vorker to process the received data
		go exch.SetWorkers(wg, dataFlows[i])

		m.Exchanges = append(m.Exchanges, exch)
	}

	if len(m.Exchanges) != 3 {
		return nil, nil, errors.New("failed to connect to 3 exchanges")
	}

	mergedCh := MergeFlows(dataFlows)

	aggregated, rawDatach := Aggregate(mergedCh)

	go func() {
		wg.Wait()
		for i := 0; i < len(m.Exchanges); i++ {
			if m.Exchanges[i] == nil {
				continue
			}
			if m.Exchanges[i].conn != nil {
				m.Exchanges[i].conn.Close()
			}
		}

		slog.Info("All workers have finished processing.")
	}()
	return aggregated, rawDatach, nil
}

func Aggregate(mergedCh chan []domain.Data) (chan map[string]domain.ExchangeData, chan []domain.Data) {
	aggregatedCh := make(chan map[string]domain.ExchangeData)
	rawDataCh := make(chan []domain.Data)

	go func() {

		for dataBatch := range mergedCh {

			// To prevent the main thread from being delayed
			go func() {
				rawDataCh <- dataBatch
			}()

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
		close(aggregatedCh)
		close(rawDataCh)
	}()

	return aggregatedCh, rawDataCh
}

func MergeFlows(dataFlows [3]chan domain.Data) chan []domain.Data {
	mergedCh := make(chan domain.Data, 15)
	ch := make(chan []domain.Data, 3)

	closedCount := 0
	var muClosed sync.Mutex

	go func() {
		defer close(mergedCh)

	mainLoop:
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
				break mainLoop
			}
			muClosed.Unlock()
		}
	}()

	t := time.NewTicker(time.Second)
	rawData := make([]domain.Data, 0)
	done := make(chan bool)
	mu := sync.Mutex{}

	go func() {
		defer close(ch)

	mainLoop:
		for {
			select {
			case tick := <-t.C:
				slog.Debug(tick.String())
				mu.Lock()

				if len(rawData) == 0 {
					mu.Unlock()
					continue
				}
				ch <- rawData
				rawData = make([]domain.Data, 0)

				mu.Unlock()
			case <-done:
				break mainLoop
			}
		}

	}()

	go func() {
		for data := range mergedCh {
			mu.Lock()
			rawData = append(rawData, data)
			mu.Unlock()
		}

		done <- true
		close(done)
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

func (exch *Exchange) Reconnect(address string) error {
	var err error
	for i := 0; i < 5; i++ {
		time.Sleep(2 * time.Second)
		exch.conn, err = net.Dial("tcp", address)
		if err == nil {
			slog.Info("Reconnected to exchange: " + exch.number)
			return nil
		}
		slog.Warn("Reconnect attempt failed: " + err.Error())
	}
	return err
}

func (exch *Exchange) FetchData(wg *sync.WaitGroup) {
	defer wg.Done()

	scanner := bufio.NewScanner(exch.conn)
	address := exch.conn.RemoteAddr().String()

	closeCh := make(chan bool, 1)
	exch.closeCh = closeCh

	reconnect := true
	mu := sync.Mutex{}

	go func() {
		<-closeCh
		mu.Lock()
		reconnect = false
		mu.Unlock()
	}()

	log.Println("Starting reading data on exchange: ", exch.number)

	for {
		for scanner.Scan() && reconnect {
			line := scanner.Text()
			exch.messageChan <- line
		}

		log.Printf("Connection lost on exchange %s. Reconnecting...\n", exch.number)

		if reconnect {
			if err := exch.Reconnect(address); err != nil {
				log.Printf("Failed to reconnect exchange %s: %v", exch.number, err)
				break
			}

			scanner = bufio.NewScanner(exch.conn)
		} else {
			break
		}
	}

	log.Println("Giving up on exchange: ", exch.number)
	close(closeCh)
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
		slog.Debug("Local workers finished work in exchange " + exch.number)
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
