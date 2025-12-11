package mongo

import (
	"fmt"
	"log"
	"math"
	"sync"

	"github.com/obfio/tmx-solver-golang/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Object struct {
	CreatedAt int64  `bson:"createdAt"`
	ExpiresAt int64  `bson:"expiresAt"`
	TotalUses int64  `bson:"totalUses"`
	MaxUses   int64  `bson:"maxUses"`
	Name      string `bson:"name"`
	Valid     bool   `bson:"valid"`
	// key = API key (random hex string)
	Key string `bson:"APIKey" json:"APIKey"`
}

var usersColl = config.DbCnx.Collection("keys")

var locker = sync.RWMutex{}
var users = make(map[string]*Object)

// Optimize for reads and force unique
func init() {
	index := mongo.IndexModel{
		Keys:    bson.D{{"APIKey", 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := usersColl.Indexes().CreateOne(config.DbCtx, index)
	if err != nil {
		panic(err)
	}
	log.Println("Optimized/Force unique userID index for profiles collection")
	users, err = GetUsers()
	if err != nil {
		panic(err)
	}
	log.Println("Loaded", len(users), "users")
	go UpdateUsers()

	// populate prints
	pipeline := mongo.Pipeline{
		{{"$project", bson.D{
			{"wgl.wglblob1", 0},
			{"wgl.wglblob2", 0},
			{"canvas.canvasblob", 0},
			{"audio.audioblob", 0},
		}}},
	}

	cursor, err := printsColl.Aggregate(config.DbCtx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(config.DbCtx)

	var results []bson.M
	if err = cursor.All(config.DbCtx, &results); err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(results))
	// decode result[0] to Print
	for _, result := range results {
		b, _ := bson.Marshal(result)
		var print Print
		if err := bson.Unmarshal(b, &print); err != nil {
			log.Fatal(err)
		}
		prints = append(prints, &print)
	}
}

func AddUser(obj *Object) error {
	_, err := usersColl.InsertOne(config.DbCtx, obj)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		}
		return err
	}
	locker.Lock()
	users[obj.Key] = obj
	locker.Unlock()
	return nil
}

func UpdateUser(APIKey string, totalNewSessions int64) error {
	filter := bson.M{"APIKey": APIKey}
	update := bson.M{
		"$set": bson.M{
			"totalUses": totalNewSessions,
		},
	}
	_, err := usersColl.UpdateOne(config.DbCtx, filter, update)
	if err != nil {
		return err
	}
	users[APIKey].TotalUses = totalNewSessions
	return nil
}

func GetUsers() (map[string]*Object, error) {
	out := make(map[string]*Object)
	cur, err := usersColl.Find(config.DbCtx, bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	for cur.Next(config.DbCtx) {
		var elem *Object
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		out[elem.Key] = elem
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteUser by key
func DeleteUser(APIKey string) error {
	filter := bson.M{"APIKey": APIKey}
	_, err := usersColl.DeleteOne(config.DbCtx, filter)
	if err != nil {
		return err
	}
	locker.Lock()
	delete(users, APIKey)
	locker.Unlock()
	return nil
}

// UpdateTime
func UpdateTime(APIKey string, time int64) error {
	filter := bson.M{"APIKey": APIKey}
	update := bson.M{
		"$set": bson.M{
			"expiresAt": time,
		},
	}
	_, err := usersColl.UpdateOne(config.DbCtx, filter, update)
	if err != nil {
		return err
	}
	locker.Lock()
	users[APIKey].ExpiresAt = time
	locker.Unlock()
	return nil
}

// UpdateMaxUses
func UpdateMaxUses(APIKey string, maxUses int64) error {
	if maxUses == -1 {
		maxUses = math.MaxInt64
	}
	filter := bson.M{"APIKey": APIKey}
	update := bson.M{
		"$set": bson.M{
			"maxUses":   maxUses,
			"totalUses": 0,
		},
	}
	_, err := usersColl.UpdateOne(config.DbCtx, filter, update)
	if err != nil {
		return err
	}
	locker.Lock()
	users[APIKey].MaxUses = maxUses
	users[APIKey].TotalUses = 0
	locker.Unlock()
	return nil
}

// UpdateName
func UpdateName(APIKey string, name string) error {
	filter := bson.M{"APIKey": APIKey}
	update := bson.M{
		"$set": bson.M{
			"name": name,
		},
	}
	_, err := usersColl.UpdateOne(config.DbCtx, filter, update)
	if err != nil {
		return err
	}
	locker.Lock()
	users[APIKey].Name = name
	locker.Unlock()
	return nil
}
