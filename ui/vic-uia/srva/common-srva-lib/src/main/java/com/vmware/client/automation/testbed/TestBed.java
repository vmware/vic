/**
 * Copyright 2012 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.testbed;

import java.io.File;
import java.net.URL;
import java.util.List;

import com.vmware.vsphere.client.automation.srv.common.spec.UserSpec;


/**
 * Definition of base test bed.
 */
public abstract class TestBed {
   private static final String NGC_URL_TEMAPLATE = "https://%s:%s/vsphere-client/?logLevel=error";

   // User specification
   private UserSpec _userSpec;

   // Ngc url
   private String _ngcUrl;

   /**
    * Gets administrator user credentials.
    *
    * @return
    *       Data model representing the vCD administrator.
    */
   public UserSpec getAdminUser() {
      if (_userSpec == null) {
         _userSpec = new UserSpec();
         _userSpec.username.set(getUserName());
         _userSpec.password.set(getUserPassword());
      }
      return _userSpec;
   }

   /**
    * Gets NGC url
    * .
    * @return
    *       The URL to the login page of the NGC client.
    */
   public String getNGCURL() {
      if (_ngcUrl == null) {
         _ngcUrl = String.format(NGC_URL_TEMAPLATE, getNgcIp(), getNgcPort());
      }
      return _ngcUrl;
   }

   /**
    * Prepares the test bed configuration
    */
   public void startUp() {
      // Nothing for testbed preparation .
   }

   /**
    * Cleans up the test bed.
    */
   public void cleanUp() {
      // Nothing to clean up.
   }

   /**
    * Finds file stored in class loader bootstrap directory.
    */
   protected String getDefaultFileName(String fileName) {
      try {
         URL url = this.getClass().getClassLoader().getResource(fileName);
         File f = new File(url.toURI());

         return f.getCanonicalPath();
      } catch (Exception e) {
         throw new RuntimeException("Cannot locate file: " + fileName, e);
      }

   }

   //---------------------------------------------------------------------------
   // NGC Settings

   /**
    * Gets user name.
    */
   protected abstract String getUserName();

   /**
    * Gets user password.
    */
   protected abstract String getUserPassword();

   /**
    * Gets NGC ip address.
    */
   protected abstract String getNgcIp();

   /**
    * Gets NGC port.
    */
   protected abstract String getNgcPort();

   //---------------------------------------------------------------------------
   // VCD Settings

   /**
    * Gets vCD IP.
    */
   public abstract String getVCD();

   /**
    * Gets vCD administrator user name.
    */
   public abstract String getVCDUsername();

   /**
    * Gets vCD administrator password.
    */
   public abstract String getVCDPassword();


   //---------------------------------------------------------------------------
   // VC Settings

   /**
    * Gets the IP of the VC
    *
    */
   public abstract String getVc();

   /**
    * Gets VC administrator user name
    *
    */
   public abstract String getVcUsername();

   /**
    * Gets VC administrator password
    */
   public abstract String getVcPassword();

   /**
    * Gets VC thumbprint
    */
   public abstract String getVcThumbprint();


   //---------------------------------------------------------------------------
   // Default ESX admin user Settings

   /**
    * Gets ESX default administrator user name.
    */
   public abstract String getESXAdminUsername();

   /**
    * Gets ESX default administrator user password.
    */
   public abstract String getESXAdminPasssword();

   //---------------------------------------------------------------------------
   // Common Inventory Objects

   /**
    * Gets the name of the CRP provided as common setup.
    */
   public abstract String getCommonCrpName();

   /**
    * Gets the name of the Datacenter provided as common inventory setup.
    */
   public abstract String getCommonDatacenterName();

   /**
    * Gets the name of the Cluster provided as common inventory setup.
    */
   public abstract String getCommonClusterName();

   /**
    * Gets the IP/domain name of the ESX provided as common inventory setup.
    */
   public abstract String getCommonHost();

   /**
    * Gets the common datastore name
    */
   public abstract String getCommonDatastoreName();

   /**
    * Gets the common datastore size
    */
   public abstract String getCommonDatastoreSize();

   /**
    * Gets the local datastore name.
    */
   public abstract String getLocalDatastoreName();

   /**
    * Gets the local datastore size
    */
   public abstract String getLocalDatastoreSize();

   /**
    * Gets the VVOL server IP.
    */
   public abstract String getVvolServerIp();

   //---------------------------------------------------------------------------
   // Inventory Objects needed for specific test scenarios(host list, network storage and etc.)

   /**
    * Return list of host IPs/domain names that can be used by the tests.
    *
    * @param hostCount    number of host to be returned
    * @return list of hosts
    */
   public abstract List<String> getHosts(int hostCount);
}
