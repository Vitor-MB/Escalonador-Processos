# Backend - Escalonador de Processos

## Visão geral

Este backend foi desenvolvido em Go para simular algoritmos de escalonamento de processos e expor os resultados por meio de uma API REST. Ele recebe uma lista de processos, aplica o algoritmo solicitado e retorna métricas como tempo total, tempo de espera, turnaround, resposta e a linha do tempo da execução.

A aplicação é voltada para integração com a interface web do projeto, permitindo que a página front-end envie os dados de entrada e receba o resultado da simulação em formato JSON.

## Tecnologias

- Go
- HTTP Server nativo da biblioteca padrão

## Estrutura do backend

```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── controllers/
│   │   └── scheduler.go
│   ├── dto/
│   │   ├── requestSimulate.go
│   │   ├── responseAlgorithms.go
│   │   └── responseSimulate.go
│   ├── response/
│   │   └── response.go
│   ├── router/
│   │   ├── router.go
│   │   └── routes.go
│   ├── scheduler/
│   │   ├── algorithms.go
│   │   ├── engine.go
│   │   ├── process.go
│   │   └── simulate.go
│   └── service/
│       └── scheduler.go
├── go.mod
└── README.md
```

## Design e arquitetura

O backend segue uma estrutura em camadas para separar responsabilidades:

1. `cmd/api/main.go`
   - Ponto de entrada da aplicação.
   - Inicializa o roteador HTTP e sobe o servidor na porta `8080`.

2. `internal/router`
   - Define as rotas HTTP e habilita CORS para permitir chamadas vindas do frontend.
   - Cada rota mapeia para uma função do controller correspondente.

3. `internal/controllers`
   - Recebe as requisições HTTP.
   - Decodifica o JSON enviado pelo cliente.
   - Chama o serviço apropriado e responde no formato JSON.

4. `internal/service`
   - Converte os DTOs de entrada em estruturas internas do scheduler.
   - Invoca a lógica de simulação e transforma os resultados em DTOs de resposta.

5. `internal/scheduler`
   - Contém o núcleo da lógica de escalonamento.
   - Implementa os algoritmos, a validação de parâmetros, o cálculo das métricas e a geração da linha do tempo.

6. `internal/dto`
   - Define as estruturas de entrada e saída da API.
   - Mantém a comunicação entre camadas estável e organizada.

### Fluxo de execução

```text
Frontend --> HTTP Request --> Controller --> Service --> Scheduler --> Resultados --> Response JSON
```

## Endpoints

A API roda em:

```text
http://localhost:8080
```

### 1) POST /simulate

Executa a simulação de escalonamento com os processos informados.

#### Requisição

```json
{
  "algorithm": "rr",
  "quantum": 2,
  "aging": 1,
  "processes": [
    { "name": "P1", "arrival": 0, "burst": 5, "priority": 2 },
    { "name": "P2", "arrival": 1, "burst": 3, "priority": 1 },
    { "name": "P3", "arrival": 2, "burst": 2, "priority": 3 }
  ]
}
```

#### Campos da requisição

- `algorithm`: algoritmo de escalonamento.
- `quantum`: quantum utilizado por algoritmos de time-sharing (ex.: `rr`).
- `aging`: valor de envelhecimento para algoritmos que utilizam prioridade com aging.
- `processes`: lista de processos.

Cada processo contém:

- `name`: nome do processo
- `arrival`: instante de chegada
- `burst`: tempo de CPU necessário
- `priority`: prioridade do processo

#### Resposta

```json
{
  "algorithm": "rr",
  "total_time": 10,
  "context_switches": 3,
  "processes": [
    {
      "name": "P1",
      "arrival": 0,
      "burst": 5,
      "priority": 2,
      "start": 0,
      "finish": 5,
      "waiting": 2,
      "turnaround": 5,
      "response": 0
    }
  ],
  "intervals": [
    {
      "process_name": "P1",
      "start": 0,
      "finish": 2
    }
  ],
  "timeline": [
    {
      "from": 0,
      "to": 1,
      "states": {
        "P1": "running",
        "P2": "ready"
      }
    }
  ],
  "averages": {
    "turnaround": 5.5,
    "waiting": 2.0,
    "response": 1.5
  }
}
```

#### Validações

A simulação exige que:

- o algoritmo exista;
- a lista de processos não esteja vazia;
- `quantum` seja maior que zero para algoritmos que usam quantum;
- `aging` seja maior que zero para algoritmos que usam aging;
- `arrival` e `burst` sejam valores válidos.

Em caso de erro, a API retorna um JSON no seguinte formato:

```json
{
  "error": "mensagem detalhada do problema"
}
```

### 2) GET /algorithm

Retorna a lista de algoritmos suportados pelo backend, além de metadados sobre seu comportamento.

#### Resposta exemplo

```json
{
  "algorithms": [
    {
      "algorithm": "fcfs",
      "preemptive": false,
      "uses_quantum": false,
      "uses_aging": false,
      "uses_priority": false
    },
    {
      "algorithm": "rr",
      "preemptive": true,
      "uses_quantum": true,
      "uses_aging": false,
      "uses_priority": false
    }
  ]
}
```

## Algoritmos suportados

A API oferece suporte aos seguintes algoritmos:

- `fcfs` — First Come, First Served
- `sjf` — Shortest Job First
- `srtf` — Shortest Remaining Time First
- `prio-np` — Priority sem preempção
- `prio-p` — Priority com preempção
- `rr` — Round Robin
- `rr-prio-aging` — Round Robin com prioridade e aging

## Tratamento de CORS

O roteador habilita CORS com `Access-Control-Allow-Origin: *`, permitindo que o frontend, mesmo em outra origem, consiga consumir a API. Também trata requisições `OPTIONS` automaticamente.

## Como executar

Dentro da pasta `backend`:

```bash
cd backend
 go mod download
 go run ./cmd/api
```

Depois, a API ficará disponível em:

```text
http://localhost:8080
```

## Observações de design

- O backend foi construído para ser simples, didático e fácil de estender.
- A lógica de escalonamento fica isolada em `internal/scheduler`, permitindo que novas políticas sejam adicionadas sem alterar diretamente a camada HTTP.
- A separação entre DTOs, serviços e scheduler ajuda a manter o código organizado e reutilizável.
- A estrutura foi pensada para receber dados do frontend e retornar resultados em formato JSON padronizado.

