fetch dependencies from go.mod :
- go mod download

run server :
- go run server.go

build docker images :
- docker build -t application-name:tag-name .

run docker :
- docker run -p 8080 -d application-name

api register user :
curl --location --request POST '127.0.0.1:8080/api/auth/register' \
--header 'Content-Type: application/json' \
--data-raw '{
"name": "Tri Budiyono",
"email": "tribudiyono93@gmail.com",
"password": "blah123"
}'

api login user :
curl --location --request POST '127.0.0.1:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--data-raw '{
"email": "tribudiyono93@gmail.com",
"password": "blah123"
}'

api get books :
curl --location --request GET '127.0.0.1:8080/api/books' \
--header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cCI6MTY0NjI0MTM3NywiaWF0IjoxNjE0NzA1Mzc3LCJpc3MiOiJkZWZhdWx0SXNzdWVyIn0.p43MNWhvRXNHsWl3Ne9y4c-ivFfDdQwmulJA3piPqJ8'

