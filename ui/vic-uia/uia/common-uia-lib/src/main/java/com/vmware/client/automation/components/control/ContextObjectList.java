/**
 * Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.control;

import java.lang.reflect.Constructor;
import java.util.Arrays;
import java.util.List;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.flexui.selenium.MethodCallUtil;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Implements the flex ContextObjectList control. Similar to AdvancedDataGrid but can not be handled by it.
 */
public class ContextObjectList {

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private final static String ID_CONTENT_HOLDER = "/className=AdvancedListBaseContentHolder";
   private final static String ID_CELL_TEXT = "/className=UITextField";

   private final String _id;

   private final List<String> _columnNames;
   private final List<String> _columnClassNames;
   private final int _columnsCount;
   private String _firstColumnClassName;
   private int _firstColumnClassNameCount;

   /**
    * Constructs the ContextObjectList
    * @param id - the ID of the list
    * @param columnNames - the name of the list columns
    * @param columnClassNames - the Flex class names for the cell controls that correspond to the columns
    */
   public ContextObjectList(String id, String [] columnNames, String [] columnClassNames) {
      _id = id;

      _columnNames = Arrays.asList(columnNames);
      _columnClassNames = Arrays.asList(columnClassNames);
      _columnsCount = columnClassNames.length;

      _firstColumnClassName = null;
      _firstColumnClassNameCount = 1;
      for (String columnClassName: _columnClassNames) {
         if (_firstColumnClassName == null) {
            _firstColumnClassName = columnClassName;
            continue;
         }

         if (columnClassName.equals(_firstColumnClassName)) {
            _firstColumnClassNameCount++;
         }
      }
   }

   /**
    * Gets the count of the columns
    * @return the count of the columns
    */
   public int getColumnsCount() {
      return _columnsCount;
   }

   /**
    * Gets the index of column by name
    * @param columnName - the name of the column
    * @return the index of column by name or -1 if not found
    */
   public int getColumnIndex(String columnName) {
      return _columnNames.indexOf(columnName);
   }

   /**
    * Gets the row count
    * @return the row count
    */
   public int getRowCount() {
      int controlCount = getVisibleControlCount(_id + ID_CONTENT_HOLDER + "/className=" + _firstColumnClassName);
      return controlCount / _firstColumnClassNameCount;
   }

   /**
    * Gets the String value in a given cell in the grid.
    * @param rwoIndex - row index
    * @param columnName - the column from which to get the value
    * @return - returns the cell value of the specified entity in the grid
    */
   public String getCellValue(int rowIndex, String columnName) {
      int columnIndex = getColumnIndex(columnName);

      String controlClassName = _columnClassNames.get(columnIndex);

      int controlIndex = rowIndex * getControlsOfClassPerRow(controlClassName) + getControlIndexInRow(controlClassName, columnIndex);
      DisplayObject control = getVisibleControl(DisplayObject.class, _id + ID_CONTENT_HOLDER + "/className="
            + controlClassName + "[%d]" + ID_CELL_TEXT, controlIndex);
      return control.getProperty("text");
   }

   /**
    * Gets the index of an item in the list
    * @param columnName - the name of the column to search in
    * @param columnValue - the column value to search for
    * @return the index of an item in the list
    */
   public int getItemIndex(String columnName, String columnValue) {
      int rowCount = getRowCount();
      for (int row = 0; row < rowCount; row++) {
         if (getCellValue(row, columnName).equals(columnValue)) {
            return row;
         }
      }

      return -1;
   }

   // ---------------------------------------------------------------------------
   // Private methods

   private int getControlCount(String id) {
      for (int i = 0; i < 200; i++) {
         DisplayObject control = new DisplayObject(id + "[" + i + "]", BrowserUtil.flashSelenium);
         if (!control.isComponentExisting()) {
            return i;
         }
      }
      return 0;
   }

   private int getVisibleControlCount(String id) {
      int controlCount = getControlCount(id);
      int result = 0;
      for (int i = 0; i < controlCount; i++) {
         if (MethodCallUtil.getVisibleOnPath(BrowserUtil.flashSelenium, id + "[" + i + "]")) {
            result ++;
         }
      }
      return result;
   }

   private <T extends DisplayObject> T getVisibleControl(Class<T> controlClass, String idFormat, int index) {
      int visibleIndex = -1;
      for (int i = 0; i < 200; i++) {
         String id = String.format(idFormat, i);
         if (MethodCallUtil.getVisibleOnPath(BrowserUtil.flashSelenium, id)) {
            visibleIndex++;
            if (visibleIndex == index) {
               Constructor<T> constructor;
               try {
                  constructor = controlClass.getConstructor(String.class, FlashSelenium.class);
                  return constructor.newInstance(id, BrowserUtil.flashSelenium);
               } catch (Exception e) {
                  throw new RuntimeException(e);
               }
            }
         }
      }

      return null;
   }

   private int getControlsOfClassPerRow(String className) {
      int result = 0;
      for (int i = 0; i < _columnsCount; i++) {
         if (_columnClassNames.get(i).equals(className)) {
            result++;
         }
      }
      return result;
   }

   private int getControlIndexInRow(String className, int columnIndex) {
      int result = -1;
      for (int i = 0; i <= columnIndex; i++) {
         if (_columnClassNames.get(i).equals(className)) {
            result++;
         }
      }
      return result;
   }
}
