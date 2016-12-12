/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.common;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.URL;
import java.util.Arrays;
import java.util.List;
import java.util.Properties;

import com.google.common.base.Strings;

/**
 * This class is an implementation of the testbed handling when
 * it is described in a properties file. It holds the constants
 * necessary to describe the inventory and the implementation to
 * read the properties file.
 */
@Deprecated
public class LocalTestBedImpl implements CommonTestBed {

   // Delimiter to be used to specify list of properties in the DEFAULT_CONFIG_FILE
   private static final String LIST_PROPERTY_DELIMITER = ",";

   private final Properties config;

   // Configuration properties
   // Property for the list of ESX-es
   private static final String ESX_LIST = "esx.list";

   // Common setup objects
   private static final String VC_COMMON_DC = "vc.common.datacenter";
   private static final String VC_COMMON_CLUSTER = "vc.common.cluster";
   private static final String VC_COMMON_HOST = "vc.common.host";
   private static final String COMMON_DATASTORE = "common.datastore";
   private static final String CONTENT_LIBRARY = "content.library";
   private static final String VDC = "vdc";
   private static final String VC_NAME = "vc.ip";
   private static final String VC_USERNAME = "vc.username";
   private static final String VC_PASSWORD = "vc.password";
   private static final String VC_SSO_USERNAME = "ngc.username";
   private static final String VC_SSO_PASSWORD = "ngc.password";
   private static final String VC_THUMBPRINT = "vc.thumprint";

   public LocalTestBedImpl(String propertyFile) {
      Properties properties = new Properties();
      InputStream in = null;

      try {
         URL url = LocalTestBedImpl.class.getClassLoader().getResource(propertyFile);
         in = new FileInputStream(new File(url.toURI()));
         properties.load(in);
      } catch (IOException e) {
         throw new RuntimeException("IOExcpetion thrown while loading properties from: "
               + propertyFile, e);
      } catch (Exception e) {
         throw new RuntimeException("Excpetion thrown while loading properties from: "
               + propertyFile, e);
      } finally {
         try {
            if (in != null) {
               in.close();
            }
         } catch (IOException e) {
            throw new RuntimeException("Failed to close the input stream for the file: "
                  + propertyFile, e);
         }
      }

      config = properties;
   }

   //---------------------------------------------------------------------------
   // Common Inventory Objects

   @Override
   public String getVcName() {
      return config.getProperty(VC_NAME);
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
   public String getVcSsoUsername() {
      return config.getProperty(VC_SSO_USERNAME);
   }

   @Override
   public String getVcSsoPassword() {
      return config.getProperty(VC_SSO_PASSWORD);
   }

   @Override
   public String getVcThumbprint() {
      return config.getProperty(VC_THUMBPRINT);
   }

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getCommonDatacenterName()
    */
   @Override
   public String getCommonDatacenterName() {
      return config.getProperty(VC_COMMON_DC);
   }

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getCommonClusterName()
    */
   @Override
   public String getCommonClusterName() {
      return config.getProperty(VC_COMMON_CLUSTER);
   }

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getCommonHostName()
    */
   @Override
   public String getCommonHostName() {
      return config.getProperty(VC_COMMON_HOST);
   }

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getCommonDatastoreName()
    */
   @Override
   public String getCommonDatastoreName(){
      return config.getProperty(COMMON_DATASTORE);
   }

   //---------------------------------------------------------------------------
   // Fixture Specific Inventory Objects

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getContentLibraryName()
    */
   @Override
   public String getContentLibraryName(){
      return config.getProperty(CONTENT_LIBRARY);
   }

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getVdcName()
    */
   @Override
   public String getVdcName(){
      return config.getProperty(VDC);
   }

   //---------------------------------------------------------------------------
   // Inventory Objects needed for specific test scenarios(host list, network storage and etc.)

   /* (non-Javadoc)
    * @see com.vmware.vsphere.client.automation.testbed.common.TestbedUtil#getHosts(int)
    */
   @Override
   public List<String> getHosts(int hostCount) {
      String hostListConfigValue = config.getProperty(ESX_LIST);
      if(Strings.isNullOrEmpty(hostListConfigValue)) {
         throw new RuntimeException("No ESX hosts were found in the local testbed configuration file!");
      }

      List<String> hostsList = Arrays.asList(hostListConfigValue.split(LIST_PROPERTY_DELIMITER));
      if(hostsList.size() < hostCount) {
         throw new RuntimeException("The configuration file provides " + hostsList.size() + " but test requests " + hostCount);
      }

      return hostsList.subList(0, hostCount);
   }

}
