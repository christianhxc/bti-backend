
```
docker build -t 302878951089.dkr.ecr.us-west-1.amazonaws.com/bti-backend:1.0 .
```

```
~/.docker/config.json
```

```
{
    "credHelpers": {
    	"302878951089.dkr.ecr.us-west-1.amazonaws.com": "ecr-login"
    }
}
```

```
sudo mv ~/environment/go/bin/docker-credential-ecr-login /usr/local/go/bin/
aws ecr create-repository --repository-name bti-backend --region us-west-1
```

```
docker push 302878951089.dkr.ecr.us-west-1.amazonaws.com/bti-backend:1.0
```