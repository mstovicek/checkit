package git

type Client interface {
	Checkout(cloneUrl string, server string, repositoryName string, commitHash string) (string, error)
}
