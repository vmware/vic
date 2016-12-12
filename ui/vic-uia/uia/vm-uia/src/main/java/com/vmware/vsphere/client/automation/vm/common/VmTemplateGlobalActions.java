/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.vm.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Class that holds the IDs of the actions in the All Actions menu of VM
 * templates.
 *
 * @deprecated  Replaced by
 *    {@link com.vmware.vsphere.client.automation.common.VmTemplateGlobalActions}}
 */
@Deprecated
public class VmTemplateGlobalActions {

   // Sub Menus
   public static final IDGroup AI_NEW_VM_FROM_THIS_TEMPLATE = IDGroup
         .toIDGroup("vsphere.core.vm.provisioning.cloneTemplateToVmAction");

}
