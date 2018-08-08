# Configuration of bridge network widths

## Issues addressed

This design addresses the following issues:

* vmware/vic#3943
* vmware/vic#3771

The following pull request should be resolved prior to completion of this work:

* vmware/vic#7183

## Requirements

The following requirements should be addressed by this design

1. Configuration of a default width for bridge networks:
    * In `vic-machine create`
    * In `vic-machine configure`
1. Inspect current default configuration via `vic-machine inspect`
1. Support the `--subnet` option to `docker network create` to allow override of default configuration
1. Apply a sane minimum width that must allow:
    * gateway - this is the endpointVM
    * broadcast address
    * containerVM - a minimum of one address for a containerVM
1. It should be possible to specify a starting offset in both the `bridge-ip-range` vic-machine option _and_ in the `--subnet` option
    * This requirement is present to simplify surrounding network configuration - it may be desirable from an infrastructure/viadmin perspective to be able to use some of a CIDR subnet for non-VCH managed addresses, and the rest for the bridge network range but without reduce the bridge width by half.

## Reconfigure (and upgrade) behaviour

The expected behaviour and limitations on reconfigure and/or upgrade:

1. Existing networks configured with default values are unchanged.
    * this includes the default bridge network that is automatically created for every VCH.
1. Existing networks configured with a specific subnet are unchanged.
1. Validation of the specified default width vs the configured bridge-ip-range.
    * the bridge range must permit at least one network of the default width.

## Test considerations

1. Default bridge network has configured default width at create time (width-X).
1. Default bridge network has original width-X after default is reconfigured to width-Y.
1. Bridge network created without subnet option has current default width (after create and after reconfigure).
1. Bridge network created with the subnet option has the specified width.
1. Containers can be powered when connected to a given network until the addresses are exhausted. Then a sane error message must be returned noting that there are no more addresses available within the configured subnet.
    * Default bridge network
    * Created bridge network without `--subnet` specified
    * Created bridge network with `--subnet` specified
    * Connected to two networks, one of which still has addresses available and one which does not.
1. `docker network inspect` shows the expected IP range and gateway for all network types.
1. Name resolution and routing over the endpointVM gateway address should be tested for all configurations.
1. vic-machine should fail if default configured range is below minimum permitted width.
1. vic-machine should fail if default width is greater than `bridge-ip-range`.

## Documentation notes

Should document:

1. reasons for the minimum width
1. when the default is applied
1. when the default is _not_ applied
1. what cannot be reconfigured
1. what occurs when you run out of addresses in a range

## Concerns

The following are items of concern that are not addressed elsewhere:

1. if using `--subnet` option a lot we potentially fragment the available CIDR address space