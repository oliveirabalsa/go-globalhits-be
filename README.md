## Back-End Go Globalhitss

### ℹ️ Sobre

Este é o back-end do challenge da GlobalHitss, desenvolvida em Go, utilizando o framework Chi para roteamento, Swagger para documentação da API, RabbitMQ para mensageria, e Gorm como ORM para PostgreSQL.

### 🚀 Como Iniciar

#### Pré-requisitos
Certifique-se de ter o Go instalado em sua máquina

1. Renomeie ou copie o arquivo `env.example` para `.env` e preencha as variáveis de ambiente necessárias:

   ```plaintext
   PORT=8082 
   
   POSTGRES_HOST=localhost
   POSTGRES_PORT=5432
   POSTGRES_USER=globalhitss
   POSTGRES_PASSWORD=globalhitss
   POSTGRES_DB=globalhitss

   RABBITMQ_USER=globalhitss
   RABBITMQ_PASSWORD=globalhitss
   RABBITMQ_HOST=localhost
   RABBITMQ_PORT=5672
   RABBITMQ_API_PORT=15672
   RABBITMQ_QUEUE=globalhitss
   ```
2. Execute o serviço docker:   
   
```bash
docker-compose up -d
```

Isso iniciará os serviços do RabbitMQ e do PostgreSQL em contêineres Docker, conforme configurado no arquivo `docker-compose.yml`.

3. Execute o script `start.sh` para iniciar o worker, criar a fila e iniciar a aplicação:

   ```bash
   sh ./start.sh
   ```

### Testando a Aplicação

Acesse o Swagger em `http://localhost:8082/swagger/index.html` testar a API ou você pode encontrar na pasta postman um json para importação.

### Informações Adicionais

- As requisições de criação, atualização e exclusão são encaminhadas através de filas RabbitMQ, enquanto as requisições de obtenção interagem diretamente com o banco de dados.
- Os dados são criptografados na inserção e descriptografados na seleção do banco de dados para segurança do usuário.
  
---

#### Ações que eu faria na aplicação a partir dessa versão:


- Filtros de busca
- Paginação
- Testes de carga

---

### 📝 Licença

Este projeto está licenciado sob a Licença MIT. Consulte o arquivo [LICENSE](LICENSE) para mais informações.

