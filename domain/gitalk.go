package domain

// Gitalk is a struct representing Gitalk's configuration
type Gitalk struct {
	Token  string
	Repo   string
	Owner  string
	Admins []string
}

func (g Gitalk) GetContext() map[string]interface{} {
	return map[string]interface{}{
		"token":  g.Token,
		"repo":   g.Repo,
		"owner":  g.Owner,
		"admins": g.Admins,
	}
}
