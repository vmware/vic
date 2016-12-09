/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec for backing tag
 */
public class BackingTagSpec extends ManagedEntitySpec {
   /**
    * Description of the object.
    */
   public DataProperty<String> description;

   /**
    * Category of the tag
    */
   public DataProperty<BackingCategorySpec> category;

   /**
    * the object tagged with this tag
    */
   public DataProperty<ManagedEntitySpec> taggedObjects;
}
