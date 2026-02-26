package worker

import (
	"log"
	"time"

	"main/internal/storage"
)

type Scheduler struct {
	wp      *WorkerPool
	storage *storage.InMemoryStorageMonitors
	quit    chan struct{}
}

func NewScheduler(wp *WorkerPool, storage *storage.InMemoryStorageMonitors) *Scheduler {
	return &Scheduler{
		wp:      wp,
		storage: storage,
		quit:    make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Second)
	log.Println("INFO", "Scheduler started 🕒")

	// !!! ИСПРАВЛЕНИЕ 1: Запускаем слушателя результатов!
	// Без этого канал results забивается и всё виснет.
	go s.processResults()

	go func() {
		for {
			select {
			case <-ticker.C:
				s.scheduleTasks()
			case <-s.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.quit)
}

func (s *Scheduler) scheduleTasks() {
	// 1. Получаем список всех мониторов
	monitors, err := s.storage.List("adminadmin1332adminadmin")
	if err != nil {
		log.Println("ERR", "Не удалось получить задачи:", err)
		return
	}

	for _, m := range monitors {
		// --- ЛОГИКА ПРОВЕРКИ ВРЕМЕНИ ---

		// 1. Если LastCheck пустой (монитор только создан) -> Проверяем сразу!
		if m.LastCheck == "" {
			log.Println("DEFAULT lastChek")
			s.storage.UpdateLastCheck(m.ID, time.Now().Format("2006-01-02 15:04:05"))
			s.sendTask(&m)
			continue
		}

		// 2. Парсим время последней проверки из строки
		// (Формат должен совпадать с тем, как ты сохраняешь! Обычно RFC3339)
		layout := "2006-01-02 15:04:05"
		lastCheckTime, err := time.ParseInLocation(layout, m.LastCheck, time.Local)
		if err != nil {
			log.Printf("ERR: Ошибка парсинга времени для ID=%d: %v", m.ID, err)
			// Если ошибка, лучше проверить на всякий случай
			s.sendTask(&m)
			continue
		}

		// 3. Вычисляем время следующей проверки
		// nextCheck = lastCheck + interval (секунд)
		nextCheckTime := lastCheckTime.Add(time.Duration(m.Interval) * time.Second)

		// 4. Если текущее время МЕНЬШЕ следующей проверки -> РАНО! Пропускаем.
		if time.Now().Before(nextCheckTime) {
			continue
		}

		// --- КОНЕЦ ЛОГИКИ ---

		// Если дошли сюда -> ПОРА ПРОВЕРЯТЬ!
		s.storage.UpdateLastCheck(m.ID, time.Now().Format("2006-01-02 15:04:05"))
		s.sendTask(&m)
	}
}

// Вспомогательная функция для отправки, чтобы не дублировать код
func (s *Scheduler) sendTask(m *storage.Monitor) { // *storage.Monitor - поменяй на свой тип

	task := Task{
		ID:   int(m.ID),
		Data: m.URL,
	}

	// Важно: запускаем в горутине, чтобы не блокировать цикл перебора,
	// если воркер пул занят.
	go func(t Task) {
		s.wp.Submit(t)
	}(task)
}

func (s *Scheduler) processResults() {
	// Читаем результаты. Цикл кончится, когда закроется канал results (при остановке пула)
	for res := range s.wp.Results() {
		var newStatus string
		if res.Err != nil {
			log.Printf("Monitor ID=%d CHECK FAILED: %v", res.TaskID, res.Err)
			newStatus = "DOWN"
		} else if res.Value >= 200 && res.Value < 300 {
			newStatus = "UP"
		} else {
			log.Printf("Monitor ID=%d BAD STATUS: %d", res.TaskID, res.Value)
			newStatus = "DOWN"
		}

		// Обновляем статус
		err := s.storage.UpdateStatusByID(int64(res.TaskID), newStatus)
		if err != nil {
			log.Printf("ERR: Failed to update status for ID=%d: %v", res.TaskID, err)
		} else {
			log.Printf("INFO: Monitor ID=%d status updated to %s", res.TaskID, newStatus)
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		err = s.storage.UpdateLastCheck(int64(res.TaskID), timestamp)
		// !!! ТУТ ВАЖНО: Обнови LastCheck (время проверки)
		// s.storage.UpdateLastCheck(int64(res.TaskID), time.Now())
		if err != nil {
			log.Printf("ERR: Failed to update LastTime for ID=%d: %v", res.TaskID, err)
		}
	}
}
