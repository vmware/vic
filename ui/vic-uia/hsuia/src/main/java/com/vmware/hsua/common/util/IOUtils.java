/** Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.hsua.common.util;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.io.OutputStream;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Properties;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * I/O utility functions.
 */
public class IOUtils {

   // logger
   private static final Logger _logger = LoggerFactory.getLogger(IOUtils.class);

   /**
    * Loads the configuration properties from the specified file. The configuration
    * should provide settings for all the parameters described here
    * https://wiki.eng.vmware.com/wiki/index.php?title=VSuiteAdminUI/NGC2014/QE/
    * UIAutomation#Common_resources_parameters
    *
    * @param configurationFile
    *       path to the configuration file.
    * @param ui
    *       UI automation tool.
    */
   public static Properties readConfiguration(String configurationFile) {
      Properties properties = new Properties();
      InputStream in = null;

      _logger.info("Loading configuration from: " + configurationFile);

      try {
         in = new FileInputStream(configurationFile);
         properties.load(in);
      } catch (IOException e) {
         throw new RuntimeException("IOExcpetion thrown while loading properties from: "
               + configurationFile, e);
      } catch (Exception e) {
         throw new RuntimeException("Excpetion thrown while loading properties from: "
               + configurationFile, e);
      } finally {
         try {
            if (in != null) {
               in.close();
            }
         } catch (IOException e) {
            throw new RuntimeException("Failed to close the input stream for the file: "
                  + configurationFile, e);
         }
      }

      return properties;
   }

   /**
    * Read object from stream.
    */
   public static Object readObject(InputStream is) throws Exception {
      ObjectInputStream ois = new ObjectInputStream(is);
      Serializable object = (Serializable) ois.readObject();

      return object;
   }

   /**
    * Write object to stream.
    */
   public static void writeObject(OutputStream os, Object object) throws IOException {
      ObjectOutputStream oos = new ObjectOutputStream(os);
      oos.writeObject(object);
   }

   /**
    * Writes the strings from the lines list to a file specified by the path.
    * @param lines
    * @param path
    * @return
    */
   public static boolean writeLinesToFile(List<String> lines, String path) {
      FileWriter writer = null;
      boolean success = true;
      try {
         ensureDirecoryTreeExists(path);
         writer = new FileWriter(path);
         for (String line : lines) {
            writer.append(line);
            writer.append('\n');
         }
         writer.flush();
         success = true;
      } catch (IOException e) {
         e.printStackTrace();
         success = false;
      } finally {
         try {
            if (writer != null) {
               writer.close();
               return success;
            }
         } catch (IOException e) {
            e.printStackTrace();
         }
      }
      return success;
   }

   /**
    * Read line by line the test file path and return list of these lines.
    * @param path
    * @return
    * @throws IOException
    */
   public static List<String> readLinesFromFile(String path) throws IOException {
      List<String> result = new ArrayList<String>();
      BufferedReader bufferedReader = null;
      String line = null;
      try {
         bufferedReader = new BufferedReader(new FileReader(path));
         while ((line = bufferedReader.readLine()) != null) {
            result.add(line);
         }
      } finally {
         if (bufferedReader != null) {
            try {
               bufferedReader.close();
            } catch (IOException e) {
               e.printStackTrace();
            }
         }
      }
      return result;
   }

   /**
    * Creates the directorCreates the directory structure for the specidies file path.
    */
   private static void ensureDirecoryTreeExists(String path) {
      File file = new File(path);
      if (!file.getParentFile().exists()) {
         file.getParentFile().mkdirs();
      }
   }
}
