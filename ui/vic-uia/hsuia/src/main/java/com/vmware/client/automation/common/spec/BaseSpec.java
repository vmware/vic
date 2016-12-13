/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.spec;

import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.hsua.common.datamodel.PropertyBox;
import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * This class will be the base class for the NGC vCD UI data model. It is based
 * on the SUITA PropertyBox class definition. For now only the Name annotation
 * is provided.
 */
public class BaseSpec extends PropertyBox {

   protected static final Logger _logger = LoggerFactory.getLogger(BaseSpec.class);

   @Retention(RetentionPolicy.RUNTIME)
   public static @interface Name {
   }

   /**
    * General tags annotating the specification.
    */
   public DataProperty<String> tag;
}
