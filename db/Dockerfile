FROM postgres:15.1

# TODO: pull this out to env var / remote secret
ENV POSTGRES_PASSWORD TfwePfOzum
ENV POSTGRES_DB chat_app

COPY schema.sql /docker-entrypoint-initdb.d/