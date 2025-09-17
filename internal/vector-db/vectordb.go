package vectordb

import (
	"context"
	"fmt"
	"log"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	g "github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/gusram01/linked-bookmarks/internal/platform/config"
)

type VectorDB struct {
	Client     chroma.Client
	Collection chroma.Collection
}

var VDB *VectorDB

func Initialize() {
	ctx := context.Background()

	var baseUrl = fmt.Sprintf("http://%s:%s", config.ENVS.VectorDBHost, config.ENVS.VectorDBPort)

	c, err := chroma.NewHTTPClient(
		chroma.WithBaseURL(baseUrl),
	)

	if err != nil {
		log.Fatalf("Failed to create Chroma client: %v", err)
	}

	ef, err := g.NewGeminiEmbeddingFunction(g.WithAPIKey(config.ENVS.GeminiApiKey), g.WithDefaultModel("gemini-embedding-001"))
	if err != nil {
		log.Fatalf("Error creating Gemini embedding function: %s \n", err)
	}

	cll, err := c.GetOrCreateCollection(ctx, config.ENVS.VectorDBCollectionName, chroma.WithEmbeddingFunctionCreate(ef))

	if err != nil {
		log.Fatalf("Failed to get or create collection: %v", err)
	}

	VDB = &VectorDB{
		Client:     c,
		Collection: cll,
	}

}

func (vdb *VectorDB) Shutdown() {

	defer func() {
		err := vdb.Client.Close()

		if err != nil {
			log.Fatalf("Failed to close Chroma client: %v", err)
		}

	}()

}
