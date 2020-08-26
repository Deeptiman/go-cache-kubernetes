TMPFILE=$(mktemp)
/usr/bin/openssl rand -base64 741 > $TMPFILE
kubectl create secret generic mongo-share-data --from-file=internal-auth-mongodb-keyfile=$TMPFILE
rm $TMPFILE
