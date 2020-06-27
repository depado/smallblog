package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Depado/smallblog/domain"
	"github.com/Depado/smallblog/models"
)

// BlogRouter holds the necessary data for the router to run properly
type BlogRouter struct {
	Host    string
	Port    int
	Debug   bool
	RootURL string

	Pages        string
	Share        bool
	AnalyticsTag string
	Gitalk       *domain.Gitalk
	r            *gin.Engine
}

// New will create a new BlogRouter with the given arguments
func New(pages, host string, port int, debug bool, rootURL string, gitalk bool,
	token, repo, owner string, admins []string, analytics bool,
	analyticsTag string, share bool) *BlogRouter {

	r := &BlogRouter{
		Pages:   pages,
		Host:    host,
		Port:    port,
		Debug:   debug,
		RootURL: rootURL,
		Share:   share,
	}

	if gitalk {
		r.Gitalk = &domain.Gitalk{
			Owner:  owner,
			Admins: admins,
			Repo:   repo,
			Token:  token,
		}
	}

	if analytics {
		r.AnalyticsTag = analyticsTag
	}

	return r
}

func (br *BlogRouter) GenerateCtx() gin.H {
	c := gin.H{
		"rootURL":     br.RootURL,
		"extra_style": models.GlobCSS,
		"share":       br.Share,
	}
	if br.AnalyticsTag != "" {
		c["analyticsTag"] = br.AnalyticsTag
	}
	if br.Gitalk != nil {
		c["gitalk"] = gin.H{
			"token":  br.Gitalk.Token,
			"repo":   br.Gitalk.Repo,
			"owner":  br.Gitalk.Owner,
			"admins": br.Gitalk.Admins,
		}
	}

	return c
}

func (br *BlogRouter) Start() {
	logrus.WithFields(logrus.Fields{
		"host": br.Host,
		"port": br.Port,
		"root": br.RootURL,
	}).Info("Starting server")
	if err := br.r.Run(fmt.Sprintf("%s:%d", br.Host, br.Port)); err != nil {
		logrus.WithError(err).Fatal("Unable to start server")
	}
}
