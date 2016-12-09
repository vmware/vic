package com.vmware.hsua.common.util;

import java.util.ArrayList;
import java.util.List;

public class KeyValuePairFileBuilder {

   public static String DEPENDENCIES_KEY = "dependencies";

   private final String _pathToFile;
   private final List<String> _fileContent = new ArrayList<String>();

   public KeyValuePairFileBuilder(String pathToFile){
      this._pathToFile = pathToFile;

   }

   /**
    * Add a comment to the internal file content array.
    *
    * @param comment
    */
   public void addComment(String comment) {
      _fileContent.add("#" + comment);
   }

   /**
    * Adds a simple key value pair to internal file content array.
    *
    * @param key
    *           non null, non empty string
    * @param value
    *           non null, non empty string
    */
   public void addSimpleKeyValuePair(String key, String value) {
      checkInput(key, value);
      _fileContent.add(key + "=" + value);
   }

   /**
    * Adds a array key value pair to internal file content array.
    *
    * @param key
    *           non null, non empty string
    * @param array
    *           must contain non null, non empty strings
    */
   public void addArrayKeyValuePair(String key, String[] array) {
      checkInput(key, array);
      StringBuilder builder = new StringBuilder();
      builder.append(key + "=[");
      for (int i = 0; i < array.length; i++) {
         builder.append(array[i]);
         if (i + 1 < array.length) {
            builder.append(", ");
         } else {
            builder.append("]");
         }
      }
      _fileContent.add(builder.toString());
   }

   /**
    * Creates a file with the added content.
    *
    * @return true if IO operation was success, false otherwise
    */
   public boolean build() {
      return IOUtils.writeLinesToFile(_fileContent, _pathToFile);
   }

   private void checkInput(String key, String[] array) {
      if (key == null || array == null) {
         throw new IllegalArgumentException(
               "KeyValuePairFileBuilder.addArrayKeyValuePair: key or array is null!");
      }
      if (key.isEmpty()) {
         throw new IllegalArgumentException(
               "KeyValuePairFileBuilder.addArrayKeyValuePair: key is empty!");
      }
      for (String string : array) {
         if (string == null || string.isEmpty()) {
            throw new IllegalArgumentException(
                  "KeyValuePairFileBuilder.addArrayKeyValuePair: array element is null or empty");
         }
      }
   }

   private void checkInput(String key, String value) {
      if (key == null || value == null) {
         throw new IllegalArgumentException(
               "KeyValuePairFileBuilder.addSimpleKeyValuePair: key or value is null!");
      }
      if (key.isEmpty() || value.isEmpty()) {
         throw new IllegalArgumentException(
               "KeyValuePairFileBuilder.addSimpleKeyValuePair: key or value is empty!");
      }
   }

}
