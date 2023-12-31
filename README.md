<h1 align="center">Hi there, I'm <a href="https://github.com/Amore14rn" target="_blank">Roman</a> 

<h3 align="center">Test question in Anadea</h3>

- [Question](#Question)
- [Answer](#Answer)

## Question
"Какие методы можно использовать для считывания данных из логов, с использованием нескольких воркеров, а затем для записи этих данных в другую систему, гарантируя, что данные будут записаны только один раз и никакие данные не будут потеряны? 

## Answer

### Я использовал следующие методы:

#### Многопоточное считывание: 
- Я использую несколько горутин для считывания данных из логов (каждая горутина представляет отдельного воркера). Каждая горутина читает данные из канала и обрабатывает их параллельно.

#### Каналы для связи:
- Я использую каналы для передачи данных между горутинами. Это позволяет координировать работу воркеров и обеспечивает безопасный доступ к данным.

#### Контекст для управления:
- Я использовал контекст. Это позволяет корректно завершить работу воркеров при получении сигнала завершения.

#### Очередь задач: 
- Каждая горутина-воркер считывает задачи из общего канала и обрабатывает их. Причем, в случае сбоя записи во внешнюю систему, она возвращает ошибку, которая обрабатывается централизованно.

#### Обработка ошибок:
- Ошибки при записи во внешнюю систему обрабатываются в отдельной горутине, что позволяет асинхронно обрабатывать ошибки, не блокируя выполнение воркеров.

#### Graceful Sutdown:
- Я использовал сигналы SIGINT и SIGTERM для обработки запросов на Graceful Sutdown. Контекст позволяет корректно остановить работу воркеров и обработать Graceful Sutdown.

#### Безопасность доступа: 
- Я использовал мьютексы для безопасного доступа к данным при считывании логов. Это гарантирует, что данные не будут изменены несколькими горутинами одновременно.

#### Что можно улучшить:
- Можно добавить шифрование данных для передачи "зашифрованных логов" во внешнее хранилище. Это сделает передачу данных более безопасносной и если это чувствительные данные, как логи транзакций их лучше шифровать, хоть это и усложнит систему, но это позволит передавать данные в зашифровонном виде. Но если все находится в одном кластере можно обойтись и без шифрования.

#### Если вкратце то:
- errgroup: использована для управления группой горутин и обработки ошибок.
- context: Использование context позволяет управлять горутинками и корректно завершать их при получении сигналов завершения.
- Мьютексы: sync.Mutex используется для обеспечения безопасного доступа к данным.
- Каналы: Каналы используются для координации работы горутин и передачи данных между ними.

#### И последнее:
- Мой вариант не гарантирует,что данные будут записаны только один раз и никакие данные не будут потеряны. Для этого можно использовать или БД, что не очень, лучше использовать Apache Kafka или RabbitMQ. Это гарантирует что логи или данные будут обработаны и доставлены именно один раз.