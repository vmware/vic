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
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorClusterSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorDatacenterSpec;
import com.vmware.vsphere.client.automation.provider.simulator.spec.SimulatorVcSpec;
import com.vmware.vsphere.client.automation.provider.simulator.util.SimulatorUtil;

public class CreateVirtualInventorySimStep implements ProviderWorkflowStep {

   private static final Logger _logger =
         LoggerFactory.getLogger(CreateVirtualInventorySimStep.class);

   public static final String TESTBED_KEY_SERVICE_ENDPOINT = "testbed.service.default.endpoint";
   public static final String TESTBED_KEY_SERVICE_USERNAME = "testbed.service.default.user";
   public static final String TESTBED_KEY_SERVICE_PASSWORD = "testbed.service.default.pass";

   public static final String TESTBED_KEY_VC_NAME = "testbed.inventory.virtualcenter.name";
   public static final String TESTBED_KEY_DC_NAME = "testbed.inventory.datacenter.name";
   public static final String TESTBED_KEY_CL_NAME = "testbed.inventory.cluster.name";

   private static final String CONTROL_KEY_RNDFAILURE = "control.simulator.allowRandomFailure";

   private SimulatorVcSpec _publisherVcSpec;
   private SimulatorDatacenterSpec _publisherDcSpec;
   private SimulatorClusterSpec _publisherClusterSpec;

   private SimulatorVcSpec _assemblerVcSpec;
   private SimulatorDatacenterSpec _assemblerDcSpec;
   private SimulatorClusterSpec _assemblerClusterSpec;

   private boolean _allowRandomFailure;

   @Override
   public void prepare(
         PublisherSpec filteredPublisherSpec, AssemblerSpec filterAssemblerSpec,
         boolean isAssembling, SettingsReader sessionSettingsReader) throws Exception {

      _allowRandomFailure = SettingsUtil.getBooleanValue(
            sessionSettingsReader, CONTROL_KEY_RNDFAILURE);

      if (isAssembling) {
         _logger.info("Loading Virtual Center, Datacenter and Cluster assembler specs");

         _assemblerVcSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filterAssemblerSpec, SimulatorVcSpec.class);
         _assemblerDcSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filterAssemblerSpec, SimulatorDatacenterSpec.class);
         _assemblerClusterSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filterAssemblerSpec, SimulatorClusterSpec.class);
      } else {
         _logger.info("Loading Virtual Center, Datacenter and Cluster publisher specs");

         _publisherVcSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filteredPublisherSpec, SimulatorVcSpec.class);
         _publisherDcSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filteredPublisherSpec, SimulatorDatacenterSpec.class);
         _publisherClusterSpec =
               SpecTraversalUtil.getRequiredSpecFromContainerNode(filteredPublisherSpec, SimulatorClusterSpec.class);
      }
   }

   @Override
   public void assemble(SettingsWriter testbedSettingsWriter) throws Exception {
      // Pretend that an inventory is built.

      _logger.info("Simulating inventory construction by waiting a second...");

      Thread.sleep(1000);

      SimulatorUtil.trySucceed(_allowRandomFailure);

      _logger.info("Saving test bed connection data...");

      // Save connection to VC
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_SERVICE_ENDPOINT,
            _assemblerVcSpec.service.get().endpoint.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_SERVICE_USERNAME,
            _assemblerVcSpec.service.get().username.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_SERVICE_PASSWORD,
            _assemblerVcSpec.service.get().password.get());

      // Save inventory identifiers
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_VC_NAME,
            _assemblerVcSpec.name.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_DC_NAME,
            _assemblerDcSpec.name.get());
      testbedSettingsWriter.setSetting(
            TESTBED_KEY_CL_NAME,
            _assemblerClusterSpec.name.get());
   }

   @Override
   public boolean checkHealth() throws Exception {
      // Pretend to check the published entities assembled here exists.
      _logger.info(
            "Simulating inventory health-checking by verifying the entity names are set.");

      // TODO: Add verification for correct parent-child relationships

      return !Strings.isNullOrEmpty(_publisherVcSpec.name.get())
            && !Strings.isNullOrEmpty(_publisherDcSpec.name.get())
            && !Strings.isNullOrEmpty(_publisherVcSpec.name.get());
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
                  _publisherClusterSpec.name.get(),
                  _publisherVcSpec.name.get()));

      _logger.info(
            String.format(
                  "Pretending to remove empty datacenter: %s from virtual center: %",
                  _publisherDcSpec.name.get(),
                  _publisherVcSpec.name.get()));
   }

}
