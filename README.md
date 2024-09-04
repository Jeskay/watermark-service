# Сервис для добавления вотермарок на изображения и организации доступа к сгенерированным картинкам
Является улучшенной реализацией проекта из статьи https://www.velotio.com/engineering-blog/build-a-containerized-microservice-in-golang
## Параметры запуска сервисов
### Watermark Service
аргументы
```
 --config "путь к yaml файлу"
 --log "путь для выгрузки логов"
```
переменные среды
```
 HTTP_PORT
 HTTP_HOST
 GRPC_PORT
 GRPC_HOST
 DB_HOST
 DB_PORT
 DB_USER
 DB_PASSWORD
 DB_DATABASE
 CLOUDINARY_CLOUD - параметры облачного хранилища cloudinary
 CLOUDINARY_API
 CLOUDINARY_SECRET
 AUTH_PORT - порт сервиса аутентикации
 AUTH_HOST
 PICTURE_PORT - порт сервиса обработки изображений
 PICTURE_HOST
 JAEGER_PORT - порт трейсинг платформы jaeger
 JAEGER_HOST
```
### Authentication Service
аргументы
```
 --config "путь к yaml файлу"
 --log "путь для выгрузки логов"
```
переменные среды
```
HTTP_PORT
HTTP_HOST
GRPC_PORT
GRPC_HOST
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_DATABASE
SECRET_KEY
```
### Picture Service
аргументы
```
 --config "путь к yaml файлу"
 --log "путь для выгрузки логов"
```
переменные среды
```
HTTP_PORT
HTTP_HOST
GRPC_PORT
GRPC_HOST
JAEGER_PORT - порт трейсинг платформы jaeger
JAEGER_HOST
```
