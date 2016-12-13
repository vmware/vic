/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.common.WorkflowStepsSequence;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.TestBedBridge;

/**
 * The purpose of the class is to simplify the workflow for the so called elemental providers.
 * The elemental providers are providers that only goal is to provide basic resource without any
 * customizations. Such providers are HostProvider which only deploys a ESX and provides its credential.
 */
public abstract class BaseElementalProvider implements ProviderWorkflow {

   @Override
   public void initAssemblerSpec(AssemblerSpec assemblerSpec, TestBedBridge testbedBridge)
         throws Exception {
   }

   @Override
   public final void composeProviderSteps(WorkflowStepsSequence<? extends WorkflowStepContext> flow) throws Exception {
      flow.appendStep("Elemental provider", new ProviderWorkflowStep() {

         private boolean _onlyDeterminesResourceVersion = false;
         private boolean _retrieveResource = false;

         @Override
         public void prepare(PublisherSpec filteredPublisherSpec,
               AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
               SettingsReader sessionSettingsReader) {
            if ("true".equals(sessionSettingsReader
                  .getSetting("control.assemble.resourceVersionCheckOnly"))) {
               _onlyDeterminesResourceVersion = true;
            }

            if ("true".equals(sessionSettingsReader
                  .getSetting("control.assemble.retrieveResource"))) {
               _retrieveResource = true;
            }

            prepareForOperations(
                  filteredPublisherSpec,
                  filteredAssemblerSpec,
                  isAssembling,
                  sessionSettingsReader);
         }

         @Override
         public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
            // "control.assemble.resourceVersionCheckOnly=true"
            // "control.assemble.retrieveResource=true"
            // "resource.version="

            String resourceVersion = determineResourceVersion();
            testbedSettingsWriter.setSetting("resource.version", resourceVersion);
            if (_onlyDeterminesResourceVersion) {
               return;
            }

            if (_retrieveResource) {
               retrieveResource();
            } else {
               deployTestbed(testbedSettingsWriter);
            }
         }

         @Override
         public boolean checkHealth() throws Exception {
            return checkTestbedHealth();
         }

         @Override
         public void disassemble() throws Exception {
            destroyTestbed();
         }
      });
   }

   public abstract void prepareForOperations(PublisherSpec filteredPublisherSpec,
         AssemblerSpec filteredAssemblerSpec, boolean isAssembling,
         SettingsReader sessionSettingsReader); // OK

   public abstract String determineResourceVersion() throws Exception; // OK

   /**
    * Create template for downloaded resource
    * The method is invoked by the HSUIA system
    * --prepareTemplate
    * How to communicate between the
    * @throws Exception
    */
   public abstract void retrieveResource() throws Exception;

   @Override
   public abstract int providerWeight();

   public abstract void deployTestbed(SettingsWriter testbedSettingsWriter)
         throws Exception; // OK

   public abstract boolean checkTestbedHealth() throws Exception; // OK

   public abstract void destroyTestbed() throws Exception; // OK
}
