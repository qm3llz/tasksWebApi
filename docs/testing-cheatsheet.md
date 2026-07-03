# Шпаргалка: тестирование HTTP-хендлеров в Go

> Цель: восстановить ход мысли за 5 минут. Завтра — попробуй пройти по шагам сам,
> подглядывая сюда только когда застрял. Запоминается то, что вывел заново, а не перечитал.

## Зачем всё это (одной фразой)

Хендлер нельзя тестировать, пока он намертво прибит к `*pgx.Conn` (реальной БД).
Решение: хендлер зависит от **интерфейса**, а в тесте на место БД подставляем **заглушку**.

```
Хендлер ──► TaskRepo (интерфейс) ──► настоящий репозиторий (в проде, идёт в БД)
                                 └──► fakeRepo (в тесте, живёт в памяти)
```

## Шаг 1. Интерфейс у потребителя

Объявляется в пакете `handler` (там, где ИСПОЛЬЗУЕТСЯ), не в `repository`.
Перечисляет ровно те методы, что хендлер реально зовёт. Сигнатуры — копия методов репозитория.

```go
type TaskRepo interface {
	Create(ctx context.Context, task models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (models.Task, error)
	GetAllByUser(ctx context.Context, userID uuid.UUID) ([]models.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, task models.Task) error
}
```

## Шаг 2. Хендлер зависит от интерфейса (а не от конкретного типа)

И поле, и аргумент конструктора — интерфейс. Тогда импорт `repository` из файла исчезает
(связь разорвана — это признак успеха).

```go
type TaskHandler struct { repo TaskRepo }

func NewTaskHandler(repo TaskRepo) *TaskHandler {   // НЕ *repository.TaskRepository
	return &TaskHandler{repo: repo}
}
```

Почему компилируется в `main.go`: настоящий `*repository.TaskRepository` имеет все 5 методов
→ он НЕЯВНО удовлетворяет `TaskRepo` (в Go нет `implements`, «крякает как утка — значит утка»).

## Шаг 3. Заглушка (fakeRepo) — в файле *_test.go

Поля под данные у каждого метода свои, ошибка `err` — общая.
Метод возвращает ПОЛЕ (`f.task`), а не захардкоженный литерал → каждый тест задаёт данные сам.
Получатель у всех методов одного типа называется одинаково (`f`).

```go
type fakeRepo struct {
	task models.Task   // что вернёт GetByID
	err  error         // nil = успех; не-nil = изобразить ошибку БД
}

func (f *fakeRepo) GetByID(ctx context.Context, id uuid.UUID) (models.Task, error) {
	return f.task, f.err
}
// + ещё 4 метода-заглушки, иначе *fakeRepo не удовлетворяет TaskRepo
// (Go проверит это в момент NewTaskHandler(&fakeRepo{}))
```

## Шаг 4. Table-driven тест

```go
func TestGetById(t *testing.T) {
	tests := []struct {
		name       string
		repoErr    error    // что вернёт заглушка
		wantStatus int      // какой статус ждём
	}{
		{"success",   nil,                      http.StatusOK},
		{"not found", errors.New("boom"),       http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// repo и h — ВНУТРИ цикла, заново на каждый случай (изоляция)
			repo := &fakeRepo{task: models.Task{Name: "x"}, err: tc.repoErr}
			h := NewTaskHandler(repo)

			body := strings.NewReader(`{"id":"11111111-1111-4111-1111-111111111111"}`)
			req := httptest.NewRequest(http.MethodGet, "/tasks", body)
			rec := httptest.NewRecorder()      // фейковый ResponseWriter, пишет в память

			h.GetById(rec, req)                // вызываем хендлер напрямую, без сети

			if rec.Code != tc.wantStatus {
				t.Errorf("%s: ждали %d, получили %d", tc.name, tc.wantStatus, rec.Code)
			}
		})
	}
}
```

Запуск: `go test ./internal/handler/ -v`

## Два правила, которые отличают надёжный тест от фикции

1. **Убедись, что тест умеет краснеть.** Временно соври в таблице (ждём не тот статус) →
   тест ДОЛЖЕН упасть. Упал → ему можно верить. Верни обратно.
2. **Подставляй поля `tc`, а не захардкоженные значения.** Если в теле цикла стоит
   `err: nil` и `http.StatusOK` вместо `tc.repoErr` / `tc.wantStatus` — таблица мёртвая,
   все строки прогоняют один случай. И сообщение `t.Errorf` печатай с `tc.wantStatus`,
   иначе диагностика соврёт при падении.

## Границы этого теста (что он НЕ проверяет)

Заглушка `GetByID` игнорирует `id` → разбор `id` из тела тест не покрывает.
`err` долетает до логики хендлера → его тест ловит. Поэтому:
- меняешь `err` в таблице → тест реагирует;
- меняешь тело/`id` → тест молчит.
Чтобы покрыть разбор `id`, нужен приём «проверка аргументов»: поле `gotID` в заглушке
запоминает, с каким id её позвали, и тест это сверяет.

## Шпаргалка по инструментам

| Инструмент | Зачем |
|---|---|
| `httptest.NewRequest(метод, путь, тело)` | собрать фейковый `*http.Request` |
| `httptest.NewRecorder()` | поймать ответ в память: `.Code`, `.Body`, `.Header()` |
| `t.Run(имя, func)` | подтест — виден отдельно, падение не глушит соседей |
| `t.Errorf` / `t.Fatalf` | пометить провал и продолжить / прекратить тест |
| `go test ./... -v` | прогнать все тесты подробно |
| `go fmt ./...` | выровнять форматирование по канону (отступы и пр.) |
