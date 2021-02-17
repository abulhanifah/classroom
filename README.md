# classroom
This is a simple program to book a class on simple class management using go language

# Initialization
Clone the repository
```
git clone https://github.com/abulhanifah/classroom.git
```
## Setup Environment
Please copy .env.sample to .env
```
DB_IS_DEBUG={set true if you want to debug query and false to disable debug}
APP_ENV={application environment, like local,development, production}
APP_PORT={application port}
APP_URL={application url}
APP_KEY={application encryption key}
OAUTH2_ACCESS_EXPIRE_IN={oauth token expiration, ex:1h, 3h or 30m}
OAUTH2_REFRESH_EXPIRE_IN={oauth token expiration, ex:720h 1d}
DB_DRIVER={database driver, supported mysql & postgres}
DB_HOST=localhost
DB_PORT=3306
DB_NAME=declassroom
DB_USER=root
DB_PASSWORD=root
```
## Run Application
```
go run main.go
```
## Endpoints
- [x] Authorization
GET /api/login
headers{Authorization:Bearer token}
Token is base64 encode "email:password" or "refresh_token:thetoken"
response
```
{
  "token":"the encode oauth token"
}
```
- [x] Get Class Room Lists
GET /api/login
headers{Authorization:Bearer token}
Token is from authorization
response
```
{
    "count": 3,
    "data": [
        {
            "class_id": 1,
            "class_name": "Class A",
            "columns": 4,
            "rows": 5
        },
        {
            "class_id": 2,
            "class_name": "Class B",
            "columns": 4,
            "rows": 5
        },
        {
            "class_id": 3,
            "class_name": "Class C",
            "columns": 4,
            "rows": 5
        }
    ]
}```
- [x] Get Class Room By Id
- [x] Check In Class
- [x] Check Out Class
