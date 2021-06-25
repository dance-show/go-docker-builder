go-docker-builder
=======

`go-docker-builder` library it is wrapper over docker client SDK that provides a set of packages that help to manage the most common docker use cases such as build or push images.

It also manages docker registry autentication, prepares docker build context to be used by docker client SDK and supports docker build context either from local path or git repository.


<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Use cases](#use-cases)
  - [Build](#build)
    - [Context](#context)
      - [Path](#path)
      - [Git](#git)
    - [Context filesystem](#context-filesystem)
  - [Push](#push)
  - [Copy](#copy)
- [Authentication](#authentication)
  - [Docker registry](#docker-registry)
  - [Git server](#git-server)
    - [Basic auth](#basic-auth)
    - [SSH key](#ssh-key)
    - [SSH agent](#ssh-agent)
- [Response](#response)
- [Examples](#examples)
- [References](#references)
- [License](#license)

<!-- /code_chunk_output -->

## Use cases
`go-docker-builder` library has been written to provide an easy way to interactuate with docker client SDK on the use cases listed below:

- **Build**: build a docker images
- **Push**: push a docker image to a docker registry
- **Copy**: copy a docker image from one Docker registry to another one

### Build
Package `build` purpose is to build docker images.
To perform a build action, must be created a `DockerBuildCmd` instance.

```go
// DockerBuilderCmd
type DockerBuildCmd struct {
	// Cli is the docker api client
	Cli types.DockerClienter
	// ImageName is the name of the image
	ImageName string
	// ImageBuildOptions from docker sdk
	ImageBuildOptions *dockertypes.ImageBuildOptions
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockertypes.ImagePushOptions
	// PushAfterBuild when is true images are automatically pushed to registry after build
	PushAfterBuild bool
	// Response manages responses from docker client
	Response types.Responser
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoveAfterPush when is true images are removed from local after push
	RemoveAfterPush bool
}
```

Below there is a recipe to build docker images using `go-docker-builder`:

1. Create a docker client from docker client SDK
```go
dockerCli, err = client.NewClientWithOpts(client.FromEnv)
if err != nil {
	return err
}
```

2. Give a name to the image
```go
registry := "registry.go-docker-builder.test"
imageName := strings.Join([]string{registry, "alpine"}, "/")
```

3. In case you need a custom response, create a response object
```go
res := response.NewDefaultResponse(
	response.WithTransformers(
		transformer.Prepend("my-custom-response"),
	),
)
```

4. Create `DockerBuildCmd` instance
```go
dockerBuilder := &build.DockerBuildCmd{
	Cli:       dockerCli,
	ImageName: imageName,
	Response:  res,
}
```

5. Create a docker build context
```go
imageDefinitionPath := filepath.Join(".", "files")
dockerBuildContext := &contextpath.PathBuildContext{
	Path: imageDefinitionPath,
}
```

6. Add the docker build Context to `DockerBuildCmd`
```go
err = dockerBuilder.AddBuildContext(dockerBuildContext)
if err != nil {
	return err
}
```

7. Include extra docker image tags, in case are needed
```go
dockerBuilder.AddTags(strings.Join([]string{imageName, "tag1"}, ":"))
```

8. Include authorization either for pull, push or both, when is required.
```go
err = dockerBuilder.AddAuth(username, password, registry)
if err != nil {
	return err
}
```

9. Include build arguments, in case are needed
```go
err = dockerBuilder.AddBuildArgs("key", "value")
if err != nil {
	return err
}
```

10. Start the build
```go
err = dockerBuilder.Run(context.TODO())
if err != nil {
	return err
}
```

#### Context
Docker build context is a set of files required to build a docker image. `go-docker-builder` library supports two kind of sources to create the Docker build context: `path` and `git`

##### Path
When files are located on local host, Docker build context must be created as `path`, and it is only required to define the local folder where files are located.

A `PathBuildContext` instance represents a Docker build context as path.
```go
// PathBuildContext creates a build context from path
type PathBuildContext struct {
	// Path is context location on the local system
	Path string
}
```

##### Git
When files are located on a git repository, Docker build context must be created as `git`.

A `GitBuildContext` instance represents a Docker build context as git.
```go
// GitBuildContext defines a build context from a git repository
type GitBuildContext struct {
	// Path must be set when docker build context is located in a subpath inside the repository
	Path string
	// Repository which will be used as docker build context
	Repository string
	// Reference is the name of the branch to clone. By default is used 'master'
	Reference string
	// Auth
	Auth auth.GitAuther
}
```

To define a `git` Docker build context, the only required attribute is `Repository`, although it accepts other configurations such as the `Reference` name (branch, commit or tag), `Auth` which is used to be authenticated over git server or `Path` which let you to define as Docker build context base a subfolder inside the repository.

`go-docker-builder` uses [go-git](https://github.com/go-git/go-git) library to manage `git` Docker build context.

#### Context filesystem
Context filesystem has been created as an intermediate filesystem between the source files and Docker build context.

Context filesystem, is build on top of [afero](https://github.com/spf13/afero) filesystem. It supports to **tar** the entire filesystem and also to **join** multiple filesystems.

### Push
Package `push` purpose is to push Docker images to a Docker registry.
To perform a push action, must be created a `DockerPushBuildCmd` instance.

```go
// DockerPushCmd is used to push images to docker registry
type DockerPushCmd struct {
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockertypes.ImagePushOptions
	// ImageName is the name of the image
	ImageName string
	// Tags is a list of the images to push
	Tags []string
	// Response manages the docker client output
	Response types.Responser
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoveAfterPush when is true the image from local is removed after push
	RemoveAfterPush bool
}
```

Below there is a recipe to build docker images using `go-docker-builder`:

1. Create a docker client from docker client SDK
```go
dockerCli, err = client.NewClientWithOpts(client.FromEnv)
if err != nil {
	return err
}
```

2. Give a name to the image
```go
registry := "registry.go-docker-builder.test"
imageName := strings.Join([]string{registry, "alpine"}, "/")
```

3. Create `DockerPushCmd` instance
```go
dockerPush := &push.DockerPushCmd{
	Cli:       dockerCli,
	ImageName: imageName,
}
```

4. Add authorization, when is required
```go
user := "myregistryuser"
pass := "myregistrypass"
dockerPush.AddAuth(user, pass)
```

5. Push the image to Docker registry
```go
err = dockerPush.Run(context.TODO())
if err != nil {
	return err
}
```


### Copy
Package `copy` can be understand such a `push` use case variation. Its purpose is to push images either from local host or from a Docker registry to another Docker registry. It can be also used to copy images from one Docker registry namespace to another namespace. 
To perform a copy action, must be created a `DockerImageCopyCmd` instance.

```go
// DockerCopyImageCmd is used to copy images to docker registry. Copy image is understood as tag an existing image and push it to a docker registry
type DockerImageCopyCmd struct {
	// Cli is the docker client to use
	Cli types.DockerClienter
	// ImagePushOptions from docker sdk
	ImagePullOptions *dockertypes.ImagePullOptions
	// ImagePushOptions from docker sdk
	ImagePushOptions *dockertypes.ImagePushOptions
	// SourceImage is the name of the image to be copied
	SourceImage string
	// TargetImage is the name of the copied image
	TargetImage string
	// Tags is a copied images tags list
	Tags []string
	// UseNormalizedNamed when is true tags are transformed to a fully qualified reference
	UseNormalizedNamed bool
	// RemoteSource when is true the source image is pulled from registry before push it to its destination
	RemoteSource bool
	// RemoveAfterPush when is true the image from local is removed after push
	RemoveAfterPush bool
	// Response manages the docker client output
	Response types.Responser
}
```

Some cases where copy can be used are:
- Copy one image from [dockerhub](https://hub.docker.com/) to a private Docker registry
- When Docker images used on your staging environments needs to be promoted somewhere else before use them on production environment.

## Authentication
`go-docker-builder` library contains a bunch of packages to manage authentication either to docker registry or git server.

### Docker registry
Package `github.com/apenella/go-docker-builder/pkg/auth/docker` provides a set of functions to create the authentication items required by Docker registry. Docker registy may require authentication either for `push` or `pull` operations.

Packages [build](#build), [push](#push) or [copy](#copy) already uses that package on its authentication methods, and is not necessary to use it directly.

### Git server
Git server authentication is needed when Docker build context is located on a git repository and the git server requires it authorize you for cloning any repository.

`go-docker-builder` supports several git authentication methods such as `basic auth`, `ssh-key` or `ssh-agent`.

All this authentication methods generates an `AuthMethod` (github.com/go-git/go-git/v5/plumbing/transport).

#### Basic auth
`Basic auth` requires a `username` and `password`. It have to be used when comunication is done through `http/s`.

```go
type BasicAuth struct {
	Username string
	Password string
}
```

#### SSH key
You can use an `ssh key` when comunication to git server is done over ssh. It requires your privete key location (`PkFile`). In case your key is being protected by password, you have to set it on `PkPassword` attribute. Finally, when git user is not `git`, you can define it on `GitSSHUser` attribute.
```go
type KeyAuth struct {
	GitSSHUser string
	PkFile     string
	PkPassword string
}
```

#### SSH agent
To authenticate to git server, `ssh-agent` method uses the ssh agent running on you host. When git user is not `git`, you can define it on `GitSSHUser` attribute.
```go
type SSHAgentAuth struct {
	GitSSHUser string
}
```

## Response
To control how Docker output is sent to the user, `go-docker-build` provides `response` package.
By default, [DockerBuildCmd](#build), [DockerPushCmd](#push) and [DockerCopyCmd](#copy) instances use a basic configuration of `response` that prints Docker output to `os.Stdout`. But you could customize docker client SDK response, defining your own `response` instance, and pass it to them.

`response` receives `ImageBuildResponse` or `ImagePushResponse` items, unmarshals into an struct and finally prepares a user-frendly output.
To customize the Docker output coming from response items, it could receive a list of transfromer functions. That function must complain the type `TransformerFunc` defined on `github.com/apenella/go-common-utils/transformer/string`.

```go
// TransformerFunc is used to enrich or update messages before to be printed out
type TransformerFunc func(string) string
```

A custom response definition could be (the snipped belongs to [build-and-push](https://github.com/apenella/go-docker-builder/tree/master/examples/build-and-push) example):
```go
res := response.NewDefaultResponse(
	response.WithTransformers(
		transformer.Prepend("buildAndPush"),
	),
	response.WithWriter(w),
)

dockerBuilder := &build.DockerBuildCmd{
	Cli:               dockerCli,
	ImageBuildOptions: &dockertypes.ImageBuildOptions{},
	ImagePushOptions:  &dockertypes.ImagePushOptions{},
	PushAfterBuild:    true,
	RemoveAfterPush:   true,
	ImageName:         imageName,
	Response:          res,
}
```

## Examples
On folder [examples](https://github.com/apenella/go-docker-builder/tree/master/examples), you could find some `go-docker-build` examples. Among those examples you could find how to build images using distinct Docker build context locations, Docker registry or git authentications, etc.

To run any example, the repository is provided with some resources that let you to start an ephemeral environment where examples can run. The environments run on `docker-compose` and start a Docker registry, a git server and a client container where example runs. That environments are also used to run functional test.

Each example is also provide by a `Makefile` which helps you to start the examples or tests.

```shell
❯ make help

 Executing example build-and-push

 help                 list allowed targets
 start                start docker registry
 cleanup              cleanup example environment
 generate-certs       generate certificate for go-docker-builder.test
 cleanup-certs        cleanup certificates
 prepare              prepare docker images required to run the example or test
 example              executes the examples
 test                 executes functional test
 logs                 show services logs
 ```

## References
- Docker engine API specifications for building an image: https://docs.docker.com/engine/api/v1.39/#operation/ImageBuild
- Taring files strategy was inspired by: https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07

## License
`go-docker-builder` is available under [MIT](https://github.com/apenella/go-docker-builder/blob/master/LICENSE) license.
