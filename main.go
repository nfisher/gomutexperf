package main

import (
	"flag"
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	defaultDuration    = 60 * 5 // Run for 5m
	defaultThreadCount = 10
)

type Int64s []int64

func (a Int64s) Len() int {
	return len(a)
}

func (a Int64s) Less(i, j int) bool {
	return a[i] < a[j]
}

func (a Int64s) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

var histogram Int64s
var lock *sync.Mutex = &sync.Mutex{}
var start = time.Now()

func main() {
	var userThreads int
	var userDuration int

	flag.IntVar(&userThreads, "threads", defaultThreadCount, "specify the number of go routines to launch.")
	flag.IntVar(&userDuration, "duration", defaultDuration, "specify the number of go routines to launch.")
	flag.Parse()

	histogram = make(Int64s, userDuration)
	wg := &sync.WaitGroup{}

	fmt.Println("Running", userThreads, "threads for", userDuration, "seconds")

	for i := 0; i < userThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for incrementBucket() {
			}
		}()
	}

	wg.Wait()

	sort.Sort(histogram)

	best := userDuration - 1
	worst := 0
	pctl99 := int(0.99 * float64(userDuration))

	fmt.Println("best", histogram[best], "worst", histogram[worst], "99PCTL", histogram[pctl99])
}

func incrementBucket() (inRange bool) {
	bucket := int(time.Since(start) / time.Second)

	inRange = bucket < len(histogram)

	if !inRange {
		return
	}

	lock.Lock()
	histogram[bucket] = histogram[bucket] + 1
	lock.Unlock()
	return
}
