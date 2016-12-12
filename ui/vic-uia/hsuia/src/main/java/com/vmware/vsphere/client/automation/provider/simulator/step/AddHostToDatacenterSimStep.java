/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.vsphere.client.automation.provider.simulator.step;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.explorer.SettingsReader;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;
import com.vmware.client.automation.workflow.explorer.SpecTraversalUtil;
import com.vmware.client.automation.workflow.provider.AssemblerSpec;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowStep;
import com.vmware.client.automation.workflow.provider.PublisherSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorDatacenterSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorHostSpec;
import com.vmware.vsphere.client.automation.provider.simulator.util.SimulatorUtil;

public class AddHostToDatacenterSimStep implements ProviderWorkflowStep {

   private static final Logger _logger = LoggerFactory.getLogger(AddHostToDatacenterSimStep.class);

   public static final String TESTBED_KEY_DCHOST_NAME = "testbed.inventory.datacenterhost.name";
   private static final String CONTROL_KEY_RNDFAILURE = "control.simulator.allowRandomFailure";

   private SimulatorDatacenterSpec _publisherDatacenterSpec;
   private SimulatorHostSpec _publisherDatacenterHostSpec;

   private SimulatorDatacenterSpec _assemblerDatacenterSpec;;
   private SimulatorHostSpec _assemblerDatacenterHostSpec;

   private boolean _allowRandomFailure;

   @Override
   public void prepare(
         PublisherSpec filteredPublisherSpec, AssemblerSpec filterAssemblerSpec,
         boolean isAssembling, SettingsReader sessionSettingsReader)  throws Exception {

      _allowRandomFailure =
            SettingsUtil.getBooleanValue(sessionSettingsReader, CONTROL_KEY_RNDFAILURE);

      if (isAssembling) {
         _logger.info("Loading datacenter and datacenter's host assembler specs");

         _assemblerDatacenterSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(
                     filterAssemblerSpec, SimulatorDatacenterSpec.class);
         _assemblerDatacenterHostSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(
                     filterAssemblerSpec, SimulatorHostSpec.class);
      } else {
         _logger.info("Loading datacenter and datacenter's host publisher specs");

         _publisherDatacenterSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(
                     filteredPublisherSpec, SimulatorDatacenterSpec.class);
         _publisherDatacenterHostSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(
                     filteredPublisherSpec, SimulatorHostSpec.class);
      }
   }

   @Override
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      // Pretend adding a host to the datacenter.

      _logger.info("Simulating inventory construction by waiting a second...");

      Thread.sleep(1000);
      pretendAddHost(_assemblerDatacenterSpec, _assemblerDatacenterHostSpec);
      SimulatorUtil.trySucceed(_allowRandomFailure);

      _logger.info("Saving the host identification");

      // Save inventory identifier
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_DCHOST_NAME,
            _assemblerDatacenterHostSpec.service.get().endpoint.get());
   }

   @Override
   public boolean checkHealth() throws Exception {
      // Pretend to check the published entities assembled here exists.
      _logger.info(
            "Simulating inventory health-checking by verifying the entity names are set.");

      // TODO: Add verification for correct parent-child relationships

      return !Strings.isNullOrEmpty(_publisherDatacenterHostSpec.name.get());
   }

   @Override
   public void disassemble() throws Exception {
      // Pretend to remove the published entities added here.
      _logger.info("Simulating removal of intentory items by waiting for a second...");

      Thread.sleep(1000);

      SimulatorUtil.trySucceed(_allowRandomFailure);

      _logger.info(
            String.format(
                  "Pretending to remove empty cluster: %s from virtual center: %",
                  _publisherDatacenterHostSpec.name.get(),
                  _publisherDatacenterSpec.name.get()));
   }

   private void pretendAddHost(SimulatorDatacenterSpec hostingEntitySpec, SimulatorHostSpec hostSpec) {
      _logger.info(
            String.format(
                  "Pretending to add host %s to entity %s",
                  hostSpec.toString(),
                  hostingEntitySpec.toString()));
   }

}
