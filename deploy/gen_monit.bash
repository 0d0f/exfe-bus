U_ID="$1"
G_ID="$2"

if [ "$U_ID" == "" ]
then
  echo Usage: $0 [UID] [GID]
  exit -1
fi

if [ "$G_ID" == "" ]
then
  echo Usage: $0 [UID] [GID]
  exit -1
fi

for f in ../bin/*
do
  NAME=${f##.*/}
  cat monit.templ | sed "s/{{bin_name}}/$NAME/g" | sed "s/{{uid}}/${U_ID}/g" | sed "s/{{gid}}/${G_ID}/g"
done
