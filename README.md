Tool for deleting unused namespaces in k8s

How to deploy:

Push new version
```
export DOCKER_IMAGE=ns-deleter-5
docker build -t $DOCKER_IMAGE . &&  docker push $DOCKER_IMAGE
```

Update deployment to the new version
```
kubectl -n monitoring set image deployment/ns-deleter ns-deleter=$DOCKER_IMAGE
```


How to add users:

Edit conf.yml following the existing format
