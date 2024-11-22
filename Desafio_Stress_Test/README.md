Construção e Execução do Docker:

Construir a imagem Docker:
docker build -t stresstest .

Executar a imagem Docker com os parâmetros necessários:
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10

O comando 
docker run stresstest --url=http://google.com --requests=1000 --concurrency=10 
executa o teste de carga na URL http://google.com, realizando 1000 requisições com uma concorrência de 10 requisições simultâneas.

O relatório final será exibido no terminal após a execução, incluindo informações sobre o tempo total, sucesso das requisições e distribuição de status HTTP.