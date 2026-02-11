Пример запросов для Windows.


* Регистрация:

curl.exe -X POST http://localhost:8080/auth/signup -H "Content-Type: application/json" -d "{\`"login\`":\`"player1\`",\`"password\`":\`"123\`"}"

Будет создан пользователь "player1" с паролем "123"


* Авторизация:

curl.exe -X POST http://localhost:8080/auth/signin -u "player1:123"


* Создать игру:

curl.exe -X POST http://localhost:8080/game -u "player1:123" -d "{\`"opponent\`":\`"Human\`"}"

"Human" - для игры "против игрока"
"AI" - для игры против компьютера.


> {gameID} и {userID} отображаются в теле ответа.


* Подключиться к созданной игре (в случае игры против игрока):

curl.exe -X POST http://localhost:8080/game/join/{gameID} -u "player2:123"


* Сделать ход:

curl.exe -X POST http://localhost:8080/game/{gameID} -u "player1:123" -H "Content-Type: application/json" -d "{\`"field\`":[[1, 0, 0],[0, 0, 0],[0, 0, 0]]}"


* Получить инфо об игре:

curl.exe -X GET http://localhost:8080/game/{gameID} -u "player1:123"


* Получить инфо об играх в ожидании (созданные игры "против игрока", но второй игрок еще не подключился):

curl.exe -X GET http://localhost:8080/waiting -u "player1:123"


* Получить инфо о пользователе:

curl.exe -X GET http://localhost:8080/user/{userID} -u "player1:123"
