/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Holder for Datacenter action ID constants.
 */
public class DatacenterGlobalActions {

   public static final IDGroup AI_ALL_VCENTER_ACTIONS =
         IDGroup.toIDGroup("afContextMenu.vCenter");

   public static final IDGroup AI_REMOVE_DVS =
         IDGroup.toIDGroup("vsphere.core.dvs.deleteAction");

   public static final IDGroup AI_CREATE_DATACENTER =
         IDGroup.toIDGroup("vsphere.core.datacenter.createAction");

   public static final IDGroup AI_ADD_HOST =
         IDGroup.toIDGroup("vsphere.core.host.addAction");
}
