# clean everything
docker compose down --rmi all --remove-orphans
docker image prune -a -f
docker builder prune --all --force



# always clean build
docker compose down
docker compose --env-file dms.env build --no-cache
docker compose --env-file dms.env up --build -d


docker compose --env-file dms.env up


docker compose down
mvn clean package -DskipTests && \
JAR_NAME="app-$(date +%s).jar" && \
mv target/dishmanagementservice-0.0.1-SNAPSHOT.jar target/$JAR_NAME && \
echo "JAR Renamed to: $JAR_NAME"

JAR_NAME=$JAR_NAME docker compose --env-file dms.env build --no-cache
JAR_NAME=$JAR_NAME docker compose --env-file dms.env up --build -d


