# GymScore Backend (Go)

Este projeto é a reestruturação completa do backend do **GymScore** — originalmente desenvolvido em Node.js — agora migrado para **Go (Golang)**. A nova versão adota Clean Architecture, o framework Fiber para alta performance e o GORM para persistência no MySQL, oferecendo uma base robusta, limpa e escalável.

## 🚀 Tecnologias Utilizadas

- **Linguagem**: [Go (Golang) 1.25+](https://go.dev/)
- **Framework Web**: [Fiber v2](https://gofiber.io/) (inspirado no Express.js, garantindo fácil transição para quem vem do Node.js)
- **Banco de Dados**: [MySQL 8+](https://www.mysql.com/)
- **ORM**: [GORM](https://gorm.io/) (com auto-migração ativada)
- **Autenticação**: JWT (JSON Web Tokens)
- **Segurança**: Bcrypt (Hash de senhas), CORS Middleware

## 📁 Estrutura do Projeto

O projeto foi organizado seguindo os princípios da **Clean Architecture** e o padrão padrão de layout de projetos em Go:

```text
gym-score-go/
├── database/
│   └── schema.sql               # DDL completo do MySQL (Tabelas e Procedures)
├── internal/
│   ├── config/                  # Configurações e conexão com DB (.env)
│   ├── controllers/             # Manipuladores HTTP (Handlers)
│   ├── middlewares/             # Interceptadores (Auth, CORS, Logger)
│   ├── models/                  # Entidades de domínio e DTOs
│   ├── repositories/            # Camada de acesso a dados (GORM)
│   ├── routes/                  # Definição das rotas da API
│   └── services/                # Regras de negócio (Use Cases)
├── pkg/
│   └── utils/                   # Utilitários compartilhados (JWT, Respostas, Validações)
├── .env.example                 # Exemplo de variáveis de ambiente
├── go.mod / go.sum              # Gerenciamento de dependências
├── main.go                      # Ponto de entrada da aplicação
└── README.md                    # Documentação do projeto
```

## ⚙️ Pré-requisitos

Para rodar este projeto localmente, você precisará de:

1. **Go** instalado (versão 1.21 ou superior recomendada)
2. **MySQL** rodando localmente ou em um container Docker
3. Opcional: **Make** para facilitar os comandos de build

## 🛠️ Como Configurar e Rodar Localmente

### 1. Configurar o Banco de Dados (MySQL)

Crie um banco de dados vazio no seu servidor MySQL:

```sql
CREATE DATABASE gymscore CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

> **Nota:** O GORM criará automaticamente as tabelas quando o projeto for iniciado (AutoMigrate). Caso prefira criar manualmente ou utilizar as Procedures originais, execute o script disponível em `database/schema.sql`.

### 2. Configurar Variáveis de Ambiente

Copie o arquivo de exemplo e crie o seu `.env`:

```bash
cp .env.example .env
```

Edite o arquivo `.env` com as suas credenciais do MySQL:

```env
APP_PORT=3000
APP_ENV=development
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=sua_senha_aqui
DB_NAME=gymscore
JWT_SECRET=uma_chave_muito_segura_e_longa
```

### 3. Instalar Dependências e Executar

Baixe as dependências do módulo Go:

```bash
go mod tidy
```

Inicie o servidor em modo de desenvolvimento:

```bash
go run cmd/api/main.go
```

O servidor estará rodando em: `http://localhost:3000`

### 4. Build para Produção

Para compilar o binário final:

```bash
go build -o bin/gymscore-api ./cmd/api/main.go
./bin/gymscore-api
```

## 🌐 Integração com Frontend

A API REST foi projetada para manter **compatibilidade máxima** com o frontend original. O CORS já está configurado para aceitar requisições de qualquer origem (`*`) em ambiente de desenvolvimento.

### Exemplos de Endpoints

**1. Criar Usuário (Público)**
- **POST** `/api/usuarios`
- **Body**:
  ```json
  {
    "nome": "João",
    "sobrenome": "Silva",
    "email": "joao@example.com",
    "senha": "senha_segura",
    "data_nascimento": "1995-03-15",
    "genero": "M"
  }
  ```

**2. Autenticação (Público)**
- **POST** `/api/login`
- **Body**:
  ```json
  {
    "email": "joao@example.com",
    "senha": "senha_segura"
  }
  ```
- **Retorno**: Retorna um token JWT que deve ser enviado no header `Authorization: Bearer <token>` nas próximas requisições.

**3. Criar Desafio (Protegido por JWT)**
- **POST** `/api/desafios`
- **Body**:
  ```json
  {
    "titulo": "Supino 100kg",
    "descricao": "Quem levanta mais repetições",
    "valor": 50.00,
    "local": "Academia Central",
    "id_criador": 1
  }
  ```

## 🧪 Testes

A lógica de negócios e as funções utilitárias possuem testes unitários. Para rodá-los:

```bash
go test ./... -v
```

## 🔄 Resumo da Migração (Node.js -> Go)

1. **Validação de Email e Saldo**: Anteriormente feitas em um microsserviço Java auxiliar, foram integradas nativamente em Go no pacote `pkg/utils/validator.go`.
2. **Procedures vs ORM**: O projeto original dependia inteiramente de Stored Procedures (`CALL criar_usuario`, etc). Nesta versão, o GORM foi introduzido para maior segurança (evitando SQL Injection nativamente) e facilidade de manutenção, mas o arquivo `schema.sql` ainda contém as procedures reescritas caso seja necessário o uso legado.
3. **Senhas**: Adicionado `bcrypt` para hash seguro das senhas, substituindo o armazenamento em texto plano implícito no projeto original.
4. **Autenticação**: Adicionado JWT nativo através de um middleware Fiber, protegendo rotas sensíveis que antes não possuíam validação de sessão na API.
