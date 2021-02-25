package chaincode

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}



// my types
type Car struct {
	ID             string `json:"ID"`
	Manufacturer   string `json:"manufacturer"`
	Model          string `json:"model"`
	Year           int `json:"year"`
	Color          string `json:"color"`
	Owner          string `json:"owner"`
	Price		   float32 `json:"price"`
}

type CarBrakedown struct {
	ID             string `json:"ID"`
	Description    string `json:"description"`
	Price          float32 `json:"price"`
	Car			   string `json:"car"`
	Fixed		   bool `json:"fixed"`
	Mechanic       string `json:"mechanic"`
}

type Person struct {
	ID             string `json:"ID"`
	Name   		   string `json:"name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	Money          float32 `json:"money"`
	Mechanic	   bool `json:"mechanic"`
}



// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	persons := []Person {
		{ID: "person1", Name: "Dusan", Surname: "Bucan", Email: "dusanbzr@gmail.com", Money: 200000, Mechanic:false},
		{ID: "person2", Name: "Sergej", Surname: "Lopatin", Email: "sergerl@gmail.com", Money: 100, Mechanic:false},
		{ID: "person3", Name: "Marko", Surname: "Blagojevic", Email: "blagoje@gmail.com", Money: 5000, Mechanic:false},
		{ID: "person4", Name: "Majstor", Surname: "Majstorovic", Email: "mmajstorovic@gmail.com", Money: 0, Mechanic:true},
	}

	cars := []Car {
		{ID: "car1", Manufacturer: "Toyota", Model: "Corola", Year: 2004, Color: "blue", Owner: "person1", Price: 3200},
		{ID: "car2", Manufacturer: "Toyota", Model: "Avensis", Year: 2006, Color: "gray", Owner: "person1", Price: 3400},
		{ID: "car3", Manufacturer: "Honda", Model: "Civic", Year: 2006, Color: "red", Owner: "person2", Price: 2000},
		{ID: "car4", Manufacturer: "Honda", Model: "Accord", Year: 2006, Color: "black", Owner: "person2", Price: 3000},
		{ID: "car5", Manufacturer: "BMW", Model: "320", Year: 2001, Color: "gray", Owner: "person3", Price: 4000},
		{ID: "car6", Manufacturer: "BMW", Model: "i340", Year: 2010, Color: "blue", Owner: "person3", Price: 6000},
	}


	for _, person := range persons {
		personJSON, err := json.Marshal(person)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(person.ID, personJSON)
		if err != nil {
			return fmt.Errorf("failed to put person to world state. %v", err)
		}
	}

	for _, car := range cars {
		carJSON, err := json.Marshal(car)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(car.ID, carJSON)
		if err != nil {
			return fmt.Errorf("failed to put car to world state. %v", err)
		}
	}

	return nil
}

// ReadPerson returns the person stored in the world state with given id.
func (s *SmartContract) ReadPerson(ctx contractapi.TransactionContextInterface, id string) (*Person, error) {
	personJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if personJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var person Person
	err = json.Unmarshal(personJSON, &person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

// ReadCar returns the car stored in the world state with given id.
func (s *SmartContract) ReadCar(ctx contractapi.TransactionContextInterface, id string) (*Car, error) {
	carJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if carJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var car Car
	err = json.Unmarshal(carJSON, &car)
	if err != nil {
		return nil, err
	}

	return &car, nil
}




// FixCarBrakedown fixes c with specfic ID
func (s *SmartContract) FixCarBrakedown(ctx contractapi.TransactionContextInterface,
	id string, 
	mechanicId string) error {

	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the carBreakdown with %s doesn't exist", id)
	} else {
		carBreakdownJSON, err := ctx.GetStub().GetState(id)
		if err != nil {
			return fmt.Errorf("failed to read from world state: %v", err)
		}
		var carBreakdown CarBrakedown
		err = json.Unmarshal(carBreakdownJSON, &carBreakdown)
		if err != nil {
			return err
		}
		// proveri da li majstor sa ID postoji
		mechanic, errMechanic := s.ReadPerson(ctx, mechanicId)
		if errMechanic != nil {
			return fmt.Errorf("failed to read mechanic from world state: %v", errMechanic)
		}
		if !mechanic.Mechanic {
			return fmt.Errorf("Person with id %s is not mechanic", mechanic.ID)
		}

		// pronadji kola i vlasnika kola
		car, errCar := s.ReadCar(ctx, carBreakdown.Car)
		if errCar != nil {
			return fmt.Errorf("failed to read mechanic from world state: %v", errCar)
		}
		owner, errOwner := s.ReadPerson(ctx, car.Owner)
		if errOwner != nil {
			return fmt.Errorf("failed to read mechanic from world state: %v", errOwner)
		}
		// fiksaj kvar ako je sve okej
		if owner.Money >= carBreakdown.Price && carBreakdown.Fixed == false {
			owner.Money = owner.Money - carBreakdown.Price;
			carBreakdown.Fixed = true;
			carBreakdown.Mechanic = mechanic.ID
			mechanic.Money = mechanic.Money + carBreakdown.Price;

			
			carBreakdownJSON, errSerializeCar := json.Marshal(carBreakdown)
			if errSerializeCar != nil {
				return errSerializeCar
			}		
			putCarBreakdownErr := ctx.GetStub().PutState(id, carBreakdownJSON)
			
			if putCarBreakdownErr != nil {
				return putCarBreakdownErr
			}

			ownerJSON, errSerializeOwner := json.Marshal(owner)
			if errSerializeOwner != nil {
				return errSerializeOwner
			}		
			putOwnerErr := ctx.GetStub().PutState(owner.ID, ownerJSON)

			if putOwnerErr != nil {
				return putOwnerErr
			}

			mechanicJSON, errSeralizeMechanic := json.Marshal(mechanic)
			if errSeralizeMechanic != nil {
				return errSeralizeMechanic
			}		
			putMechanicErr := ctx.GetStub().PutState(mechanic.ID, mechanicJSON)
				
			if putMechanicErr != nil {
				return putMechanicErr
			}
			
			// ako nema gresaka vrati nil
			return nil


		} else {
			return fmt.Errorf("failed to fix carBreakdown with id %s", id)
		}
	}

}



// ReadAllCarBreakDown returns the all carBreakdown for specific car stored in the world state with given id.
func (s *SmartContract) ReadAllCarBreakDown(ctx contractapi.TransactionContextInterface, 
	carId string) ([]*CarBrakedown, 
	error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var carBrakedowns []*CarBrakedown
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var carBrakedown CarBrakedown
		err = json.Unmarshal(queryResponse.Value, &carBrakedown)
		if err == nil {
			// ako je to sto naidje CarBrakedown onda gledaj da li ima id kola
			if carBrakedown.Car == carId {
				carBrakedowns = append(carBrakedowns, &carBrakedown)	
			}
		}
		
	}
	return carBrakedowns, nil
}


// GetAllCars returns the all cars in world state
func (s *SmartContract) GetAllCars(ctx contractapi.TransactionContextInterface) ([]*Car, 
	error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var cars []*Car
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var car Car
		err = json.Unmarshal(queryResponse.Value, &car)
		if err == nil {
			if strings.HasPrefix(car.ID, "car") {
				cars = append(cars, &car)
			}
		}
		
	}
	return cars, nil
}


func (s *SmartContract) GetAllCarsByCollor(
	ctx contractapi.TransactionContextInterface, color string) ([]*Car, error) {

	cars, getCarrErr := s.GetAllCars(ctx)
	if getCarrErr != nil {
		return nil, getCarrErr
	}
	var carsFiltered []*Car
	for _, car := range cars {
		if car.Color == color {
			carsFiltered = append(carsFiltered, car)
		}
	}
	return carsFiltered, nil
}

func (s *SmartContract) GetAllCarsByCollorAndOwner(
	ctx contractapi.TransactionContextInterface, color string, owner string) ([]*Car, error) {

	cars, getCarrErr := s.GetAllCarsByCollor(ctx, color)
	if getCarrErr != nil {
		return nil, getCarrErr
	}
	var carsFiltered []*Car
	for _, car := range cars {
		if car.Owner == owner {
			carsFiltered = append(carsFiltered, car)
		}
	}
	return carsFiltered, nil
}


// GetAllUsers returns the all persons in world state
func (s *SmartContract) GetAllUsers(ctx contractapi.TransactionContextInterface) ([]*Person, 
	error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var persons []*Person
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var person Person
		err = json.Unmarshal(queryResponse.Value, &person)
		if err == nil {
			if strings.HasPrefix(person.ID, "person") {
				persons = append(persons, &person)
			}	
		}
		
	}
	return persons, nil
}


func (s *SmartContract) GetCarBreakDownUnfixedCost(ctx contractapi.TransactionContextInterface, carId string) (float32, error) {
	var retVal float32 = 0.0
	breakDowns, err:= s.ReadAllCarBreakDown(ctx, carId)
	if err != nil {
		return retVal, err
	}
	
	for _, breadDown := range breakDowns {
		if breadDown.Fixed == false {
			retVal += breadDown.Price
		}
	}
	return retVal, nil
}

func (s *SmartContract) UpdateCarColor(ctx contractapi.TransactionContextInterface, 
	id string, color string, cost float32, mechanicId string) error {

	car, carErr := s.ReadCar(ctx, id)
	if carErr != nil {
		return carErr
	}
	owner, ownerErr := s.ReadPerson(ctx, car.Owner)
	if ownerErr != nil {
		return ownerErr
	}

	machanic, machanicErr := s.ReadPerson(ctx, mechanicId)
	if machanicErr != nil {
		return machanicErr
	}
	if !machanic.Mechanic {
		return fmt.Errorf("Person with id %s is not mechanic", mechanicId)
	}

	if owner.Money >= cost {
		owner.Money = owner.Money - cost;
		car.Color = color;
		machanic.Money = machanic.Money + cost;

		mechanicJSON, errSeralizeMechanic := json.Marshal(machanic)
		if errSeralizeMechanic != nil {
			return errSeralizeMechanic
		}		
		putMechanicErr := ctx.GetStub().PutState(machanic.ID, mechanicJSON)
				
		if putMechanicErr != nil {
			return putMechanicErr
		}

		ownerJSON, errSeralizeOwner := json.Marshal(owner)
		if errSeralizeOwner != nil {
			return errSeralizeOwner
		}		
		putOwnerErr := ctx.GetStub().PutState(owner.ID, ownerJSON)
				
	
		if putOwnerErr != nil {
			return putOwnerErr
		}

		carJSON, errSeralizeCar := json.Marshal(car)
		if errSeralizeCar != nil {
			return errSeralizeCar
		}		
		putCarErr := ctx.GetStub().PutState(car.ID, carJSON)
				
		if putCarErr != nil {
			return putCarErr
		}
		return nil;
	} else {
		return fmt.Errorf("Owner with id %s has no money for change color of car", owner.ID)
	}
}

func (s *SmartContract) UpdateCarOwner(ctx contractapi.TransactionContextInterface, 
	id string, newOwnerId string, buyWithErrors bool) error {

	carBroken := false
	carBreadownCost, carBrokenErr := s.GetCarBreakDownUnfixedCost(ctx, id)
	if carBrokenErr != nil {
		return carBrokenErr
	}
	if carBreadownCost > 0 {
		carBroken = true
	}

	car, carErr := s.ReadCar(ctx, id)
	if carErr != nil {
		return carErr
	}
	owner, ownerErr := s.ReadPerson(ctx, car.Owner)
	if ownerErr != nil {
		return ownerErr
	}

	newOwner, newOwnerErr := s.ReadPerson(ctx, newOwnerId)
	if newOwnerErr != nil {
		return newOwnerErr
	}
	
	if newOwner.Money >= car.Price && newOwnerId != owner.ID {

		if !buyWithErrors && !carBroken{
			owner.Money = owner.Money + car.Price;
			car.Owner = newOwnerId;
			newOwner.Money = newOwner.Money - car.Price;
		} 
		if !buyWithErrors && carBroken {
			return fmt.Errorf("Car has unfixed erros of total cost %f and buy with errors is not allowed", carBreadownCost)
		}
		if buyWithErrors {
			owner.Money = owner.Money + (car.Price - carBreadownCost);
			car.Owner = newOwnerId;
			newOwner.Money = newOwner.Money - (car.Price - carBreadownCost);
		}

		newOwnerJSON, errSeralizeNewOwner := json.Marshal(newOwner)
		if errSeralizeNewOwner != nil {
			return errSeralizeNewOwner
		}		
		putNewOwnerErr := ctx.GetStub().PutState(newOwnerId, newOwnerJSON)
				
		if putNewOwnerErr != nil {
			return putNewOwnerErr
		}

		ownerJSON, errSeralizeOwner := json.Marshal(owner)
		if errSeralizeOwner != nil {
			return errSeralizeOwner
		}		
		putOwnerErr := ctx.GetStub().PutState(owner.ID, ownerJSON)
				
	
		if putOwnerErr != nil {
			return putOwnerErr
		}

		carJSON, errSeralizeCar := json.Marshal(car)
		if errSeralizeCar != nil {
			return errSeralizeCar
		}		
		putCarErr := ctx.GetStub().PutState(car.ID, carJSON)
				
		if putCarErr != nil {
			return putCarErr
		}
		return nil;
	} else {
		return fmt.Errorf("New owner with id %s has no money to buy car", newOwner.ID)
	}
}

// CreateCarBrakedown issues a new carBreakdown to the world state with given details.
func (s *SmartContract) CreateCarBrakedown(ctx contractapi.TransactionContextInterface,
	description string, price float32, carId string) error {
		   

	car, carErr := s.ReadCar(ctx, carId)
	if carErr != nil {
		return carErr
	}
   
   var newBreadkDownIndex int = 0
   breakDowns, breakDownsErr := s.ReadAllCarBreakDown(ctx, carId);
   if breakDownsErr != nil {
	   return breakDownsErr
   }
   if len(breakDowns) > 0 {
	   newBreadkDownIndex = len(breakDowns)
   }
   var newBreakDownId string = "breakDown" + strconv.Itoa(newBreadkDownIndex) + carId

   // uslov da se brise auto ako je cena poravke troskova veca od cene auta
   totalFixCost, errCalcTotalFixCost := s.GetCarBreakDownUnfixedCost(ctx, carId)
   if errCalcTotalFixCost != nil {
	   return errCalcTotalFixCost
   }
   totalFixCost += price
   // obrisi auto
   if totalFixCost > car.Price {
	   deleteCarErr := s.DeleteCar(ctx, carId)
	   if deleteCarErr != nil {
		   return deleteCarErr
	   }
   }




   brakedown := CarBrakedown{
	   ID: newBreakDownId,
	   Description: description,
	   Price: price,
	   Car: carId,
	   Fixed: false,
	   Mechanic: "",
   }

   brakedownJSON, err := json.Marshal(brakedown)
   if err != nil {
	   return err
   }

   return ctx.GetStub().PutState(newBreakDownId, brakedownJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteCar(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the car %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}
