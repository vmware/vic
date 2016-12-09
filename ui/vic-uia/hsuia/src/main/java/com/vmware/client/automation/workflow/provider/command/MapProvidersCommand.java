/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.provider.command;

import java.util.ArrayList;
import java.util.List;
import java.util.Set;

import org.apache.commons.collections4.CollectionUtils;

import com.vmware.client.automation.workflow.command.CommandException;
import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.client.automation.workflow.provider.ProviderWorkflow;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowContext;
import com.vmware.client.automation.workflow.provider.ProviderWorkflowController;
import com.vmware.hsua.common.util.KeyValuePairFileBuilder;

// Sat: provider map, provider list, test list
// Sun: Code fit & finish
// Mon-Wed: Wiring with the provisioner
// ? Downloading?

@WorkflowCommandAnnotation(commandName = "provider-map")
public class MapProvidersCommand extends BaseProviderCommand {

   private String _providerClassName;
   private String _providerMapFilePath;

   @Override
   public void prepare(String[] commandParams) {
      validateParameters(
            commandParams,
            new String[] {"providerClassName", "providerMapFilePath"});

      _providerClassName = commandParams[0];
      _providerMapFilePath = commandParams[1];
   }

   @Override
   public void execute() throws CommandException {
      // Create context for all workflows. This will be the map in the registry.
      // Get the provider class
      Class<? extends ProviderWorkflow> providerClass = getProviderClass(_providerClassName);

      // Run custom validation

      WorkflowRegistry registry = getRegistry();
      ProviderWorkflowContext context = null;
      ProviderWorkflowController controller = null;
      try {
         //   context = registry.registerWorkflowInstance(providerClassName);
         context = registry.registrerWorkflowContext(providerClass);
         controller = ProviderWorkflowController.create(context, null);
         controller.analyze();

         boolean success = createOutputFile(registry, context);
         if (success == false) {
            throw new Exception("Error in writing output to file!");
         }
      } catch (Exception e) {
         throw new CommandException(
               String.format("Execution failed on command %s %s", getCommandName(), e.getMessage()), e);
         /* } catch (IOException e) {
       // TODO: add decent message on the exception. Need to decide where the issues will be triaged for the reporting.
       throw new Exception(
       String.format(
             "Error writing testbed settings file: {0}", providerMapFilePath),
             e);

          */
      } finally {
         // Make sure the context is always properly removed.
         registry.unregistrerWorkflowContext(context);
      }
      // Make sure prepare() phase is launched in all steps, so the test bed types could be generated. Types are required only for primary providers
      // Read the registry and build the map.
      // Dispose everything

      // initProviderSpec cannot consume related specs
      // Unique type is set for the provider
      // There's download phase or allocate physical resource
   }

   private boolean createOutputFile(WorkflowRegistry registry,
         ProviderWorkflowContext rootProviderContext) {
      Set<ProviderWorkflowContext> componentContexts = registry
            .getRegisteredConsumerProviders(rootProviderContext);
      KeyValuePairFileBuilder fileBuilder = new KeyValuePairFileBuilder(_providerMapFilePath);
      fileBuilder.addComment("provider-map command result file");
      if (CollectionUtils.isEmpty(componentContexts)) {
         // Elemental provider
         fileBuilder.addArrayKeyValuePair(
               KeyValuePairFileBuilder.DEPENDENCIES_KEY,
               new String[] { rootProviderContext.getProviderWorkflow().getClass()
                     .getSimpleName() });
         addProviderToOutput(fileBuilder, rootProviderContext.getProviderWorkflow());
      } else {
         List<String> dependencies = new ArrayList<String>();
         for (ProviderWorkflowContext providerWorkflowContext : componentContexts) {
            dependencies.add(providerWorkflowContext.getProviderWorkflow().getClass()
                  .getSimpleName());
         }
         fileBuilder.addArrayKeyValuePair(
               KeyValuePairFileBuilder.DEPENDENCIES_KEY,
               dependencies.toArray(new String[dependencies.size()]));
         for (ProviderWorkflowContext providerWorkflowContext : componentContexts) {
            addProviderToOutput(fileBuilder, providerWorkflowContext.getProviderWorkflow());
         }
      }
      return fileBuilder.build();
   }

   private void addProviderToOutput(KeyValuePairFileBuilder fileBuilder, ProviderWorkflow provider) {
      fileBuilder.addSimpleKeyValuePair(
            provider.getClass().getSimpleName() + ".testClass",
            provider.getClass().getCanonicalName());
      fileBuilder.addSimpleKeyValuePair(
            provider.getClass().getSimpleName()+ ".testWeight",
            provider.providerWeight() + "");
   }

   @SuppressWarnings("unchecked")
   private Class<? extends ProviderWorkflow> getProviderClass(
         String providerClassName) throws CommandException {
      try {
         return
               (Class<? extends ProviderWorkflow>) getRegistry().getRegisteredWorkflowClass(
                     providerClassName, ProviderWorkflow.class);
      } catch (ClassNotFoundException e) {
         throw new CommandException(
               String.format(
                     "Cannot find provider class: %s",
                     providerClassName),
                     e);
      }
   }
}
