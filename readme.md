## tempfiles-backend

frontend - https://github.com/tempfiles-Team/tempfiles-frontend

tempfiles-backend is a backend for tempfiles-frontend.

## How to run - docker

### 1. build docker image

```bash
docker build -t tempfiles-backend .
```

or pull image from docker hub

```bash
docker pull minpeter/tempfiles-backend
```

### 2. run docker container

```bash
docker run -dp 5000:5000 \
  -e BACKEND_BASEURL=http://localhost:5000 \
  -e JWT_SECRET=<your secret> \
  -e DB_TYPE=sqlite \
  -v $(pwd)/backend-data:/tmp \
  tempfiles-backend
```

## How to run - local

1. config .env file

### nessary

|       key       |         value         |   description    |
| :-------------: | :-------------------: | :--------------: |
| BACKEND_BASEURL | http://localhost:5000 | backend base url |
|   JWT_SECRET    |     <your secret>     |    jwt secret    |
|     DB_TYPE     |  sqlite or postgres   |    select db     |

### optional

|     key      |   value   |                  description                   |
| :----------: | :-------: | :--------------------------------------------: |
| BACKEND_PORT |   5000    |     If you want to change the backend port     |
|   DB_HOST    | localhost | If postgres is selected, its db ip or hostname |
|   DB_PORT    |   5432    |      If postgres is selected, its db port      |
|   DB_NAME    |  tempdb   |   If postgres is selected, its db table name   |
|   DB_USER    |  tempdb   |   If postgres is selected, its db user name    |
| DB_PASSWORD  |  tempdb   | If postgres is selected, its db user password  |

2. run server

```bash
go run .
```

## test server

https://api.tempfiles.ml

- 파일 목록 조회 [Get]
  https://api.tempfiles.ml/list

- 파일 업로드 [Post]
  multipart/form-data "file" 필드에 업로드할 파일을 넣어서 요청
  https://api.tempfiles.ml/upload

- 파일 다운로드 [Get]
  https://api.tempfiles.ml/dl/(file_id)

- 파일 삭제 [Delete]
  https://api.tempfiles.ml/del/(file_id)
