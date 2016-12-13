/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Properties;
import java.util.Set;

import org.apache.commons.collections4.CollectionUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.common.spec.BaseSpec;
import com.vmware.client.automation.common.spec.EntitySpec;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.connector.TestbedConnector;
import com.vmware.client.automation.workflow.common.StepPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowController;
import com.vmware.client.automation.workflow.explorer.ProviderBridgeImpl;
import com.vmware.client.automation.workflow.explorer.SpecTraversalUtil;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.vsphere.client.automation.common.spec.ManagedEntitySpec;

/**
 * Controller that implementation the provider workflow stages.
 */
public class ProviderWorkflowController extends WorkflowController {
   private static final String SETTINGS_KEY_PROVIDER_ID = "provider.id";

   private ProviderWorkflowContext _context;
   private String _testBedConfigFilePath;

   private ProviderBridgeImpl _providerBridge;

   // logger
   private static final Logger _logger =
         LoggerFactory.getLogger(ProviderWorkflowController.class);

   @Override
   public WorkflowRegistry getRegistry() {
      return WorkflowRegistry.getRegistry();
   }

   public static ProviderWorkflowController create(ProviderWorkflowContext context,
         String testBedConfigFilePath) {
      ProviderWorkflowController controller = new ProviderWorkflowController();
      controller._context = context;
      controller._testBedConfigFilePath = testBedConfigFilePath;
      controller._providerBridge = new ProviderBridgeImpl(controller._context);
      return controller;
   }

   public void analyze() throws ProviderControllerException, ProviderWorkflowException {
      guardSingleRequest();

      runWorkflowIntantiate();

      // runWorkflowStepsCompose();
      // runStepsSequencePrepare(true);

      this._providerBridge.analyze();
   }

   public PublisherSpec getPublisherSpec()
         throws ProviderControllerException, ProviderWorkflowException {
      // Could this be static i.e. find provider,

      // allocate context
      guardSingleRequest();

      runWorkflowIntantiate();

      // 5. Return the provider. The returned 'publisher' spec only?
      return _context.getPublisherSpec();
   }

   public void assignTestBed() throws ProviderControllerException, ProviderWorkflowException {

      guardSingleRequest();

      runWorkflowTestbedSettingsLoad();

      runWorkflowAssignTestBed();

      runWorkflowAssignConnectors();
   }

   // initSpec -> bridge -> controller -> registry -> context (map)
   // search and find only top level fixtures

   public void connect() throws ProviderControllerException, ProviderWorkflowException {
      guardSingleRequest();

      runEstablishConnections();
   }

   /**
    * Run this scenario to validate that provider is able to successfully assign
    * the provided test bed settings to the specs.
    *
    * @throws ProviderWorkflowException
    * @throws ProviderControllerException
    */
   public void register() throws ProviderWorkflowException, ProviderControllerException {
      guardSingleRequest();

      // Builds new workflow instance
      runWorkflowIntantiate();

      // 2. Load configuration file
      runWorkflowTestbedSettingsLoad();

      // Validate
      runWorkflowAssignTestBed();

      // Make the registration with the registry?
      getRegistry().addTestBed(
            _context.getProviderWorkflow(),
            _testBedConfigFilePath,
            _context.getTestbedSettings());
   }

   public void assemble() throws ProviderControllerException, ProviderWorkflowException {
      guardSingleRequest();

      // Clear context or throw exception if this method is called twice

      // Instantiate the workflow for assembling
      runWorkflowIntantiate();

      // 0. Init spec - done
      // 1. Compose steps - done
      runWorkflowStepsCompose();

      // 2. Claim test bed instances as participating in this test bed /
      // workflow / run (call call initSpecs(), give connections settings(),
      // /open connections/)
      _providerBridge.claimTestbeds();

      // 2.1. Automatic assign of the respective service specs based on
      // the hierarchy defined by the parent property.
      // TODO: rkovachev - move the logic to ProviderWorkflow
      runWorkflowAssignServices(_context.getAssemblerSpec());

      // 3. Execute prepare() for all steps - done
      runStepsSequencePrepare(true);

      // 4. Establish the connections as listed in the corresponding ServiceSpec
      // classes./ If needed the connections might be stored in the registry./
      _providerBridge.connectToTestbeds();

      // 5. Execute assemble()
      try {
         runStepsSequenceAssemble();
      } catch (ProviderWorkflowException pwe) {
         pwe.printStackTrace();
         try {
            collectNewSettings();
            runWorkflowAssemblePhaseRollback(getNotExecutedAssemblingStepsCount());
         } catch (ProviderWorkflowException e) {
            // TODO: set
         }
      }

      // 5. Execute check() for all steps; On failure, rollback by calling
      // disassemble. - done
      try {
         collectNewSettings();
         runWorkflowIntantiate();
         runWorkflowStepsCompose();
         runWorkflowAssignTestBed();
         runWorkflowAssignConnectors();
         runStepsSequencePrepare(false);
         runEstablishConnections();

         runStepsSequenceCheckHealth(); // ->r re-inite, check, on failure
         // rollback

      } catch (ProviderWorkflowException pwe) {
         pwe.printStackTrace();
         runWorkflowAssemblePhaseRollback(0);
      }

      // 7. Save settings file - done
      runWorkflowTestbedSettingsSave();

      // Close the connections - this is quite a mess at this point - should me
      // moved inside the sequencers.
      _providerBridge.releaseTestbeds();

      // return true or throw?
   }

   public void checkHealth() throws ProviderControllerException, ProviderWorkflowException {
      guardSingleRequest();

      // Builds new workflow instance
      runWorkflowIntantiate();

      // 0. Init spec - done
      // 1. Compose steps - done
      runWorkflowStepsCompose();

      // 2. Load configuration file - done
      runWorkflowTestbedSettingsLoad();

      runWorkflowAssignTestBed();

      runWorkflowAssignConnectors();

      // 3. Execute prepare() for all steps - done
      runStepsSequencePrepare(false);


      // 4. Establish the connections as listed in the corresponding ServiceSpec
      // classes.
      runEstablishConnections();

//      getRegistry().claimTestbed(_context);

      // 6. Execute check() for all steps - done
      runStepsSequenceCheckHealth();
      // 7. Close connections

      // return true or throw?
   }

   public void disassemble() throws ProviderControllerException, ProviderWorkflowException {
      guardSingleRequest();

      // Builds new workflow instance
      runWorkflowIntantiate();

      // 0. Init spec - done
      // 1. Compose steps - done
      runWorkflowStepsCompose();
      // 2. Load configuration file - done
      runWorkflowTestbedSettingsLoad();

      runWorkflowAssignTestBed();

      runWorkflowAssignConnectors();

      // 3. Execute prepare() for all steps - done
      runStepsSequencePrepare(false);
      // 4. Establish the connections as listed in the corresponding ServiceSpec
      // classes.
      runEstablishConnections();
      // 6. Execute disassemble() for all steps in the reverse order - done
      runStepsSequenceDisassemble();
      // 7. Close connections
      // 8. Delete the content of the file?

      // return true or throw?
   }

   /**
    * Direct instantiation is not enabled. Use the factory method instead.
    */
   protected ProviderWorkflowController() {
   }

   // Workflow runners

   /**
    * Instantiates the workflow and related data and stores it in context.
    *
    * @throws ProviderControllerException
    *            The exception is triggered if a system error occurs.
    */
   private void runWorkflowIntantiate()
         throws ProviderControllerException, ProviderWorkflowException {

      try {
         _context.createWorkflowInstance();
      } catch (InstantiationException | IllegalAccessException e1) {
         throw new ProviderControllerException(
               "The controller cannot instantiate the provider workflow", e1);
      }

      ProviderWorkflow workflow = _context.getProviderWorkflow();

      try {
         workflow.initPublisherSpec(_context.getPublisherSpec());
         // TODO: rkovachev - we do not need to init assemble spec if not in
         // assemble command - sync to Nacho
         workflow.initAssemblerSpec(_context.getAssemblerSpec(), _providerBridge);
      } catch (Exception e) {
         throw new ProviderWorkflowException("Error initializing provider specs.", e);
      }
   }

   /**
    * Calls the workflow to build the list of its steps.
    */
   private void runWorkflowStepsCompose()
         throws ProviderWorkflowException {
      ProviderWorkflow workflow = _context.getProviderWorkflow();

      try {
         workflow.composeProviderSteps(_context.getFlow());
         // TODO: Add validation that the workflow has at least one step
      } catch (Exception e) {
         throw new ProviderWorkflowException("Error composing provider steps.", e);
      }
   }

   /**
    * Assign the test bed settings to the specs.
    */
   private void runWorkflowAssignTestBed() throws ProviderWorkflowException {
      try {
         _context.getProviderWorkflow().assignTestbedSettings(
               _context.getPublisherSpec(),
               _context.getTestbedSettings());

         _context.getProviderWorkflow().assignTestbedSettings(
               _context.getAssemblerSpec(),
               _context.getTestbedSettings());
      } catch (Exception e) {
         throw new ProviderWorkflowException(
               "Error assigning testbed settings the provider specs", e);
      }
   }

   private void runEstablishConnections() throws ProviderWorkflowException {
      Map<String, TestbedConnector> connectionMap = _context.getKeyConncetionMap();

      for (Map.Entry<String, TestbedConnector> mapEntry : connectionMap.entrySet()) {
         mapEntry.getValue().connect();
      }
   }

   // TODO: rkovachev Move to the Providers
   private void runWorkflowAssignServices(BaseSpec containerSpec)
         throws ProviderWorkflowException {
      //      AssemblerSpec assemblerSpec = _context.getAssemblerSpec();
      List<ManagedEntitySpec> specList =
            containerSpec.links.getAll(ManagedEntitySpec.class);

      // Elemental assembler workflow - no ManagedEntitySpec objects
      for (ManagedEntitySpec spec : specList) {
         assignParentService(spec);
      }
   }

   private void assignParentService(ManagedEntitySpec spec) {
      if (!spec.parent.isAssigned() || spec.parent.get() == null) {
         // Container spec- vc, host and etc.
         if (spec.service.get() == null) {
            throw new RuntimeException("Spec " + spec
                  + " is container and has no service assigned!");
         }
      } else {
         assignParentService(spec.parent.get());
         spec.service.set(spec.parent.get().service.get());
      }
   }

   private void runWorkflowAssignConnectors() throws ProviderWorkflowException {
      // TODO: Establish connections and store them in the provider context

      // 0. Collect all entity specs
      PublisherSpec publisherSpec = _context.getPublisherSpec();

      Set<EntitySpec> entitySpecs =
            SpecTraversalUtil.getAllSpecsFromContainerNode(
                  publisherSpec,
                  EntitySpec.class);
      // entitySpecs.toArray()
      if (CollectionUtils.isEmpty(entitySpecs)) {
         // There are no entity specs defined, so there is no need for
         // establishing connections.
         return;
      }

      // 0. Collect all unique service specs. One and the same specs are
      // considered those with the same key.
      // 1. Clone each unique serice spec and the map keys from them.

      Map<String, ServiceSpec> connectionKeyServiceSpecMap =
            new HashMap<String, ServiceSpec>();
      for (EntitySpec entitySpec : entitySpecs) {
         // TODO: rkovachev talk to Nacho about assigning services to not
         // elemental TB. And the fact that the elemental TB provide invokes the
         // runWorkflowAssignConnectors before the one that request them
         if (entitySpec.service == null || entitySpec.service.get() == null) {
            continue;
         }
         String connectionKey = entitySpec.service.toString();
         if (!connectionKeyServiceSpecMap.containsKey(connectionKey)) {
            ServiceSpec newServiceSpec = null;
            try {
               Class<? extends ServiceSpec> serviceSpecClass =
                     entitySpec.service.get().getClass();
               newServiceSpec = serviceSpecClass.newInstance();
            } catch (Exception e) {
               throw new ProviderWorkflowException("Cannot instantiate service spec",
                     null); // TODO: Polish
               // the message
            }
            newServiceSpec.copy(entitySpec.service.get());
            _logger.info("ADD SERVICE SPEC WITH KEY: " + connectionKey);
            connectionKeyServiceSpecMap.put(connectionKey, newServiceSpec);
         }
      }

      Map<ServiceSpec, TestbedConnector> serviceConnectorsMap =
            new HashMap<ServiceSpec, TestbedConnector>();
      for (String connectionKey : connectionKeyServiceSpecMap.keySet()) {
         serviceConnectorsMap.put(connectionKeyServiceSpecMap.get(connectionKey), null);
      }

      // 2. Ask the workflow to build the service-connection map.
      try {
         _context.getProviderWorkflow().assignTestbedConnectors(serviceConnectorsMap);
      } catch (Exception e) {
         throw new ProviderWorkflowException(
               "Error assigning connections for the provider", e);
      }

      // 3. Verify there're connectors for each service
      for (Entry<ServiceSpec, TestbedConnector> entry: serviceConnectorsMap.entrySet()) {
         TestbedConnector connector = entry.getValue();
         if (connector == null) {
            ServiceSpec serviceSpec = entry.getKey();
            ProviderWorkflowException exception = new ProviderWorkflowException(
                  "Not all connectors are set. Unset serviceSpec: "
                        + serviceSpec.toString(), null);
            _logger.info(exception.getMessage());
         }
      }

      // 4. Set the new connectors in the registry or make sure old ones are
      // reused? Why this is stored in the registry, should it be
      // better to stay in the provider context. It's the latter as there's no
      // value from the connection pooling in this situation.
      // It's gives a few seconds per provider (test respective, but maintaining
      // it will be quite expensive. It also makes sense to
      // build it in way that that provider contexts are recycled instead
      // maintaining separate connection pools.
      Map<String, TestbedConnector> connectionMap = _context.getKeyConncetionMap();
      connectionMap.clear();
      for (ServiceSpec serviceSpec : serviceConnectorsMap.keySet()) {
         String connectionKey = serviceSpec.toString();
         connectionMap.put(connectionKey, serviceConnectorsMap.get(serviceSpec));
      }
   }

   /**
    * Load the key-value configuration into the <code>Properties</code> set and
    * invoke the workflow's configuration loaded to dispatch them to specs.
    *
    * @throws ProviderWorkflowException
    */
   private void runWorkflowTestbedSettingsLoad() throws ProviderControllerException, ProviderWorkflowException {
      Properties testbedSettings = new Properties();

      try (FileReader fileReader = new FileReader(_testBedConfigFilePath);
            BufferedReader bufferedReader = new BufferedReader(fileReader)) {
         testbedSettings.load(bufferedReader);
      } catch (FileNotFoundException e) {
         throw new ProviderControllerException(String.format(
               "Testbed settings file not available: {0}",
               _testBedConfigFilePath), e);
      } catch (IOException e) {
         throw new ProviderControllerException(String.format(
               "Error reading testbed settings file: {0}",
               _testBedConfigFilePath), e);
      }

      _context.getTestbedSettings().putAll(testbedSettings);
   }

   /**
    * Save into file the key-value configuration available in the
    * <code>Properties</code> after calling work's configuration saved to dump
    * the settings from the specs.
    *
    * @throws ProviderWorkflowException
    */
   private void runWorkflowTestbedSettingsSave() throws ProviderControllerException, ProviderWorkflowException {
      Properties newSettings = new Properties();

      try {
         newSettings.put(SETTINGS_KEY_PROVIDER_ID, _context.getProviderWorkflow()
               .getClass().getCanonicalName());
         newSettings.putAll(_context.getTestbedSettings());
      } catch (Exception e) {
         throw new ProviderWorkflowException(
               "Error saving provider specs  configuration data into key-value testbed settings.",
               e);
      }

      try (FileWriter fileWriter = new FileWriter(_testBedConfigFilePath);
            BufferedWriter bufferedWriter = new BufferedWriter(fileWriter)) {
         newSettings.store(bufferedWriter, null);
      } catch (IOException e) {
         throw new ProviderControllerException(String.format(
               "Error writing testbed settings file: {0}",
               _testBedConfigFilePath), e);
      }
   }

   /**
    *
    * @param numberOfStepsToRemove
    * @throws ProviderWorkflowException
    */
   private void runWorkflowAssemblePhaseRollback(int numberOfNotStartedSteps)
         throws ProviderControllerException, ProviderWorkflowException {
      runWorkflowIntantiate();
      runWorkflowStepsCompose();
      if (numberOfNotStartedSteps > 0) {
         _context.getFlow().removeLastSteps(numberOfNotStartedSteps);
      }
      runWorkflowAssignTestBed();
      runWorkflowAssignConnectors();
      runStepsSequencePrepare(false);
      runEstablishConnections();
      runStepsSequenceDisassemble();
   }

   // Steps sequences runners

   /**
    * Run the prepare() methods of all tests in the workflow.
    *
    * Each method is given isolated copy of the workflow spec.
    *
    * @throws ProviderControllerException
    * @throws ProviderWorkflowException
    */
   private void runStepsSequencePrepare(boolean forAssembling)
         throws ProviderControllerException, ProviderWorkflowException {
      List<ProviderWorkflowStepContext> stepsContext = _context.getFlow().getAllSteps();

      for (ProviderWorkflowStepContext stepContext : stepsContext) {
         enableStepPhasePrepare(stepContext);
         runStepPhasePrepare(stepContext, forAssembling); // TODO: Handle
         // exception
      }
   }

   private void runStepsSequenceAssemble() throws ProviderWorkflowException, ProviderControllerException {
      List<ProviderWorkflowStepContext> stepsContext = _context.getFlow().getAllSteps();

      for (ProviderWorkflowStepContext stepContext : stepsContext) {
         enableStepPhaseAssemble(stepContext);
         enableStepPhaseDisassemble(stepContext);
         runStepPhaseAssemble(stepContext);
      }
   }

   private void runStepsSequenceCheckHealth() throws ProviderWorkflowException, ProviderControllerException {
      List<ProviderWorkflowStepContext> stepsContext = _context.getFlow().getAllSteps();

      for (ProviderWorkflowStepContext stepContext : stepsContext) {
         enableStepPhaseCheckHealth(stepContext);
         runStepPhaseCheckHealth(stepContext);
      }
   }

   private void runStepsSequenceDisassemble() throws ProviderWorkflowException, ProviderControllerException {
      List<ProviderWorkflowStepContext> straightSteps = _context.getFlow().getAllSteps();
      List<ProviderWorkflowStepContext> reversedSteps =
            new ArrayList<ProviderWorkflowStepContext>(straightSteps);

      Collections.reverse(reversedSteps);
      for (ProviderWorkflowStepContext stepContext : reversedSteps) {
         // Note: Important - keep on disassembling to the first one.
         enableStepPhaseDisassemble(stepContext);
         runStepPhaseDisassemble(stepContext);
      }
   }

   // Single step phase

   private void enableStepPhasePrepare(ProviderWorkflowStepContext stepContext)
         throws ProviderControllerException {
      if (stepContext.getPrepareState() == StepPhaseState.BLOCKED) {
         stepContext.setPrepareState(StepPhaseState.READY_TO_START);
      } else {
         throw new ProviderControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private void runStepPhasePrepare(ProviderWorkflowStepContext stepContext,
         boolean forAssembling) throws ProviderWorkflowException,
         ProviderControllerException {

      if (stepContext.getPrepareState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      PublisherSpec filteredPublisherSpec = null;
      AssemblerSpec filteredAssemblerSpec = null;
      try {
         filteredPublisherSpec =
               (PublisherSpec) getDeepClonedSpec(
                     _context.getPublisherSpec(),
                     stepContext.getSpecTagsFilter());
         filteredAssemblerSpec =
               (AssemblerSpec) getDeepClonedSpec(
                     _context.getAssemblerSpec(),
                     stepContext.getSpecTagsFilter());
      } catch (Exception e) {
         throw new ProviderControllerException("", e);
      }

      try {
         stepContext.setPrepareState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().prepare(
               filteredPublisherSpec,
               filteredAssemblerSpec,
               forAssembling,
               getRegistry().getSessionSettingsReader());
      } catch (Exception e) {
         stepContext.setPrepareState(StepPhaseState.FAILED);
         throw new ProviderWorkflowException("", e);
      }

      stepContext.setPrepareState(StepPhaseState.DONE);
   }

   private void enableStepPhaseAssemble(ProviderWorkflowStepContext stepContext)
         throws ProviderControllerException {
      if (stepContext.getAssembleState() == StepPhaseState.BLOCKED) {
         stepContext.setAssembleState(StepPhaseState.READY_TO_START);
      } else {
         throw new ProviderControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private void runStepPhaseAssemble(ProviderWorkflowStepContext stepContext)
         throws ProviderWorkflowException {

      if (stepContext.getAssembleState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      try {
         stepContext.setAssembleState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().assemble(stepContext.getSettingsWriter());
      } catch (Exception e) {
         stepContext.setAssembleState(StepPhaseState.FAILED);
         throw new ProviderWorkflowException("", e);
      }
      stepContext.setAssembleState(StepPhaseState.DONE);
   }

   private void enableStepPhaseCheckHealth(ProviderWorkflowStepContext stepContext)
         throws ProviderControllerException {
      if (stepContext.getCheckHealthState() == StepPhaseState.BLOCKED) {
         stepContext.setCheckHealthState(StepPhaseState.READY_TO_START);
      } else {
         throw new ProviderControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private void runStepPhaseCheckHealth(ProviderWorkflowStepContext stepContext)
         throws ProviderWorkflowException {

      if (stepContext.getCheckHealthState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      try {
         stepContext.setCheckHealthState(StepPhaseState.IN_PROGRESS);
         // TODO: rkovachev - how do we handle the failure of a step?
         if (stepContext.getStep().checkHealth()) {
            stepContext.setCheckHealthState(StepPhaseState.DONE);
         } else {
            stepContext.setCheckHealthState(StepPhaseState.FAILED);
         }
      } catch (Exception e) {
         stepContext.setCheckHealthState(StepPhaseState.FAILED);
         throw new ProviderWorkflowException("", e);
      }

   }

   private void enableStepPhaseDisassemble(ProviderWorkflowStepContext stepContext)
         throws ProviderControllerException {
      if (stepContext.getDisassembleState() == StepPhaseState.BLOCKED) {
         stepContext.setDisassembleState(StepPhaseState.READY_TO_START);
      } else {
         throw new ProviderControllerException(
               "The controller attempted invalid enabling of step phase.");
      }
   }

   private void runStepPhaseDisassemble(ProviderWorkflowStepContext stepContext)
         throws ProviderWorkflowException {
      if (stepContext.getDisassembleState() != StepPhaseState.READY_TO_START) {
         // This check prevents the phase double run or run before other phases
         // on which it depends are completed successfully.
         return;
      }

      try {
         stepContext.setDisassembleState(StepPhaseState.IN_PROGRESS);
         stepContext.getStep().disassemble();
      } catch (Exception e) {
         stepContext.setDisassembleState(StepPhaseState.FAILED);
         throw new ProviderWorkflowException("", e);
      }
      stepContext.setDisassembleState(StepPhaseState.DONE);
   }

   // Helper methods

   /**
    * Stores in the workflow context test bed settings specified during
    * assembling phase of each step.
    */
   private void collectNewSettings() {
      List<ProviderWorkflowStepContext> stepsContext = _context.getFlow().getAllSteps();

      for (ProviderWorkflowStepContext stepContext : stepsContext) {
         if (stepContext.getAssembleState() == StepPhaseState.DONE
               || stepContext.getAssembleState() == StepPhaseState.FAILED) {
            _context.getTestbedSettings().putAll(stepContext.getSettings());
         }
      }
   }

   /**
    * Returns the number of steps on which the assembling was not performed.
    *
    * @return Number of steps on which the assembling was not started.
    */
   private int getNotExecutedAssemblingStepsCount() {
      List<ProviderWorkflowStepContext> stepsContext = _context.getFlow().getAllSteps();
      int count = 0;
      for (ProviderWorkflowStepContext stepContext : stepsContext) {
         if (stepContext.getAssembleState() == StepPhaseState.BLOCKED) {
            // TODO: Review the step phase states.
            count++;
         }
      }
      return count;
   }
}
