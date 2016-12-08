/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import com.vmware.client.automation.workflow.common.WorkflowStep;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;

/**
 * The interface defines the API of provider workflow step.
 *
 * Provider steps should implement full life-cycle of resource management from 
 * assembling, through verification until the the resource disposal. Therefore
 * steps should be designed in a way that this phases are present in all of them.
 */
public interface ProviderWorkflowStep extends WorkflowStep {

   /**
    * Represents general preparation phase of the step.
    * 
    * Use this method to retrieve the required specs from the provider's root
    * spec and save them in member variables so they can be used later in the
    * assemble, verification, or clean-up phases.
    * 
    * filteredProviderSpec
    *    Reference to the <code>ProviderSpec</code> of the provider's workflow.
    * @throws Exception 
    */
   // public void prepare(ProviderSpec filteredProviderSpec, SettingsReader sessionSettingsReader);
   public void prepare(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filterAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader) throws Exception;

   /**
    * Represents general assemble phase of the step.
    * This phase has to be triggered if there are no fixtures available.
    *
    * @throws Exception    if something goes wrong in the assemble phase
    */
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception;

   /**
    * Represents health checking step for configuration this step cares about.
    *
    * @return              true if the verification of this step was successful,
    *                      false otherwise
    * @throws Exception    if something goes wrong in the verification phase
    */
   public boolean checkHealth() throws Exception;

   /**
    * Represents general disassemble phase.
    *
    * @throws Exception    if disassemble of the fixture workflow step is not possible
    */
   public void disassemble() throws Exception;
}
