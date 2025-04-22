package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	//task1_2()

	//task3()
	//task4()
	//task5()
	//task6()
	//task7()
	task8()
}

func task1_2() {
	var wg sync.WaitGroup
	wg.Add(5)

	for i := 1; i <= 5; i++ {
		go printWorker(i, &wg)

	}
	wg.Wait()
	fmt.Println("Все воркеры завершены")

}

func printWorker(worker int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Hello from goroutine!", worker)
	time.Sleep(time.Second)
}

func task3() {

	ch := make(chan int)
	go func() {
		for i := 1; i < 6; i++ {
			fmt.Println("to ch ", i)
			ch <- i
		}
		close(ch)

	}()

	var sum int
	for i := range ch {
		sum += i
	}
	fmt.Println(sum)

}

func task4() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	countWorker := 10
	var counter int
	wg.Add(countWorker)
	for i := 1; i <= countWorker; i++ {
		go func() {
			defer wg.Done()
			defer mu.Unlock()
			mu.Lock()
			counter++
			fmt.Println("counter:", counter)
		}()
	}
	wg.Wait()
	fmt.Println("comp")
}

func task5() {
	var counter int32
	countWorker := 10
	var wg sync.WaitGroup

	for i := 1; i <= countWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt32(&counter, 1)
			fmt.Println(atomic.LoadInt32(&counter))
		}()
	}
	wg.Wait()
	fmt.Println("comp", counter)
}

func task6() {
	testURLs := []string{
		// Обычные веб‑страницы
		"https://www.example.com",
		"https://www.google.com",

		// Простой JSON‑API
		"https://jsonplaceholder.typicode.com/posts/1",

		// Симуляция медленного ответа (задержка 3 сек)
		"https://httpstat.us/200?sleep=3000",

		// Статус 404
		"https://httpstat.us/404",

		// Некорректный хост (должно вернуть "error")
		"http://invalid.nonexistent.domain",

		// HTTPS‑страница с большим телом
		"https://www.gutenberg.org/files/1342/1342-0.txt",
	}

	results := FetchURLs(testURLs)
	for url, body := range results {
		fmt.Printf("URL: %s\nResponse: %s\n\n", url, body)
	}

}

func FetchURLs(urls []string) map[string]string {
	results := make(map[string]string, len(urls))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()

			resp, err := http.Get(u)
			if err != nil {
				mu.Lock()
				results[u] = "error"
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				mu.Lock()
				results[u] = "error"
				mu.Unlock()
				return
			}

			snippet := string(body)
			if len(snippet) > 100 {
				snippet = snippet[:100]
			}

			mu.Lock()
			results[u] = snippet
			mu.Unlock()
		}(url)
	}

	wg.Wait()
	return results
}

// 1, 2, 3, 10, 11
func task7() {
	c1 := make(chan int, 3)
	c2 := make(chan int, 2)
	for i := 1; i <= 3; i++ {
		c1 <- i
	}
	close(c1)
	for i := 10; i <= 11; i++ {
		c2 <- i
	}
	close(c2)

	merged := Merge(c1, c2)

	for i := range merged {
		fmt.Println("Merged: ", i)
	}

}

func Merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func task8() {
	tasks := make(chan int)

	workers := Split(tasks, 3)

	var wg sync.WaitGroup
	for id, ch := range workers {
		wg.Add(1)
		go func(id int, c <-chan int) {
			defer wg.Done()
			for v := range c {
				fmt.Printf("wr %d t %d\n", id, v)
			}
		}(id, ch)
	}

	go func() {
		for i := 1; i <= 10; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	wg.Wait()

}

func Split(ch <-chan int, n int) []<-chan int {
	outs := make([]chan int, n)
	for i := range outs {
		outs[i] = make(chan int)
	}

	go func() {
		defer func() {
			for _, o := range outs {
				close(o)
			}
		}()

		i := 0
		for v := range ch {
			outs[i%n] <- v
			i++
		}
	}()

	result := make([]<-chan int, n)
	for i, o := range outs {
		result[i] = o
	}
	return result
}
