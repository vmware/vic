/* Copyright 2013 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.common.chunkbuilder;

import com.vmware.client.automation.testbed.TestBed;
import com.vmware.client.automation.workflow.WorkflowComposition;
import com.vmware.hsua.common.datamodel.PropertyBoxLinks;

/**
 * Common test chunk builder interface.
 * It allows defining specs and steps in one place.
 */
public interface TestChunkBuilder {

   /**
    * Gives possibility to append specs to the provided links.
    *
    * @param testBed    test bed that can be used for retrieving test bed elements
    * @param links      links to which created specs will be added.
    */
   public void appendSpecs(TestBed testBed, PropertyBoxLinks links);

   /**
    * Appends prerequisite steps to the workflow.
    *
    * @param composition   composition to which prerequisite steps will be appended
    */
   public void appendPrereqSteps(WorkflowComposition composition);

   /**
    * Appends test steps to the workflow.
    *
    * @param composition   composition to which test steps will be appended
    */
   public void appendTestSteps(WorkflowComposition composition);
}
