package com.vmware.client.automation.util;

import com.vmware.client.automation.servicespec.SeleniumServiceSpec;
import com.vmware.vsphere.client.automation.provider.connector.SeleniumConnector;

public class SeleniumUtil {

   // Fake singleton implementation
   // TODO: make it work with connection registry
   private static SeleniumConnector connector = null;

   private static SeleniumServiceSpec _seleniumSpec;

   // helper methods for getting selenium service spec before fixture functionality is introduced
   // TODO : remove it once fixture for selenium server is provided
   public static void setSeleniumSpec(SeleniumServiceSpec seleniumSpec) {
      SeleniumUtil._seleniumSpec = seleniumSpec;
   }

   // helper methods for getting selenium service spec before fixture functionality is introduced
   // TODO : remove it once fixture for selenium server is provided
   @Deprecated
   public static SeleniumServiceSpec getSeleniumSpec() {
      return _seleniumSpec;
   }

   // Fake implementation
   // TODO: make it work with connection registry
   public static SeleniumConnector getSeleniumConnector(SeleniumServiceSpec _seleniumSpec) {
      if (connector == null) {
         connector = new SeleniumConnector(_seleniumSpec);
      }
      return connector;
   }
}
