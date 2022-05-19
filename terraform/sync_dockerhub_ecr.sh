#!/bin/sh

n=0
until [ "$n" -ge $RETRIES ]
do
  docker pull $SOURCE_REPO:$IMAGE_TAG && break
  n=$((n+1))
  sleep $SLEEP
done

aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com

docker tag $SOURCE_REPO:$IMAGE_TAG $ECR_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_NAME:$IMAGE_TAG

docker push $ECR_ACCOUNT.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_NAME:$IMAGE_TAG