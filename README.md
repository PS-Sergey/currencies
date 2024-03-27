# currencies

## Описание
Cервис котировок валютных курсов.  
Сервис предоставляет асинхронный интерфейс, т.е. пользователь сначала
выполняет запрос на обновление котировки, а затем, через некоторое время,
запрос на получение котировки. При этом непосредственно обновление
котировки происходит в фоновом режиме.  
За значениями цены сервис ходит в [ExchangeRate-API](https://www.exchangerate-api.com/).  
В качестве базы данных используется PostgreSQL.  
Обновление котировки в фоновом режиме реализовано по триггеру через канал.  
Сервисный слой покрыт Unit-тестами.
    
## Endpoints
Обновить котировку: POST /currency/v1/rate?pair={валютная_пара}  
В query param pair указывается код валютной пары (напр. EUR/MXN ).   
Доступные валюты: USD, EUR, MXN.  
В ответ вернется код идентификатора обновления в формате UUID ("currency_rate_id").

Получить котировку по идентификатору: GET /currency/v1/rate/{id}  
В запросе указывается идентификатор обновления.   
В ответе сервис отдает значение цены, время обновления и статус запроса.  

Получить последнее значение котировки: GET /currency/v1/rate/{base}/{target}  
В запросе указываются коды валют.   
В ответе сервис отдает значение цены, время обновления и статус последнего успешного запроса.
    
## Запуск
Для запуска сервиса с окружением (postgres, migrator, swagger-ui) в Docker, выполнить команду:

    docker-compose up

в корне проекта.  
При добавлении конфига, указать путь до него в переменной окружения `CONFIG_PATH`, 
дефолтный путь поиска конфига - `./config/config.yml`.  
Сервис доступен на порту 8080, swagger на порту 8085.
    
    
## Примечания
Запросы через swagger-ui доходят до сервиса и успешно выполняются, но в ui отображатеся ошибка "Failed to fetch.", не успел разобраться в причинах такого поведения.
