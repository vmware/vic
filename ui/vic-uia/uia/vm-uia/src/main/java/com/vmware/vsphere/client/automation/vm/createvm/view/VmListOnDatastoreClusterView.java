/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.createvm.view;

import com.vmware.vsphere.client.automation.srv.common.view.VmListView;

/**
 * The class represents the Virtual Machines list view at Datastore Cluster ->
 * VMs -> Virtual Machines
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.vm.lib.createvm.view.VmListOnDatastoreClusterView}
 */
@Deprecated
public class VmListOnDatastoreClusterView extends VmListView {

   @Override
   protected String getGridId() {
      return "vsphere.core.dscluster.relatedVMs/vmsForDatastoreCluster/list";
   }
}
