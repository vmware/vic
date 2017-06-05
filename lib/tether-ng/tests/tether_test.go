package tetherng_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	tether "github.com/vmware/vic/lib/tether-ng"
	"github.com/vmware/vic/lib/tether-ng/interaction"
	mocker "github.com/vmware/vic/lib/tether-ng/mocks"
	"github.com/vmware/vic/lib/tether-ng/network"
	"github.com/vmware/vic/lib/tether-ng/toolbox"
	"github.com/vmware/vic/lib/tether-ng/types"
)

var (
	config = &types.ExecutorConfig{
		Sessions: map[string]*types.SessionConfig{
			"attach": &types.SessionConfig{
				Env:        []string{},
				Cmd:        []string{"/bin/ls", "-l"},
				WorkingDir: "/",

				Attach:    true,
				OpenStdin: true,
				RunBlock:  true,
				Tty:       true,
				Restart:   false,
			},
		},
	}
)

func TestRegister(t *testing.T) {
	var err error

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tether := tether.NewTether(ctx)

	callctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interaction := mocker.NewMockPlugin(ctrl)
	network := mocker.NewMockPlugin(ctrl)
	toolbox := mocker.NewMockPlugin(ctrl)

	interaction.EXPECT().UUID(callctx).Return(uuid.New())
	err = tether.Register(callctx, interaction)
	assert.Nil(t, err)

	network.EXPECT().UUID(callctx).Return(uuid.New())
	err = tether.Register(callctx, network)
	assert.Nil(t, err)

	toolbox.EXPECT().UUID(callctx).Return(uuid.New())
	err = tether.Register(callctx, toolbox)
	assert.Nil(t, err)
}

func TestRegisterMock(t *testing.T) {
	var err error

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tether := mocker.NewMockPluginRegistrar(ctrl)

	callctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interaction := interaction.NewInteraction(callctx)
	network := network.NewNetwork(callctx)
	toolbox := toolbox.NewToolbox(callctx)

	tether.EXPECT().Register(callctx, interaction)
	err = tether.Register(callctx, interaction)
	assert.Nil(t, err)

	tether.EXPECT().Register(callctx, network)
	err = tether.Register(callctx, network)
	assert.Nil(t, err)

	tether.EXPECT().Register(callctx, toolbox)
	err = tether.Register(callctx, toolbox)
	assert.Nil(t, err)
}

func TestRegisterFailure(t *testing.T) {
	var err error

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tether := tether.NewTether(ctx)

	callctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interaction := mocker.NewMockPlugin(ctrl)

	uuid := uuid.New()

	interaction.EXPECT().UUID(callctx).Return(uuid)
	err = tether.Register(callctx, interaction)
	assert.Nil(t, err)

	interaction.EXPECT().UUID(callctx).Return(uuid)
	err = tether.Register(callctx, interaction)
	assert.NotNil(t, err)
}

func TestConfigure(t *testing.T) {
	var err error

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	tether := tether.NewTether(ctx)

	interaction := mocker.NewMockPlugin(ctrl)
	network := mocker.NewMockPlugin(ctrl)
	toolbox := mocker.NewMockPlugin(ctrl)

	interaction.EXPECT().UUID(ctx).Return(uuid.New())
	err = tether.Register(ctx, interaction)
	assert.Nil(t, err)

	network.EXPECT().UUID(ctx).Return(uuid.New())
	err = tether.Register(ctx, network)
	assert.Nil(t, err)

	toolbox.EXPECT().UUID(ctx).Return(uuid.New())
	err = tether.Register(ctx, toolbox)
	assert.Nil(t, err)

	interaction.EXPECT().Configure(ctx, config).Return(nil)
	network.EXPECT().Configure(ctx, config).Return(nil)
	toolbox.EXPECT().Configure(ctx, config).Return(nil)
	for _, i := range tether.Plugins(ctx) {
		i.Configure(ctx, config)
	}
}
