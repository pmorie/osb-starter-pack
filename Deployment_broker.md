# Documentation for the Deployment of Dataverse-broker on Openshift


##Generate Docker Image

build:
```actionscript
go build -i github.com/SamiSousa/dataverse-broker/cmd/dataverse-broker
```

test: `## Runs the tests`
```actionscript
go test -v $(shell go list ./... | grep -v /vendor/ | grep -v /test/)
```

linux: `## Builds a Linux executable`
```actionscript
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
go build -o dataverse-broker-linux --ldflags="-s" github.com/SamiSousa/dataverse-broker/cmd/dataverse-broker
```

image: `linux ## Builds a Linux based image`
```actionscript
cp dataverse-broker-linux image/dataverse-broker
$(SUDO_CMD) docker build image/ -t "$(IMAGE):$(TAG)"
```


##Push image to Quay.io


First you should register on Quay.io and create a repositoty for your images.
Next, login via the docker login command in terminal.
`$ docker login quay.io`

Push the image to Quay.io
`$ docker push quay.io/brutto/dataverse-broker`

All of your pushed images are here:![]
![](https://github.com/bruttozz/airmules/blob/master/3.png)

![](https://img.shields.io/github/stars/pandao/editor.md.svg)

Your teammates can pull the image they need anytime by this command:
`$ docker pull quay.io/brutto/dataverse-broker`


##Deployment on Openshift Container Platform

Go to https://openshift.massopen.cloud/
Login with your username and password or use other authentication.

Deployment method One:
In the web console, Click the blue button to create a new project.

![](https://github.com/bruttozz/airmules/blob/master/4.png)

Name it Big Data Containers.
Deploy the dataverse-broker through the Docker image.
Your can finish this procedure by following this tutorial:
https://learn.openshift.com/introduction/deploying-images/

Method Two:
Copy your login command by clicking “Copy Login Command” which is under the username.

![](https://github.com/bruttozz/airmules/blob/master/2.png)

Go to terminal and paste the command:
`$ oc login https://openshift.massopen.cloud --token=(your own token)`

Run the command like that and you can talk to the Openshift Container Platform directly!

*$ oc login https://openshift.massopen.cloud --token=ZL0cHbHDno7nJS0Bk-NgiqACbA82sxtgX6sf77ibpio
Logged into "https://openshift.massopen.cloud:443" as "brutto" using the token provided.

You have one project on this server: "bdc"

Using project "bdc".*

You can see what project you have and which project you are in.
Then use oc commands in terminal, which will go into effect automatically on the Openshift cluster.

After you finish your work, just run `$ oc logout` to disconnect your terminal with Openshift.


##Debugging

If there is anything wrong in your code, it may cause the deployed broker failed.

If that happens, you should go to the Overview page, click your application’s name. For our project, it’s dataverse-broker.

![](https://github.com/bruttozz/airmules/blob/master/1.png)

First, you can click “View Log” to check logs. If that doesn’t work, you should use the “Edit YAML” under the Actions, to find out if any arguments or commands missed in the YAML file.
