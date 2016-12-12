/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.workflow.command;

import java.io.BufferedReader;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.text.MessageFormat;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.workflow.explorer.WorkflowRegistry;
import com.vmware.hsua.common.util.IOUtils;

/**
 * The <code>BaseCommand</code> implements common utilities and validation
 * methods required by most of the commands.
 *
 * The class implements <code>WorkflowCommand</code>, which is required for
 * all classes that must be recognized as command classes.
 *
 * Concrete command implementations could extend this class for convenience,
 * but this is not required by the system.
 */
public abstract class BaseCommand implements WorkflowCommand {

   private static final Logger _logger = LoggerFactory
         .getLogger(BaseCommand.class);

   @Override
   public final String getCommandName() {
      WorkflowCommandAnnotation annotation = this.getClass().getAnnotation(
            WorkflowCommandAnnotation.class);
      if (annotation == null) {
         throw new RuntimeException(MessageFormat.format(
               "Command '{0}' isn't annotated with WorkflowCommandAnnotation",
               this.getClass().getCanonicalName()));
      }

      return annotation.commandName();
   }

   /**
    * Provides reference to the global registry.
    */
   protected WorkflowRegistry getRegistry() {
      return WorkflowRegistry.getRegistry();
   }

   /**
    * Perform basic validation of the input parameters.
    *
    * The number of command parameters i
    * An exception is thrown if the validation doesn't check.
    *
    * Check if the numbers are rig
    * @param params
    *    String parameters array.
    *
    * @param paramTitles
    *    Holds the title of each parameter as ordered in the <code>params</code> array.
    *    The length of this array is used to infer the expected number of command
    *    parameters.
    */
   protected void validateParameters(
         String[] commandParams,
         String[] paramTitles) {

      if (commandParams == null || commandParams.length == 0) {
         throw new IllegalArgumentException(
               "Illegal call to validateParameters. The commandParams array must be filled in.");
      }

      if (paramTitles == null || paramTitles.length == 0) {
         throw new IllegalArgumentException(
               "Illegal call to validateParameters. The paramTitles array must be filled in.");
      }


      if (commandParams.length != paramTitles.length) {
         throw new IllegalArgumentException(
               String.format(
                     "(0) command parameters are provided for command %s, but %s are expected. Provide the required number of command parameters.",
                     commandParams.length,
                     this.getCommandName(),
                     paramTitles.length));
      }


      for (int i=0;i<commandParams.length;++i) {
         if (Strings.isNullOrEmpty(commandParams[i])) {

            throw new IllegalArgumentException(
                  String.format(
                        "{0} command parameter is empty. Set the parameter.",
                        paramTitles[i]));
         }
      }
   }

   /**
    * Validate that a file path exists and the file can be opened.
    * If this behavior doesn't check, an exception is thrown.
    *
    * TODO: Consider if it's good to structurally validate the file here - e.g. KV, CSV
    *
    * @param filePath
    *    Path to the file
    */
   protected void validateInputFile(String filePath) {
      try {
         IOUtils.readLinesFromFile(filePath);
      } catch (IOException e) {
         _logger.error("Error during reading: " + filePath);
         throw new RuntimeException(e);
      }
   }

   /**
    * Create a new empty file and close it.
    *
    * The method is used to validate that the specified file path is accessible.
    * If this behavior doesn't check, an exception is thrown.
    *
    * @param filePath
    *    Path to the file
    */
   protected void createAndValidateOutputFile(String filePath) {
      List<String> emptyLines = new ArrayList<String>();
      IOUtils.writeLinesToFile(emptyLines, filePath);
   }

   /**
    * Validate that actual workflow class for the given name and workflow type exists.
    *
    * Throw an exception if it could not be found.
    *
    * @param workflowClassName
    *      Name of the workflow class.
    * @param workflowType
    *      Super class for the given workflow class. If null, the check is omitted.s
    */
//   protected void validateWorkflowExists(
//         String workflowClassName, Class<? extends Workflow> workflowType) {
//
//      if (workflowClassName == null) {
//         throw new IllegalArgumentException(
//               "Illegal call to validateParameters. The workflowClassName should be set.");
//      }
//
//      if (workflowType == null) {
//         throw new IllegalArgumentException(
//               "Illegal call to validateParameters. The workflowType should be set.");
//      }
//
//      if (!getRegistry().checkWorkflowExists(workflowClassName, workflowType)) {
//         throw new IllegalArgumentException(
//               String.format(
//                     "The workflow class {0} cannot be found or is not of type {1}.",
//                     workflowClassName,
//                     workflowType.getCanonicalName()));
//      }
//   }

   /**
    * Load the key-value settings from a file into <code>Properties</code>.
    * @param filePath
    *    File path
    *
    * @return
    *    <code>Properties</code>
    */
   protected Properties loadKVSettings(String filePath) throws CommandException {
      Properties properties = new Properties();

      // TODO: Consider moving the generic implementation in the IO util.

      try (FileReader fileReader = new FileReader(filePath);
            BufferedReader bufferedReader = new BufferedReader(fileReader) ) {
         properties.load(bufferedReader);
      } catch (FileNotFoundException e) {
         throw new CommandException(filePath + " was not found!", e);
      } catch (IOException e) {
         throw new CommandException("", e);
      }

      return properties;
   }
}
