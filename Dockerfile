FROM public.ecr.aws/lambda/provided:al2 as build

RUN yum install -y golang
RUN go env -w GOPROXY=direct

ADD src/go.mod src/go.sum ./
RUN go mod download

ADD src/. .
RUN go build -o /main

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
ENTRYPOINT [ "/main" ]