sudo apt-get install gnupg curl
VERSION="7.0"
curl -fsSL https://pgp.mongodb.com/server-$VERSION.asc | \
sudo gpg -o /usr/share/keyrings//mongodb-server-$VERSION.gpg \
--dearmor

echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-$VERSION.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/$VERSION multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-$VERSION.list

sudo apt-get update

sudo apt-get install -y mongodb-org

