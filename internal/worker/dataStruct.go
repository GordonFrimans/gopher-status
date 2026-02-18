package worker

type Task struct {
	ID   int    // Уникальный ID задачи
	Data string // URL
}

// Result - результат обработки
type Result struct {
	TaskID int   // ID задачи, к которой относится результат
	Value  int   // Результат проверки сайта(statusCODE)
	Err    error // Ошибка возникшая при работе воркера
}
