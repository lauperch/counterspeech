FROM golang

# if left blank app will run with dev settings
# to build pprod image run:
# $ docker build ./api --build-args app_env=prod
ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/github.com/user/dthought/crawl
WORKDIR /go/src/github.com/user/dthought/crawl

RUN go get ./
RUN go build

# if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
# if prod setting will build binary
CMD if [ ${APP_ENV} = prod ]; \
	then \
	crawl; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi
	
EXPOSE 3030
