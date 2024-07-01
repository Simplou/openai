package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Simplou/goxios"
)

type EmbeddingRequest[Input string | []string] struct {
	// Input text to embed, encoded as a string or array of tokens.
	Input Input `json:"input"`
	// ID of the model to use.
	Model string `json:"model"`
	// The format to return the embeddings in. Can be either float or base64.
	Encoding string `json:"encoding_format,omitempty"`
	// The number of dimensions the resulting output embeddings should have. Only supported in text-embedding-3 and later models.
	Dimensions int `json:"dimensions,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

type (
	Base64 string
	
	Embedding[Encoding []float64 | Base64] struct {
		Object    string   `json:"object"`
		Embedding Encoding `json:"embedding"`
		Index     int      `json:"index"`
	}

	EmbeddingResponse[Encoding []float64 | Base64] struct {
		Object string                `json:"object"`
		Data   []Embedding[Encoding] `json:"data"`
		Model  string                `json:"model"`
		Usage  Usage                 `json:"usage"`
	}
)

// CreateEmbedding sends a request to create embeddings for the given input.
func CreateEmbedding[Input string | []string, Encoding []float64 | Base64](api OpenAIClient, httpClient HTTPClient, body *EmbeddingRequest[Input]) (*EmbeddingResponse[Encoding], *OpenAIErr) {
	api.AddHeader(contentTypeJSON)
	b, err := json.Marshal(body)
	if err != nil {
		return nil, errCannotMarshalJSON(err)
	}
	options := goxios.RequestOpts{
		Headers: api.Headers(),
		Body:    ioReader(b),
	}
	res, err := httpClient.Post(api.BaseURL()+"/embeddings", &options)
	if err != nil {
		return nil, errCannotSendRequest(err)
	}
	if res.StatusCode >= http.StatusBadRequest {
		return nil, openaiHttpError(res)
	}
	response := new(EmbeddingResponse[Encoding])
	if err := goxios.DecodeJSON(res.Body, response); err != nil {
		return nil, errCannotDecodeJSON(err)
	}
	if err := res.Body.Close(); err != nil {
		return nil, errCloseBody(err)
	}
	return response, nil
}

type ChunkTextOpts struct {
	Text      string
	ChunkSize int
}

// ChunkText splits the input text into chunks of specified size.
func ChunkText(opts ChunkTextOpts) []string {
	chunkSize := opts.ChunkSize
	if chunkSize <= 0 {
		chunkSize = 512
	}

	words := strings.Fields(opts.Text)
	var chunks []string

	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)
	}

	return chunks
}

// ChunksSummary returns the summary of a randomly selected relevant chunk
func ChunksSummary(client OpenAIClient, httpClient HTTPClient, relevantChunks []string, query string) (string, error) {
	if len(relevantChunks) == 0 {
		return "", errors.New("no relevant chunks provided")
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := rand.Intn(len(relevantChunks))
	randomChunk := relevantChunks[randomIndex]

	summary, err := chunkSummary(client, httpClient, randomChunk, query)
	if err != nil {
		return "", err
	}

	return summary, nil
}

// chunkSummary generates a summary for the given chunk and query.
func chunkSummary(client OpenAIClient, httpClient HTTPClient, chunk, query string) (string, error) {
	response, err := ChatCompletion(client, httpClient, &CompletionRequest[DefaultMessages]{
		Model: "gpt-3.5-turbo",
		Messages: DefaultMessages{
			{Role: "system", Content: "Você deve resumir a resposta para a pergunta do usuário usando o conteúdo de forma clara e concisa."},
			{Role: "user", Content: fmt.Sprintf("Desenvolva uma resposta curta e clara para a pergunta %s baseada no seguinte conteúdo: %s", query, chunk)},
		},
	})
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}

// FindMostRelevantEmbeddings finds the most relevant embeddings.
// q is the query embedding, and e is the matrix of embeddings to search in.
func FindMostRelevantEmbeddings(q, e [][]float64) ([]int, error) {
	if len(q) == 0 || len(e) == 0 || len(q[0]) == 0 || len(e[0]) == 0 {
		return nil, fmt.Errorf("input matrices cannot be empty")
	}

	qNorm := normalize(q[0])
	eNorm := make([][]float64, len(e))

	for i := range e {
		eNorm[0] = normalize(e[i])
	}

	similarities := make([]float64, len(e))
	for i := range eNorm {
		similarities[i] = dotProduct(qNorm, eNorm[i])
	}

	type similarityIndex struct {
		similarity float64
		index      int
	}
	similarityIndices := make([]similarityIndex, len(similarities))
	for i, sim := range similarities {
		similarityIndices[i] = similarityIndex{similarity: sim, index: i}
	}

	sort.Slice(similarityIndices, func(i, j int) bool {
		return similarityIndices[i].similarity > similarityIndices[j].similarity
	})

	topIndices := make([]int, len(similarities))
	for i := 0; i < len(similarities); i++ {
		topIndices[i] = similarityIndices[i].index
	}

	return topIndices, nil
}

// normalize computes the Euclidean norm of a slice and returns the normalized slice.
func normalize(vec []float64) []float64 {
	normVal := norm(vec)
	normalizedVec := make([]float64, len(vec))
	for i, v := range vec {
		normalizedVec[i] = v / normVal
	}
	return normalizedVec
}

// norm computes the Euclidean norm of a slice.
func norm(vec []float64) float64 {
	var sum float64
	for _, v := range vec {
		sum += v * v
	}
	return math.Sqrt(sum)
}

// dotProduct computes the dot product of two slices.
func dotProduct(vec1, vec2 []float64) float64 {
	var sum float64
	for i, value := range vec1 {
		sum += value * vec2[i]
	}
	return sum
}
