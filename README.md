### Имитатор терминалов ГЛОНАСС/GPS и анализ канала передачи телеметрии

# Аргументы CLI
- ```cfg``` - путь к файлу конфига (значение по умолчанию: ```рабочая директория```)

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
            "sats": 1,
            "stepDistanceM": 100,
            "stepMillis": 500
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
- ```sats``` - количество спутников
- ```stepDistanceM``` - шаг интерполяции маршрута (м)
- ```stepMillis``` - шаг времени (мс)

## Архитектура проекта
Архитектурой проекта является "Слоистая архитектура" (Layered Architecture)

