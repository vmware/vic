/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Container spec class that holds backing tag rename properties.
 * The category to which the target tag belongs, the tag itself
 * and the new configs that should be applied to the tag.
 * Change in category to which the tag belongs is not supported.
 */
public class UpdateBackingTagSpec extends BaseSpec {

   /**
    * Holds spec of the backing tag category
    * to which the target tag belongs.
    */
   public DataProperty<BackingCategorySpec> category;

   /**
    * The spec of the tag that will be updated.
    * It has to belong to the category above.
    */
   public DataProperty<BackingTagSpec> targetTag;

   /**
    * The spec that contains the new tag configurations.
    */
   public DataProperty<BackingTagSpec> newTargetConfigs;
}
