# zimplebank


## Resources used in this project:

### Database Related:
* #### DB: PostgreSQL 16 docker image
* #### DB Design Tool: [dbdiagram.io](https://dbdiagram.io/)
* #### DB Design: https://dbdiagram.io/d/Zimple_Bank-6563f8823be1495787c588f4
* #### DB Docs Generation Tool: [dbdocs.io](https://dbdocs.io/)
* #### DB Docs: https://dbdocs.io/ZhangZhihuiAAA/zimple_bank
* #### DB SQL Code Generation Tool: [dbml2sql](https://dbml.dbdiagram.io/cli)
* #### SQL GUI Client: Tableplus(https://tableplus.com)
* #### CRUD Go Code Generation Tool: SQLC(https://sqlc.dev)

### Non-Standard Library Go modules/packages:
* #### PostgreSQL Driver: github.com/jackc/pgx/v5
* #### DB Migration: github.com/golang-migrate/migrate/v4
* #### Unit Test: github.com/stretchr/testify/require
* #### Unit Test: github.com/uber-go/mock
* #### Web Framework: github.com/gin-gonic/gin
* #### Config Management: github.com/spf13/viper
* #### Validator: github.com/go-playground/validator/v10
* #### UUID: github.com/google/uuid
* #### Token Maker: github.com/o1egl/paseto
* #### GRPC: google.golang.org/grpc
* #### Protobuf: google.golang.org/protobuf
* #### GRPC Gateway: github.com/grpc-ecosystem/grpc-gateway
* #### Google APIs: github.com/googleapis/googleapis (needed by generating reverse-proxy using protoc-gen-grpc-gateway)
* #### Static files to binary: github.com/rakyll/statik
* #### Logging: github.com/rs/zerolog
* #### Sending Email: github.com/jordan-wright/email

### CI/CD:
* #### Github Actions
* #### Github Marketplace Actions:
    * Amazon ECR "Login" Action for GitHub Actions
    * Kubectl tool installer

### Cloud services:
* #### AWS Free Tier Account
* #### AWS Identity and Access Management (IAM)
* #### AWS Elastic Container Registery (ECR)
* #### AWS Relational Database Service (RDS)
* #### AWS Secrets Manager
* #### AWS Key Management Service (KMS)
* #### AWS Elastic Kubernetes Service (EKS)
* #### AWS Elastic Compute Cloud (EC2)
* #### AWS CloudWatch
* #### AWS Route 53
* #### AWS Simple Email Service (SES)
* #### AWS ElastiCache

### Others:
* #### VSCode
* #### Makefile
* #### migrate
* #### SQLC
* #### Postman
* #### Docker
* #### AWS CLI
* #### kubectl
* #### k9s
* #### Kubernetes Ingress
* #### Kubernetes Add-on: cert-manager
* #### protoc (Protobuf Complier)
* #### GRPC Gateway (a plugin of protoc)
* #### Evans (a GRPC client tool)
* #### Swagger (https://swagger.io)
    * Display API Docs on SwaggerHub: (https://app.swaggerhub.com/apis-docs/ZHANGZHIHUIAAA/zimple-bank/1.2)
* #### SwaggerUI (https://github.com/swagger-api/swagger-ui)
    * Free tool to display API Docs on your own server
* #### statik
* #### Redis