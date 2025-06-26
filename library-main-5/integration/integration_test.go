//go:build integration

package integration

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAuthorGrpc(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	const authorName = "Test testovich"

	registerRes, err := client.RegisterAuthor(ctx, &RegisterAuthorRequest{
		Name: authorName,
	})

	require.NoError(t, err)
	authorID := registerRes.GetId()

	author, err := client.GetAuthorInfo(ctx, &GetAuthorInfoRequest{
		Id: authorID,
	})
	require.NoError(t, err)

	require.Equal(t, authorName, author.GetName())
	require.Equal(t, authorID, author.GetId())

	_, err = client.ChangeAuthorInfo(ctx, &ChangeAuthorInfoRequest{
		Id:   authorID,
		Name: authorName + "123",
	})
	require.NoError(t, err)

	newAuthor, err := client.GetAuthorInfo(ctx, &GetAuthorInfoRequest{
		Id: authorID,
	})
	require.NoError(t, err)

	require.Equal(t, authorName+"123", newAuthor.GetName())
	require.Equal(t, authorID, newAuthor.GetId())
}

func TestBookGrpc(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	const (
		authorName = "Test testovich"
		bookName   = "go"
	)

	registerRes, err := client.RegisterAuthor(ctx, &RegisterAuthorRequest{
		Name: authorName,
	})

	require.NoError(t, err)
	authorID := registerRes.GetId()

	response, err := client.AddBook(ctx, &AddBookRequest{
		Name:     bookName,
		AuthorId: []string{authorID},
	})
	require.NoError(t, err)

	book := response.GetBook()

	require.Equal(t, bookName, book.GetName())
	require.Equal(t, 1, len(book.GetAuthorId()))
	require.Equal(t, authorID, book.GetAuthorId()[0])

	_, err = client.UpdateBook(ctx, &UpdateBookRequest{
		Id:   book.GetId(),
		Name: bookName + "-2024",
	})
	require.NoError(t, err)

	newBook, err := client.GetBookInfo(ctx, &GetBookInfoRequest{
		Id: book.GetId(),
	})

	require.NoError(t, err)
	require.Equal(t, bookName+"-2024", newBook.GetBook().GetName())
	require.Equal(t, 1, len(newBook.GetBook().GetAuthorId()))
	require.Equal(t, authorID, newBook.GetBook().GetAuthorId()[0])

	books := getAllAuthorBooks(t, authorID, client)

	require.NoError(t, err)
	require.Equal(t, 1, len(books))

	require.Equal(t, newBook.GetBook().GetName(), books[0].GetName())
	require.Equal(t, newBook.GetBook().GetAuthorId(), books[0].GetAuthorId())
}

func TestBookManyAuthorsGrpc(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	var (
		authorBasicName = "Donald Knuth"
		authorsCount    = 10
		bookName        = "The Art of Computer Programming"
	)

	authorIds := make([]string, authorsCount)
	for i := range authorsCount {
		author, err := client.RegisterAuthor(ctx, &RegisterAuthorRequest{
			Name: authorBasicName + strconv.Itoa(rand.N[int](10e9)),
		})
		require.NoError(t, err)
		authorIds[i] = author.Id
	}

	bookAdded, err := client.AddBook(ctx, &AddBookRequest{
		Name:     bookName,
		AuthorId: authorIds,
	})
	require.NoError(t, err)
	require.ElementsMatch(t, bookAdded.Book.AuthorId, authorIds)

	bookReceived, err := client.GetBookInfo(ctx, &GetBookInfoRequest{
		Id: bookAdded.Book.Id,
	})
	require.NoError(t, err)
	require.EqualExportedValues(t, bookAdded.Book, bookReceived.Book)
}

func TestConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	var (
		authorName = "Test testovich" + strconv.Itoa(rand.N[int](10e9))
		totalBooks = 1234
		workers    = 50
	)

	registerRes, err := client.RegisterAuthor(ctx, &RegisterAuthorRequest{
		Name: authorName,
	})

	require.NoError(t, err)
	authorID := registerRes.GetId()

	books := make([]string, 0, totalBooks)
	for i := range totalBooks {
		books = append(books, strconv.Itoa(i))
	}

	perWorker := totalBooks / workers
	start := 0

	wg := new(sync.WaitGroup)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(s int) {
			defer wg.Done()

			right := s + perWorker
			if i == workers-1 {
				right = len(books)
			}

			for b := s; b < right; b++ {
				_, err := client.AddBook(ctx, &AddBookRequest{
					Name:     books[b],
					AuthorId: []string{authorID},
				})
				require.NoError(t, err)
			}
		}(start)

		start += perWorker
	}

	wg.Wait()

	authorBooks := lo.Map(getAllAuthorBooks(t, authorID, client), func(item *Book, index int) string {
		return item.GetName()
	})

	slices.Sort(authorBooks)
	slices.Sort(books)

	require.Equal(t, books, authorBooks)
}

func TestAuthorNotFound(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	_, err := client.GetAuthorInfo(ctx, &GetAuthorInfoRequest{
		Id: uuid.New().String(),
	})

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

func TestAuthorInvalidArgument(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	_, err := client.GetAuthorInfo(ctx, &GetAuthorInfoRequest{
		Id: "123",
	})

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

func TestBookNotFound(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	_, err := client.GetBookInfo(ctx, &GetBookInfoRequest{
		Id: uuid.New().String(),
	})

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, s.Code())
}

func TestBookInvalidArgument(t *testing.T) {
	ctx := context.Background()
	client := newGRPCClient(t)

	_, err := client.GetBookInfo(ctx, &GetBookInfoRequest{
		Id: "123",
	})

	s, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, s.Code())
}

func TestGrpcGateway(t *testing.T) {
	type RegisterAuthorResponse struct {
		ID string `json:"id"`
	}

	type GetAuthorResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	registerUrl := fmt.Sprintf("http://127.0.0.1:%s/v1/library/author",
		cmp.Or(os.Getenv("GRPC_GATEWAY_PORT"), "8080"))

	request, err := http.NewRequest("POST", registerUrl, strings.NewReader(`{"name": "Name"}`))
	require.NoError(t, err)

	response, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	var registerAuthorResponse RegisterAuthorResponse

	err = json.Unmarshal(data, &registerAuthorResponse)
	require.NoError(t, err)

	require.NotEmpty(t, registerAuthorResponse)

	getUrl := fmt.Sprintf("http://127.0.0.1:%s/v1/library/author_info/%s",
		cmp.Or(os.Getenv("GRPC_GATEWAY_PORT"), "8080"), registerAuthorResponse.ID)

	getRequest, err := http.NewRequest("GET", getUrl, nil)
	require.NoError(t, err)

	getResponse, err := http.DefaultClient.Do(getRequest)
	require.NoError(t, err)

	getData, err := io.ReadAll(getResponse.Body)
	require.NoError(t, err)

	var author GetAuthorResponse
	err = json.Unmarshal(getData, &author)
	require.NoError(t, err)

	require.Equal(t, author.ID, registerAuthorResponse.ID)
	require.Equal(t, author.Name, "Name")
}

func TestGrpcGatewayUnknownUrl(t *testing.T) {
	unknownUrl := fmt.Sprintf("http://127.0.0.1:%s/v0/not_library/not_author_info",
		cmp.Or(os.Getenv("GRPC_GATEWAY_PORT"), "8080"))

	response, err := http.Get(unknownUrl)

	require.NoError(t, err)
	require.Equal(t, response.StatusCode, 404)
}

func newGRPCClient(t *testing.T) LibraryClient {
	t.Helper()

	addr := "127.0.0.1:" + cmp.Or(os.Getenv("GRPC_PORT"), "9090")
	c, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	return NewLibraryClient(c)
}

func getAllAuthorBooks(t *testing.T, authorID string, client LibraryClient) []*Book {
	t.Helper()
	ctx := context.Background()

	result := make([]*Book, 0)
	stream, err := client.GetAuthorBooks(ctx, &GetAuthorBooksRequest{
		AuthorId: authorID,
	})
	require.NoError(t, err)

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			return result
		}

		require.NoError(t, err)

		result = append(result, resp)
	}
}
