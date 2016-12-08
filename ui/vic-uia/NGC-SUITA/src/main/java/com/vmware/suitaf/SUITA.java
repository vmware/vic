package com.vmware.suitaf;

import static com.vmware.suitaf.util.CommonUtils.smartEqual;

import java.util.HashMap;
import java.util.Map;
import java.util.logging.Level;

import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.AutomationPlatformLinkExt;
import com.vmware.suitaf.apl.HostLogger;
import com.vmware.suitaf.util.Logger;

/**
 * This class serves as the main access point for the SUITA framework. It
 * includes a factory subclass that allows creation one instance of the
 * framework's UI Tool and allows access to it. The factory method creates also
 * the {@link AutomationPlatformLink} instance through which the framework
 * connects to the chosen Automation Agent and Automation Platform.<br>
 * <i><b>NOTE:</b> The SUITA framework at its current state is designed for
 * sequential test execution. The parallel execution of tests must be
 * implemented by running multiple instances of the UI Test Automation
 * application against different test environments.</i> <br>
 * <br>
 *
 * @author dkozhuharov
 */
public final class SUITA {
   /**
    * This class plays the role of a factory for the SUITA framework components.
    * It takes care for creation of one {@link UIAutomationTool} instance that
    * is provided to all framework users. It also has the tool methods for
    * creation, configuration and release of an APL interface implementation,
    * which is the back-end of the automation tool. <br>
    * <br>
    *
    * @author dkozhuharov
    */
   public static final class Factory {
      // =================================================================
      // UI Automation Tool instantiation
      // =================================================================
      public static final UIAutomationTool UI_AUTOMATION_TOOL;
      // Prepare an ad-hoc HostLogger implementation that redirects
      // log messages from the APL implementation to the default Logger
      private static final HostLogger HOST_LOGGER;
      static {
         HOST_LOGGER = new HostLogger() {
            @Override
            public void debug(String message) {
               Logger.debug(message);
            }

            @Override
            public void error(String message) {
               Logger.error(message);
            }

            @Override
            public void info(String message) {
               Logger.info(message);
            }

            @Override
            public void warn(String message) {
               Logger.warn(message);
            }

            @Override
            public void dump(String message) {
               Logger.info(message);
            }
         };

         // Create the new SUITA framework instance
         UI_AUTOMATION_TOOL = new UIAutomationTool(HOST_LOGGER);
      };

      // =================================================================
      // APL Base Control methods
      // =================================================================
      private static String instClassName = null;
      private static boolean isInstantiated = false;

      private static Map<String, String> initParams = null;
      private static boolean isInitialized = false;

      private static AutomationPlatformLink apl = null;
      private static AutomationPlatformLinkExt aplx = null;

      public static void aplSetup(final String instClassName,
            final Map<String, String> initParams) {
         if (!smartEqual(Factory.initParams, initParams)) {
            if (isInitialized) {
               aplClose();
            }
            Factory.initParams = initParams;
         }

         if (!smartEqual(Factory.instClassName, instClassName)) {
            if (isInstantiated) {
               apl = null;
               aplx = null;
               isInstantiated = false;
            }
            Factory.instClassName = instClassName;
         }

         if (!isInstantiated) {
            Logger.info("INSTANTIATING A.P.L. WITH: " + instClassName);

            // Instantiation of the APL implementation
            try {
               // Runtime class loading
               Class<?> aplClass = Class.forName(instClassName);
               // Instantiation of the APL implementation
               apl = (AutomationPlatformLink) aplClass.newInstance();
               // Check if extended APL interface is supported too
               if (apl instanceof AutomationPlatformLinkExt) {
                  aplx = (AutomationPlatformLinkExt) apl;
               }
               // Instantiation completed
               isInstantiated = true;
            } catch (Exception e) {
               // Instantiation failed
               throw new RuntimeException(
                     "Failed to instantiate APL implementation for"
                           + " classname '" + instClassName + "'.", e);
            }
         }
      }

      private static AutomationPlatformLink aplInst() {
         if (!isInstantiated) {
            throw new RuntimeException("APL not instantiated."
                  + " Call method [aplSetup} first.");
         }
         return apl;
      }

      private static synchronized AutomationPlatformLinkExt aplxInst() {
         if (!isInstantiated) {
            throw new RuntimeException("APL not instantiated."
                  + " Call method [aplSetup} first.");
         }
         return aplx;
      }

      public static synchronized AutomationPlatformLink apl() {
         if (!isInitialized) {
            aplReset();
         }
         return apl;
      }

      public static synchronized AutomationPlatformLinkExt aplx() {
         if (aplxInst() == null) {
            throw new RuntimeException("APL implementation"
                  + " does not support extended functionality!");
         }
         if (!isInitialized) {
            aplReset();
         }
         return aplx;
      }

      public static void aplReset() {
         isInitialized = true;
         Logger.info("INITIALIZING S.U.I.T.A. WITH: " + initParams);
         aplInst().resetLink(initParams, HOST_LOGGER);
      }

      public static void aplRepair() {
         isInitialized = true;
         Logger.info("INITIALIZING S.U.I.T.A. WITH: " + initParams);
         aplxInst().attachLink(initParams, HOST_LOGGER);
      }

      public static void aplClose() {
         if (isInitialized) {
            aplInst().resetLink(new HashMap<String, String>(), null);
            isInitialized = false;
         }
      }
   }

   /**
    * This class serves the role of a static environment for the operations of
    * the SUITA framework. The state of this environment must be made up to date
    * before each UI-test is executed. <br>
    * <br>
    *
    * @author dkozhuharov
    */
   public static final class Environment {
      // =============================================================
      // Configuration parameters initialization
      // =============================================================
      public static long getBackendJobSmall() {
         return PAGE_LOAD_TIMEOUT * 4;
      }

      public static long getBackendJobMid() {
         return PAGE_LOAD_TIMEOUT * 8;
      }

      public static long getBackendJobLarge() {
         return PAGE_LOAD_TIMEOUT * 20;
      }

      public static long getPageLoadTimeout() {
         return PAGE_LOAD_TIMEOUT;
      }

      public static long getUIOperationTimeout() {
         return PAGE_LOAD_TIMEOUT / 3;
      }

      private static long PAGE_LOAD_TIMEOUT = 30000;

      private static boolean IS_DEBUG_ENABLED = false;

      public static boolean isDebugEnabled() {
         return IS_DEBUG_ENABLED;
      }

      private static String SCREENSHOT_DIR = "<Not INITIALIZED>";

      static String getErrorImageDir() {
         return SCREENSHOT_DIR;
      }

      // Web server where the screen shots can be viewed.
      private static String SCREENSHOT_WEB_SERVER = "";

      static String getErrorImageWebURL() {
         return SCREENSHOT_WEB_SERVER;
      }

      public static void set(Boolean debugEnabled, String hostSnapshotFolder,
            Long pageLoadTimeout, String snapshotWebServer) {
         IS_DEBUG_ENABLED = debugEnabled;
         SCREENSHOT_DIR = hostSnapshotFolder;
         SCREENSHOT_WEB_SERVER = snapshotWebServer;
         PAGE_LOAD_TIMEOUT = pageLoadTimeout;

         // Disable java.util.logging if the test does not run in debug mode.
         if (IS_DEBUG_ENABLED) {
            java.util.logging.Logger.getLogger("").setLevel(Level.FINE);
         } else {
            java.util.logging.Logger.getLogger("").setLevel(Level.OFF);
         }
      }
   }

}
