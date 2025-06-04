## Curl command

```
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com", "password":"password123"}' | jq -r .token)

curl -X GET "http://localhost:8080/api/users?page=1&limit=5" \
  -H "Authorization: Bearer $TOKEN"

curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN"

curl -X PUT http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Name","email":"newemail@example.com","password":"newpassword123"}'

curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer $TOKEN"

curl -X GET http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer $TOKEN"
```
