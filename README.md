# Домашнее задание #7

1. Логи уже были структурные (zerolog), сейчас только немного изменил интерфейс и добавил больше логов.
2. Трейсами покрыл все, кроме кафки.
   ![trace-1](./readmeImages/img.png)
   ![trace-2](./readmeImages/img_1.png)
   ![trace-3](./readmeImages/img_2.png)

3. Метрики на gRPC серверы, клиенты и БД есть.
   ![grafana-1](./readmeImages/img_5.png)
   ![grafana-2](./readmeImages/img_6.png)

4. Сделал алерты в тг на падение сервисов, высокий процент ошибок и высокую задержку ответов.
   ![alerts-1](./readmeImages/img_3.png)
   ![alerts-2](./readmeImages/img_4.png)

5. Сделал wrapper для клиентов и интерцепторы для серверов (см. `./lib/grpc`).
6. Сделал wrapper для клиента к БД (см. `./lib/db`).