# PIDASP_2020

Projekat iz predmeta PIDASP aplikacija koja upravlja automobilima.<br>
Aplikacija koristi Hyperledger Fabric mrežu i NodeSDK za rad sa mrežom.

## Hyperledger Fabric mreža

Mreža se sastoji iz 3 organizacije sa po 3 peer-a i jednim orderer peer-om.
Tabela ispod prikazuje bitne informacije o mreži.

| Organizacija  | Naziv | docker port | Anchor peer | kanal     |
|---------------|-------|-------------|-------------|-----------|
| Organization1 | peer0 | 7051        | jeste       | myChannel |
| Organization1 | peer1 | 8051        | nije        | myChannel |
| Organization1 | peer2 | 9051        | nije        | myChannel |
| Organization2 | peer0 | 10051       | jeste       | myChannel |
| Organization2 | peer1 | 11051       | nije        | myChannel |
| Organization2 | peer2 | 12051       | nije        | myChannel |
| Organization3 | peer0 | 13051       | jeste       | myChannel |
| Organization3 | peer1 | 14051       | nije        | myChannel |
| Organization3 | peer2 | 15051       | nije        | myChannel |
| Orderer       | peer0 | 7050, 7053  | jeste       | myChannel |
