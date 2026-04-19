# 🎬 Movie Data Pipeline & Dashboard

Este projeto demonstra um pipeline de dados completo: desde a ingestão de dados CSV para um banco de dados PostgreSQL até a exibição de métricas em um dashboard web desenvolvido em Go.

## 🏗️ Arquitetura do Sistema

A arquitetura foi desenhada para ser modular e conteinerizada, garantindo que cada componente cumpra uma função específica:

1.  **Camada de Dados (PostgreSQL):** Um banco de dados relacional que armazena informações de filmes, usuários e avaliações.
2.  **Camada de Ingestão (Python):** Um script que verifica a existência das tabelas, processa arquivos CSV e popula o banco de dados.
3.  **Camada de Aplicação (Go):** Um servidor web que consulta o banco de dados e renderiza páginas HTML dinâmicas com paginação e estatísticas.

## 📂 Estrutura de Pastas

* `app/`: Contém o código fonte do servidor Go e os templates HTML.
* `scripts/to_postgres/`: Contém o script de carga de dados e os arquivos CSV originais.
* `.github/workflows/`: Automação de CI/CD para build e push das imagens para o Docker Hub.

---

## 🚀 Como Executar

```bash
make
```

Acesse em: `http://localhost:8080`

Para remover containers e redes, basta executar

```bash
make clean
```

---

## 📊 Funcionalidades do Dashboard

* **Listagem Geral:** Tabela paginada com todos os filmes cadastrados.
* **Top 5 Filmes:** Ranking baseado na quantidade de interações dos usuários.
* **Análise de Gênero:** Identificação do gênero com a melhor média de avaliações.
* **Geolocalização de Usuários:** Estatística de qual país possui o maior volume de avaliações registradas.

## 🛠️ Tecnologias Utilizadas

* **Linguagens:** Go, Python 
* **Banco de Dados:** PostgreSQL 
* **Bibliotecas Python:** Pandas, SQLAlchemy, Psycopg2
* **Bibliotecas Go:** html/template, lib/pq
* **DevOps:** Docker, GitHub Actions (CI/CD)

---

### 📝 Notas de Implementação
O projeto utiliza **Docker Networks** para permitir que os serviços se comuniquem usando nomes de host internos (ex: `postgre:5432`). O workflow do GitHub está configurado para realizar o build automático das imagens sempre que houver um push na branch `main`.
