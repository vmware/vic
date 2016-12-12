/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.model;

/**
 * Specifies the different types of entities in the inventory. Used to build up the Fixtures item
 * managed entities. Each item is prefixed with the Fixtures name of the item to which it belongs.
 */
@Deprecated
public enum FixtureEntities {
   VC,
   NGC_COMMON_CLUSTERED_HOST,
   NGC_COMMON_CLUSTER,
   NGC_COMMON_DATASTORE,
   NGC_COMMON_DATACENTER,
   CONTENT_LIBRARY_LOCAL,
   CONTENT_LIBRARY_PUBLISHED,
   CONTENT_LIBRARY_SUBSCRIBED,
   VDC_VDC,
   VM_DEFAULT,
}