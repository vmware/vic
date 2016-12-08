/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.suitaf.apl.sele.sikulix;

import org.sikuli.script.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.vmware.suitaf.apl.sele.sikulix.ConnectionType.LOCAL;
import static com.vmware.suitaf.apl.sele.sikulix.ConnectionType.REMOTE_VNC;

/**
 * Class used to initialize SikuliX and get a Region object.  The Region
 * object is a representation of the screen and is used to do all Sikuli
 * operations.
 */
public class SikuliInitializer {
   private static final Logger _logger = LoggerFactory.getLogger(
         SikuliInitializer.class);
   private static final String DEFAULT_VNC_PORT = "5900";
   private static final String LOCALHOST_NAME = "localhost";
   private static final String LOCALHOST_IP = "127.0.0.1";
   private static Region screen;

   /**
    * Setting the constructor as private, in order to prevent instantiations
    */
   private SikuliInitializer() {
   }

   /**
    * Initializes Sikuli and returns a valid Region object.
    *
    * @param hostName
    *
    * @return
    */
   public static Region initSikuliAndGetRegion(String hostName) {
      _logger.debug(String.format("Init Sikuli with '%s'", hostName));
      verifyHostNameInput(hostName);
      initializeSikuli(hostName);
      _logger.debug("Done initializing Sikuli.");
      return screen;
   }

   /**
    * Verify that the hostname is not Null or Empty.
    *
    * @param hostName
    */
   private static void verifyHostNameInput(String hostName) {
      if (hostName == null || hostName.isEmpty()) {
         throw new IllegalArgumentException("Host name is empty!");
      }
   }

   /**
    * initializes Sikuli using the the target hostname
    *
    * @param hostName - the host name of the target machine
    */
   private static void initializeSikuli(String hostName) {
      ConnectionType connectionType = getConnectionType(hostName);
      final String[] parameters = new String[]{hostName, DEFAULT_VNC_PORT};
      _logger.debug(String.format("ConnectionType is: '%s'", connectionType));
      initializeScreen(connectionType, parameters);
   }

   /**
    * @param hostName
    *
    * @return connection type based on the hostname
    */
   private static ConnectionType getConnectionType(String hostName) {
      ConnectionType connectionType;
      if (isHostLocal(hostName)) {
         connectionType = LOCAL;
      } else {
         connectionType = REMOTE_VNC;
      }
      _logger.debug(String.format("Connection type is %s", connectionType));
      return connectionType;
   }

   /**
    * Checks whether the host name is local (127.0.0.1 or localhost).
    *
    * @param hostName
    *
    * @return
    */
   private static boolean isHostLocal(String hostName) {
      boolean isHostLocal;
      if (isHostNameLocal(hostName)) {
         isHostLocal = true;
      } else {
         isHostLocal = false;
      }

      return isHostLocal;
   }

   /**
    * Checks whether the passed in host name matches one of the formats for
    * local host
    *
    * @param hostName
    *
    * @return
    */
   private static boolean isHostNameLocal(String hostName) {
      return LOCALHOST_NAME.equals(hostName) || LOCALHOST_IP.equals(hostName);
   }

   /**
    * Initializes Sikuli by making a call to SikuliFactory, which returns a
    * valid Region object based on passed in connectionType and parameters.
    *
    * @param connectionType
    * @param parameters
    */
   private static void initializeScreen(
         ConnectionType connectionType, String[] parameters) {
      _logger.debug("Initializing screen");
      if (screen == null) {
         _logger.debug("Sikuli screen is null, creating a new one.");
         screen = SikuliFactory.getScreenRegion(connectionType, parameters);
      } else {
         _logger.debug("Sikuli screen is NOT null. NOT creating a new one");
      }
   }
}
