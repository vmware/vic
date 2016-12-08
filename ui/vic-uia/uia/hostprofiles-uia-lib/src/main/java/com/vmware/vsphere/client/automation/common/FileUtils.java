/*
 *  Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential
 */

package com.vmware.vsphere.client.automation.common;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Path;

import static org.apache.commons.io.FileUtils.openInputStream;
import static org.apache.commons.io.FileUtils.readFileToString;

/**
 * Utilities for file management
 */
public class FileUtils {

    private static final Logger logger =
            LoggerFactory.getLogger(FileUtils.class);
    private static FileUtils INSTANCE = null;

    /**
     * Private method. Use getInstance() instead.
     */
    private FileUtils() {

    }

    /**
     * Get instance of HostProfileSrvApi.
     *
     * @return created instance
     */
    public static FileUtils getInstance() {
        if (INSTANCE == null) {
            synchronized (FileUtils.class) {
                if (INSTANCE == null) {
                    logger.info("Initializing FileUtils.");
                    INSTANCE = new FileUtils();
                }
            }
        }

        return INSTANCE;
    }

    /**
     * Write a text string to a text file. The file is placed in the working
     * project folder.
     *
     * @param text     the text to be written to file
     * @param fileName the file name
     */
    public void writeStringToTextFile(String text, String fileName) {
        try (Writer writer = new BufferedWriter(new OutputStreamWriter(new FileOutputStream(
                fileName), "utf-8"))) {
            writer.write(text);
        } catch (IOException e) {
            String errorText = "Error writing text %s to file %s";
            String errorMessage = String.format(errorText, text, fileName);
            throw new RuntimeException(errorMessage, e);
        }
    }

    /**
     * Deletes a file from the working project folder
     *
     * @param fileName the file name to use for file deletion
     */
    public void deleteFile(String fileName) {
       logger.info(String.format("Deleting file with name: %s", fileName));
        File file = new File(fileName);
        Path filePath = file.toPath();
        deleteFile(filePath);
    }

   /**
    * Reads a file from the working project folder
    * @param fileName the file name to use for reading
    * @return the content of the file
    */
   public String readFile(String fileName) {
      File file = new File(fileName);
      String fileContents;
      try {
         fileContents = readFileToString(file);
      } catch (IOException e) {
         String message = String.format("Error reading file: %s",fileName);
         throw new RuntimeException(message, e);
      }

      return fileContents;
   }

   /**
    * Reads a file from the working project folder and returns a
    * FileInputStream.
    * @param fileName the file name to use
    * @return FileInputStream
    */
   public FileInputStream readFileToStream(String fileName) {
      File file = new File(fileName);
      FileInputStream fileInputStream;
      try {
         fileInputStream = openInputStream(file);
      } catch (IOException e) {
         String message = String.format("Error reading file: %s", fileName);
         throw new RuntimeException(message, e);
      }

      return fileInputStream;
   }

    /**
     * Deletes a file
     *
     * @param filePath the path to the file
     */
    public void deleteFile(Path filePath) {
        try {
            System.gc();
            Files.deleteIfExists(filePath);
        } catch (IOException e) {
            throw new RuntimeException("Error deleting file", e);
        }
    }

    /**
     * Checks whether a file exists. The file should be in the working
     * project folder.
     *
     * @param fileName the file name to search for
     * @return true if it exists, false is it doesn't
     */
    public boolean isFileExists(String fileName) {
        File file = new File(fileName);
        boolean fileExists = file.exists();
        return fileExists;
    }
}
