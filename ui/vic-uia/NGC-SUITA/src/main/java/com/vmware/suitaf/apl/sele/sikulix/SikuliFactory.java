/*
 * Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.suitaf.apl.sele.sikulix;

import org.sikuli.script.ImagePath;
import org.sikuli.script.Region;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import sun.reflect.Reflection;

import java.net.URI;
import java.net.URISyntaxException;
import java.net.URL;

/**
 * A factory, which creates Sikuli Regions objects.
 */
class SikuliFactory {
   private static final Logger _logger =
         LoggerFactory.getLogger(SikuliFactory.class);

   /**
    * Gets a Sikuli Region based on connection type and parameters
    *
    * @param type Connection type
    * @param parameters Parameters if any are necessary (local doesn't use any)
    *
    * @return implementation of {@link org.sikuli.script.Region} object.
    */
   public static Region getScreenRegion(
         ConnectionType type, String[] parameters) {
      verifyConnectionTypeIsNotNull(type);
      setupImagePath();
      return doGetScreenRegion(type, parameters);
   }

   private static void verifyConnectionTypeIsNotNull(ConnectionType type) {
      if (type == null) {
         throw new IllegalArgumentException("ConnectionType is null.");
      }
   }

   /**
    * Get the image path and pass them to Sikuli. After Sikuli knows where
    * our images are, we only need to pass in the file names.
    */
   private static void setupImagePath() {
      //TODO mkalinov Find a proper way to get the image path
      // This method gets the resource folder of the current Maven module.
      // This is a temporary solution and not optimal.
      final ClassLoader classLoader = SikuliFactory.class.getClassLoader();
      final URL resource = classLoader.getResource("sikuli");
      String path = "";
      if (resource != null) {
         path = resource.getPath();
      }
      _logger.debug(String.format("Setting the image path to %s", path));
      ImagePath.add(path);
   }

   /**
    * Gets a Sikuli Region based on connection type and parameters
    *
    * @param connectionType
    * @param parameters
    *
    * @return
    */
   private static Region doGetScreenRegion(
         ConnectionType connectionType, String[] parameters) {
      switch (connectionType) {
         case LOCAL:
            return getLocalRegion();
         case REMOTE_VNC:
            return getRemoteRegion(parameters);
         default:
            throw new RuntimeException("The connection type is not recognized");
      }
   }

   /**
    * Returns a new local region
    *
    * @return an initialized Region object
    */
   private static Region getLocalRegion() {
      return new LocalRegionFactory().getRegion();
   }

   /**
    * @param parameters
    *
    * @return an initialized Region object
    */
   private static Region getRemoteRegion(final String[] parameters) {
      validateRemoteScreenParameters(parameters);
      final RemoteVncRegionFactory factory = new RemoteVncRegionFactory(
            parameters[0], parameters[1]);
      return factory.getRegion();
   }

   /**
    * Validate remote screen parameters
    *
    * @param parameters
    */
   private static void validateRemoteScreenParameters(final String[] parameters) {
      validateParametersNotNullAndCorrectCount(parameters);
      validateHostAndPortFormat(parameters);
   }

   /**
    * Validate that parameters value is not null and the parameter count is
    * correct
    *
    * @param parameters
    */
   private static void validateParametersNotNullAndCorrectCount(String[] parameters) {
      if (parameters == null || parameters.length != 2) {
         throw new IllegalArgumentException(
               "parameters must be set and must be exactly two");
      }
   }

   /**
    * Validate host and port values
    *
    * @param parameters
    */
   private static void validateHostAndPortFormat(String[] parameters) {
      try {
         //validate by using the URI class. It does automatic validation
         URI uri = new URI("vnc://" + parameters[0] + ":" + parameters[1]);
         final String host = uri.getHost();
         final int port = uri.getPort();

         validateHostParameter(host);
         validatePortParameter(port);
      } catch (URISyntaxException e) {
         throw new IllegalArgumentException("Host or port format is not valid");
      }
   }

   /**
    * Validate port value is valid
    *
    * @param port
    */
   private static void validatePortParameter(int port) {
      if (port == -1) {
         throw new IllegalArgumentException("The port format is not valid");
      }
   }

   /**
    * Validate host name value is valid
    *
    * @param host
    */
   private static void validateHostParameter(String host) {
      if (host == null || host.isEmpty()) {
         throw new IllegalArgumentException("The host format is not valid.");
      }
   }
}
