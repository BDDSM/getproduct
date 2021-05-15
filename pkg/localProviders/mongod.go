package localProviders

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/korableg/getproduct/pkg/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "getproduct"
const collectionName = "products"

type MongoDB struct {
	hostname string
	port     int
	username string
	password string
}

func NewMongoDB(hostname string, port int, username, password string) (*MongoDB, error) {
	return newMongo(hostname, port, username, password, collectionName)
}

func newMongo(hostname string, port int, username, password, colName string) (*MongoDB, error) {

	m := MongoDB{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
	}

	err := m.initDB(colName)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *MongoDB) initDB(colName string) error {

	ctx := context.Background()

	connect, err := m.connect(ctx)
	if err != nil {
		return err
	}

	err = connect.Database(dbName).CreateCollection(ctx, colName)

	if err != nil && err.(mongo.ServerError).HasErrorCode(48) {
		return nil
	} else if err != nil {
		return err
	}

	unicue := true
	idxName := fmt.Sprintf("%s_barcode", colName)
	idxOpts := options.IndexOptions{
		Name:   &idxName,
		Unique: &unicue,
	}
	idx := mongo.IndexModel{Keys: bson.M{"barcode": 1}, Options: &idxOpts}
	collection := connect.Database(dbName).Collection(colName)
	_, err = collection.Indexes().CreateOne(ctx, idx)
	if err != nil {
		return err
	}

	return nil

}

func (m *MongoDB) AddProduct(ctx context.Context, p *product.Product) error {
	return m.addProduct(ctx, p, collectionName)
}

func (m *MongoDB) addProduct(ctx context.Context, p *product.Product, colName string) error {
	client, err := m.connect(ctx)
	if err != nil {
		return err
	}
	defer m.disconnect(ctx, client)

	pJson, err := json.Marshal(p)
	if err != nil {
		return err
	}

	collection := client.Database(dbName).Collection(colName)

	upsert := true
	returnDocument := options.After
	opts := options.FindOneAndReplaceOptions{
		Upsert:         &upsert,
		ReturnDocument: &returnDocument,
	}
	filter := bson.M{"barcode": p.Barcode()}
	data := bson.M{"barcode": p.Barcode(), "value": pJson}

	res := collection.FindOneAndReplace(ctx, filter, data, &opts)

	return res.Err()
}

func (m *MongoDB) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {
	return m.getProduct(ctx, barcode, collectionName)
}

func (m *MongoDB) getProduct(ctx context.Context, barcode string, colName string) (*product.Product, error) {

	client, err := m.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer m.disconnect(ctx, client)

	filter := bson.M{"barcode": barcode}
	collection := client.Database(dbName).Collection(colName)

	mongoData := make(map[string]interface{})

	err = collection.FindOne(ctx, filter).Decode(&mongoData)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	p := product.Product{}
	err = json.Unmarshal(mongoData["value"].(primitive.Binary).Data, &p)
	if err != nil {
		return nil, err
	}

	return &p, nil

}

func (m *MongoDB) DeleteProduct(ctx context.Context, p *product.Product) error {
	return m.deleteProduct(ctx, p, collectionName)
}

func (m *MongoDB) deleteProduct(ctx context.Context, p *product.Product, colName string) error {

	client, err := m.connect(ctx)
	if err != nil {
		return err
	}

	defer m.disconnect(ctx, client)

	filter := bson.M{"barcode": p.Barcode()}

	_, err = client.Database(dbName).Collection(colName).DeleteOne(ctx, filter)

	return err

}

func (m *MongoDB) connect(ctx context.Context) (*mongo.Client, error) {

	var userInfo *url.Userinfo
	if m.username != "" {
		userInfo = url.UserPassword(m.username, m.password)
	}

	url := url.URL{
		Scheme: "mongodb",
		User:   userInfo,
		Host:   fmt.Sprintf("%s:%d", m.hostname, m.port),
	}

	urlRaw := url.String()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(urlRaw))
	if err != nil {
		return nil, err
	}

	return client, err

}

func (m *MongoDB) disconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		log.Println(err)
	}
}
