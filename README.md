# Spotify Integration

This project is an implementation of the spotify API 
In includes a frontend done in React + redux with a backend in Goland with microservices
The microservices are served using nginx as a gateway through docker


## Requisites

You will need docker / docker-compose installed

## Start the project

To start the whole project you need to run:
```
docker-compose up
```

Once started you should have something like this: 
```
Creating spotify-app_client_1 ...
Creating spotify-app_user_1   ...
Creating spotify-app_player_1 ...
Creating spotify-app_nginx_1  ...
Creating spotify-app_playlist_1 ...
```

This will start multiple containers
It will build each microservices and start the unit tests before starting them.
Nginx will be the gateway to redirect requests to the correct microservices.
Nginx will serve them through the host from the port 8080.
The client will be built and served on the port 3000.

Once the docker-compose is up and running you'll able to access the client from `127.0.0.1:3000`

## Todo:

- Update the method to retrieve a token to be able to have a refresh token and refresh it.

- Adding metrics microservices with grpc to keep track of country / mail / product / birthday etc. (outside of the public network) while requesting user info.

- Adding unit tests on client side

- Many other improvements
