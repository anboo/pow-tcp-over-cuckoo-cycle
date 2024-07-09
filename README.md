***Запуск через***
```
make
```

Задача:
Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge

Объяснение:

Пробовал Hashcash, Argon2, в итоге проанализировал интернет и выбрал CuckooCycle.
Hashcash, Argon2 при нужной сложности, которая действительно замедлит DDOS атаку и сделает ее нерентабельной, потребует примерно 10 мб и 60 мб по памяти и есть вариант, что клиент решит либо слишком быстро задачу, либо слишком медленно и придется обвязываться таймаутами и повторными попытками получить новую задач.
Cuckoo Cycle отличается высокими требованиями к памяти, в среднем около 160 МБ, что усложняет и удорожает разработку специализированных ASIC устройств и в целом удорожает DDOS атаку.