### Имитатор терминалов ГЛОНАСС/GPS и анализ канала передачи телеметрии

# Аргументы CLI
- ```cfg``` - путь к файлу конфига (значение по умолчанию: ```рабочая директория```)
- ```qos``` - путь к директории таблиц QoS (значение по умолчанию: ```рабочая директория/qos/```)

# Файл конфига
Пример полного конфига:
```
{
    "simulations": [
        {
            "creds": {
                "host": "0.0.0.0",
                "port": "12345",
                "imei": "1234567890123456789",
                "password": "67890"
            },
            "netCfg": {
                "maxTimeout": 5,
                "maxReconnectTries": 5,
                "maxResendTries": 5,
                "resendWaitTime": 5
            },
            "protocolType": "wialonips",
            "hasControlChannel": false,
            "channelCfg": {
                "jitter": 25,
                "delay": 50,
                "lossPercent": 0.3,
                "connBreakPercent": 0.1,
                "maxLossCount": 5,
                "maxConnBreakCount": 5
            },
            "autoStartPos": false,
            "startPos": {
                "x": 51.546529,
                "y": 46.037097
            },
            "autoEndPos": false,
            "endPos": {
                "x": 51.536257,
                "y": 46.02405
            },
            "boundingBoxes": [
                {
                    "_comment": "Moscow",
                    "minLat": 55.49,
                    "maxLat": 56.02,
                    "minLon": 36.80,
                    "maxLon": 38.10
                }
            ],
            "sats": 1,
            "stepDistanceM": 100,
            "stepMillis": 500,
            "maxDist": 200,
            "minDist": 100
        }
    ]
}
```

### Параметры конфига
```simulations``` - массив структур SimulationConfig
- ```creds``` - данные для подключения к серверу
- ```netCfg``` - параметры обслуживания TCP-соедиенения
    - ```maxTimeout``` - таймаут на I/O (с)
    - ```maxReconnectTries``` - максимальное количество попыток переподключения
    - ```maxResendTries``` - максимальное количество попыток переотправления пакета
    - ```resendWaitTime``` - задержка при отправке (мс)
- ```protocolType``` - тип протокола (возможные значения: ```wialonips```, ```egts```)
- ```hasControlChannel``` - флаг для включения управляемого канала
- ```channelCfg``` - параметры управляемого канала
- ```autoStartPos``` - флаг для автоматической генерации стартовой позиции
- ```startPos``` - стартовая позиция
- ```autoEndPos``` - флаг для автоматической генерации конечной позиции
- ```endPos``` - конечная позиция
- ```boundingBoxes``` - коллекция координатных границ
- ```sats``` - количество спутников
- ```stepDistanceM``` - шаг интерполяции маршрута (м)
- ```stepMillis``` - шаг времени (мс)
- ```minDist``` - минимальное расстояние
- ```maxDist``` - максимальное расстояние

## Архитектура проекта
Архитектурой проекта является "Слоистая архитектура" (Layered Architecture)

## Переменные окружения для тестов

- EMULATOR_TEST_UNIT_ID
- EMULATOR_TEST_UNIT_ID_2
- EMULATOR_TEST_PASSWORD
- EMULATOR_TEST_HOST
- EMULATOR_TEST_PORT