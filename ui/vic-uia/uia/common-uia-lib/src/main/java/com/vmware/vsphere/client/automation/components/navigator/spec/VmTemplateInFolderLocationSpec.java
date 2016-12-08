/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.components.navigator.spec;

import com.vmware.vsphere.client.automation.components.navigator.NGCNavigator;

/**
 * This class is a location spec which represents navigation to
 * vCenter>VM Templates in Folders>VM_template.
 */
public class VmTemplateInFolderLocationSpec extends NGCLocationSpec {

   /**
    * Navigation to vCenter>VM Templates in Folders.
    */
   public VmTemplateInFolderLocationSpec() {
      this(null);
   }

   /**
    * Navigation to vCenter>VM Templates in Folders>VM_template.
    */
   public VmTemplateInFolderLocationSpec(String vmTemplateName) {
      super(NGCNavigator.NID_HOME_VCENTER,
            NGCNavigator.NID_VCENTER_TEMPLATES_IN_FOLDERS,
            vmTemplateName);
   }
}
