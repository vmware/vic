/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common;

import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Utilities for host tests
 */
public class HostUtil {

    public static final int HOST_PORT_DEFAULT = 443;

    public enum HostStates {
        CONNECTED, DISCONNECTED, ENTER_MAINTENANCE_MODE, EXIT_MAINTENANCE_MODE, NOT_RESPONDING;
    }

    /**
     * Method to build the host specification
     *
     * @param name - name of the host
     * @param username - name of host user
     * @param password - password of host user
     * @param port - port for connections, usually 443
     * @param parent - datacenter or cluster specification
     * @return The HostSpec for the host
     */
    public static HostSpec buildHostSpec(String name, String username, String password, int port,
        ManagedEntitySpec parent) {
        HostSpec hostSpec = new HostSpec();
        hostSpec.name.set(name);
        hostSpec.userName.set(username);
        hostSpec.password.set(password);
        hostSpec.port.set(port);
        hostSpec.parent.set(parent);
        return hostSpec;
    }
}
