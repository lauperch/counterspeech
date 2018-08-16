FROM python:2.7
MAINTAINER Christoph Lauper "christoph.lauper@gmail.com"
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
ENTRYPOINT ["python"]
CMD ["app.py"]
