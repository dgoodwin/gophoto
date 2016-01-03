# Properly rebuilds the environment, brings up the db, populates schema, and then
# starts the gotphoto container.
set -e
sudo docker-compose stop
sudo docker-compose rm -f
sudo docker-compose build
sudo docker-compose up -d db
sleep 5
sudo docker-compose run --rm gophoto goose up
sudo docker-compose up -d gophoto

