/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;


/**
 * A spec used in the RenameEntitySpec class
 */
public class RenameEntitySpec extends BaseSpec {
   /**
    * Old name of the entity
    */
   public DataProperty<String> oldName;

   /**
    * New name of the entity
    */
   public DataProperty<String> newName;

   /**
    * Does the new entity name already exist?
    */
   public DataProperty<Boolean> isNewNameExisting;

   /**
    * Task name
    */
   public DataProperty<String> taskName;
}
