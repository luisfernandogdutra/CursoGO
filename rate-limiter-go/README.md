# Rate Limiter Go

Este projeto implementa um middleware rate limiter em Go que pode limitar requisições por IP ou token de acesso. Ele usa o Redis para persistência e suporta configuração via variáveis de ambiente ou arquivo `.env`.

## Como rodar

### 1. Subir o Redis com Docker

Execute o seguinte comando para iniciar o Redis usando Docker Compose:

```bash
docker-compose up -d
