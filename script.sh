docker compose -f ./docker-composer.yml kill;
docker compose -f ./docker-composer.yml rm -f;
docker compose -f ./docker-composer.yml build;
docker compose -f ./docker-composer.yml up -d;
