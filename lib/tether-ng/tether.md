# Tether

This document describes tether re-organization (this is not a complete re-design). The reasons of re-organizing the code base could be summarized as below;

- Making concurrent parts of the code little bit easy to follow by sharing the data via communicating,
- Making tether code easy to test,
- Aggregating guestinfo operations,
- Minimizing the code duplication between tether and vic-init,
- Making the code base more approachable,

With this re-organization, tether itself will be as small as and as portable as possible. Its main responsibility will be passing the configuration data to its plugins, 
starting them (in specific order) and then reporting back their status and diagnostics data.

The main use case of tether is being the parent of all sessions it creates and acting like a init system (aka pid 1) 
but it could potentially run under any other init system with minimal effort (e.g. disabling some plugins via flags)

As a rule of thumb every exported method now receives context.Context as its first parameter - even if the receiver doesn't need it.

There should be two different context scopes - Global and call specific. Global context gets created by the main function passed around when we create
the tether and its plugins. Its intended use is for tether wide cancellation.

The call specific context should derive from this global context and should be passed to the individual calls for their lifecycle management.

```go
	ctx := context.Background()

	tether := tether.NewTether(ctx)

	callctx, cancel := context.WithCancel(ctx)
	defer cancel()
    ...
```

## Directory structure

Proposed directory structure follows with each plugin having its own package.

```
cmd
├── tether-ng
│   ├── main.go
│
lib
├── tether-ng
│   ├── doc.go
│   ├── interaction
│   │   └── plugin.go
│   ├── mocks
│   │   └── tether.go
│   ├── network
│   │   └── plugin.go
│   ├── tests
│   │   └── tether_test.go
│   ├── tether.go
│   ├── tether.md
│   ├── toolbox
│   │   └── plugin.go
│   └── types
│       └── types.go
```

## Tether

Tether implements following interface

```go
type Tetherer interface {
	Register(ctx context.Context, plugin Plugin) error
	Unregister(ctx context.Context, plugin Plugin) error
	Plugins(ctx context.Context) []Plugin

    // FIXME(caglar10ur)
    TODO
}
```

which can be used like following;

```go
type Tether struct {
	Tetherer

	ctx context.Context

	m       sync.RWMutex
	plugins map[uuid.UUID]Plugin
}

func NewTether(ctx context.Context) Tetherer {
	return &Tether{
		ctx:     ctx,
		plugins: make(map[uuid.UUID]Plugin),
	}
}

...
```

plugins can provide additional capabilities via implementing following interfaces

```go
// Interaction calls this to release Waiter
type Releaser interface {
	Release(ctx context.Context, out chan<- chan struct{})
}

// Process calls this to wait Releaser to release
type Waiter interface {
	Wait(ctx context.Context, in <-chan chan struct{})
}

// Process calls this to mutate writers/readers
type Interactor interface {
	Interact(ctx context.Context, in <-chan *types.Session) <-chan struct{}

	Close(ctx context.Context, in <-chan *types.Session) <-chan struct{}
}

//
type Reaper interface {
	Reap(ctx context.Context) error
}
```

## Configuration

Tether and its plugins uses following structures as their configuration. Tether uses extraconfig package to serialize/deserialize this and uses vmx as the persistent datastore. 

> This mechanism can be replaced by NamespaceDB

```go
type ExecutorConfig struct {
    // Sessions is the set of sessions currently hosted by this executor
    // These are keyed by session ID
    Sessions map[string]*SessionConfig `vic:"0.1" scope:"read-only" key:"sessions"`
}


type SessionConfig struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	// List of environment variable to set in the container
	Env []string
	// Command to run when starting the container
	Cmd []string
	// Current directory (PWD) in the command will be launched
	WorkingDir string

	// Allow attach or not
	Attach bool `vic:"0.1" scope:"read-only" key:"attach"`

	// Open stdin or not
	OpenStdin bool `vic:"0.1" scope:"read-only" key:"openstdin"`

	// Delay launching the Cmd until an attach request comes
	RunBlock bool `vic:"0.1" scope:"read-only" key:"runblock"`

	// Allocate a tty or not
	Tty bool `vic:"0.1" scope:"read-only" key:"tty"`

	// Restart controls whether a process gets relaunched if it exists
	Restart bool `vic:"0.1" scope:"read-only" key:"restart"`

	// StopSignal is the signal name or number used to stop a container
	StopSignal string `vic:"0.1" scope:"read-only" key:"stopSignal"`

	// User and group for setuid programs
	User  string `vic:"0.1" scope:"read-only" key:"user"`
	Group string `vic:"0.1" scope:"read-only" key:"group"`
}

type Session struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	*SessionConfig

	Interaction

	ID string

	// The primary process for the session
	Cmd exec.Cmd `vic:"0.1" scope:"read-only" key:"cmd" recurse:"depth=2,nofollow"`

	// The exit status of the process, if any
	ExitStatus int `vic:"0.1" scope:"read-write" key:"status"`

	// This indicates the launch status
	Started string `vic:"0.1" scope:"read-write" key:"started"`

	// RessurectionCount is a log of how many times the entity has been restarted due
	// to error exit
	ResurrectionCount int `vic:"0.1" scope:"read-write" key:"resurrections"`
}

type Interaction struct {
	// Exclusive access to the structure
	sync.Mutex `vic:"0.1" scope:"read-only" recurse:"depth=0"`

	Done <-chan struct{}

	Pty       *os.File               `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Outwriter dio.DynamicMultiWriter `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Errwriter dio.DynamicMultiWriter `vic:"0.1" scope:"read-only" recurse:"depth=0"`
	Reader    dio.DynamicMultiReader `vic:"0.1" scope:"read-only" recurse:"depth=0"`
}
```

## Plugins

Plugins are isolated independent functionalities controlled by the tether. Depending on the the target guest operating system some of them 
could be a NOOP. As a rule of thumb they don't share data with each other - but they could call methods on each other.

Plugins implement following interface

```go
type Plugin interface {
	Configure(ctx context.Context, config *types.ExecutorConfig) error

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	UUID(ctx context.Context) uuid.UUID
    
    // FIXME(caglar10ur)
    TODO
}
```

* Process
    * Launch processes: Needs to order multiple processes to solve dependencies between them.
    * Monitor processes: Needs to monitor the lifecycle of the process and apply the dictated policies (restart etc.)

* Interaction
    * Capture and forward process output for logging/debugging purposes. The logging medium needs to be persistent (accessible when powered off), 
    should provide timestamps and possibly differentiate different OS streams.
    * Inject process input: Needs to control/own the process input/output to enable interaction shells.

* Tools
    * Dispatch signals to processes [calls Signal@Process plugin]
    * Handle toolbox commands (Power/shutdown operations etc.)
    * Report back to infrastructure like toolbox  (IP address publishing etc.)
* Storage
    * Bidirectional file copy (support `docker cp` and similar)
    * Mount file systems:
        * block device mount
        * network share mount
    
* Network
    *Handle network configurations
        * IP address
        * Firewall
* Hardware
    * Handle runtime hardware events (like hotadd of mem/cpu/vnic - support `docker network connect` and similar)
    ``` Can't we just listen reconfigure events and trigger tether via toolbox and bring cpu/mem/nic online```
* Capabilities
    * Handle vmfork/instant clone
        * pre/post handling
        * fork trigger
    * Handle system limits (e.g. ulimits)

## Initialization order

- Toolbox
- Capabilities
- Network
- Storage
- Hardware
- Process
- Interaction (on demand)

## Non-functional requirements
* limited filesystem dependencies
* immutable configuration - accessible when powered off

VMX, NamespaceDB, or something else?

* publishing of runtime data - e.g. process status, DHCP-assigned address

Toolbox?

* idempotent application of configuration
* more portable the better
* no in-guest network dependency

Serial, VMCI or something else?

* containerVM and endpointVM enforce mutual authentication

## Tests

Now that we use interfaces we should consider using gomock to generate the mocks automatically.

```shell

# cat doc.go
package tetherng

//go:generate mockgen github.com/vmware/vic/lib/tether-ng Plugin,Tetherer

# go generate > mock/tether.go && go test -v ./tests/
=== RUN   TestRegister
--- PASS: TestRegister (0.00s)
=== RUN   TestRegisterMock
--- PASS: TestRegisterMock (0.00s)
=== RUN   TestRegisterFailure
--- PASS: TestRegisterFailure (0.00s)
=== RUN   TestConfigure
--- PASS: TestConfigure (0.00s)
PASS
ok      github.com/vmware/vic/lib/tether-ng/tests       0.012s
```