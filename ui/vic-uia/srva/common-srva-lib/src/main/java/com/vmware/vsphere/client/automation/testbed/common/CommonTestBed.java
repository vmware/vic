/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.common;

import java.util.List;

/**
 * Interface that represents the utility methods to get
 */
@Deprecated
public interface CommonTestBed {

   /**
    * Returns name of the virtual center.
    */
   public abstract String getVcName();

   /**
    * Returns login username.
    */
   public abstract String getVcUsername();

   /**
    * Returns login password.
    */
   public abstract String getVcPassword();

   /**
    * Returns login username.
    */
   public abstract String getVcSsoUsername();

   /**
    * Returns login password.
    */
   public abstract String getVcSsoPassword();

   /**
    * Return virtual center thumbprint.
    */
   public abstract String getVcThumbprint();

   /**
    * Method that returns the name of the datacenter in the common
    * inventory setup
    * @return String with the name of the datacenter
    */
   public abstract String getCommonDatacenterName();

   /**
    * Method that returns the name of the cluster in the common
    * inventory setup
    * @return String with the name of the cluster
    */
   public abstract String getCommonClusterName();

   /**
    * Method that returns the name of the clustered host in the
    * common inventory setup
    * @return String with the name of the clustered host
    */
   public abstract String getCommonHostName();

   /**
    * Method that returns the name of the shared datastore in the
    * common inventory setup
    * @return String with the name of the shared datastore
    */
   public abstract String getCommonDatastoreName();

   /**
    * Method that returns the name of the content library in the
    * inventory setup common for content library tests
    * @return String with the name of the content library
    */
   public abstract String getContentLibraryName();

   /**
    * Method that returns the name of the vdc in the inventory
    * setup common to vdc tests
    * @return String with the name of the vdc
    */
   public abstract String getVdcName();

   /**
    * Method that returns the names of available free hosts in
    * the inventory setup
    * @param hostCount
    * @return List of the names of the available free hosts
    */
   public abstract List<String> getHosts(int hostCount);
}