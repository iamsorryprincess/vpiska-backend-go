# vpiska-backend
[![Deploy](https://github.com/iamsorryprincess/vpiska-backend-go/workflows/build-deploy/badge.svg)](https://github.com/iamsorryprincess/vpiska-backend-go/actions)
____

## WebSocket эвента
____

Url для подключения - **wss://vp1ska.ru/api/v1/websockets/event?accessToken=qweasd&eventId=E9F6D9A2-2FF4-4A15-96EB-7C13F47F9CA8**    
access_token - jwt token юзера, если он есть, если нет (юзер не зареган), то не передавать параметр токена    
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
