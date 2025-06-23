package loaders

// import vikstrous/dataloadgen with your other imports
import (
	"context"
	"graphql-backend/app"
	"graphql-backend/graph/model"
	trans "graphql-backend/transport"
	"net/http"
	"time"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	UserLoader    *dataloadgen.Loader[string, *model.User]
	ProductLoader *dataloadgen.Loader[string, *model.Product]
}

type reader struct {
	repo app.Repo
}

// getUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *reader) getUsers(ctx context.Context, ids []string) ([]*model.User, []error) {
	users, err := u.repo.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	res := trans.UsersRes{}
	res.Bind(users)

	return res.Res, nil
}

// getProducts implements a batch function that can retrieve many products by ID,
// for use in a dataloader
func (u *reader) getProducts(ctx context.Context, ids []string) ([]*model.Product, []error) {
	products, err := u.repo.GetProductsByIDs(ctx, ids)
	if err != nil {
		return nil, []error{err}
	}

	res := trans.ProductsRes{}
	res.Bind(products)

	return res.Res, nil
}

func NewLoaders(repo app.Repo) *Loaders {
	// define the data loader
	ur := &reader{repo: repo}
	return &Loaders{
		UserLoader:    dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
		ProductLoader: dataloadgen.NewLoader(ur.getProducts, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(next http.Handler, repo app.Repo) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(repo)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

func GetUser(ctx context.Context, id string) (*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, id)
}

func GetProduct(ctx context.Context, id string) (*model.Product, error) {
	loaders := For(ctx)
	return loaders.ProductLoader.Load(ctx, id)
}

func GetUsers(ctx context.Context, userIDs []string) ([]*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}

func GetProducts(ctx context.Context, ids []string) ([]*model.Product, error) {
	loaders := For(ctx)
	return loaders.ProductLoader.LoadAll(ctx, ids)
}
