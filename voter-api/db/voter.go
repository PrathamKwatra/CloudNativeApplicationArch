package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type voterPoll struct {
	PollID   uint
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type VoterMap map[uint]Voter //A map of VoterIDs as keys and Voter structs as values

type Voters struct {
	voterMap VoterMap
}

func New() (*Voters, error) {

	voter := &Voters{
		voterMap: make(map[uint]Voter),
	}

	return voter, nil
}

func (t *Voters) AddItem(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	_, ok := t.voterMap[item.VoterID]
	if ok {
		return errors.New("item already exists")
	}

	//Now that we know the item doesn't exist, lets add it to our map
	t.voterMap[item.VoterID] = item

	//If everything is ok, return nil for the error
	return nil
}

func (t *Voters) DeleteItem(id uint) error {

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(t.voterMap, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /voters endpoint
func (t *Voters) DeleteAll() error {
	t.voterMap = make(map[uint]Voter)

	return nil
}

func (t *Voters) UpdateItem(item Voter) error {

	_, ok := t.voterMap[item.VoterID]
	if !ok {
		return errors.New("item does not exist")
	}

	t.voterMap[item.VoterID] = item

	return nil
}

func (t *Voters) GetItem(id uint) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	item, ok := t.voterMap[id]
	if !ok {
		return Voter{}, errors.New("item does not exist")
	}

	return item, nil
}

func (t *Voters) ChangeVoterID(id uint, value bool) error {

	return errors.New("not implemented")
}

func (t *Voters) GetAllItems() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voters []Voter

	//Now lets iterate over our map and add each item to our slice
	for _, item := range t.voterMap {
		voters = append(voters, item)
	}

	//Now that we have all of our items in a slice, return it
	return voters, nil
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (t *Voters) PrintItem(item Voter) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (t *Voters) PrintAllItems(itemList []Voter) {
	for _, item := range itemList {
		t.PrintItem(item)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (t *Voters) JsonToItem(jsonString string) (Voter, error) {
	var item Voter
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}
