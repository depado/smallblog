# Smallblog

![Go Version](https://img.shields.io/badge/go-1.9-brightgreen.svg)
![Go Version](https://img.shields.io/badge/go-1.10-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/smallblog)](https://goreportcard.com/report/github.com/Depado/smallblog)
[![Build Status](https://drone.depado.eu/api/badges/Depado/smallblog/status.svg)](https://drone.depado.eu/Depado/smallblog)
[![codecov](https://codecov.io/gh/Depado/smallblog/branch/master/graph/badge.svg)](https://codecov.io/gh/Depado/smallblog)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/smallblog/blob/master/LICENSE)

A simple self-hosted markdown flat files blog

## What is Smallblog

The main goal of this project is to show how easily you can develop a flat file
blog with markdown as the primary writing language. It's not perfect, it will
never be, some people are already doing great things based on that idea, like
[Hugo](https://gohugo.io/) for instance. Let's note though, that's **not** a
static website generator.

## Disclaimer

Smallblog is a quick project. The goal is to take some markdown files with
front matter headers (written in yaml), parse them, render them and store the
resulting HTML in memory. Which means, the more markdown files you have, the
more memory smallblog will consume.

## Features

- Filesystem monitoring. Drop a new file or modify a file, bam, your blog is
updated.
- Automatic syntax highlighting using [bfchroma](https://github.com/Depado/bfchroma)
(which uses [chroma](https://github.com/alecthomas/chroma) under the hood)
- No CGO dependencies
- Tag system
- Simple and customizable template and CSS for easy reading
- Comments using [gitalk](https://github.com/gitalk/gitalk)

## TODO

- [] Add a search engine

## Tutorial

First of all you'll need a `pages` directory where you'll store all your
articles, so let's create that.

```sh
$ mkdir pages
```

Without configuration, smallblog will try to find a `pages/` directory in your 
current working directory. So let's create your very first article, and for that 
we'll be using the `new` command :

```sh
$ mkdir pages
$ smallblog new first-article --title "My very first Smallblog article" --draft
INFO[0000] Successfully generated new article            file=pages/first-article.md
```

This commands takes a single argument : The name of the file this command will
write. It will automatically append the `.md` suffix if you don't specify it
and will place it in the configured (or default) `pages` directory. Also we're
specifying that this article is a draft. A draft can only be visible if the
server was started with the `blog.draft` is set to true.

To check if this is working, simply run the following command :

```sh
$ smallblog serve --blog.draft                
INFO[0000] Generated files                               files=1 took="201.207µs"
INFO[0000] Starting server                               host=127.0.0.1 port=8080
```

That's it, you successfully generated and served your first article.

## Configure

You can add a `conf.yml` or `conf.json` in your working directory to make your
configuration persistent. You may have noticed the warning about the
configuration file that wasn't found earlier. Here is an example `conf.yml`
file :

```yaml
server:
  host: 127.0.0.1
  port: 8080
  debug: false
blog:
  title: Depado's Blog
  description: A simple blog from a developer who does things.
  pages: pages/
  code:
    style: monokai
  draft: true
gitalk:
  enabled: true
  client_id: xxxxx
  client_secret: xxxx
  repo: articles
  owner: Depado
  admins: [Depado]
log:
  format: text
  level: info
  line: true
```

Explanations of these fields will follow. Just note that all these values
can be customized using command line flags which have a higher priority than
the configuration file. For example you might want to turn on the `debug` mode
of the server without changing your configuration file. In which case you can
just pass the flag `--server.debug` when starting smallblog.

### Server

Server related configuration. Defines on which host/port the server should
listen as well as debug mode.

| Field   | Description                            | Default     |
|---------|----------------------------------------|-------------|
| `host`  | Host on which the server should listen | "127.0.0.1" |
| `port`  | Port on which the server should listen | 8080        |
| `debug` | Enable debug mode for the router       | false       |

The `debug` value is especially useful if you modify the HTML templates of
smallblog since you won't need a restart of the service to see your changes
applied. 

### Blog

General blog configuration such as the blog title, description, pages directory
and code highlighting style.

| Field         | Description                               | Default   |
|---------------|-------------------------------------------|-----------|
| `title`       | Title of your blog                        | -         |
| `description` | Description of your blog                  | -         |
| `pages`       | Directory where your articles are stored  | "pages/"  |
| `code.style`  | Style of the syntax highlighting          | "monokai" |
| `draft`       | Display articles that are marked as draft | false     |

### Gitalk

Please see the documentation of [gitalk](https://github.com/gitalk/gitalk) to
see how you can configure gitalk for smallblog.

### Log

| Field    | Description                                                                                                | Default |
|----------|------------------------------------------------------------------------------------------------------------|---------|
| `format` | Log format, either "text" or "json"                                                                        | "text"  |
| `level`  | Defines the minimum level of logging that is displayed. One of "debug", "info", "warn", "error" or "fatal" | "info"  |
| `line`   | Show where the log happened in the source code (filename, line number)                                     | false   |

## Write Posts

There is no naming convention for file names. You can name them whatever you
want, it won't change the server's behavior. A post (or page/article) file is
divided in two parts. The first part is yaml data. The second part is the actual
content of your article. The two parts are separated by a blank line.

Here is the list of yaml values you can fill

| Key         | Description                                                                           | Mandatory |
| ----------- | ------------------------------------------------------------------------------------- | --------- |
| title       | The title of your article.                                                            | **Yes**   |
| description | The description of your article (sub-title)                                           | No        |
| slug        | The link you want for your article. If left empty, will be generated from title.      | No        |
| banner      | URL of the banner (image at the top level of your article)                            | No        |
| author      | Author of the article                                                                 | No        |
| date        | The date of writing/publication of your article.                                      | **Yes**   |
| tags        | A list of tags you want to apply on the article (useless right now, but still pretty) | No        |

If any of the two mandatory values (`date` and `title`) are omitted, the parser will complain and simply ignore the file.

## Example Post

`pages/first-article`

```
title: First Article
description: The reasons I made SmallBlog
banner: ""
slug: first-article
author: Depado
date: 2016-05-06 11:22:00
draft: false
tags:
    - inspiration
    - dev

# Actual Markdown Content
Notice the blank line right after the `tags` list.
That's how you tell the parser that you are done with yaml format.
```

This article will be parsed, and available at `example.com/post/first-article`.
It will also be listed at `example.com/`.

## Filesystem Monitoring

The directory you define in your `conf.yml` file is constantly watched by the
server. Which means several things :
 - If you create a new file, it will be parsed and added to your site.
   (Also if you `mv` a file inside the directory)
 - If you modify an exisiting file, it will be parsed and modified.
 - If you delete an existing file, the article will be removed. (Also if you
   `mv` a file out of the directory)

All these changes are instant. Usually a file takes ~250µs to be parsed. When
you restart the server, all the files will be parsed again so they are stored in
RAM (which is really efficient unless you have 250Mo of markdown file).
