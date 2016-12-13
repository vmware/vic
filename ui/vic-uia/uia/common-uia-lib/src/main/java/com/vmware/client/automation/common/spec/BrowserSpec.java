/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * Data spec containing browser connection and settings data.
 */
public class BrowserSpec extends BaseSpec {

   /**
    * The URL to open.
    */
   public DataProperty<String> url;
}
