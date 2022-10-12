# syntax=docker/dockerfile:1

FROM python:3.8-slim-buster

WORKDIR /app

COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

RUN apt-get update
RUN apt-get install ffmpeg libsm6 libxext6  -y

COPY . .

RUN pip3 install --ignore-installed --upgrade tensorflow-2.8.0-cp38-cp38-linux_x86_64.whl 

CMD [ "python", "./client.py"]
