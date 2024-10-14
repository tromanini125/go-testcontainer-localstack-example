#Create docker container
docker run --rm -it -p 4566:4566 -p 4510-4559:4510-4559 localstack/localstack

export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"
export AWS_DEFAULT_REGION="us-east-1"

##Criar fila
aws --endpoint-url=http://localhost:4566 awslocal sqs create-queue --queue-name test-queue

##Listar filas
aws --endpoint-url=http://localhost:4566 sqs list-queues

##postar evento na fila
aws --endpoint-url=http://localhost:4566 sqs send-message --queue-url http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/test-queue --message-body "{ \"cardId\": 1, \"cardHolderName\": \"Thiago\", \"cardNumber\": \"1234321056789876\", \"cvv\": \"123\", \"expiryDate\": \"2024-01-30\" }"