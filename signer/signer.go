package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const TH_COUNT = 6

func ExecutePipeline(jobs ...job) {
	wg := sync.WaitGroup{}

	inCh := make(chan interface{})
	outCh := make(chan interface{})
	wg.Add(len(jobs))
	for _, job := range jobs {
		go runWork(&wg, job, inCh, outCh)
		inCh = outCh
		outCh = make(chan interface{})
	}

	wg.Wait()
	close(outCh)
}

func runWork(wg *sync.WaitGroup, job job, inCh, outCh chan interface{}) {
	job(inCh, outCh)

	wg.Done()
	close(outCh)
}

func SingleHash(inCh, outCh chan interface{}) {
	wgSingleHash := sync.WaitGroup{}

	mutex := sync.Mutex{}

	for input := range inCh {

		wgSingleHash.Add(1)

		go operateSingleHash(&wgSingleHash, outCh, strconv.Itoa(input.(int)), &mutex)

	}

	wgSingleHash.Wait()

}

func operateSingleHash(wgSingleHash *sync.WaitGroup, outCh chan interface{}, data string, mutex *sync.Mutex) {
	crc32Ch := make(chan string)
	crc32md5Ch := make(chan string)

	go func(crc32Ch chan string, data string) {
		crc32Ch <- DataSignerCrc32(data)
	}(crc32Ch, data)

	go func(crc32md5Ch chan string, data string) {
		mutex.Lock()
		md5 := DataSignerMd5(data)
		mutex.Unlock()
		crc32md5Ch <- DataSignerCrc32(md5)
		close(crc32md5Ch)
	}(crc32md5Ch, data)

	crc32 := <-crc32Ch
	crc32md5 := <-crc32md5Ch
	singleHash := crc32 + "~" + crc32md5

	outCh <- singleHash

	wgSingleHash.Done()
}

func MultiHash(inCh, outCh chan interface{}) {
	wgMultiHash := sync.WaitGroup{}

	for data := range inCh {

		wgMultiHash.Add(1)

		go operateMultiHash(&wgMultiHash, outCh, data.(string))

	}

	wgMultiHash.Wait()
}

func operateMultiHash(wgMultiHash *sync.WaitGroup, outCh chan interface{}, SingleHash string) {
	fmt.Println(SingleHash)
	slice := make([]chan string, 6, 6)

	for i := 0; i < TH_COUNT; i++ {
		newCh := make(chan string)
		slice[i] = newCh
		go func(newCh chan string, i int, SingleHash string) {
			newCh <- DataSignerCrc32(strconv.Itoa(i) + SingleHash)
		}(newCh, i, SingleHash)
	}
	result := ""
	for i := 0; i < TH_COUNT; i++ {
		result += <-slice[i]
		close(slice[i])
	}
	outCh <- result
	wgMultiHash.Done()
}

func CombineResults(inCh, outCh chan interface{}) {
	var results []string
	for input := range inCh {
		strInput := input.(string)
		results = append(results, strInput)
	}
	sort.Strings(results)
	outCh <- strings.Join(results, "_")
}
