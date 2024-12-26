1 - docker-compose up --build
2- curl -X POST http://localhost:8080/cep -d '{"cep": "00000000"}' alterando para o cep desejado
3- Acesse Zipkin em http://localhost:9411