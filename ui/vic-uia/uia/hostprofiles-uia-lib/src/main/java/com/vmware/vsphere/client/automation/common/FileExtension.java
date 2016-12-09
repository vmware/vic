/*
 *  Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common;

/**
 * Enum representing a file extension. Several tests deal with files and this
 * class centralizes the file extension suffix details.
 */
public enum FileExtension {
   CSV("csv"),
   VPF("vpf");

   private static final String FILENAME_DELIMETER = ".";
   private String extension;

   FileExtension(String extension) {
      this.extension = extension;
   }

   public String getSuffix(){
      return FILENAME_DELIMETER + extension;
   }
}
