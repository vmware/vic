/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Container class for datastore properties.
 */
public class DatastoreSpec extends ManagedEntitySpec {

   /**
    * Datastore type.
    */
   public DataProperty<DatastoreType> type;

   /**
    * Remote host.
    * This property applies only if NFS type is used.
    */
   public DataProperty<String> remoteHost;

   /**
    * Remote path.
    * This property applies only if NFS type is used.
    */
   public DataProperty<String> remotePath;



}
