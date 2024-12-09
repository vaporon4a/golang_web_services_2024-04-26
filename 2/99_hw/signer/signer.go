package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, j := range jobs {
		wg.Add(1)

		out := make(chan interface{})
		go JobWorker(j, in, out, wg)
		in = out
	}
	wg.Wait()
}

func JobWorker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	job(in, out)

}

func SingleHash(data, result chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for d := range data {
		wg.Add(1)
		go SingleHashWorker(d, result, wg, mu)
	}

	wg.Wait()

}

func SingleHashWorker(in interface{}, out chan interface{}, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	data := strconv.Itoa(in.(int))

	mu.Lock()
	md5Data := DataSignerMd5(data)
	mu.Unlock()

	crc32Chan := make(chan string)
	crc32Md5Chan := make(chan string)

	go asyncCrc32(data, crc32Chan)
	go asyncCrc32(md5Data, crc32Md5Chan)

	crc32Data := <-crc32Chan
	crc32Md5Data := <-crc32Md5Chan

	out <- crc32Data + "~" + crc32Md5Data

}

func asyncCrc32(data string, out chan string) {
	defer close(out)
	out <- DataSignerCrc32(data)
}

func MultiHash(data, result chan interface{}) {
	const num = 6
	wg := &sync.WaitGroup{}

	for d := range data {
		wg.Add(1)
		go multiHashWorker(d, result, num, wg)
	}
	wg.Wait()

}

func multiHashWorker(d interface{}, out chan interface{}, num int, wg *sync.WaitGroup) {
	defer wg.Done()

	jobWg := &sync.WaitGroup{}
	jobMu := &sync.Mutex{}

	jobSlice := make([]string, num)

	for i := 0; i < num; i++ {
		jobWg.Add(1)
		data := strconv.Itoa(i) + d.(string)
		go func(jobSlice []string, i int, data string, jobMu *sync.Mutex, jobWg *sync.WaitGroup) {
			defer jobWg.Done()
			data = DataSignerCrc32(data)
			jobMu.Lock()
			jobSlice[i] = data
			jobMu.Unlock()

		}(jobSlice, i, data, jobMu, jobWg)
	}

	jobWg.Wait()
	sb := strings.Builder{}
	for _, s := range jobSlice {
		sb.WriteString(s)
	}

	out <- sb.String()

}

func CombineResults(datas, result chan interface{}) {
	//defer close(result)
	//defer close(datas)
	sl := make([]string, 0)
	for d := range datas {
		sl = append(sl, d.(string))
	}
	sort.Strings(sl)
	result <- strings.Join(sl, "_")

}

func main() {
	// input := make(chan interface{})
	// out := make(chan interface{})
	// inputData := []int{0, 1}
	inputData := []int{0, 1, 1, 2, 3, 5, 8}

	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			for _, fibNum := range inputData {
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			dataRaw := <-in
			data, ok := dataRaw.(string)
			if !ok {
				panic("cant convert result data to string")
			}
			fmt.Println(data)
		}),
	}
	ExecutePipeline(hashSignJobs...)

}
