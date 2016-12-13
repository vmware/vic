/**
 * Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vim.binding.dataservice.tagging.CategoryInfo.Cardinality;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec for backing category
 */
public class BackingCategorySpec extends ManagedEntitySpec {
   /**
    * Description of the object.
    */
   public DataProperty<String> description;

   /**
    * Associated objects
    */
   public DataProperty<String> associatedObjects;

   /**
    * Cardinality - multiple tags per object or one tag per object - multiple or
    * single cardinality.
    */
   public DataProperty<Cardinality> cardinality;

}
