/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.testbed;

import java.util.Arrays;
import java.util.List;
import java.util.Properties;

import com.google.common.base.Strings;
import com.vmware.hsua.common.util.IOUtils;

/**
 * This class is used to load and keep all the common resources used by the
 * tests.
 */
public class LocalTestBed extends TestBed {

   private static final String DEFAULT_CONFIG_FILE = "localTestbed.properties";
   // Delimiter to be used to specify list of properties in the DEFAULT_CONFIG_FILE
   private static final String LIST_PROPERTY_DELIMITER = ",";

   private final Properties config;

   // Configuration properties

   // NGC client credentials
   private static final String NGC_USERNAME = "ngc.username";
   private static final String NGC_PASSWORD = "ngc.password";
   private static final String NGC_IP = "ngc.ip";
   private static final String NGC_PORT = "ngc.port";

   // vCD credentials
   private static final String VCD_IP = "vcd.ip";
   private static final String VCD_USERNAME = "vcd.username";
   private static final String VCD_PASSWORD = "vcd.password";

   // VC Settings
   private static final String VC_IP = "vc.ip";
   private static final String VC_USERNAME = "vc.username";
   private static final String VC_PASSWORD = "vc.password";
   private static final String VC_THUMBPRINT = "vc.thumprint";

   // ESX Settings
   private static final String ESX_ADMIN_USER_NAME = "esx.admin.username";
   private static final String ESX_ADMIN_USER_PASSWORD = "esx.admin.password";
   // Property for the list of ESX-es
   private static final String ESX_LIST = "esx.list";

   // Common setup objects
   private static final String VCD_COMMON_CRP = "vcd.common.crp";
   private static final String VC_COMMON_DC = "vc.common.datacenter";
   private static final String VC_COMMON_CLUSTER = "vc.common.cluster";
   private static final String VC_COMMON_HOST = "vc.common.host";
   private static final String COMMON_DATASTORE = "common.datastore";
   private static final String COMMON_DATASTORE_SIZE = "common.datastore.size";
   private static final String LOCAL_DATASTORE = "local.datastore";
   private static final String LOCAL_DATASTORE_SIZE = "local.datastore.size";
   private static final String VVO_SERVER_IP = "vvol.datastore.ip";


   public LocalTestBed(String configurationFile) {
      if (configurationFile.length() == 0) {
         configurationFile = getDefaultFileName(DEFAULT_CONFIG_FILE);
      }

      config = IOUtils.readConfiguration(configurationFile);
   }


   //---------------------------------------------------------------------------
   // NGC Settings

   @Override
   protected String getUserName() {
      return config.getProperty(NGC_USERNAME);
   }

   @Override
   protected String getUserPassword() {
      return config.getProperty(NGC_PASSWORD);
   }

   @Override
   protected String getNgcIp() {
      return config.getProperty(NGC_IP);
   }

   @Override
   protected String getNgcPort() {
      return config.getProperty(NGC_PORT);
   }


   //---------------------------------------------------------------------------
   // VCD Settings

   @Override
   public String getVCD() {
      return config.getProperty(VCD_IP);
   }

   @Override
   public String getVCDUsername() {
      return config.getProperty(VCD_USERNAME);
   }

   @Override
   public String getVCDPassword() {
      return config.getProperty(VCD_PASSWORD);
   }


   //---------------------------------------------------------------------------
   // VC Settings

   @Override
   public String getVc() {
      return config.getProperty(VC_IP);
   }

   @Override
   public String getVcUsername() {
      return config.getProperty(VC_USERNAME);
   }

   @Override
   public String getVcPassword() {
      return config.getProperty(VC_PASSWORD);
   }

   @Override
   public String getVcThumbprint() {
      return config.getProperty(VC_THUMBPRINT);
   }


   //---------------------------------------------------------------------------
   // Default ESX admin user settings
   @Override
   public String getESXAdminUsername() {
      return config.getProperty(ESX_ADMIN_USER_NAME);
   }


   @Override
   public String getESXAdminPasssword() {
      return config.getProperty(ESX_ADMIN_USER_PASSWORD);
   }

   //---------------------------------------------------------------------------
   // Common Inventory Objects

   @Override
   public String getCommonCrpName() {
      return config.getProperty(VCD_COMMON_CRP);
   }

   @Override
   public String getCommonDatacenterName() {
      return config.getProperty(VC_COMMON_DC);
   }

   @Override
   public String getCommonClusterName() {
      return config.getProperty(VC_COMMON_CLUSTER);
   }

   @Override
   public String getCommonHost() {
      return config.getProperty(VC_COMMON_HOST);
   }

   @Override
   public String getCommonDatastoreName(){
      return config.getProperty(COMMON_DATASTORE);
   }

   @Override
   public String getCommonDatastoreSize() {
      return config.getProperty(COMMON_DATASTORE_SIZE);
   }

   @Override
   public String getLocalDatastoreName() {
      return config.getProperty(LOCAL_DATASTORE);
   }

   @Override
   public String getLocalDatastoreSize() {
      return config.getProperty(LOCAL_DATASTORE_SIZE);
   }
   @Override
   public String getVvolServerIp() {
      return config.getProperty(VVO_SERVER_IP);
   }

   //---------------------------------------------------------------------------
   // Inventory Objects needed for specific test scenarios(host list, network storage and etc.)
   @Override
   public List<String> getHosts(int hostCount) {
      String hostListConfigValue = config.getProperty(ESX_LIST);
      if(Strings.isNullOrEmpty(hostListConfigValue)) {
         throw new RuntimeException("No ESX hosts were found in the local testbed configuration file!");
      }

      List<String> hostsList = Arrays.asList(hostListConfigValue.replace(
            LIST_PROPERTY_DELIMITER + " ", LIST_PROPERTY_DELIMITER).split(
                  LIST_PROPERTY_DELIMITER));
      if(hostsList.size() < hostCount) {
         throw new RuntimeException("The configuration file provides " + hostsList.size() + " but test requests " + hostCount);
      }

      return hostsList.subList(0, hostCount);

   }
}
