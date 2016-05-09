version = 1.0.2

.PHONY: release clean

release:
	-mkdir -p releases/smallblog-linux-amd64-${version}/pages
	-mkdir -p releases/smallblog-linux-arm-${version}/pages
	cp -r assets/ releases/smallblog-linux-amd64-${version}
	cp -r assets/ releases/smallblog-linux-arm-${version}
	cp -r templates/ releases/smallblog-linux-amd64-${version}
	cp -r assets/ releases/smallblog-linux-arm-${version}
	cp examples/article.md.example releases/smallblog-linux-amd64-${version}/pages/article.md
	cp examples/article.md.example releases/smallblog-linux-arm-${version}/pages/article.md
	cp examples/conf.yml.example releases/smallblog-linux-amd64-${version}/conf.yml
	cp examples/conf.yml.example releases/smallblog-linux-arm-${version}/conf.yml
	go build -o releases/smallblog-linux-amd64-${version}/smallblog
	GOARCH=arm go build -o releases/smallblog-linux-arm-${version}/smallblog
	tar czf releases/smallblog-linux-amd64-${version}.tar.gz releases/smallblog-linux-amd64-${version}
	tar czf releases/smallblog-linux-arm-${version}.tar.gz releases/smallblog-linux-arm-${version}

clean:
	-rm -r releases/
