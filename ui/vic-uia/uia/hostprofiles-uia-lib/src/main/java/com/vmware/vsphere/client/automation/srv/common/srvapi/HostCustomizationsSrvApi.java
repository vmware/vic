/*
 *  Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.vsphere.client.automation.srv.common.srvapi;

import com.vmware.client.automation.common.TestSpecValidator;
import com.vmware.client.automation.common.spec.ServiceSpec;
import com.vmware.client.automation.exception.VcException;
import com.vmware.client.automation.util.VcServiceUtil;
import com.vmware.vim.binding.vim.HostSystem;
import com.vmware.vim.binding.vim.ServiceInstanceContent;
import com.vmware.vim.binding.vim.Task;
import com.vmware.vim.binding.vim.TaskInfo;
import com.vmware.vim.binding.vim.profile.DeferredPolicyOptionParameter;
import com.vmware.vim.binding.vim.profile.host.AnswerFile;
import com.vmware.vim.binding.vim.profile.host.ProfileManager;
import com.vmware.vim.binding.vim.profile.host.ProfileManager.ApplyHostConfigSpec;
import com.vmware.vim.binding.vim.profile.host.ProfileManager.CustomizationData;
import com.vmware.vim.binding.vim.profile.host.ProfileManager.EntityCustomizations;
import com.vmware.vim.binding.vim.profile.host.ProfileManager.FormattedCustomizations;
import com.vmware.vim.binding.vmodl.ManagedObjectReference;
import com.vmware.vise.search.util.Strings;
import com.vmware.vise.vim.commons.vcservice.VcService;
import com.vmware.vsphere.client.automation.common.CommonHostProfileConstants;
import com.vmware.vsphere.client.automation.common.FileUtils;
import com.vmware.vsphere.client.automation.common.HostProfilesUtil;
import com.vmware.vsphere.client.automation.common.spec.AnswerFileUpdateSpec;
import com.vmware.vsphere.client.automation.common.spec.HostPolicyUpdateSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostProfileSpec;
import com.vmware.vsphere.client.automation.srv.common.spec.HostSpec;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.List;

import static com.vmware.vsphere.client.automation.srv.common.srvapi.ManagedEntityUtil.getManagedObjectFromMoRef;

/**
 * API server wrapper for host profile customization related calls.
 */
public class HostCustomizationsSrvApi {

   private static final Logger log = LoggerFactory
         .getLogger(HostProfileSrvApi.class);
   private static final int HEADER_LINE_INDEX = 0;
   private static final String QUOTE = "\"";

   private static HostCustomizationsSrvApi INSTANCE = null;
   private static final String NEW_LINE_DELIMITER = "\n";
   private static final String ADVANCED_CONFIG_PREFIX =
         "Advanced configuration option ";

   // Disallow the direct creation of the class. use getInstance() instead
   private HostCustomizationsSrvApi() {
   }

   /**
    * Get instance of HostCustomizationsSrvApi.
    *
    * @return created instance
    */
   public static HostCustomizationsSrvApi getInstance() {
      if (INSTANCE == null) {
         synchronized (HostCustomizationsSrvApi.class) {
            if (INSTANCE == null) {
               log.info("Initializing HostCustomizationsSrvApi.");
               INSTANCE = new HostCustomizationsSrvApi();
            }
         }
      }

      return INSTANCE;
   }

   /**
    * Modify csv answer file values
    *
    * @param filename        the csv file name
    * @param fileUpdateSpecs the update specs used to modify the values
    */
   public void modifyAnswerFileValues(String filename,
                                      List<AnswerFileUpdateSpec> fileUpdateSpecs) {
      Path path = Paths.get(filename);

      if (fileUpdateSpecs == null || fileUpdateSpecs.isEmpty()) {
         writeAnswerFile(path, "");
      } else {
         modifyAnswerFileValues(path, fileUpdateSpecs);
      }
   }

   /**
    * Modify csv answer file values
    *
    * @param filePath        the csv file path
    * @param fileUpdateSpecs the update specs used to modify the values
    */
   public void modifyAnswerFileValues(Path filePath,
                                      List<AnswerFileUpdateSpec> fileUpdateSpecs) {
      List<String> inLines = readAnswerFile(filePath);

      String resultContent;
      resultContent = updateAnswerFileValues(fileUpdateSpecs, inLines);

      writeAnswerFile(filePath, resultContent);
   }

   /**
    * Exports host customizations to a csv text file with the specified
    * filename. The filename should NOT include the file extension, as it is
    * set to csv.
    * <p/>
    * The file is placed in the project folder
    *
    * @param hostProfileSpec host profile spec to export customizations from
    * @param fileName        the name of the file to be exported
    */
   public void exportHostCustomizationsToFile(HostProfileSpec hostProfileSpec,
                                              String fileName) {
      String answerFileString = exportHostCustomizations(hostProfileSpec);

      FileUtils fileUtils = FileUtils.getInstance();
      fileUtils.writeStringToTextFile(answerFileString, fileName);
   }

   /**
    * Retrieve host customizations and return the result as a String. The
    * String represents the exported csv text file.
    *
    * @param hostProfileSpec host profile spec to export customizations from
    * @return String representation of the exported csv file
    */
   public String exportHostCustomizations(HostProfileSpec hostProfileSpec) {
      ManagedObjectReference taskReference;
      taskReference = startExportAndGetTaskMor(hostProfileSpec);

      waitForTaskSuccess(hostProfileSpec, taskReference);

      String taskResult = getTaskResult(hostProfileSpec, taskReference);
      return taskResult;
   }

   /**
    * Method that remediate the supplied hosts
    *
    * @param hostSpecs - hosts to remediate
    * @return true if successful, false otherwise
    * @throws Exception - if there is error in VC communication or actions
    */
   public boolean remediate(HostSpec[] hostSpecs) throws Exception {
      TestSpecValidator.ensureNotEmpty(Arrays.asList(hostSpecs),
                                       "Supply entity to check compliance!");

      ServiceSpec svcSpec = hostSpecs[0].service.get();
      ProfileManager profileManager = getHostProfileManager(svcSpec);
      ApplyHostConfigSpec[] specs = getApplyConfigSpecs(hostSpecs);
      // When host is remediated it may require reboot
      // Set the reboot flag so host to be rebooted during remediate.
      // The reboot flag is ignored if the host doesn't require reboot.
      for (ApplyHostConfigSpec spec : specs) {
         spec.setRebootHost(true);
      }
      ManagedObjectReference
            taskMor =
            profileManager.applyEntitiesConfiguration(specs);
      return VcServiceUtil.waitForTaskSuccess(taskMor, hostSpecs[0]);
   }

   /**
    * Read a file and return a list of text lines.
    *
    * @param fileName The filename
    * @return a List of text lines
    */
   public List<String> readAnswerFile(String fileName) {
      Path path = Paths.get(fileName);
      return readAnswerFile(path);
   }

   /**
    * Method used to reset the host customizations to their initial state
    *
    * @param hostSpec - host, whose customizations to be reset
    * @return - true if successful, false otherwise
    * @throws Exception - if vc connection exception happens, or host is not
    *                   present or not connected to a host profile
    */
   public boolean resetHostCustomizations(HostSpec hostSpec) throws Exception {
      HostSystem host = HostBasicSrvApi.getInstance().getHostSystem(hostSpec);
      ServiceSpec serviceSpec = hostSpec.service.get();
      ProfileManager profileManager = getHostProfileManager(serviceSpec);

      ProfileManager.AnswerFileOptionsCreateSpec resetSpec =
         new ProfileManager.AnswerFileOptionsCreateSpec();
      resetSpec.userInput = new DeferredPolicyOptionParameter[0];

      ManagedObjectReference taskMor =
         profileManager.updateAnswerFile(host._getRef(), resetSpec);
      return VcServiceUtil.waitForTaskSuccess(taskMor, serviceSpec);
   }

   /**
    * Method used just to accept the customizations with default values so that
    * host becomes not requiring host customizations for teh pruposes of
    * testing
    *
    * @param hostSpec - host whose customizations to set
    * @return - true if successful, false otherwise
    * @throws Exception - if vc connection exception happens, or host is not
    *                   present or not connected to a host profile
    */
   public boolean setHostCustomizations(HostSpec hostSpec) throws Exception {
      HostSystem host = HostBasicSrvApi.getInstance().getHostSystem(hostSpec);
      ServiceSpec serviceSpec = hostSpec.service.get();
      ProfileManager profileManager = getHostProfileManager(serviceSpec);
      AnswerFile answerFile = profileManager.retrieveAnswerFile(host._getRef());
      HostProfilesUtil.ensureNotNull(answerFile,
         "There is no answer file for the host: " + hostSpec.name.get());
      ProfileManager.AnswerFileOptionsCreateSpec createSpec =
         new ProfileManager.AnswerFileOptionsCreateSpec();
      createSpec.userInput = answerFile.getUserInput();

      ManagedObjectReference taskMor =
         profileManager.updateAnswerFile(host._getRef(), createSpec);
      return VcServiceUtil.waitForTaskSuccess(taskMor, serviceSpec);
   }

   // Private methods

   /**
    * Method that retrieves the config spec necessary fo remediation
    *
    * @param hostSpecs array of host specs
    * @return array of apply host config specs
    * @throws Exception
    */
   private ProfileManager.ApplyHostConfigSpec[] getApplyConfigSpecs(
         HostSpec[] hostSpecs) throws Exception {
      TestSpecValidator.ensureNotEmpty(Arrays.asList(hostSpecs),
                                       "Supply at least one HostSpec!");
      ServiceSpec serviceSpec = hostSpecs[0].service.get();
      VcService vcService = VcServiceUtil.getVcService(serviceSpec);
      ProfileManager profileManager = getHostProfileManager(serviceSpec);
      ManagedObjectReference preCheckRemediateTaskMor = profileManager
            .generateHostConfigTaskSpec(getHostCustSpecs(hostSpecs));

      VcServiceUtil.waitForTaskCompletion(preCheckRemediateTaskMor,
                                          serviceSpec);

      Task preCheckRemediateTask =
            vcService.getManagedObject(preCheckRemediateTaskMor);

      TaskInfo taskInfo = preCheckRemediateTask.getInfo();

      if (taskInfo.getError() != null) {
         throw new RuntimeException(taskInfo.getError());
      }

      return (ProfileManager.ApplyHostConfigSpec[]) taskInfo.getResult();
   }

   /**
    * Method that retrieves the host customization specs necessary for
    * remediation
    *
    * @param hostSpecs array of host specs
    * @return array of structured customizations
    * @throws Exception
    */
   private ProfileManager.StructuredCustomizations[] getHostCustSpecs(
         HostSpec[] hostSpecs) throws Exception {
      HostSystem[] hosts = new HostSystem[hostSpecs.length];
      for (int j = 0; j < hostSpecs.length; j++) {
         HostBasicSrvApi api = HostBasicSrvApi.getInstance();
         hosts[j] = api.getHostSystem(hostSpecs[j]);
      }
      ProfileManager.StructuredCustomizations[] hostCustSpecs;
      hostCustSpecs =
            new ProfileManager.StructuredCustomizations[hosts.length];

      for (int i = 0; i < hosts.length; i++) {
         ProfileManager.StructuredCustomizations hostCustSpec;
         hostCustSpec = new ProfileManager.StructuredCustomizations();
         hostCustSpec.setEntity(hosts[i]._getRef());
         hostCustSpecs[i] = hostCustSpec;
      }

      return hostCustSpecs;
   }

   private List<String> readAnswerFile(Path path) {
      List<String> inLines;
      try {
         inLines = Files.readAllLines(path, StandardCharsets.UTF_8);
      } catch (IOException e) {
         String format = "Error reading file %s";
         String message = String.format(format, path.getFileName());
         throw new RuntimeException(message, e);
      }
      return inLines;
   }

   private String updateAnswerFileValues(List<AnswerFileUpdateSpec> updateSpecs,
                                         List<String> fileLines) {
      for (AnswerFileUpdateSpec updateSpec : updateSpecs) {
         updateValuesForHost(updateSpec, fileLines);
      }

      return getLinuxFormattedTextContent(fileLines);
   }

   private void writeAnswerFile(Path path, String textLines) {
      try {
         Files.write(path, textLines.getBytes());
      } catch (IOException e) {
         String format = "Error writing file '%s'";
         String message = String.format(format, path.getFileName());
         throw new RuntimeException(message, e);
      }
   }

   /**
    * Updates the values for a specific host
    *
    * @param updateSpec the answer file update spec
    * @param fileLines  the file lines
    */
   private void updateValuesForHost(AnswerFileUpdateSpec updateSpec,
                                    List<String> fileLines) {
      final int hostRowIndex = getHostRowIndex(updateSpec, fileLines);
      String[] hostLineValues = getValuesByLineIndex(hostRowIndex, fileLines);
      String[] headerLine = getHeaderValues(fileLines);

      List<HostPolicyUpdateSpec> values;
      values = updateSpec.policyUpdateSpec.getAll();

      findAndUpdatePolicyValues(values, headerLine, hostLineValues);

      fileLines.set(hostRowIndex, joinValuesToCsvLine(hostLineValues));
   }

   private String[] getHeaderValues(List<String> fileLines) {
      String[] headerValues = getValuesByLineIndex(HEADER_LINE_INDEX,
                                                   fileLines);
      return HostProfilesUtil.trimLeadingTrailingChars(headerValues, QUOTE);
   }

   /**
    * Returns the index of host value line.
    *
    * @param updateSpec the answer file update spec
    * @param fileLines  the file lines
    * @return index of the host value line, -1 if not found
    */
   private int getHostRowIndex(AnswerFileUpdateSpec updateSpec,
                               List<String> fileLines) {
      final HostSpec host = updateSpec.host.get();
      return getHostRowIndex(host, fileLines);
   }

   /**
    * Returns an array of values for a line with index lineIndex
    *
    * @param lineIndex line index
    * @param lines     a list of lines
    * @return array of values
    */
   private String[] getValuesByLineIndex(int lineIndex,
                                         List<String> lines) {
      String line = lines.get(lineIndex);
      return line.split(CommonHostProfileConstants.CSV_DELIMITER, -1);
   }

   /**
    * Returns the index of host value line.
    *
    * @param host      the host spec
    * @param fileLines the file lines
    * @return index of the host value line, -1 if not found
    */
   public int getHostRowIndex(HostSpec host, List<String> fileLines) {
      final String hostIp = host.name.get();
      int result = -1;
      for (int i = 0; i < fileLines.size(); i++) {
         final String[] columns = getValuesByLineIndex(i, fileLines);
         //the host ip value is the first column
         final String firstColumnValue = columns[0];
         if (hostIp.equals(firstColumnValue)) {
            result = i;
            break;
         }
      }
      return result;
   }

   private void findAndUpdatePolicyValues(List<HostPolicyUpdateSpec> hostProfileUpdateSpecs,
                                          String[] headerLine,
                                          String[] valuesLine) {
      for (HostPolicyUpdateSpec hostProfileUpdateSpec : hostProfileUpdateSpecs) {
         final String columnName = getFullColumnName(hostProfileUpdateSpec);
         final int columnIndex = findColumnIndex(headerLine, columnName);

         String newValue = hostProfileUpdateSpec.newPropertyValue.get();
         valuesLine[columnIndex] = QUOTE + newValue + QUOTE;
      }
   }

   private String joinValuesToCsvLine(String[] policyValues) {
      String updatedValuesLine;
      updatedValuesLine =
            Strings.join(CommonHostProfileConstants.CSV_DELIMITER,
                         policyValues);
      return updatedValuesLine;
   }

   public String getLinuxFormattedTextContent(List<String> inLines) {
      String[] newEmptyArray = new String[inLines.size()];
      String[] outLinesArray = inLines.toArray(newEmptyArray);
      String outLines;
      outLines = Strings.join(NEW_LINE_DELIMITER, outLinesArray);
      return outLines;
   }

   private String getFullColumnName(HostPolicyUpdateSpec hostProfileUpdateSpec) {
      String policyName = hostProfileUpdateSpec.newPropertyName.get();
      return ADVANCED_CONFIG_PREFIX + policyName;
   }

   private int findColumnIndex(String[] headerLine, String columnName) {
      int columnIndex = searchForColumnIndex(columnName, headerLine);
      verifyColumnFound(columnName, columnIndex);

      return columnIndex;
   }

   private int searchForColumnIndex(String columnName, String[] headerLine) {
      int columnIndex = -1;

      for (int i = 0; i < headerLine.length; i++) {
         if (headerLine[i].equals(columnName)) {
            columnIndex = i;
            break;
         }
      }

      return columnIndex;
   }

   private void verifyColumnFound(String columnName, int columnIndex) {
      if (columnIndex == -1) {
         String message;
         message = String.format("Could not find column '%s'", columnName);
         throw new RuntimeException(message);
      }
   }

   /**
    * Starts the export API call and returns the newly started task reference.
    *
    * @param hostProfileSpec host profile spec to export customizations from
    * @return task reference
    */
   private ManagedObjectReference startExportAndGetTaskMor(HostProfileSpec hostProfileSpec) {
      ManagedObjectReference[] attachedHosts;
      attachedHosts = getHosts(hostProfileSpec);

      ProfileManager profileManager = getHostProfileManager(hostProfileSpec);

      ManagedObjectReference taskRef;
      taskRef = profileManager.exportCustomizations(attachedHosts, "csv");
      return taskRef;
   }

   /**
    * Wait for the task to complete. If the task fails to complete, or is
    * unsuccessful, a RuntimeException is thrown.
    *
    * @param entitySpec the target entity spec associated with the task
    * @param taskRef    the task reference
    */
   private void waitForTaskSuccess(HostProfileSpec entitySpec,
                                   ManagedObjectReference taskRef) {
      boolean taskCompletedSuccessfully;
      try {
         taskCompletedSuccessfully =
               VcServiceUtil.waitForTaskSuccess(taskRef, entitySpec);
      } catch (VcException e) {
         throw new RuntimeException("Task failed.", e);
      }

      if (!taskCompletedSuccessfully) {
         throw new RuntimeException("Tasks did not succeed.");
      }
   }

   /**
    * Retrieves the task result as String from a task reference and the
    * associated host profile spec
    *
    * @param hostProfileSpec host profile spec
    * @param taskRef         task reference
    * @return String representation of the task result
    */
   private String getTaskResult(HostProfileSpec hostProfileSpec,
                                ManagedObjectReference taskRef) {
      final Task task = getTask(hostProfileSpec, taskRef);

      return getHostCustomizationTaskResult(task);
   }

   /**
    * Retrieves the array of hosts attached to the host profile
    *
    * @param hostProfileSpec host profile spec
    * @return attachedHosts
    */
   private ManagedObjectReference[] getHosts(HostProfileSpec hostProfileSpec) {
      HostProfileSrvApi hostProfileApi = HostProfileSrvApi.getInstance();
      ManagedObjectReference[] attachedHosts;
      attachedHosts = hostProfileApi.getAttachedHosts(hostProfileSpec);

      return attachedHosts;
   }

   private ProfileManager getHostProfileManager(
         HostProfileSpec hostProfileSpec) {
      ProfileManager manager;
      try {
         ServiceSpec serviceSpec = hostProfileSpec.service.get();
         manager = getHostProfileManager(serviceSpec);
      } catch (VcException e) {
         throw new RuntimeException("Could not get host profile manager", e);
      }
      return manager;
   }

   /**
    * Retrieve the Task from a HostProfileSpec and a task reference.
    *
    * @param hostProfileSpec host profile spec
    * @param taskRef         task reference
    * @return task
    */
   private Task getTask(HostProfileSpec hostProfileSpec,
                        ManagedObjectReference taskRef) {
      VcService vcService = getVcService(hostProfileSpec);
      Task task = vcService.getManagedObject(taskRef);
      return task;
   }

   /**
    * Retrieve the task result as a String
    *
    * @param task the task to retrieve result from
    * @return a String representation of the task result
    */
   private String getHostCustomizationTaskResult(Task task) {
      TaskInfo info = task.getInfo();

      CustomizationData customizationData = getHostCustomizationData(info);

      EntityCustomizations entityCustomization;
      entityCustomization = customizationData.entityCustomizations[0];
      FormattedCustomizations formattedCustomizations;
      formattedCustomizations = (FormattedCustomizations) entityCustomization;

      String answerFileString;
      answerFileString = formattedCustomizations.formattedCustomizations;
      return answerFileString;
   }

   /**
    * Retrieve HostProfileManager from a ServiceSpec
    *
    * @param serviceSpec service Spec
    * @return profileManager
    * @throws VcException
    */
   private ProfileManager getHostProfileManager(ServiceSpec serviceSpec)
         throws VcException {
      ServiceInstanceContent serviceContent;
      serviceContent = getServiceInstanceContent(serviceSpec);

      ManagedObjectReference hostProfileManager;
      hostProfileManager = serviceContent.getHostProfileManager();

      ProfileManager profileManager;
      profileManager = getManagedObjectFromMoRef(hostProfileManager,
                                                 serviceSpec);

      return profileManager;
   }

   /**
    * Retrieve VcService from a HostProfileSpec
    *
    * @param hostProfileSpec host profile spec
    * @return vcService
    */
   private VcService getVcService(HostProfileSpec hostProfileSpec) {
      VcService vcService;
      ServiceSpec serviceSpec = hostProfileSpec.service.get();
      vcService = VcServiceUtil.getVcService(serviceSpec);

      return vcService;
   }

   /**
    * Retrieves the host customization data from the TaskInfo object
    *
    * @param taskInfo the task info object
    * @return the customization data result
    */
   private CustomizationData getHostCustomizationData(TaskInfo taskInfo) {
      Object result = taskInfo.getResult();
      CustomizationData hostCustomizationData = (CustomizationData) result;

      if (hostCustomizationData == null) {
         throw new RuntimeException("Host Customization data it null.");
      }

      return hostCustomizationData;
   }

   /**
    * Retrieve service instance content from service spec
    *
    * @param serviceSpec service spec
    * @return service instance content
    * @throws VcException
    */
   private ServiceInstanceContent getServiceInstanceContent(ServiceSpec serviceSpec)
         throws VcException {
      VcService vcService = getVcService(serviceSpec);

      ServiceInstanceContent serviceContent;
      serviceContent = doGetServiceInstanceContent(vcService);
      return serviceContent;
   }

   /**
    * Retrieve VcService from ServiceSpec
    *
    * @param serviceSpec service spec
    * @return vcService
    * @throws VcException
    */
   private VcService getVcService(ServiceSpec serviceSpec) throws VcException {
      return VcServiceUtil.getVcService(serviceSpec);
   }

   /**
    * Retrieves the service instance content from a VcService
    *
    * @param vcService vc service
    * @return serviceInstanceContent
    */
   private ServiceInstanceContent doGetServiceInstanceContent(VcService vcService) {
      ServiceInstanceContent serviceInstanceContent;
      serviceInstanceContent = vcService.getServiceInstanceContent();

      if (serviceInstanceContent == null) {
         throw new RuntimeException("ServiceInstanceContent is null");
      }
      return serviceInstanceContent;
   }
}
