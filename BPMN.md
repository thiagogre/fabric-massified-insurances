# Entendendo o Negócio do Projeto

-   Área de Domínio: Mercado de Seguros
-   Projeto: Sistema para promover transparência, eficiência e conformidade regulatória na operação de um Seguro Massificado.

# Descrição de Processos de Negócio

-   Seguradora
    -   Registro da contratação
    -   Gestão e faturamento
    -   Análise da solicitação de sinistro
-   Parceiros de distribuição (Corretores, Representantes e Estipulantes)
    -   Oferta do seguro para o consumidor
    -   Registro da contratação
-   Consumidor
    -   Contratação do seguro, através de um dos parceiros de distribuição
    -   Acionamento do seguro

# BPMN

## Contratação

![BPMN Contratação](assets/bpmn_contratacao.png)

### Lista de Tarefas

1.  Solicitar seguro através do parceiro de distribuição
2.  Verificar possíveis coberturas de seguro
3.  Selecionar cobertura do seguro
4.  Registrar contratação
5.  Faturar prêmios

### Detalhe das Tarefas

<table>
    <tr>
        <th></th>
        <th>Nome da tarefa</th>
        <th>Dados de entrada</th>
        <th>Dados de saída</th>
        <th>Detalhamento da tarefa</th>
    </tr>
    <tr>
        <td>1</td>
        <td>Solicitar seguro através do parceiro de distribuição</td>
        <td>-</td>
        <td>Dados do produto</td>
        <td>O consumidor solicita uma cobertura de seguro para seu produto e informa dados do produto</td>
    </tr>
    <tr>
        <td>2</td>
        <td>Verificar possíveis coberturas de seguro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>3</td>
        <td>Selecionar cobertura do seguro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>4</td>
        <td>Registrar contratação</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>5</td>
        <td>Faturar prêmios</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
</table>

## Acionamento

![BPMN Acionamento](assets/bpmn_acionamento.png)

### Lista de Tarefas

1.  Solicitar acionamento do seguro
2.  Solicitar evidências do sinistro
3.  Anexar evidências do sinistro
4.  Analisar sinistro
5.  Aprovar/rejeitar solicitação de sinistro

### Detalhe das Tarefas

<table>
    <tr>
        <th></th>
        <th>Nome da tarefa</th>
        <th>Dados de entrada</th>
        <th>Dados de saída</th>
        <th>Detalhamento da tarefa</th>
    </tr>
    <tr>
        <td>1</td>
        <td>Solicitar acionamento do seguro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>2</td>
        <td>Solicitar evidências do sinistro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>3</td>
        <td>Anexar evidências do sinistro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>4</td>
        <td>Analisar sinistro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
    <tr>
        <td>5</td>
        <td>Aprovar/rejeitar solicitação de sinistro</td>
        <td>-</td>
        <td>-</td>
        <td>-</td>
    </tr>
</table>
