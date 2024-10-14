

##Listar filas
aws --endpoint-url=http://localhost:4566 sqs list-queues

##postar evento na fila
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/test-queue --message-body "{ \"cardId\": 1, \"cardHolderName\": \"Thiago\", \"cardNumber\": \"1234321056789876\", \"cvv\": \"123\", \"expiryDate\": \"2024-01-30\" }"