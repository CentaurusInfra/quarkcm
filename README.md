# quarkcm

quarkcm is Connection Manager for Quark CDMA Networking Solution.

## push to docker hub
sudo docker login -u quarkcm

sudo docker image build -t quarkcm/quarkcm:v0.1.0 -f Dockerfile .
sudo docker image push quarkcm/quarkcm:v0.1.0


## deploy to cluster
kubectl apply -f deploy-quarkcm.yaml
