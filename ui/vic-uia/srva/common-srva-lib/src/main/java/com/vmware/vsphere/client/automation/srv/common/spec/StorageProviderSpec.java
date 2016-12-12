/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * A spec for a storage provider.
 */
public class StorageProviderSpec extends ManagedEntitySpec {

   /**
    * The URL of the provider.
    */
   public DataProperty<String> providerUrl;

   /**
    * The username for the storage provider URL.
    */
   public DataProperty<String> username;

   /**
    * The password for the storage provider URL.
    */
   public DataProperty<String> password;
}
