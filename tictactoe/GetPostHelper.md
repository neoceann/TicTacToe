Регистрация:

curl.exe -X POST http://localhost:8080/auth/signup -H "Content-Type: application/json" -d "{\`"login\`":\`"player1\`",\`"password\`":\`"123\`"}"


Авторизация:

curl.exe -X POST http://localhost:8080/auth/signin -u "player1:123"


Создать игру:

curl.exe -X POST http://localhost:8080/game -u "player1:123" -d "{\`"opponent\`":\`"Human\`"}"


Подключиться:

curl.exe -X POST http://localhost:8080/game/join/123 -u "player2:123"


Сделать ход:

curl.exe -X POST http://localhost:8080/game/123 -u "player1:123" -H "Content-Type: application/json" -d "{\`"field\`":[[1, 0, 0],[0, 0, 0],[0, 0, 0]]}"


Получить инфо об игре:

curl.exe -X GET http://localhost:8080/game/123 -u "player1:123"


Получить инфо об играх в ожидании:

curl.exe -X GET http://localhost:8080/waiting -u "player1:123"


Получить инфо о пользователе:

curl.exe -X GET http://localhost:8080/user/25f61bdf-5fd9-434d-a772-fe4248307fee -u "player1:123"