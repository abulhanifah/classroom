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
GET /api/get_class_list
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
}
```

- [x] Get Class Room By Id
GET /api/get_class_by_id/:id
headers{Authorization:Bearer token}
Token is from authorization
Id is class room id
response
```
{
    "available_seats": [
        "1A",
        "1B",
        "1C",
        "1D",
        "2A",
        "2B",
        "2C",
        "2D",
        "3A",
        "3B",
        "3C",
        "3D",
        "4A",
        "4B",
        "4C",
        "4D",
        "5A",
        "5B",
        "5C",
        "5D"
    ],
    "class_id": 2,
    "columns": 4,
    "occupied_seats": [],
    "rows": 5,
    "teacher": "out"
}
```

- [x] Create class
POST /api/create_class
headers{Authorization:Bearer token}
Token is from authorization
body
```
{
    "name": "Class D",
    "rows": 5,
    "columns": 4
}
```
response
```
{
    "available_seats": [
        "1A",
        "1B",
        "1C",
        "1D",
        "2A",
        "2B",
        "2C",
        "2D",
        "3A",
        "3B",
        "3C",
        "3D",
        "4A",
        "4B",
        "4C",
        "4D",
        "5A",
        "5B",
        "5C",
        "5D"
    ],
    "class_id": 4,
    "columns": 4,
    "occupied_seats": [],
    "rows": 5,
    "teacher": "out"
}
```
- [x] Check In Class
POST /api/check_in
headers{Authorization:Bearer token}
Token is from authorization
body
```
{
    "class_id": 4,
    "user_id": "0edd2b6c-87cd-4007-a74e-5fcf86437bf5"
}
```
response
```
{
    "available_seats": [
        "1B",
        "1C",
        "1D",
        "2A",
        "2B",
        "2C",
        "2D",
        "3A",
        "3B",
        "3C",
        "3D",
        "4A",
        "4B",
        "4C",
        "4D",
        "5A",
        "5B",
        "5C",
        "5D"
    ],
    "class_id": 4,
    "columns": 4,
    "message": "Hi Naruto, your seat is 1A",
    "occupied_seats": [
      {
        "seat":"1A",
        "student_name":"Naruto"
      }
    ],
    "rows": 5,
    "teacher": "out"
}
```

- [x] Check Out Class
POST /api/check_in
headers{Authorization:Bearer token}
Token is from authorization
body
```
{
    "class_id": 4,
    "user_id": "0edd2b6c-87cd-4007-a74e-5fcf86437bf5"
}
```
response
```
{
    "available_seats": [
        "1A",
        "1B",
        "1C",
        "1D",
        "2A",
        "2B",
        "2C",
        "2D",
        "3A",
        "3B",
        "3C",
        "3D",
        "4A",
        "4B",
        "4C",
        "4D",
        "5A",
        "5B",
        "5C",
        "5D"
    ],
    "class_id": 4,
    "columns": 4,
    "message": "Hi Naruto, 1A is now available for other students",
    "occupied_seats": [],
    "rows": 5,
    "teacher": "out"
}
```
