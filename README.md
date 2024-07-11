***Запуск через***
```
# соберет сервер, клиент и запустит сервер и клиент
make
# можно несколько раз запустить созданный клиент
make run-client
```

```go
func performPoW(challenge string, difficulty int) (string, int) {
    nonce := 0
    var hash string
    
    // Используем параметр N, который увеличивается с ростом сложности
    N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
    r := 8
    p := 1
    
    for {
        nonce++
        record := fmt.Sprintf("%s%d", challenge, nonce)
        h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
        hash = hex.EncodeToString(h[:])
        
        // Простая проверка: хэш должен начинаться с двух нулей для усложнения задачи
        if strings.HasPrefix(hash, "00") {
            break
        }
    }
    
    return hash, nonce
}
```

Server verify:
```go
func verifyPoW(challenge string, nonce int, hash string, difficulty int) bool {
	record := fmt.Sprintf("%s%d", challenge, nonce)
	N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
	r := 8
	p := 1
	h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
	calculatedHash := hex.EncodeToString(h[:])

	return strings.HasPrefix(calculatedHash, "00") && calculatedHash == hash
}
```

Попытка добиться линейного увелчения времени в зависимости от сложности, которая меняется в зависимости от скорости перебора хешей на основе статистики за последние 5 расчетов хэша для защиты от DDOS:
```go
Server is listening on port :12345
2024/07/10 07:34:33 INFO handle conn currentDifficulty=1 nonce=407 size=101
2024/07/10 07:34:35 INFO handle conn currentDifficulty=1 nonce=233 size=101
2024/07/10 07:34:38 INFO handle conn currentDifficulty=1 nonce=280 size=101
2024/07/10 07:34:41 INFO handle conn currentDifficulty=1 nonce=655 size=101
2024/07/10 07:34:45 INFO handle conn currentDifficulty=1 nonce=34 size=100
2024/07/10 07:34:47 INFO handle conn currentDifficulty=2 nonce=186 size=101
2024/07/10 07:34:51 INFO handle conn currentDifficulty=3 nonce=326 size=101
2024/07/10 07:34:56 INFO handle conn currentDifficulty=4 nonce=168 size=101
2024/07/10 07:35:02 INFO handle conn currentDifficulty=5 nonce=83 size=100
2024/07/10 07:35:20 INFO handle conn currentDifficulty=6 nonce=165 size=101
2024/07/10 07:35:31 INFO handle conn currentDifficulty=5 nonce=145 size=101
2024/07/10 07:35:37 INFO handle conn currentDifficulty=4 nonce=104 size=101
2024/07/10 07:35:43 INFO handle conn currentDifficulty=3 nonce=309 size=101
```

при дефолтной цели расчета времени хэша - 500 ms, окна последних расчетов = 5 элементов, сложности - 1
```go
const (
    initialDifficulty   = 1
    targetCalculateTime = 500 * time.Millisecond
    calculateWindow     = 5
)
```

***UPD***
Есть нюанс, что за время расчета берется время, когда клиент отправит данные с решением, но тут прибавляется ко времени расчета еще и время задержки по сети.
Подразумевается, что задержка будет минимальной и клиент будет коммуницировать с сервером, который максимально близко к нему расположен.
