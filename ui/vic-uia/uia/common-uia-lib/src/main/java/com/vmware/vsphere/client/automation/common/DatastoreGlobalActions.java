/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common;

import com.vmware.suitaf.apl.IDGroup;

/**
 * Class that encapsulates all datastore actions.
 */
public class DatastoreGlobalActions {

   //---------------------------------------------------------------------------
   // IDs of datastore actions

   // Create new datastore
   public static final IDGroup AI_CREATE_DATASTORE =
         IDGroup.toIDGroup("vsphere.core.datastore.addActionGlobal/button");
}
