// Пакет для воркеров которые будут мониторить сайты и тп
package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

var (
	ErrPoolOverload = errors.New("pool is overloaded")
	ErrPoolClosed   = errors.New("worker pool is closed")
	ErrPanicWorker  = errors.New("worker recover")
)

// Тестовая функ для проверки статуса сайта по URL (UP, DOWN, PENDING)
func SimpleCheck(url string) (int, error) {
	UserAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11.1; rv:84.0) Gecko/20100101 Firefox/84.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.89 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; U; Android 9; Redmi Note 7 Build/PKQ1.180904.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/77.0.3865.92 Mobile Safari/537.36 OPR/44.1.2254.143214",
		"Mozilla/5.0 (Linux; Android 11; SM-A705FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.105 Mobile Safari/537.36 OPR/63.3.3216.58675",
		"Mozilla/5.0 (Linux; Android 8.1.0; 16th) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.2; Redmi 4X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Mobile Safari/537.36",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	defer cancel()

	// Создаем запрос с контекстом
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Добавляем User-Agent, так как некоторые защиты (Cloudflare, etc)
	// блокируют пустые или дефолтные Go-агенты

	randUserAgent := UserAgents[rand.Int64N(int64(len(UserAgents))-1)]
	fmt.Println("used UserAgent = ", randUserAgent)
	req.Header.Set("User-Agent", randUserAgent)

	// Используем свой клиент, чтобы контролировать параметры транспорта, если нужно
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// Функция которая выполняет полезную работу (!проверку сайта и тп!)
type Handler func(data string) (int, error)

type WorkerPool struct {
	// Конфигурация
	numWorkers int
	handler    Handler
	//---Logger---

	// Каналы
	// tasks - буферизированный канал для входящих задач
	tasks chan Task

	// results - буферизированный канал для ответов

	results chan Result

	// Синхронизация
	// wg - нужна, чтобы Close() ждал, пока все воркеры доделают работу

	wg sync.WaitGroup

	// mu и closed - нужны для thread-safe Submit() и Close()
	// (тот же паттерн, что мы учили в логгере)
	mu     sync.RWMutex
	closed bool
}

func NewWorkerPool(numWorkers int, handler Handler) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		handler:    handler,
		// Буфер tasks ставим хотя бы = кол-ву воркеров, чтобы они сразу схватили задачи
		tasks: make(chan Task, numWorkers),
		// Буфер results тоже полезен, чтобы воркеры не блочились, если никто быстро не читает
		results: make(chan Result, numWorkers),
	}
}

func (w *WorkerPool) Start() error {
	log.Println("INFO", "Start WorkerPool")
	for i := 0; i < w.numWorkers; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			for task := range w.tasks {
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Println("ERR", fmt.Sprintf("Panic в задаче %d: %v", task.ID, r))
							w.results <- Result{
								TaskID: task.ID,
								Value:  0,
								Err:    fmt.Errorf("panic recovered: %v", r),
							}
						}
					}()

					// Теперь безопасно вызываем handler
					res, err := w.handler(task.Data)
					w.results <- Result{
						TaskID: task.ID,
						Value:  res,
						Err:    err,
					}
				}()
			}
		}()

	}
	return nil
}

func (w *WorkerPool) Submit(t Task) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// 1. Проверяем, не закрыт ли пул, пока держим лок
	if w.closed {
		log.Println("ERR", "Пул закрыт")
		return ErrPoolClosed
	}

	// 2. Идиоматичный non-blocking send
	w.tasks <- t

	log.Println("INFO", fmt.Sprintf("task %d успешно отправлен", t.ID))
	return nil
}

func (w *WorkerPool) Results() <-chan Result {
	return w.results
}

func (w *WorkerPool) Close() error {
	w.mu.Lock() // ← Полный лок
	if w.closed {
		w.mu.Unlock()
		return errors.New("already closed")
	}
	w.closed = true
	close(w.tasks)
	w.mu.Unlock() // ← Не держи лок во время Wait!

	w.wg.Wait()
	close(w.results)
	return nil
}
