/* Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

/**
 * This enumeration represents the different Failover policies for HA Admission
 * Control
 */
public enum AdmissionControlFailoverPolicy {
   DISABLED, SLOT_POLICY, CLUSTER_RESOURCE_PERCENTAGE, DEDICATED_FAILOVER_HOSTS
}