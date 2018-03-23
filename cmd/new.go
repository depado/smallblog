package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"

	"github.com/Depado/smallblog/models"
)

var helpmsg = `Generate a new article. The first argument is the filename of the 
article created and should not exist in the configured pages directory. If the
title isn't provided with the appropriate --title flag, it should be added by
hand in the generated file.`

var newCmd = &cobra.Command{
	Use:   "new [filename]",
	Short: "Generate a new article in the pages directory.",
	Long:  helpmsg,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var fd *os.File

		// Check filename first, add .md suffix if necessary and checks the
		// existence of said file
		fname := args[0]
		if !strings.HasSuffix(fname, ".md") {
			fname = fname + ".md"
		}
		fp := filepath.Join(viper.GetString("blog.pages"), fname)
		if _, err := os.Stat(fp); err == nil {
			logrus.WithField("file", fp).Fatal("File already exists")
		}
		if fd, err = os.Create(fp); err != nil {
			logrus.WithError(err).Fatal("Couldn't create file")
		}
		defer fd.Close()

		// Generate struct
		o := models.MetaData{
			Date:        time.Now().Format("2006-01-02 15:04:05"),
			Tags:        viper.GetStringSlice("tags"),
			Author:      viper.GetString("author"),
			Title:       viper.GetString("title"),
			Description: viper.GetString("description"),
		}
		o.Slug = o.GenerateSlug()

		if err = yaml.NewEncoder(fd).Encode(o); err != nil {
			logrus.WithError(err).Fatal("Couldn't write to file")
		}
		logrus.WithField("file", fp).Info("Successfully generated new article")
	},
}

func init() {
	newCmd.Flags().StringSlice("tags", []string{}, "tags for the article")
	newCmd.Flags().String("title", "", "title of the article")
	newCmd.Flags().String("description", "", "description of the article")
	newCmd.Flags().String("slug", "", "slug of the article")
	newCmd.Flags().String("author", "", "author of the article")
	newCmd.Flags().String("banner", "", "banner URL of the article")
	viper.BindPFlags(newCmd.Flags())
}
