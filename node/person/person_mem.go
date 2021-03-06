package person

import (
	"context"
	"fmt"
	"log"

	"github.com/AntanasMaziliauskas/grpc/api"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//DataFromMem structure holds values of Data and ID
type DataFromMem struct {
	Data map[bson.ObjectId]*Person
	ID   string
}

//Init function does nothing
func (d *DataFromMem) Init() error {

	d.Data = make(map[bson.ObjectId]*Person)
	return nil
}

//ListPersons function returns a list of all persons
func (d *DataFromMem) ListPersons(ctx context.Context, in *api.Empty) (*api.MultiPerson, error) {
	listOfPersons := &api.MultiPerson{}

	for _, v := range d.Data {
		listOfPersons.Persons = append(listOfPersons.Persons, &api.Person{
			Id:         v.ID.Hex(),
			Name:       v.Name,
			Age:        v.Age,
			Profession: v.Profession,
			Node:       d.ID,
		})
	}
	if len(listOfPersons.Persons) < 1 {
		fmt.Println("No Data located.")
	}
	return listOfPersons, nil
}

//GetOnePerson function looks for person and returns it if found
func (d *DataFromMem) GetOnePerson(ctx context.Context, in *api.Person) (*api.Person, error) {
	if !bson.IsObjectIdHex(in.Id) {
		log.Println("Provided ID is invalid")

		return &api.Person{}, nil
	}

	if v, ok := d.Data[bson.ObjectIdHex(in.Id)]; ok {
		fmt.Println("Person sucessfully located.")

		return &api.Person{Id: in.Id, Name: v.Name, Age: v.Age, Profession: v.Profession, Node: d.ID}, nil
	}
	fmt.Println("Person not located.")
	return &api.Person{}, nil
}

//GetMultiPerson function looks for multiple persons and returns if found
func (d *DataFromMem) GetMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.MultiPerson, error) {
	listOfData := &api.MultiPerson{}

	for _, k := range in.Persons {
		if !bson.IsObjectIdHex(k.Id) {
			log.Println("Provided ID is invalid")

			continue
		}
		if v, ok := d.Data[bson.ObjectIdHex(k.Id)]; ok {
			listOfData.Persons = append(listOfData.Persons, &api.Person{Id: v.ID.Hex(), Name: v.Name, Age: v.Age, Profession: v.Profession, Node: d.ID})
		}
	}
	if len(listOfData.Persons) < 1 {
		fmt.Println("Unable to locate given persons")
	}

	return listOfData, nil
}

//DropOnePerson removes given person from the slice
func (d *DataFromMem) DropOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	if !bson.IsObjectIdHex(in.Id) {
		log.Println("Provided ID is invalid")

		return &api.Empty{}, nil
	}

	if _, ok := d.Data[bson.ObjectIdHex(in.Id)]; ok {
		delete(d.Data, bson.ObjectIdHex(in.Id))
		log.Println("Person deleted")

		return &api.Empty{}, nil
	}

	fmt.Println("Unable to locate given person")

	return &api.Empty{}, nil
}

//DropMultiPerson removes given persons from the slice
func (d *DataFromMem) DropMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	var success bool

	for _, k := range in.Persons {
		if !bson.IsObjectIdHex(k.Id) {
			log.Println("Provided ID is invalid")

			continue
		}
		if _, ok := d.Data[bson.ObjectIdHex(k.Id)]; ok {
			delete(d.Data, bson.ObjectIdHex(k.Id))
			success = true
		}
	}
	if !success {
		fmt.Println("Unable to locate given persons")

		return &api.Empty{}, nil
	}
	fmt.Println("Persons successfully dropped")

	return &api.Empty{}, nil
}

//UpsertOnePerson adds person to slice
func (d *DataFromMem) UpsertOnePerson(ctx context.Context, in *api.Person) (*api.Empty, error) {
	if !bson.IsObjectIdHex(in.Id) {
		log.Println("Provided ID is invalid")

		return &api.Empty{}, nil
	}

	id, _ := primitive.ObjectIDFromHex(in.Id)

	if _, ok := d.Data[bson.ObjectIdHex(in.Id)]; ok {
		d.Data[bson.ObjectIdHex(in.Id)].Name = in.Name
		d.Data[bson.ObjectIdHex(in.Id)].Age = in.Age
		d.Data[bson.ObjectIdHex(in.Id)].Profession = in.Profession
		log.Println("Data Updated")
	} else {
		d.Data[bson.ObjectIdHex(in.Id)] = &Person{
			ID:         id,
			Name:       in.Name,
			Age:        in.Age,
			Profession: in.Profession,
		}
		log.Println("Data inserted")
	}

	return &api.Empty{}, nil
}

//UpsertMultiPerson adds multiple persons to a slice
func (d *DataFromMem) UpsertMultiPerson(ctx context.Context, in *api.MultiPerson) (*api.Empty, error) {
	for _, v := range in.Persons {
		if !bson.IsObjectIdHex(v.Id) {
			log.Println("Provided ID is invalid")

			continue
		}

		id, _ := primitive.ObjectIDFromHex(v.Id)

		if _, ok := d.Data[bson.ObjectIdHex(v.Id)]; ok {
			d.Data[bson.ObjectIdHex(v.Id)].Name = v.Name
			d.Data[bson.ObjectIdHex(v.Id)].Age = v.Age
			d.Data[bson.ObjectIdHex(v.Id)].Profession = v.Profession
			log.Println("Data Updated")
		} else {
			d.Data[bson.ObjectIdHex(v.Id)] = &Person{
				ID:         id,
				Name:       v.Name,
				Age:        v.Age,
				Profession: v.Profession,
			}
			log.Println("Data Inserted")
		}
	}

	return &api.Empty{}, nil
}

//Ping function does nothing
func (d *DataFromMem) Ping(ctx context.Context, in *api.PingMessage) (*api.Empty, error) {
	return &api.Empty{}, nil
}
