# DevOps Challenge 2025

Infraestrutura com duas aplicações em linguagens diferentes, camada de cache HTTP, observabilidade e execução simplificada via Docker Compose.

---

## Tecnologias

| Componente | Tecnologia |
|---|---|
| release-service | Python 3.12 · FastAPI |
| infra-service | Go 1.22 · net/http |
| Cache | Nginx proxy_cache |
| Observabilidade | Prometheus + Grafana |
| Infraestrutura | Docker Compose |

---

## Estrutura do projeto

```
devops-challenge-2025/
├── release-service/
│   ├── main.py
│   ├── requirements.txt
│   └── Dockerfile
├── infra-service/
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── nginx/
│   └── nginx.conf
├── observability/
│   └── prometheus.yml
├── docs/
│   ├── architecture.svg
│   └── update-flow.svg
├── demo.sh
└── docker-compose.yml
```

---

## Como executar

### Pré-requisitos

- Docker
- Docker Compose

### Subir toda a infraestrutura

```bash
docker compose up --build
```

Um único comando sobe todos os serviços: as duas aplicações, o Nginx com cache, o Prometheus e o Grafana.

### Parar

```bash
docker compose down
```

---

## Endpoints

Todos os requests passam pelo **Nginx na porta 80**, que aplica o cache antes de rotear para as aplicações.

### release-service

| Rota | Descrição | Cache |
|---|---|---|
| `GET /release/` | Informações do último release | 10 segundos |
| `GET /release/time` | Timestamp do último deploy | 10 segundos |

Exemplo de resposta — `GET /release/`:
```json
{
  "service": "release-service",
  "version": "1.4.2",
  "environment": "production",
  "status": "healthy"
}
```

### infra-service

| Rota | Descrição | Cache |
|---|---|---|
| `GET /infra/` | Status dos componentes de infra | 60 segundos |
| `GET /infra/time` | Horário atual do servidor | 60 segundos |

Exemplo de resposta — `GET /infra/`:
```json
{
  "nginx": "healthy",
  "redis": "healthy",
  "kubernetes": "healthy",
  "status": "all systems operational"
}
```

---

## Cache

O cache é gerenciado pelo **Nginx** via `proxy_cache`, sem nenhuma dependência de código nas aplicações. Essa abordagem mantém as apps simples e o cache centralizado na camada HTTP.

Cada resposta inclui o header `X-Cache-Status` com os valores:

- `MISS` — primeira chamada, Nginx buscou na aplicação
- `HIT` — resposta servida direto do cache
- `EXPIRED` — cache expirado, Nginx buscou versão atualizada

### Demonstração do cache

```bash
bash demo.sh
```

O script executa chamadas sequenciais e exibe o status do cache em tempo real, incluindo a expiração do cache de 10s do release-service enquanto o cache de 60s do infra-service ainda está ativo.

---

## Observabilidade

| Serviço | URL |
|---|---|
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

Credenciais padrão do Grafana: `admin` / `admin`

O Prometheus coleta métricas dos serviços a cada 15 segundos. O Grafana está configurado com o Prometheus como datasource e inclui um dashboard com status dos serviços e duração dos scrapes.

---

## Arquitetura

![Arquitetura](docs/architecture.svg)

O Client realiza requisições HTTP na porta 80. O Nginx atua como reverse proxy e aplica o cache antes de rotear para os serviços internos. Prometheus e Grafana operam de forma independente, coletando métricas sem impacto no fluxo principal.

---

## Fluxo de atualização

![Fluxo de atualização](docs/update-flow.svg)

Para atualizar um serviço sem afetar os demais:

```bash
# Rebuild e redeploy de um serviço específico
docker compose build release-service
docker compose up -d --no-deps release-service
```

Para atualizar a configuração do Nginx:

```bash
docker compose exec nginx nginx -s reload
```

---

## Decisões arquiteturais

**Por que Nginx para cache e não Redis?**
O cache na camada HTTP via `proxy_cache` é transparente para as aplicações. Não exige nenhuma lógica de cache no código, funciona para qualquer tipo de resposta e é o padrão em ambientes de produção com reverse proxy. Redis seria mais adequado se as aplicações precisassem de cache compartilhado com lógica de negócio.

**Por que Go no infra-service?**
Go compila para um binário estático, o que permite um Dockerfile com multi-stage build — a imagem final contém apenas o binário, sem runtime adicional. Isso resulta em imagem menor e superfície de ataque reduzida.

**Por que as apps não expõem porta diretamente?**
Os serviços não têm portas mapeadas no `docker-compose.yml`. Todo o tráfego passa obrigatoriamente pelo Nginx, garantindo que o cache seja sempre aplicado e que nenhum cliente consiga acessar as apps diretamente.

---

## Pontos de melhoria

- **CI/CD pipeline** — automatizar o build e deploy via GitHub Actions ou similar ao receber um push na branch principal
- **Healthchecks** — adicionar `healthcheck` no Docker Compose para garantir que o Nginx só suba após as apps estarem prontas
- **Imagens base** — atualizar periodicamente ou migrar para imagens distroless para reduzir vulnerabilidades
- **HTTPS** — adicionar certificado TLS no Nginx para produção
- **Métricas customizadas** — expor endpoint `/metrics` nas aplicações com contadores de requests e latência por rota
- **Alertas** — configurar alertas no Grafana para notificar quando algum serviço ficar indisponível
- **Cache invalidation** — implementar estratégia de invalidação de cache sob demanda via `proxy_cache_purge`