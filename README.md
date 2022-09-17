# IoT Platform

MVP курсового проекта по курсу "Архитектура высоких нагрузок"

![image info](./doc/otus_system_design.png)


## Описание проекта

Для демонстрации работы платформы реализован следующий основной функционал:

- Имитатор IoT устройства. Приложение cmd/emulator имитирует работу 10 000 gps-маяков. Каждый маяк выходит 
на связь каждые 10-20 минут - производит подключение к tcp-серверу и производит передачу пакета данных, 
в котором содержится следующая информация:
  - id (номер сим-карты) маяка
  - гео-координаты маяка
  - скорость
  - направление передвижения
  - заряд батареи
  - температуру
  - и т.д.
  
  Все значения имитируемых параметров устройства постоянно меняются.

- TCP-сервер для приема и обработки данных. Приложение cmd/server_tcp обрабатывает входящие подключения, парсит пришедший пакет данных и 
помещает данные в очередь rabbitMQ.
- Обработчик данных rabbitMQ. Приложение cmd/processor_rmq читает канал поступления данных от rabbitMQ, производит 
сохранение или обновление данных в БД PostgreSQL.
- web-сервер для отображения данных на карте. Приложение cmd/server_api выполняет запросы к БД PostgreSQL и отображает 
на карте перемещение маяков.

## Настройка и запуск
- Выполняем git clone проекта на локальный диск. Переходим в папку проекта и выполняем `go mod tidy` и `go mod vendor`
- Производим переименование файла c env-переменными из `configs/app.env_example` в `configs/app.env`.
- Выполняем сборку приложений командами `make build_emulator`,`make build_server`,`make build_rmq`,`make build_api`
- Производим запуск БД и RabbitMQ командой `docker-compose up` 
- Выполняем запуск приложений командами `make run_emulator`,`make run_server`,`make run_rmq`,`make run_api` 
- Открываем браузер по адресу `localhost:8082` и наслаждаемся картинкой движения маяков.

![image info](./doc/screenshot.png)

![alt text](./doc/beacon.gif)