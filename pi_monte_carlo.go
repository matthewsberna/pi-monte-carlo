package main

import (
    "fmt"
    "math"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// Config
const (
	NumWorkers = 5
	TotalPoints = 1000000000
	BufferSize = 1000000
)

type Point struct {
	X, Y float64
}

func main() {
	startTime := time.Now()

	// Channels
	requestChan := make(chan Point, BufferSize)
	responseChan := make(chan int, BufferSize)

	// PRNG
	seed := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(seed)

	// WaitGroup
	wg := &sync.WaitGroup{}

	// Spawn the workers.
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		go isInCircle(i, requestChan, responseChan, wg)
	}
		
	// Feed the workers.
	go func(){
		for i := 0; i < TotalPoints; i++ {
			p := Point{X: rng.Float64(), Y: rng.Float64()}
			requestChan <- p
		}
		close(requestChan)
	}()

	// Clean up when all responses have been sent.
	go func(){
		wg.Wait()
		close(responseChan)
	}()

	// Handle responses and calculate pi!
	responseCount := 0
	pointsInCircle := 0
	for response := range responseChan {
		responseCount++
		pointsInCircle += response
	}
	fmt.Printf("Received %v responses.\n", responseCount)
	fmt.Printf("%v points are inside circle.\n", pointsInCircle)
	// I'm using a big float because I have nothing else to do during quarantine.
	pi := new(big.Float).SetPrec(20).SetFloat64(float64(4)*(float64(pointsInCircle)/float64(responseCount)))
	fmt.Printf("\nEstimated pi = %.20f\n", pi)
	fmt.Printf("\nElapsed time:  %v", time.Since(startTime))
}


func isInCircle(workerNum int, requestChan chan Point, responseChan chan int, wg *sync.WaitGroup){
	pointCount := 0
	for p := range requestChan {
		val := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
		ret := 0
		if val <= 1.0 {
			ret = 1
		}
		responseChan <- ret
		pointCount++
	}
	fmt.Printf("Worker %v processed %v points.\n", workerNum, pointCount)
	wg.Done()
}

