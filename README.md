# Smallblog
A simple self-hosted markdown flat files blog

## What is Smallblog
The main goal of this project is to show how easily you can develop a flat file
blog with markdown as the primary writing language. It's not perfect, it will
never be, some people are already doing great things based on that idea, like
[Hugo](https://gohugo.io/) for instance. Let's note though, that's **not** a
static website generator.

## Disclaimer
Smallblog is a quick project. It's efficient while there aren't too much files,
but please keep in mind that because of its design, the more articles you'll
have the more RAM the server will take. This will be fixed someday. Maybe.

As for now, Smallblog is designed to be fast. And that's all. (Which means no
markdown parsing for each request, no file opening, everything in RAM)

## Configure
Put a `conf.yml` file next to your `smallblog` binary. Here are the options you
can customize

| Key         | Description                                                               | Default     |
| ----------- | ------------------------------------------------------------------------- | ----------- |
| host        | Interface on which the server should listen.                              | "127.0.0.1" |
| port        | Port on which the server should listen.                                   | 8080        |
| debug       | Activates the router's debug mode.                                        | false       |
| pages_dir   | Local or absolute path to the directory in which your articles are stored | "pages"     |
| title       | Blog title (front page)                                                   | ""          |
| description | Blog Description (front page)                                             | ""          |

## Write Posts
There is no naming convention for file names. You can name them whatever you
want, it won't chage the server's behaviour. A post (or page/article) file is
divided in two parts. The first part is yaml data. The second part is the actual
content of your article. The two parts are separated by a blank line.

Here is the list of yaml values you can fill

| Key         | Description                                                                           | Mandatory |
| ----------- | ------------------------------------------------------------------------------------- | --------- |
| title       | The title of your article.                                                            | **Yes**   |
| description | The description of your article (sub-title)                                           | No        |
| slug        | The link you want for your article. If left empty, will be generated from title.      | No        |
| author      | Author of the article                                                                 | No        |
| date        | The date of writing/publication of your article.                                      | **Yes**   |
| tags        | A list of tags you want to apply on the article (useless right now, but still pretty) | No        |

If any of the two mandatory values (`date` and `title`) are omitted, the parser will complain and simply ignore the file.

## Example Post
`pages/first-article`
```
title: First Article
description: The reasons I made SmallBlog
slug: first-article
author: Depado
date: 2016-05-06 11:22:00
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
 - If you modify an exisiting file, it will be parsed and modified on your site
   if necessary (e.g if the slug changes).
 - If you delete an existing file, the article will be removed. (Also if you
   `mv` a file out of the directory)

All these changes are instant. Usually a file takes ~250µs to be parsed. When
you restart the server, all the files will be parsed again so they are stored in
RAM (which is really efficient unless you have 250Mo of markdown file).
