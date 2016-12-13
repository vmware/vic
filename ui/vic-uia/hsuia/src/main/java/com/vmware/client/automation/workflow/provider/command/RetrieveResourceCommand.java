package com.vmware.client.automation.workflow.provider.command;

import com.vmware.client.automation.workflow.command.WorkflowCommand.WorkflowCommandAnnotation;
import com.vmware.client.automation.workflow.explorer.SettingsRWImpl;
import com.vmware.client.automation.workflow.explorer.SettingsUtil;
import com.vmware.hsua.common.util.KeyValuePairFileBuilder;

// THIS IS A STUB!
@WorkflowCommandAnnotation(commandName = "retrieve-resource")
public class RetrieveResourceCommand extends BaseProviderCommand {

   private String _providerClassName;
   private String _resultFilePath;

   @Override
   public void prepare(String[] commandParams) throws Exception {
      validateParameters(commandParams, new String[] { "providerClassName", "resultFilePath" });
      _providerClassName = commandParams[0];
      _resultFilePath = commandParams[1];
   }

   // TODO rkovachev: this is only a stub. Implementation is due.
   @Override
   public void execute() throws Exception {
      KeyValuePairFileBuilder fileBuilder = new KeyValuePairFileBuilder(_resultFilePath);
      fileBuilder.addComment(this.getClass().getSimpleName() + " result file");
      fileBuilder.addComment("Resource retrived for provider: " + _providerClassName);
      if (!_providerClassName.contains("NimbusNfsStorageProvider")
            || !_providerClassName.contains("SeleniumNodeProvider")) {
         String buildNum = SettingsUtil.getRequiredValue(new SettingsRWImpl(),
               getProviderClassSimpleName(_providerClassName) + ".buildNumber");
         fileBuilder.addComment("Build number is: " + buildNum);
      }
      fileBuilder.build();
   }

   public String getProviderClassSimpleName(String providerCannonicalClassName) {
      String[] providerClassnameSplittedAtDots = providerCannonicalClassName.split("\\.");
      return providerClassnameSplittedAtDots[providerClassnameSplittedAtDots.length - 1];
   }

}
