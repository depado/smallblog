package filesystem

import (
	"log"
	"path/filepath"
	"sort"

	"github.com/Depado/smallblog/models"
	"github.com/rjeczalik/notify"
)

func Watch(dir string) {
	m := make(map[string]*models.Page, len(models.SPages))
	for _, p := range models.SPages {
		m[p.File] = p
	}
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	if err := notify.Watch(dir, c, notify.Remove, notify.Write); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	var ei notify.EventInfo
	var err error
	for {
		switch ei = <-c; ei.Event() {
		case notify.Write:
			if p, ok := m[filepath.Base(ei.Path())]; ok {
				// Existing file, needs to be parsed.
				if err = p.ParseFile(ei.Path()); err != nil {
					log.Printf("[SB] [ ERR] [%s] Could not process modified file : %s\n", filepath.Base(ei.Path()), err)
					continue
				}
				sort.Sort(models.SPages)
				log.Printf("[SB] [INFO] [%s] Modified [/post/%s] %s\n", p.File, p.Slug, p.Title)
			} else {
				// This is a new file. Needs to be processed and inserted.
				p := models.Page{}
				if err = p.ParseFile(ei.Path()); err != nil {
					log.Printf("[SB] [ ERR] [%s] Could not process new file : %s\n", filepath.Base(ei.Path()), err)
					continue
				}
				if err = p.Insert(false); err != nil {
					log.Printf("[SB] [ ERR] [%s] Could not add new file : %s\n", filepath.Base(ei.Path()), err)
					continue
				}
				log.Printf("[SB] [INFO] [%s] New file [/post/%s] %s\n", p.File, p.Slug, p.Title)
				m[filepath.Base(ei.Path())] = &p
			}
		case notify.Remove:
			log.Println("Removed", ei.Path())
		}
	}

}
