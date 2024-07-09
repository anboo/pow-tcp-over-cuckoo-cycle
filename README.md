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

***Объяснение почему Blake2:***
PoW должен быть достаточно сложным, чтобы затруднить массовую отправку запросов злоумышленником, но при этом не слишком сложным для легитимных пользователей. Высокая производительность Blake2b позволяет эффективно решать эту задачу.
Но blake2 базируется на вычислительной сложности, но не учитывает использование специализированного оборудования (ASIC) для ускорения вычислений, что может быть недостатком в некоторых сценариях.