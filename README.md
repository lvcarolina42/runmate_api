# runmate_api

## Estrutura

```
.
├── cmd             => Ponto de entrada (declarações de dependências e execução do programa)
│   └── main.go
├── config          => Configurações da API
│   └── env.go
├── http
│   ├── handler     => Tratamentos da API
│   │   ├── api.go
│   │   └── chat.go
│   └── model       => Representação dos modelos da API
│       ├── activity.go
│       ├── challenge.go
│       ├── event.go
│       ├── message.go
│       └── user.go
├── internal
│   ├── chat        => Tratamentos para o chat
│   │   ├── hub.go      => Gerenciamento das conexões do chat
│   │   └── kafka.go    => Consumidor e publicador do Kafka
│   ├── entity      => Representação dos modelos do banco
│   │   ├── activity.go
│   │   ├── challenge.go
│   │   ├── event.go
│   │   ├── message.go
│   │   └── user.go
│   ├── repository  => Interface com o banco
│   │   ├── activity.go
│   │   ├── challenge.go
│   │   ├── event.go
│   │   ├── message.go
│   │   └── user.go
│   └── service     => Casos de uso
│       ├── activity.go
│       ├── challenge.go
│       ├── event.go
│       ├── message.go
│       └── user.go
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Configuração da API

Para configurar a API corretamente, é necessário declarar as seguintes variáveis de ambiente:
```
export ENV="dev"
export API_PORT="3000"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="runmate"
export DB_PASSWORD="runmate"
export DB_NAME="runmate"
export KAFKA_HOST="localhost"
export KAFKA_PORT="9092"
export KAFKA_ACCESS_KEY_NAME=""
export KAFKA_ACCESS_KEY=""
```

## Casos de Uso "Complexos"

Explicação dos casos de uso mais complexos. Os casos de uso que não aparecem aqui são considerados intuitivos.

### Criar atividade (runmate_api/internal/service/activity.go(.Create))

1. Verifica se o usuário existe
1. Cria a atividade no banco
1. Atualiza a XP do usuário com base na distância percorrida (1 metro = 1 XP)
1. Busca os desafios ativos que o usuário participa
1. Cria um evento, no banco, em cada desafio com a distância percorrida
1. Se a soma total da distância percorrida for maior que a distância do desafio, o desafio é encerrado

## Criação de desafios

Existem dois tipos de desafios:

1. Desafios com meta de distância (ChallengeTypeDistance)
    a. Não existe uma data de fim para o desafio. Encerra com o primeiro usuário que atingir a meta
1. Desafios com meta de data (ChallengeTypeDate)
    a. Não existe uma distância para o desafio. Encerra quando a data de fim do desafio for atingida

### Ranking dos desafios (runmate_api/internal/service/challenge.go(.Ranking))

1. Busca pelos eventos do desafio
1. Sumariza as distâncias percorridas
1. Ordena as distâncias percorridas

### Cálculo do nível do usuário

Cada atividade realizada gera XP (pontos de experiência) para o usuário. Cada metro equivale a 1 ponto.
Para definir o nível do usuário, a partir do XP, é utilizado uma expressão clássica de jogos, para progressão
linear do nível:

`Sqrt(1000 * (2 * XP + 250)) + 500 / 1000`

Essa expressão define que o usuário deve percorrer 1.000 metros a mais que percorreu no nível anterior, ou seja:

| Nível |   XP  | Diferença |
|-------|-------|-----------|
|   1	|     0 | 	    -   |
|   2	|  1000 | 	 1000   |
|   3	|  3000 | 	 2000   |
|   4	|  6000 | 	 3000   |
|   5	| 10000 | 	 4000   |
|   6	| 15000 | 	 5000   |
|   7	| 21000 | 	 6000   |
|   8	| 28000 | 	 7000   |

## Chat

### Características

- Atrelado ao desafio
- Tempo real
- Utiliza o Kafka para publicar e consumir mensagens
    - Um tópico por desafio
- Utiliza websocket para manter conexão dos usuários

### Funcionamento

1. Usuários se conectam ao hub do desafio pelo websocket
    a. A partir desse momento, terá acesso a todas as mensagens enviadas no hub
    a. Quando o usuário se conecta, é enviada uma mensagem, exclusiva para o sistema (`type = 1`), para a criação do
    tópico do Kafka, caso ele não exista. Essa mensagem não deve ser exibida para os usuários
1. Ao enviar uma mensagem, ela é publicada no tópico do Kafka pelo Publicador (runmate_api/internal/chat/kafka.go(.Publisher))
1. O consumidor recebe as mensagens do tópico (runmate_api/http/handler/chat.go(.Consumer.Start))
    a. Interpreta a mensagem
    a. Salva no banco, para histórico
    a. Constrói um modelo mais claro para o cliente (app)
    a. Envia a mensagem para os usuários conectados ao hub (Broadcast)

### Limitações

Por ser um broadcast da mensagem num tópico geral para o desafio. O usuário que enviou a mensagem também a recebe.
É importante o cliente ter um tratamento para evitar duplicação de mensagens.
