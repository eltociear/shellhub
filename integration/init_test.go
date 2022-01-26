package integration

import (
	"context"
	"fmt"
	"log"
    "flag"
	"os"
	"testing"

	cmd "github.com/shellhub-io/shellhub/integration/commands"
	"github.com/shellhub-io/shellhub/integration/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Cloud = iota
	Enterprise
	Community
)

var DBTest *mongo.Database
const BaseURL = "http://127.0.0.1/api/"
var Logins map[string]*utils.Login
var version string

func InitializeDatabase() {
	ip := cmd.MongoIP()

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", ip))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		log.Println(err)
	}

	DBTest = client.Database("test")
}

func init() {
    flag.StringVar(&version, "version", "", "shellhub instance")
	Logins = make(map[string]*utils.Login, 0)
}

func LoadAPIService() {
    service := "api"
	params := []string{}

    versionMap := map[string]int {
        "cloud": Cloud,
        "community": Community,
        "enterprise": Enterprise,
    }

    numVersion, ok := versionMap[version]

    if !ok {
        log.Println("Choose a shellhub instance for test evaluation an param")
        log.Println("-version=version_chosen")
        log.Println("Available options: ")
        log.Println("\tcloud")
        log.Println("\tenterprise")
        log.Println("\tcommunity")
        os.Exit(1)
    }

	if !cmd.IsInstanceAlive(service, numVersion, params) {
		cmd.InitializeTestAPI(numVersion, params)
	} else {
		log.Println("test instance is already running")
	}


    // if the instance is cloud or enterprise manage to 
    // up required private project apis communication
    // path the cloud project path related apis
}

func InitData() { // Generate logins map with fake data
    user1 := utils.NewUser(DBTest, 1)
    user2 := utils.NewUser(DBTest, 2)

    namespace1 := utils.NewNamespace(DBTest, 1, user1, []utils.User{*user2}, []string{"administrator"})

    login1, err := utils.NewLogin(DBTest, namespace1, user2)
    if err != nil {
        log.Println(err)
        return
    }

	Logins[login1.Name] = login1
}

func setup() {
	log.Println("setting up...")

    // get version from given context
    flag.Parse()

    // initialize api service
    LoadAPIService()

	// initializes mongodb connection
	InitializeDatabase()

	// populates database
	InitData()
}

func teardown() {
	log.Println("tearing down...")

	collections := []string{"namespaces", "devices", "users"}
	for _, c := range collections {
		_, _ = DBTest.Collection(c, nil).DeleteMany(context.TODO(), bson.M{})
	}
}

func TestMain(m *testing.M) {
	setup()
	val := m.Run()
	teardown()
	os.Exit(val)
}
