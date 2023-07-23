package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

var voterDbFile = "./data/voter.json"

type voterPoll struct {
	PollID   uint64
	VoteDate time.Time
}

type Voter struct {
	VoterID     uint64
	FirstName   string
	LastName    string
	VoteHistory []voterPoll
}

type VoterMap map[uint64]Voter //A map of VoterIDs as keys and Voter structs as values

type Voters struct {
	voterMap VoterMap
}

func New() (*Voters, error) {

	if _, err := os.Stat(voterDbFile); err != nil {
		//If the file doesn'v exist, create it
		err := initDB(voterDbFile)
		if err != nil {
			return nil, err
		}
	}

	voter := &Voters{
		voterMap: make(map[uint64]Voter),
	}

	return voter, nil
}

func (v *Voters) AddItem(item Voter) error {

	err := v.loadDB()
	if err != nil {
		return err
	}

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	_, ok := v.voterMap[item.VoterID]
	if ok {
		return errors.New("item already exists")
	}

	//Now that we know the item doesn'v exist, lets add it to our map
	v.voterMap[item.VoterID] = item

	//If everything is ok, return nil for the error
	return nil
}

func (v *Voters) DeleteItem(id uint64) error {

	err := v.loadDB()
	if err != nil {
		return err
	}

	// we should if item exists before trying to delete it
	// this is a good practice, return an error if the
	// item does not exist

	//Now lets use the built-in go delete() function to remove
	//the item from our map
	delete(v.voterMap, id)

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /voters endpoint
func (v *Voters) DeleteAll() error {
	v.voterMap = make(map[uint64]Voter)
	// saveDb()
	return nil
}

func (v *Voters) UpdateItem(item Voter) error {

	err := v.loadDB()
	if err != nil {
		return err
	}

	_, ok := v.voterMap[item.VoterID]
	if !ok {
		return errors.New("item does not exist")
	}

	v.voterMap[item.VoterID] = item

	return nil
}

func (v *Voters) GetItem(id uint64) (Voter, error) {

	err := v.loadDB()
	if err != nil {
		return Voter{}, err
	}

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	item, ok := v.voterMap[id]
	if !ok {
		return Voter{}, errors.New("item does not exist")
	}

	return item, nil
}

func (v *Voters) ChangeVoterID(id uint64, value bool) error {

	return errors.New("not implemented")
}

func (v *Voters) GetAllItems() ([]Voter, error) {

	err := v.loadDB()
	if err != nil {
		return nil, err
	}

	//Now that we have the DB loaded, lets crate a slice
	var voters []Voter

	//Now lets iterate over our map and add each item to our slice
	for _, item := range v.voterMap {
		voters = append(voters, item)
	}

	//Now that we have all of our items in a slice, return it
	return voters, nil
}

func (v *Voters) GetPoll(id uint64, pollsid uint64) (voterPoll, error) {

	err := v.loadDB()
	if err != nil {
		return voterPoll{}, err
	}

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	item, ok := v.voterMap[id]
	if !ok {
		return voterPoll{}, errors.New("item does not exist")
	}

	for _, poll := range item.VoteHistory {
		if poll.PollID == pollsid {
			return poll, nil
		}
	}

	return voterPoll{}, errors.New("poll does not exist")
}

// Helper Functions

// Print functions
// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (v *Voters) PrintItem(item Voter) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (v *Voters) PrintAllItems(itemList []Voter) {
	for _, item := range itemList {
		v.PrintItem(item)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (v *Voters) JsonToItem(jsonString string) (Voter, error) {
	var item Voter
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return Voter{}, err
	}

	return item, nil
}

// DB functions
func initDB(dbFileName string) error {
	f, err := os.Create(dbFileName)
	if err != nil {
		return err
	}

	// Given we are working with a json array as our DB structure
	// we should initialize the file with an empty array, which
	// in json is represented as "[]
	_, err = f.Write([]byte("[]"))
	if err != nil {
		return err
	}

	f.Close()

	return nil
}

func (v *Voters) saveDB() error {
	//1. Convert our map into a slice
	//2. Marshal the slice into json
	//3. Write the json to our file

	//1. Convert our map into a slice
	var voters []Voter
	for _, item := range v.voterMap {
		voters = append(voters, item)
	}

	//2. Marshal the slice into json, lets pretty print it, but
	//   this is not required
	data, err := json.MarshalIndent(voters, "", "  ")
	if err != nil {
		return err
	}

	//3. Write the json to our file
	err = os.WriteFile(voterDbFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (v *Voters) loadDB() error {
	data, err := os.ReadFile(voterDbFile)
	if err != nil {
		return err
	}

	//Now let's unmarshal the data into our map
	var voters []Voter
	err = json.Unmarshal(data, &voters)
	if err != nil {
		return err
	}

	//Now let's iterate over our slice and add each item to our map
	for _, item := range voters {
		v.voterMap[item.VoterID] = item
	}

	return nil
}
