/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider;

import java.util.Properties;

import com.vmware.client.automation.workflow.common.StepPhaseState;
import com.vmware.client.automation.workflow.common.WorkflowStepContext;
import com.vmware.client.automation.workflow.explorer.SettingsRWImpl;
import com.vmware.client.automation.workflow.explorer.SettingsWriter;

public class ProviderWorkflowStepContext extends WorkflowStepContext {

   private ProviderWorkflowStep _providerStep;
   private SettingsRWImpl _settings;

   private StepPhaseState _prepareState;
   private StepPhaseState _assembleState;
   private StepPhaseState _checkHealthState;
   private StepPhaseState _disassembleState;

   public ProviderWorkflowStepContext(ProviderWorkflowStep providerStep) {
      _providerStep = providerStep;

      _prepareState = StepPhaseState.BLOCKED;
      _assembleState = StepPhaseState.BLOCKED;
      _checkHealthState = StepPhaseState.BLOCKED;
      _disassembleState = StepPhaseState.BLOCKED;

      _settings = new SettingsRWImpl();
   }

   public ProviderWorkflowStep getStep() {
      return _providerStep;
   }

   // Configuration data

   public SettingsWriter getSettingsWriter() {
      return _settings;
   }

   public Properties getSettings() {
      return _settings;
   }

   // Runtime data

   public StepPhaseState getPrepareState() {
      return _prepareState;
   }

   public void setPrepareState(StepPhaseState state) {
      _prepareState = state;
   }

   public StepPhaseState getAssembleState() {
      return _assembleState;
   }

   public void setAssembleState(StepPhaseState state) {
      _assembleState = state;
   }

   public StepPhaseState getCheckHealthState() {
      return _checkHealthState;
   }

   public void setCheckHealthState(StepPhaseState state) {
      _checkHealthState = state;
   }

   public StepPhaseState getDisassembleState() {
      return _disassembleState;
   }

   public void setDisassembleState(StepPhaseState state) {
      _disassembleState = state;
   }

}
