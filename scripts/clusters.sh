#/bin/sh

source settings.source

CLUSTERS=(`gcloud container clusters list | column -t | sed "1 d" | awk '{printf "%s=%s\n", $1,$4}'`)
SIZE=${#CLUSTERS[@]}
for i in "${!CLUSTERS[@]}"
do
  #temp= "value: ${CLUSTERS[$i]}"
  if [ $i = 0 ]; then
    var=$var${CLUSTERS[$i]}
  else
    var=$var,${CLUSTERS[$i]}    
  fi
done

FEDERATION_IP=`kubectl --namespace=federation-system get svc | column -t | sed "1 d" |awk '{print $3}'`

sed -e "s/\CLUSTERS_PLACEHOLDER/$var/" -e "s/\FEDERATION_IP_PLACEHOLDER/$FEDERATION_IP/" ../manifests/admin_template.yaml > ../manifests/geoserver-admin.yaml

echo "DONE - You can now deploy the admin via: make deploy-admin"

