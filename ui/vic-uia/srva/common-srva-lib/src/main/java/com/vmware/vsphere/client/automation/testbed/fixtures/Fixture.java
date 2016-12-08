/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed.fixtures;

import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;
import com.vmware.vsphere.client.automation.testbed.model.FixtureEntities;

/**
 * Interface that represents a way to get a specified Inventory Entity
 * that is contained in a specified fixture
 */
@Deprecated
public interface Fixture {

   /**
    * Method that returns the spec of a specified resource from a specified fixture setup
    * @param specClass - the class of the spec that will be generated
    * @param entity - the resource for which the spec will be generated
    * @return - generated object spec
    */

   public <T extends ManagedEntitySpec> T getFixtureResource(
         Class<T> specClass, FixtureEntities entity);

}
