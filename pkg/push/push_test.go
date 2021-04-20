package push

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	errors "github.com/apenella/go-common-utils/error"
	mockclient "github.com/apenella/go-docker-builder/test/mock"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	var w bytes.Buffer
	writer := io.Writer(&w)
	reader := ioutil.NopCloser(io.Reader(&bytes.Buffer{}))

	tests := []struct {
		desc              string
		dockerPushCmd     *DockerPushCmd
		ctx               context.Context
		pushOptions       dockertypes.ImagePushOptions
		mock              *mockclient.DockerClient
		prepareAssertFunc func(context.Context, *mockclient.DockerClient, *DockerPushCmd)
		assertFunc        func(*mockclient.DockerClient) bool
		err               error
	}{
		{
			desc:          "Testing error when DockerPushCmd is undefined",
			ctx:           context.TODO(),
			pushOptions:   dockertypes.ImagePushOptions{},
			dockerPushCmd: nil,
			err:           errors.New("(push::Run)", "DockerPushCmd is undefined"),
		},
		{
			desc:        "Testing error when ImagePushOptions is undefined",
			ctx:         context.TODO(),
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				ImagePushOptions: nil,
			},
			err: errors.New("(push::Run)", "Image push options is undefined"),
		},
		{
			desc: "Testing push a single image",
			ctx:  context.TODO(),
			dockerPushCmd: &DockerPushCmd{
				Writer:           io.Writer(writer),
				ImageName:        "test_image",
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				ExecPrefix:       "",
			},
			mock: new(mockclient.DockerClient),
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 1)
			},
			err: &errors.Error{},
		},
		{
			desc: "Testing push a single image with auth",
			ctx:  context.TODO(),
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				Writer:    io.Writer(writer),
				ImageName: "test_image",
				ImagePushOptions: &dockertypes.ImagePushOptions{
					RegistryAuth: "auth",
				},
				ExecPrefix: "",
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 1)
			},
			mock: new(mockclient.DockerClient),
			err:  &errors.Error{},
		},
		{
			desc: "Testing push a single image with tags",
			ctx:  context.TODO(),
			prepareAssertFunc: func(ctx context.Context, mock *mockclient.DockerClient, cmd *DockerPushCmd) {
				mock.On("ImagePush", ctx, cmd.ImageName, *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag1", *cmd.ImagePushOptions).Return(reader, nil)
				mock.On("ImagePush", ctx, "tag2", *cmd.ImagePushOptions).Return(reader, nil)
				cmd.Cli = mock
			},
			pushOptions: dockertypes.ImagePushOptions{},
			dockerPushCmd: &DockerPushCmd{
				Writer:           io.Writer(writer),
				ImageName:        "test_image",
				Tags:             []string{"tag1", "tag2"},
				ImagePushOptions: &dockertypes.ImagePushOptions{},
				ExecPrefix:       "",
			},
			assertFunc: func(mock *mockclient.DockerClient) bool {
				return mock.AssertNumberOfCalls(t, "ImagePush", 3)
			},
			mock: new(mockclient.DockerClient),
			err:  &errors.Error{},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Log(test.desc)

			if test.prepareAssertFunc != nil {
				test.prepareAssertFunc(test.ctx, test.mock, test.dockerPushCmd)
			}

			err := test.dockerPushCmd.Run(test.ctx)

			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				assert.True(t, test.assertFunc(test.mock))
			}
		})
	}
}

// func TestAddAuth(t *testing.T) {

// 	type args struct {
// 		username string
// 		password string
// 	}
// 	tests := []struct {
// 		name    string
// 		options *DockerPushOptions
// 		args    *args
// 		err     error
// 		res     string
// 	}{
// 		{
// 			name: "Test add user-password auth",
// 			options: &DockerPushOptions{
// 				ImageName:    "test image",
// 				RegistryAuth: nil,
// 			},
// 			args: &args{
// 				username: "user",
// 				password: "AqSwd3Fr",
// 			},
// 			err: nil,
// 			res: "eyJ1c2VybmFtZSI6InVzZXIiLCJwYXNzd29yZCI6IkFxU3dkM0ZyIn0=",
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			err := test.options.AddAuth(test.args.username, test.args.password)
// 			if err != nil {
// 				assert.Equal(t, test.err, err)
// 			} else {
// 				assert.Equal(t, test.res, *test.options.RegistryAuth, "Unexpected auth result")
// 			}
// 		})
// 	}
// }
