PROJECT = montcoio
NAME = mpubsub
TAG = latest


docker-build:
	docker build --no-cache -t us.gcr.io/$(PROJECT)/$(NAME):$(TAG) -f Dockerfile .


push:
	docker push us.gcr.io/$(PROJECT)/$(NAME):$(TAG) 

pull:
	docker pull us.gcr.io/$(PROJECT)/$(NAME):$(TAG) 


run:
	docker run -p 3000:3000 --rm -it -d --name $(NAME) us.gcr.io/$(PROJECT)/$(NAME):$(TAG) 



runnod:
	docker run -p 3000:3000 --rm -it --name $(NAME) us.gcr.io/$(PROJECT)/$(NAME):$(TAG) 

stop:
	docker stop $(NAME)

logs:
	docker logs $(NAME)



