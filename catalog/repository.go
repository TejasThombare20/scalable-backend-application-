package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	ErrNotFound = errors.New("entity Not Found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, product Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewEsasticRepository(url string) (Repository, error) {

	url1 := "http://catalog_db:9200"
	client, err := elastic.NewClient(
		elastic.SetURL(url1),
		elastic.SetSniff(false),
	)

	if err != nil {
		return nil, err
	}

	return &elasticRepository{client}, nil

}

func (r *elasticRepository) Close() {

}

func (r *elasticRepository) PutProduct(ctx context.Context, product Product) error {
	_, err := r.client.Index().Index("catalog").Type("product").Id(product.ID).BodyJson(productDocument{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}).Do(ctx)

	return err
}

func (r *elasticRepository) GetProductByID(ctx context.Context, id string) (*Product, error) {

	response, err := r.client.Get().Index("catalog").Type("product").Id(id).Do(ctx)

	if err != nil {
		return nil, err
	}

	if !response.Found {
		return nil, err
	}

	product := productDocument{}

	if err := json.Unmarshal(*response.Source, &product); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, err

}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	respose, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMatchAllQuery()).From(int(skip)).Size(int(take)).Do(ctx)

	if err != nil {
		return nil, err
	}

	products := []Product{}

	for _, hit := range respose.Hits.Hits {
		product := productDocument{}

		if err := json.Unmarshal(*hit.Source, &product); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}

	return products, err
}

func (r *elasticRepository) ListProductWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	items := []*elastic.MultiGetItem{}

	for _, id := range ids {

		items = append(items, elastic.NewMultiGetItem().Index("catalog").Type("product").Id(id))

		log.Println("items in catlog repository", items)
	}

	response, err := r.client.MultiGet().Add(items...).Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}

	log.Println("response", response)

	for _, doc := range response.Docs {
		product := productDocument{}

		if err = json.Unmarshal(*doc.Source, &product); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
		if err != nil {
			log.Println("error while appending products in Es", err)
		}
	}

	return products, nil

}
func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	response, err := r.client.Search().Index("catalog").Type("product").Query(elastic.NewMultiMatchQuery(query, "name", "description ")).From(int(skip)).Size(int(take)).Do(ctx)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []Product{}

	for _, hit := range response.Hits.Hits {
		product := productDocument{}

		if err = json.Unmarshal(*hit.Source, &product); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}

	return products, err

}
