# Регистрация:

curl.exe -X POST http://localhost:8080/auth/signup -H "Content-Type: application/json" -d "{\`"login\`":\`"player1\`",\`"password\`":\`"123\`"}"


# Авторизация:

curl.exe -X POST http://localhost:8080/auth/signin -d "{\`"login\`":\`"player1\`",\`"password\`":\`"123\`"}"


# Создать игру:

curl.exe -X POST http://localhost:8080/game -H "Authorization: Bearer 777" -H "Content-Type: application/json" -d "{\`"opponent\`":\`"Human\`"}"


# Подключиться:

curl.exe -X POST http://localhost:8080/game/join/123 -H "Authorization: Bearer 777"


# Сделать ход:

curl.exe -X POST http://localhost:8080/game/123 -H "Authorization: Bearer 777" -H "Content-Type: application/json" -d "{\`"field\`":[[1, 0, 0],[0, 0, 0],[0, 0, 0]]}"


# Получить инфо об игре:

curl.exe -X GET http://localhost:8080/game/123 -H "Authorization: Bearer 777"


# Получить инфо об играх в ожидании:

curl.exe -X GET http://localhost:8080/waiting -H "Authorization: Bearer token"


# Получить инфо о пользователе по ID:

curl.exe -X GET http://localhost:8080/user/25f61bdf-5fd9-434d-a772-fe4248307fee -H "Authorization: Bearer 777"


# Обновить access token (используя refresh token):

curl.exe -X POST http://localhost:8080/auth/refresh_access -H "Content-Type: application/json" -d "{\`"refresh_token\`":\`"777\`"}"


# Обновить все токены (используя refresh token):

curl.exe -X POST http://localhost:8080/auth/refresh_tokens -H "Content-Type: application/json" -d "{\`"refresh_token\`":\`"777\`"}"


# Получить инфо о пользователе по токену:

curl.exe -X GET http://localhost:8080/user/info_by_access_token -H "Authorization: Bearer 777"


# Получить историю игр по токену:

curl.exe -X GET http://localhost:8080/history -H "Authorization: Bearer 777"


# Получить список лидеров:

curl.exe -X GET http://localhost:8080/leaderboard?limit=10 -H "Authorization: Bearer 777"