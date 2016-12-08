/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.srv.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

public class CdDvdDriveSpec extends ManagedEntitySpec {

   /**
    * Indicates the cdrom source type
    */
   public enum CdDvdDriveType {
      CLIENT_DEVICE("Client Device", "Connect to host CD device..."), DATASTORE_ISO_FILE(
            "Datastore ISO File", "Connect to CD/DVD image on a datastore..."), CONTENT_LIBRARY_ISO_FILE(
            "Content Library ISO File",
            "Connect to ISO image from a Content Library...");

      private String value;
      private String popupMenuString;

      private CdDvdDriveType(String value, String popupMenuString) {
         this.value = value;
         this.popupMenuString = popupMenuString;
      }

      public String value() {
         return value;
      }

      public String getPopupMenuString() {
         return popupMenuString;
      }
   }

   /**
    * Specifies if cd/dvd drive will be added, edited, or deleted
    */
   public DataProperty<CdDvdDriveType> cdDvdDriveType;

   /**
    * Specifies the cdroms name
    */
   public DataProperty<String> isoName;
   /**
    * Specifies the cdroms iso type
    */
   public DataProperty<String> isoType;
   /**
    * Specifies the cdroms content library
    */
   public DataProperty<String> isoContentLibrary;

   /**
    * Indicates whether the cdrom will be added, edited, or deleted
    */
   public enum CdDvdDriveActionType {
      EDIT("Edit"), ADD("Add"), DELETE("Delete");

      private String value;

      private CdDvdDriveActionType(String value) {
         this.value = value;
      }

      public String value() {
         return value;
      }
   }

   /**
    * Specifies if cd/dvd drive will be added, edited, or deleted
    */
   public DataProperty<CdDvdDriveActionType> cdDvdDriveAction;

   /**
    * Specifies if cd/dvd drive will be connected at power on
    */
   public DataProperty<Boolean> connectAtPowerOn;
}
