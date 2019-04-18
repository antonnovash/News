FROM golang

COPY ./bin /go/news-reader/

CMD news-reader/news

EXPOSE 3006