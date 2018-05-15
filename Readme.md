# สร้าง Microservice สำหรับ Authen 

> ขอบคุณบทความดีๆจาก [acoshift](https://acoshift.me/token-based-authentication-golang-307f18b03dc3)

## 1. ใช้ Datastore เก็บข้อมูล

> จะต้องเข้าไปสร้างโปรเจคก่อนที่ https://console.cloud.google.com

- เสร็จแล้วไปที่ IAM & Admin เพื่อสร้าง Service Account สำหรับ connect ไป Datastore

- แล้วก็ไปที่ Service Account > Create Service Account ใส่ชื่อ account และ เปลี่ยน Role เป็น Project > Editor นะครับ อย่าลืมกด Furnish a new private key ด้วย

## 2. Generate RSA256 เพื่อใช้ในการ Sign Token

```bash
# สำหรับ Mac/Linux
$ openssl genrsa -out key.rsa 1024
$ openssl rsa -in key.rsa -pubout > key.pub
```

- เราจะได้ไฟล์มา 2 ไฟล์ คือ

- key.rsa เป็น private key สำหรับ Sign Token

- key.pub เป็น public key สำหรับให้ Verify Token

- เราจะให้แค่ไฟล์ key.pub กับ Microservice เพื่อใช้ในการ Verify Token แต่จะไม่ให้ key.rsa เพราะเราไม่ต้องการให้ Microservice สร้าง Token

## 3. ตั้งค่าภูมิภาค
- https://console.cloud.google.com/datastore/setup?project=golang-jwt-auth

## 4. เมื่อ register มันจะสร้าง Entity เพื่อเก็บข้อมูล user และ token ให้
- User 
- Token
- https://console.cloud.google.com/datastore/entities/query?project=golang-jwt-auth&ns=&kind=User

## 5. Register
- http://localhost:9000/auth/register
- Header
```
Content-Type: application/json
```
- Body
```
{
	"username": "prongbang",
	"password": "test"
}
```
- Response 
```
Header: Status 201
Body: Created
```

## 6. ขอ Refresh Token
- http://localhost:9000/auth
- Header
```
Content-Type: application/json
```
- Body
```
{
	"grant_type": "password",
	"username": "prongbang",
	"password": "test"
}
```
- Response 
```
Header: Status 200
Body: {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTYzOTQ0NTYwNDcyODgzMiwidHlwZSI6MiwiZXhwIjoxNTA3NjM3MDEyLCJpYXQiOjE1MDc2MzY3MTJ9.kUQKXAvJvkAfktylnff2D2qds0k_GgWfZvgKUS992lgvLrVTFzGyfBdxBpk4vSxfdFbIgSzAU91oeb7FLuJYR9Tk8I3EiZ-FdY6KyR0mIL2l6fL8nP1LbisSQdLzhZRxofQgMjj1az_T2B3mMWawHdGr0ks4zVLTyE7FFTPM4Ww",
    "token_type": "bearer",
    "expires_in": 300,
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTYzOTQ0NTYwNDcyODgzMiwidHlwZSI6MSwiaWF0IjoxNTA3NjM2NzEyfQ.rJay4ozGudduUWKakGP3GOp5BUfzsjUb8GzrwoIvk6qFZtRO7B1dknKKnheXRJ4WCcscYdrTiR4FI_snnSfTZSAq_TfArzjx90mbfFEVu5L1SlR50V97L1h4W3LsDBoLC8MJMPOXi1QpEPIIpXTBOE7qThfnWTs0nG9Tga9kHyY",
    "uid": 5639445604728832
}
```

## 7. เมื่อ Access Token หมดอายุ เราก็ต้องส่ง Refresh Token ไปขอ Access Token อันใหม่
- http://localhost:9000/auth
- Header
```
Content-Type: application/json
```
- Body
```
{
	"grant_type": "refresh_token",
	"refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTYzOTQ0NTYwNDcyODgzMiwidHlwZSI6MSwiaWF0IjoxNTA3NjM2NzEyfQ.rJay4ozGudduUWKakGP3GOp5BUfzsjUb8GzrwoIvk6qFZtRO7B1dknKKnheXRJ4WCcscYdrTiR4FI_snnSfTZSAq_TfArzjx90mbfFEVu5L1SlR50V97L1h4W3LsDBoLC8MJMPOXi1QpEPIIpXTBOE7qThfnWTs0nG9Tga9kHyY"
}
```
- Response 
```
Header: Status 200
Body: {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6NTYzOTQ0NTYwNDcyODgzMiwidHlwZSI6MiwiZXhwIjoxNTA3NjM3MzMyLCJpYXQiOjE1MDc2MzcwMzJ9.fFkUYFuvPqkfKuvzVjTLZ3z3BNpMPGAcBDBqCX_b1ZbJ-dWC_8jr6vqU2mJ-SfYirNrjD1SpZ7_qS_tq1l_M490-EcwN0hgekhHetzm8X73YqsrXcObPP1Cb1x8RQfX_CPyzOs1oG13GhVshXvp4yPw4jaMp-45LXGP3M9xqrlg",
    "token_type": "bearer",
    "expires_in": 300,
    "uid": 5639445604728832
}
```

## 8. Revoke Refresh Token
- http://localhost:9000/auth/revoke
- Header
```
Authorization: Bearer ACCESS_TOKEN
Content-Type: application/json
```
- Body
```
{
	"username": "prongbang",
	"password": "test"
}
```
- Response 
```
Header: Status 201
Body: Created
```