package filesystem

import (
	"log"
	"path/filepath"

	"github.com/rjeczalik/notify"

	"github.com/Depado/smallblog/models"
)

// Watch monitors the dir directory for changes and executes actions when a file
// is either created, written, removed, moved to or whatever.
func Watch(dir string) {
	m := make(map[string]*models.Page, len(models.SPages))
	for _, p := range models.SPages {
		m[p.File] = p
	}
	c := make(chan notify.EventInfo, 100)

	if err := notify.Watch(dir, c, notify.Remove, notify.Write, notify.InMovedTo, notify.InMovedFrom); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	var ei notify.EventInfo
	var err error
	for {
		switch ei = <-c; ei.Event() {
		case notify.Write, notify.InMovedTo:
			if p, ok := m[filepath.Base(ei.Path())]; ok {
				// Existing file, needs to be parsed.
				p.UpdateFromFile(ei.Path())
			} else {
				// New file. Needs to be parsed and inserted.
				p := new(models.Page)
				if p, err = models.NewPageFromFile(ei.Path()); err != nil {
					continue
				}
				m[filepath.Base(ei.Path())] = p
			}
		case notify.Remove, notify.InMovedFrom:
			if p, ok := m[filepath.Base(ei.Path())]; ok {
				// The removed file is an active article, needs to be removed.
				p.Pop()
				delete(m, filepath.Base(ei.Path()))
			}
		}
	}

}
