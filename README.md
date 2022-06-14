#### Описание
    Создать сервис, который обращается к сервису курсов валют, и отображает gif:
    если курс по отношению к USD за сегодня стал выше вчерашнего, 
    то отдаем рандомную отсюда https://giphy.com/search/rich
    если ниже - отсюда https://giphy.com/search/broke
    Ссылки
    REST API курсов валют - https://docs.openexchangerates.org/
    REST API гифок - https://developers.giphy.com/docs/api#quick-start-guide
    Must Have
    Сервис на Spring Boot 2 + Java / Kotlin
    Запросы приходят на HTTP endpoint (должен быть написан в соответствии с rest conventions), 
    туда передается код валюты по отношению с которой сравнивается USD
    Для взаимодействия с внешними сервисами используется Feign
    Все параметры (валюта по отношению к которой смотрится курс, адреса внешних сервисов и т.д.) 
    вынесены в настройки
    На сервис написаны тесты (для мока внешних сервисов можно использовать @mockbean или WireMock)
    Для сборки должен использоваться Gradle
    Результатом выполнения должен быть репо на GitHub с инструкцией по запуску
    Nice to Have
    Сборка и запуск Docker контейнера с этим сервисом
    
#### Endpoints
    Список доступных валют: GET /api/currencies
    Получить gif: GET /api/rate{code}
    
#### Запуск
    jar: 
    В корне проекта выполнить 
    java -jar currencies-0.0.1-SNAPSHOT.jar
    
    Docker:
    Для создания образа в корне проекта выполнить 
    docker build -t currencies .
    Для запуска созданного контейнера выполнить
    docker run -p 8080:8080 currencies
