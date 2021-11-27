FROM python:3.6

WORKDIR /code

COPY ./requirements.txt .
RUN pip install -r requirements.txt

COPY . .

ENTRYPOINT ["python", "db_watcher.py"]
