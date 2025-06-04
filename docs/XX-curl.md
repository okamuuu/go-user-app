## 1. ユーザー一覧取得（GET /users）

```bash
curl -X GET http://localhost:8080/users
```

## 2. ユーザー登録（POST /users）

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "Test User",
    "Email": "testuser@example.com",
    "Password": "password123"
  }'
```

## 3. ユーザー更新（PUT /users/\:id）

`id` は更新したいユーザーの ID に置き換えてください。
例では `1` としています。

```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "ID": 1,
    "Name": "Updated User",
    "Email": "updateduser@example.com",
    "Password": "newpassword"
  }'
```

## 4. ユーザー削除（DELETE /users/\:id）

```bash
curl -X DELETE http://localhost:8080/users/1
```
