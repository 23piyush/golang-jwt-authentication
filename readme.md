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

To run the project : 
1. go mod tidy => this command will install all the packages that we included but not installed
2. go run main.go



":=" symbol is used for short variable declaration and assignment. It allows you to declare and initialize a variable in a concise manner.

Golang gives very accurate errors. Really helpful in debugging.
string is immutable. To change them, use bytes