/**
 * Copyright 2013 VMware, Inc.  All rights reserved. -- VMware Confidential
 */
package com.vmware.vsphere.client.automation.testbed;

import sun.reflect.generics.reflectiveObjects.NotImplementedException;

import com.vmware.vsphere.client.automation.testbed.fixtures.Fixture;
import com.vmware.vsphere.client.automation.testbed.fixtures.FixtureImpl;
import com.vmware.vsphere.client.automation.testbed.model.FixtureEntities;
import com.vmware.vsphere.client.automation.testbed.model.Fixtures;
import com.vmware.vsphere.client.automation.testbed.spec.StandaloneResourceSpec;

/**
 * Class that is used to get an inventory setup, i.e. fixture, that corresponds
 * to a specific setup or a standalone resource
 */
@Deprecated
public class TestBed {
   private final String propertyFile;

   public TestBed() {
      propertyFile = "primaryVcFixture.properties";
   }

   public TestBed(String propertyFile) {
      this.propertyFile = propertyFile;
   }

   /**
    * Gets a specific Fixture according to supplied parameter
    * @param fixture - fixture to get
    * @return - Fixture
    */
   public Fixture getFixture(Fixtures fixture) {
      return FixtureImpl.createFixture(fixture, propertyFile);
   }

   /**
    * Method that returns the spec of a specified standalone resource
    * @param specClass - the class of the spec that will be generated
    * @param entity - the resource for which the spec will be generated
    * @return - generated object spec
    * @throws Exception
    */
   public <T extends StandaloneResourceSpec> T getStandaloneResource(
         Class<T> specClass, FixtureEntities entity) throws Exception {
      throw new NotImplementedException();
   }

}
