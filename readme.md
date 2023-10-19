jwt go package - to create entire authentication logic
gin-gonic - to craete api's
mongoDB - database
log in the cmd
validate the request from postman

code . - open visual studio code in current folder
go mod init github.com/23piyush/golang-jwt-project   ====similar to====== npm init
go get github.com/gin-gonic/gin    ====similar to======= npm install

What does this mean?
warning: in the working copy of 'go.sum', LF will be replaced by CRLF the next time Git touches it

Project flow: models => main.go => routes folder => database => controllers

This project structure is best for monilithic projects. It can be scaled to large number of users.
We will see how to create microservice project structure when we learn microservices using GRPC.

