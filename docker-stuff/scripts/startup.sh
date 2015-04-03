#!/bin/bash

# configure postgres!
su postgres <<'EOF'
  echo "configuring postgres"
  sed -i 's/local   all             all                                     peer/local   all             all                                     trust/g' /etc/postgresql/9.3/main/pg_hba.conf

  /etc/init.d/postgresql start

  psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker111';" &&\
      createdb -O docker docker

  echo "configuring postgres successful"
EOF

mkdir --parents /home/docker/go/src/github.com/$APP
cd /home/docker/go/src/github.com/$APP
git clone https://github.com/$APP.git /home/docker/go/src/github.com/$APP
git checkout -q $COMMIT

ret=$?
echo $ret

if [ $ret -eq 0 ]; then
  echo "go get ./..."
  go get ./...
  ret=$?
  echo $ret
fi

if [ $ret -eq 0 ]; then
  echo "go build ./..."
  go build ./...
  ret=$?
  echo $ret
fi

if [ $ret -eq 0 ]; then
  echo "go test -v ./..."
  go test -v ./...
  ret=$?
  echo $ret
fi

if [ $ret -eq 0 ]; then
  echo "minideploy"
  minideploy
  ret=$?
  echo $ret
fi
