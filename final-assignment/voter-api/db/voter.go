package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

// var voterDbFile = "./data/voter.json"

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

// type VoterMap map[uint64]Voter //A map of VoterIDs as keys and Voter structs as values

type Voters struct {
	cache
}

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

func New() (*Voters, error) {

	redisUrl := os.Getenv("REDIS_URL")
	//This handles the default condition
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}

	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*Voters, error) {

	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	//By default, redis manages keys and values, where the values
	//are either strings, sets, maps, etc.  Redis has an extension
	//module called ReJSON that allows us to store JSON objects
	//however, we need a companion library in order to work with it
	//Below we create an instance of the JSON helper and associate
	//it with our redis connnection
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &Voters{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

func (v *Voters) AddVoter(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	redisKey := redisKeyFromId(int(item.VoterID))
	var existingItem Voter
	if err := v.getItemFromRedis(redisKey, &existingItem); err == nil {
		return errors.New("item already exists")
	}

	//Add item to database with JSON Set
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (v *Voters) UpdateVoter(item Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does exist, if it does not, return an error
	redisKey := redisKeyFromId(int(item.VoterID))
	var existingItem Voter
	if err := v.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("item does not exist")
	}

	if _, err := v.jsonHelper.JSONSet(redisKey, ".", item); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (v *Voters) DeleteVoter(id uint64) error {

	pattern := redisKeyFromId(int(id))
	numDeleted, err := v.cacheClient.Del(v.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

// DeleteAll removes all items from the DB.
// It will be exposed via a DELETE /voters endpoint
func (v *Voters) DeleteAll() error {
	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	//Note delete can take a collection of keys.  In go we can
	//expand a slice into individual arguments by using the ...
	//operator
	numDeleted, err := v.cacheClient.Del(v.context, ks...).Result()
	if err != nil {
		return err
	}

	if numDeleted != int64(len(ks)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}

func (v *Voters) GetVoter(id uint64) (Voter, error) {
	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	redisKey := redisKeyFromId(int(id))
	var voterItem Voter
	if err := v.getItemFromRedis(redisKey, &voterItem); err != nil {
		return Voter{}, errors.New("item does not exist")
	}

	return voterItem, nil
}

func (v *Voters) ChangeVoterID(id uint64, value bool) error {

	return errors.New("not implemented")
}

func (v *Voters) GetAllItems() ([]Voter, error) {
	//Now that we have the DB loaded, lets crate a slice
	var voters []Voter
	var voter Voter

	//Now lets iterate over our map and add each item to our slice
	pattern := RedisKeyPrefix + "*"
	ks, _ := v.cacheClient.Keys(v.context, pattern).Result()
	for _, key := range ks {
		err := v.getItemFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voters = append(voters, voter)
	}

	//Now that we have all of our items in a slice, return it
	return voters, nil
}

func (v *Voters) GetPoll(id uint64, pollsid uint64) (voterPoll, error) {
	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	redisKey := redisKeyFromId(int(id))
	var voterItem Voter
	if err := v.getItemFromRedis(redisKey, &voterItem); err != nil {
		return voterPoll{}, errors.New("item does not exist")
	}

	for _, poll := range voterItem.VoteHistory {
		if poll.PollID == pollsid {
			return poll, nil
		}
	}

	return voterPoll{}, errors.New("poll does not exist")
}

func (v *Voters) AddPoll(voterID uint64, pollID uint64) error {
	redisKey := redisKeyFromId(int(voterID))
	var voterItem Voter
	if err := v.getItemFromRedis(redisKey, &voterItem); err != nil {
		return errors.New("item does not exist")
	}

	voteHistory := voterItem.VoteHistory
	for _, poll := range voteHistory {
		if poll.PollID == pollID {
			return errors.New("poll already exists")
		}
	}

	poll := voterPoll{
		PollID:   pollID,
		VoteDate: time.Now(),
	}

	voterItem.VoteHistory = append(voteHistory, poll)

	//Add item to database with JSON Set
	if _, err := v.jsonHelper.JSONSet(redisKey, ".", voterItem); err != nil {
		return err
	}

	return nil
}

func (v *Voters) UpdatePoll(voterID uint64, pollID uint64) error {

	redisKey := redisKeyFromId(int(voterID))
	var voterItem Voter
	if err := v.getItemFromRedis(redisKey, &voterItem); err != nil {
		return errors.New("item does not exist")
	}

	voteHistory := voterItem.VoteHistory
	for i, poll := range voteHistory {
		if poll.PollID == pollID {
			poll.VoteDate = time.Now()
			voteHistory[i] = poll
			voterItem.VoteHistory = voteHistory

			if _, err := v.jsonHelper.JSONSet(redisKey, ".", voterItem); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("poll does not exist")
}

func (v *Voters) DeletePoll(voterID uint64, pollID uint64) error {

	redisKey := redisKeyFromId(int(voterID))
	var voter Voter
	if err := v.getItemFromRedis(redisKey, &voter); err != nil {
		return errors.New("Voter does not exist")
	}

	voteHistory := voter.VoteHistory
	for i, poll := range voteHistory {
		if poll.PollID == pollID {
			newPolls := make([]voterPoll, len(voteHistory)-1)
			copy(newPolls, voteHistory[:i])
			copy(newPolls[i:], voteHistory[i+1:])

			voter.VoteHistory = newPolls
			if _, err := v.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
				return err
			}

			// saveErr := v.saveDB()
			// if saveErr != nil {
			// 	return err
			// }

			return nil
		}
	}

	return errors.New("poll does not exist")
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
// We will use this later, you can ignore for now
func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

// In redis, our keys will be strings, they will look like
// todo:<number>.  This function will take an integer and
// return a string that can be used as a key in redis
func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Helper to return a ToDoItem from redis provided a key
func (t *Voters) getItemFromRedis(key string, item *Voter) error {

	//Lets query redis for the item, note we can return parts of the
	//json structure, the second parameter "." means return the entire
	//json structure
	itemObject, err := t.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	//JSONGet returns an "any" object, or empty interface,
	//we need to convert it to a byte array, which is the
	//underlying type of the object, then we can unmarshal
	//it into our ToDoItem struct
	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return err
	}

	return nil
}
