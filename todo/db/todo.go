package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ToDoItem is the struct that represents a single ToDo item
type ToDoItem struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"done"`
}

// DbMap is a type alias for a map of ToDoItems.  The key
// will be the ToDoItem.Id and the value will be the ToDoItem
type DbMap map[int]ToDoItem

// ToDo is the struct that represents the main object of our
// todo app.  It contains a map of ToDoItems and the name of
// the file that is used to store the items.
//
// Notice how the fields in the struct are not exported
//
//	   	 (they are lowercase).  Describe why you think this is
//		 a good design decision.
//
// ANSWER: Since these fields are only used within this package,
// they do not need to be capitalized or exported. Moreover, this
// encapsulates the data and prevents other packages from accessing
// it directly. Making sure that the structs are correctly used and
// only controlled by designated functions. Lastly, if this application
// were lauched as a service, the encapsulation would ensure data
// integrity and consistency. Hence, this is a good design decision.
type ToDo struct {
	toDoMap    DbMap
	dbFileName string
}

// New is a constructor function that returns a pointer to a new
// ToDo struct.  It takes a single string argument that is the
// name of the file that will be used to store the ToDo items.
// If the file doesn't exist, it will be created.  If the file
// does exist, it will be loaded into the ToDo struct.
func New(dbFile string) (*ToDo, error) {

	//Check if the database file exists, if not use initDB to create it
	//In go, you use the os.Stat function to get information about a file
	//In this case, we are only checking the error, because if we get an
	//error we can safely assume that this file does not exist.
	if _, err := os.Stat(dbFile); err != nil {
		//If the file doesn't exist, create it
		err := initDB(dbFile)
		if err != nil {
			return nil, err
		}
	}

	//Now that we know the file exists, at at the minimum we have
	//a valid empty DB, lets create the ToDo struct
	toDo := &ToDo{
		toDoMap:    make(map[int]ToDoItem),
		dbFileName: dbFile,
	}

	// We should be all set here, the ToDo struct is ready to go
	// so we can support the public database operations
	return toDo, nil
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// AddItem accepts a ToDoItem and adds it to the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must not already exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if so, return an error
//
// Postconditions:
//
//	 (1) The item will be added to the DB
//		(2) The DB file will be saved with the item added
//		(3) If there is an error, it will be returned
func (t *ToDo) AddItem(item ToDoItem) error {
	err := t.loadDB()
	if err != nil {
		return err
	}

	if _, ok := t.toDoMap[item.Id]; ok {
		return errors.New("item already exists")
	}
	t.toDoMap[item.Id] = item

	err = t.saveDB()
	if err != nil {
		return err
	}

	return nil
}

// DeleteItem accepts an item id and removes it from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be removed from the DB
//		(2) The DB file will be saved with the item removed
//		(3) If there is an error, it will be returned
func (t *ToDo) DeleteItem(id int) error {
	err := t.loadDB()
	if err != nil {
		return err
	}

	if _, ok := t.toDoMap[id]; !ok {
		return errors.New("item doesn't exist")
	}

	delete(t.toDoMap, id)

	err = t.saveDB()
	if err != nil {
		return err
	}

	return nil
}

// UpdateItem accepts a ToDoItem and updates it in the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be updated in the DB
//		(2) The DB file will be saved with the item updated
//		(3) If there is an error, it will be returned
func (t *ToDo) UpdateItem(item ToDoItem) error {
	err := t.loadDB()
	if err != nil {
		return err
	}

	if _, ok := t.toDoMap[item.Id]; !ok {
		return errors.New("item doesn't exist")
	}
	t.toDoMap[item.Id] = item

	err = t.saveDB()
	if err != nil {
		return err
	}

	return nil
}

// GetItem accepts an item id and returns the item from the DB.
// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The item will be returned, if it exists
//		(2) If there is an error, it will be returned
//			along with an empty ToDoItem
//		(3) The database file will not be modified
func (t *ToDo) GetItem(id int) (ToDoItem, error) {
	err := t.loadDB()
	if err != nil {
		return ToDoItem{}, err
	}

	if _, ok := t.toDoMap[id]; !ok {
		return ToDoItem{}, errors.New("item doesn't exist")
	}

	return t.toDoMap[id], nil
}

// GetAllItems returns all items from the DB.  If successful it
// returns a slice of all of the items to the caller
// Preconditions:   (1) The database file must exist and be a valid
//
// Postconditions:
//
//	 (1) All items will be returned, if any exist
//		(2) If there is an error, it will be returned
//			along with an empty slice
//		(3) The database file will not be modified
func (t *ToDo) GetAllItems() ([]ToDoItem, error) {
	err := t.loadDB()
	if err != nil {
		return nil, err
	}

	var toDoList []ToDoItem
	for _, value := range t.toDoMap {
		toDoList = append(toDoList, value)
	}

	return toDoList, nil
}

// PrintItem accepts a ToDoItem and prints it to the console
// in a JSON pretty format. As some help, look at the
// json.MarshalIndent() function from our in class go tutorial.
func (t *ToDo) PrintItem(item ToDoItem) {
	jsonBytes, _ := json.MarshalIndent(item, "", "  ")
	fmt.Println(string(jsonBytes))
}

// PrintAllItems accepts a slice of ToDoItems and prints them to the console
// in a JSON pretty format.  It should call PrintItem() to print each item
// versus repeating the code.
func (t *ToDo) PrintAllItems(itemList []ToDoItem) {
	for _, item := range itemList {
		t.PrintItem(item)
	}
}

// JsonToItem accepts a json string and returns a ToDoItem
// This is helpful because the CLI accepts todo items for insertion
// and updates in JSON format.  We need to convert it to a ToDoItem
// struct to perform any operations on it.
func (t *ToDo) JsonToItem(jsonString string) (ToDoItem, error) {
	var item ToDoItem
	err := json.Unmarshal([]byte(jsonString), &item)
	if err != nil {
		return ToDoItem{}, err
	}

	return item, nil
}

// ChangeItemDoneStatus accepts an item id and a boolean status.
// It returns an error if the status could not be updated for any
// reason.  For example, the item itself does not exist, or an
// IO error trying to save the updated status.

// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The items status in the database will be updated
//		(2) If there is an error, it will be returned.
//		(3) This function MUST use existing functionality for most of its
//			work.  For example, it should call GetItem() to get the item
//			from the DB, then it should call UpdateItem() to update the
//			item in the DB (after the status is changed).
func (t *ToDo) ChangeItemDoneStatus(id int, value bool) error {
	err := t.loadDB()
	if err != nil {
		return err
	}

	item, err := t.GetItem(id)
	if err != nil {
		return err
	}

	item.IsDone = value
	if err := t.UpdateItem(item); err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------
// THESE ARE HELPER FUNCTIONS THAT ARE NOT EXPORTED AKA PRIVATE
//------------------------------------------------------------

// initDB is a helper function that creates a new file with an
// empty json array.  This is used to make sure that the DB
// file exists for operations on our ToDo struct.  This function
// should be called by the New() function if the DB file doesn't
// exist.  Notice this function does not have a receiver as its
// used by New() to create the DB file
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

func (t *ToDo) saveDB() error {
	//1. Convert our map into a slice
	//2. Marshal the slice into json
	//3. Write the json to our file

	//1. Convert our map into a slice
	var toDoList []ToDoItem
	for _, item := range t.toDoMap {
		toDoList = append(toDoList, item)
	}

	//2. Marshal the slice into json, lets pretty print it, but
	//   this is not required
	data, err := json.MarshalIndent(toDoList, "", "  ")
	if err != nil {
		return err
	}

	//3. Write the json to our file
	err = os.WriteFile(t.dbFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (t *ToDo) loadDB() error {
	data, err := os.ReadFile(t.dbFileName)
	if err != nil {
		return err
	}

	//Now let's unmarshal the data into our map
	var toDoList []ToDoItem
	err = json.Unmarshal(data, &toDoList)
	if err != nil {
		return err
	}

	//Now let's iterate over our slice and add each item to our map
	for _, item := range toDoList {
		t.toDoMap[item.Id] = item
	}

	return nil
}
