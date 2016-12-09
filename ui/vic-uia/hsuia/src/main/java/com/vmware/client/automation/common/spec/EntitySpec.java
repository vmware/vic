/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.common.spec;

import com.vmware.hsua.common.datamodel.BasePBox.DataProperty;

/**
 * A generic entity class used to represent all types of external objects and
 * systems the tests and test bed providers deal with.
 */
public class EntitySpec extends BaseSpec {

   /**
    * Specification for accessing the service that supports the entity.
    */
   public DataProperty<ServiceSpec> service;
}
