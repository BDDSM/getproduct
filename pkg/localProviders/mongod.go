package localProviders

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/korableg/getproduct/pkg/product"
	"go.mongodb.org/mongo-driver/bson"
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

func New(hostname string, port int, username, password string) (*MongoDB, error) {

	m := MongoDB{
		hostname: hostname,
		port:     port,
		username: username,
		password: password,
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// connect, err := m.connect(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// err = connect.Database(dbName).CreateCollection(ctx, collectionName)
	// if err != nil {
	// 	return nil, err
	// }

	return &m, nil

}

func (m *MongoDB) AddProduct(ctx context.Context, p *product.Product) error {
	client, err := m.connect(ctx)
	if err != nil {
		return err
	}
	defer m.disconnect(ctx, client)

	collection := client.Database(dbName).Collection(collectionName)
	_, err = collection.InsertOne(ctx, bson.D{{"barcode", p.Barcode()}, {"value", p}})

	return err

}

func (m *MongoDB) GetProduct(ctx context.Context, barcode string) (*product.Product, error) {

	client, err := m.connect(ctx)
	if err != nil {
		return nil, err
	}

	defer m.disconnect(ctx, client)

	filter := bson.D{{"barcode", barcode}}
	collection := client.Database(dbName).Collection(collectionName)

	p := product.Product{}

	err = collection.FindOne(ctx, filter).Decode(&p)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &p, nil

}

func (m *MongoDB) DeleteProduct(ctx context.Context, product *product.Product) error {
	return nil
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
