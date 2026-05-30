

# ---------------- SQLBOILER
MYSQL_DBNAME="aicha" MYSQL_HOST="127.0.0.1" MYSQL_PORT=3306 MYSQL_USER="root" MYSQL_PASS="root" MYSQL_SSLMODE="false" sqlboiler mysql --output="models" --wipe


# ---------------- REDIS CACHING
# Stop the container first

sudo systemctl stop redis


sudo docker stop redis-server

# Remove the container
sudo docker rm redis-server

# Now create a new one
sudo docker run -d --name redis-server -p 6379:6379 redis:latest

# Directly open Redis CLI
sudo docker exec -it redis-server redis-cli
