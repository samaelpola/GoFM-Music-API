# GoFM Music API

Creating a rest api in golang with mux, sqlite to stock music data and s3 to store image and file audio of music.

To run this application, you must have docker installed on your machine
#

## Getting Started

### Launch app

```
make
```

### after launch app run
## run the command below to generate the swagger doc
```
make swag
```

Open http://localhost:8083/swagger/index.html with your browser to have access to the swagger doc.

link useful:
minio -> http://localhost:9001
phpMyAdmin -> http://localhost:8081

#

## Command useful

### Build project

```
make build
```

### Up project

```
make up
```

### Down container

```
make down
```

### Doc swagger

#### generate doc swagger

```
make swag
```
