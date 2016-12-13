/*
 *  Copyright 2016 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.common.spec;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;
import com.vmware.vsphere.client.automation.common.FileExtension;

/**
 * Spec representing a file extension.
 */
public class FileExtensionSpec extends BaseSpec {
   public DataProperty<FileExtension> extension;
}
