# PIDASP_2020

Projekat iz predmeta PIDASP aplikacija koja upravlja automobilima.<br>
Aplikacija koristi Hyperledger Fabric mrežu i NodeSDK za rad sa mrežom.

## Hyperledger Fabric mreža

Mreža se sastoji iz 3 organizacije sa po 3 peer-a i jednim orderer peer-om.
Tabela ispod prikazuje bitne informacije o organizacijama na mreži.

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

### Certificate Authority na mreži
Za svaku organizaciju kao i orderer servis kreiran je pojedan CA koji generiše potrebne kriptomaterijale <br>
za tu odganizaciju. U tabeli ispod su prikazane informacije o CA kontejnerima.

| Organizacija | docker port |
|--------------|-------------|
| CA_Org1      | 7054        |
| CA_Org2      | 8054        |
| CA_Org3      | 10054       |
| CA_Orderer   | 9054        |


### Pokretanje mreže
Za kreiranje same mreže, kanala, postavljanje chain code-a na mrežu i inicializaciju stanja sveta potrebno <br> 
je pokrenuti skriptu **all.sh** koja se nalazi u folderu test-network.


## Serverska aplikacija

Za pokretanje serverske aplikacije potrebno je pozicionirati se u back-end-app direktorijum i instalirati zavisnosti<br>
preko komande **npm install**, nakon toga aplikacija se moze pokrentu komandom **npm start**. <br>
**Preduslov za uspesno pokretanje serverske aplikacije je da je mreza podignuta.**

### EndPoint-i aplikacije

| Operacija                                 | Metoda | URL                                                     | RequestBody                     |
|-------------------------------------------|--------|---------------------------------------------------------|---------------------------------|
| dobavljanje svih person-a                 | GET    | http://localhost:3000/person/all                        | /                               |
| dobavljanje person-a po id                | GET    | http://localhost:3000/person/{personId}                 | /                               |
| dobavljanje svih automobila               | GET    | http://localhost:3000/car/all                           | /                               |
| dobavljanje automobila po id              | GET    | http://localhost:3000/car/{carId}                       | /                               |
| dobavljanje svih gresaka automobila sa id | GET    | http://localhost:3000/car/errors/{carId}                 | /                               |
| kreiranje greske za automobil             | POST   | http://localhost:3000/car/makeBreakdown                 | {description, price, carId}     |
| popravka greske automobila sa id-jem      | POST   | http://localhost:3000/car/fixBreakdown                  | {id, mechanicId}                |
| promena boje automobila                   | POST   | http://localhost:3000/car/changeColor                   | {ID, Color, Cost, mechanicId}   |
| promena vlasnistva automobila             | POST   | http://localhost:3000/car/changeOwner                   | {ID, newOwnerId, buyWithErrors} |
| pretraga po boji                          | GET    | http://localhost:3000/car/filterColor/{color}           | /                               |
| pretraga po boji i vlasniku               | GET    | http://localhost:3000/car/filterColor/{color}/{ownerId} | /                               |

