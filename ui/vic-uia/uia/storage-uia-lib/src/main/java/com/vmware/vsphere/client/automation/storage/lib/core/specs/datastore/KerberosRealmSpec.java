package com.vmware.vsphere.client.automation.storage.lib.core.specs.datastore;

import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * EntitySpec implementation for a kerberos realm
 *
 */
public class KerberosRealmSpec extends EntitySpec {

   /**
    * The DNS server for the kerberos realm.
    */
   public DataProperty<String> dnsServer;

   /**
    * The kerberos realm.
    */
   public DataProperty<String> activeDirectoryDomain;

   /**
    * The kerberos database server.
    */
   public DataProperty<String> activeDirectoryServer;

   /**
    * The username for the kerberos realm.
    */
   public DataProperty<String> activeDirectoryUserName;

   /**
    * The password for the kerberos realm.
    */
   public DataProperty<String> activeDirectoryPassword;

   /**
    * The username for datastore in the kerberos realm.
    */
   public DataProperty<String> storageUserName;

   /**
    * The password for datastore in the kerberos realm.
    */
   public DataProperty<String> storagePassword;

}
