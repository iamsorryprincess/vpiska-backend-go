# vpiska-backend
[![Deploy](https://github.com/iamsorryprincess/vpiska-backend-go/workflows/build-deploy/badge.svg)](https://github.com/iamsorryprincess/vpiska-backend-go/actions)
### swagger
https://vp1ska.ru/swagger/index.html
____
### configuration
**Переменные окружения**:</br>
SERVER_PORT - порт</br>
DB_CONNECTION - строка подключения к бд</br>
DB_NAME - имя бд</br>
JWT_KEY - ключ шифрования jwt токена</br>
JWT_ISSUER - издатель jwt токена</br>
JWT_AUDIENCE - клиент jwt токена</br>
JWT_LIFETIME_DAYS - кол-во дней через которое jwt токен становится невалидным</br>
HASH_KEY - ключ хэширования для паролей пользователей</br>
LOGGING_TRACE_REQUESTS - булевый флаг, указывающий логировать ли http запросы и ответы</br>

Инфраструктуру для дебага можно поднять в докере командой make infrastructure</br>
Собрать сервис в образ докера командой make build</br>
Поднять сервис вместе со всей необходимой инфраструктурой make run</br>
____

## WebSocket эвента
____

Url для подключения - **wss://vp1ska.ru/api/v1/websockets/event?accessToken=qweasd&eventId=E9F6D9A2-2FF4-4A15-96EB-7C13F47F9CA8**    
accessToken - jwt token юзера, если он есть, если токен пуст, то в ответ получите 401  
eventId - id эвента    

### сообщения приходящие с бэка
____

**eventUpdated/{"eventId":"c5208ba0-17fa-4aab-b627-4b8ccc6060bb","name":"string","address":"string","usersCount":1,"coordinates":{"x":0,"y":0}}**: 
обновление евентеа, после слэша json
____

**chatMessage/{"userId": "qweasd", "userName": "qweasd", "userImageId": "qweasd", "message": "lol ahaha"}**: 
сообщение в чате, после слэша json
____

**mediaAdded/E9F6D9A2-2FF4-4A15-96EB-7C13F47F9CA8**: 
был добавлен медиа контент (фото или видео), после слэша идет id медиа контента в firebase
____

**mediaRemoved/E9F6D9A2-2FF4-4A15-96EB-7C13F47F9CA8**: 
был удален медиа контент, после слэша идет id медиа контента в firebase
____

**closeEvent/**: 
эвент был закрыт, после слэша ничего нет, чисто нотификация что эвент был закрыт, после этого бэк разрывает соединение
____

### сообщения для отправки в бэк
____

**chatMessage/message**: 
сообщение в чат где message - строка
____
