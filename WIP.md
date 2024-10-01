### Eliminação dos intermediários

-   Durante a banca, o prof comentou sobre a possibilidade de automatizar o processo de oferta, eliminando os intermediários. Se tratando de seguros massificados isso não é viável para o negócio, pois esse modelo de negócio funciona através da oferta a massa por meio de estipulantes e representantes, em geral a sua base de clientes. Ex: Varejo e bancos.

### Criptografia de dados vs criptografia de acesso (Pra LGPD pensar em criptografia de dados)

Estamos pensando sobre o armazenamento de dados sensíveis na blockchain.

-   Em primeiro lugar, no escopo específico deste trabalho: seguro massificados, precisamos armazenar dados do usuário? Essa dúvida surge, pois os seguros massificados não dependem de uma análise individualizada do segurado, então não estamos enxergando a necessidade de salvar esses dados. Os dados do segurado são importantes durante a etapa de análise de evidências, para validar a ocorrência do sinistro.

-   Dados do segurado poderiam ser obtidos através da plataforma do parceiro de negócio (OFF-CHAIN), visto que esse intermediário que fará a oferta do seguro ao segurado? `R: Nesse caso, a nossa PoC só armazenaria o contrato de seguro, incluindo todas etapas do seu lifecycle. Então a nossa proposta de centralização da informação entre os participantes se enfraqueceria.`

-   Se tratando de uma solução blockchain seria viável o segurado ser o portador das credenciais geradas, não importando o CPF em posse delas? Similar ao portador das seeds em criptomoedas. No caso, salvariamos o hash das seeds no contrato para relacionamento com o segurado portador delas.

-   A autenticação do segurado não é feita ON-CHAIN, mas sim OFF-CHAIN. Ou seja, para que o smart contract de acionamento do seguro seja invocado o segurado precisa se autenticar na plataforma da seguradora.

-   As identidades de acesso a blockchain são necessárias apenas para as aplicações clientes, elas se conectam por um gateway com TLS. Não é necessário criar identidades para os usuários da aplicação cliente (segurados, colaboradores etc.), a autenticação dos usuários deve ser feita OFF-CHAIN. A hyperledge recomenda a retenção da conexão criada com o gateway, o estabelecimento frequente de conexões ao gateway resultará em problemas de performance.

-   A Hyperledger Fabric possuí uma solução para lidar com dados sensíveis, visando atender a GDPR (LGPD da União Europeia). Ainda não entendemos qual o poder dessa solução, enxergamos que foi adicionada mais uma camada de controle de acesso e permissionamento ao dado sensível, porém os dados privados podem ser eliminados deixando um hash dos dados na ledger. Em nossa visão atual, isso seria equivalente a ter os dados sensíveis em um DB off-chain em que o segurado poderia exigir a exclusão de seus dados. Além disso, não conseguimos enxergar valor nesse hash dos dados salvos na ledger, isso apenas mostra que em um dado momento os dados existiram...

```
Solução Hyperleder: Private Data

For very sensitive data, even the parties sharing the private data might want — or might be required by government regulations — to periodically “purge” the data on their peers, leaving behind a hash of the data on the blockchain to serve as immutable evidence of the private data.
In some of these cases, the private data only needs to exist in the peer’s private database until it can be replicated into a database external to the peer’s blockchain. The data might also only need to exist on the peers until a chaincode business process is done with it (trade settled, contract fulfilled, etc).
To support these use cases, private data can be purged so that it is not available for chaincode queries, not available in block events, and not available for other peers requesting the private data.
```

-   A abordagem sugerida pela banca em TCC I foi utilizar uma chave simétrica para decode dos encoded senstive data salvos na ledger, porém uma vez que essa chave simétrica for excluída, apenas uma encoded string ficará na rede. Algo parecido com a proposta de Private Data da Hyperledger, porém tratando no client-side (OFF-CHAIN).

-   O nosso questionamento principal é que não conseguimos enxergar valor nesse hash dos dados salvos na ledger, isso apenas mostra que em um dado momento os dados existiram. Se os dados sensíveis do segurado forem importantes apenas até o pagamento do prêmio ou fim do contrato, então tanto a solução Private Data, como a chave simétrica para decoding dos dados são possíveis soluções para garantir a não persistência de dados sensíveis.

### Dúvidas sobre Monografia

-   Para implementação da PoC utilizamos diversas ferramentas, citamos algumas abaixo, para criar a rede de testes e as aplicações. Teremos que adicionar seções na revisão bibliográfica sobre essas ferramentas utilizadas? Se sim, qual o nível de detalhamento esperado?
    -   Hyperledger Binaries
    -   Contêineres
    -   BashScript
    -   Golang
    -   ReactJS
    -   Bancos de dados
