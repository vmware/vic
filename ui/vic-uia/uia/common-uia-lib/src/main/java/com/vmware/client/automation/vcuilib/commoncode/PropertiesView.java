/*
 * ************************************************************************
 *
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 *
 * ************************************************************************
 */
package com.vmware.client.automation.vcuilib.commoncode;

import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_DEFAULT_PROPERTY_VIEW;
import static com.vmware.client.automation.vcuilib.commoncode.IDConstants.ID_PROPERTIES_VIEW_TAB;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.verifySafely;
import static com.vmware.client.automation.vcuilib.commoncode.TestBaseUI.verifyTrueSafely;
import static com.vmware.flexui.selenium.BrowserUtil.flashSelenium;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map.Entry;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * Class to provide data from properties view
 *
 * NOTE: this class is a partial copy of the one from VCUI-QE-LIB
 */
public class PropertiesView {

   private static final Logger logger = LoggerFactory.getLogger(PropertiesView.class);

   private static final String PROPERTY_VIEW_PROP_VALUE_CLASS_NAME =
         "PropertyViewPropValue";
   private static final String PROPERTY_VIEW_PROP_MULTI_VALUE_CLASS_NAME =
         "PropertyViewPropMultiValue";
   private static final String PROPERTY_VIEW_PROP_MESSAGE_CLASS_NAME =
         "PropertyViewMessage";
   private static final String PROPERTY_VIEW_HEADING_CLASS_NAME = "PropertyViewHeading";
   private static final String PROPERTY_VIEW_PROPERTIES =
         "properties.[].uid,properties.[].className,properties.[].visible";
   private static final String VALUES_PROPERTY = "values.[]";
   private final int UID_PROPERTIES_POS = 0;
   private final int CLASS_NAME_PROPERTIES_POS = 1;
   private final int VISIBLE_PROPERTIES_POS = 2;

   List<PropertyViewObject> list;
   private String uiID;
   private DisplayObject viewObject;

   // The following constant will be used as default delimiters for arrays in
   // returns from the Selenium Flex request
   private static final String DEFAULT_ARRAY_DELIMITER = "^^";

   // The following constants will be used as default splitters for arrays/properties when
   // parsing the output of the Selenium Flex request.
   private static final String DEFAULT_ARRAY_SPLITTER = "\\^\\^";
   private static final String DEFAULT_PROPERTIES_SPLITTER = ",";

   // Row and column delimiters and splitters that will be used in parsing the
   // output from Selenium Flex call
   private String arrayDelimiter;
   private String arraySplitter;

   public PropertiesView() {
      this(ID_DEFAULT_PROPERTY_VIEW);
   }

   public PropertiesView(String uiID) {
      this(uiID, DEFAULT_ARRAY_DELIMITER, DEFAULT_ARRAY_SPLITTER);
   }

   public PropertiesView(String uiID, String arrayDelimiter, String arraySplitter) {
      this.uiID = uiID;
      this.arrayDelimiter = arrayDelimiter;
      this.arraySplitter = arraySplitter;
      viewObject = new DisplayObject(this.uiID, BrowserUtil.flashSelenium);
      init(uiID);
   }

   /**
    * Get the id of the property view
    * @return
    */
   public String getUid() {
      return this.uiID;
   }

   /**
    * Sets array delimiter
    *
    * @param arrayDelimiter
    */
   public void setArrayDelimiter(String arrayDelimiter) {
      this.arrayDelimiter = arrayDelimiter;
   }

   /**
    * Sets array splitter
    *
    * @param arraySplitter
    */
   public void setArraySplitter(String arraySplitter) {
      this.arraySplitter = arraySplitter;
   }

   /**
    * Change properties view tab
    *
    * @param tabName - tab name
    */
   public void changeTab(String tabName) {
      UIComponent tabBarTab =
            new UIComponent(uiID + "/" + ID_PROPERTIES_VIEW_TAB + tabName, flashSelenium);
      tabBarTab.leftMouseClick(false, false, false, null);
      refresh();
   }

   /**
    * Initialize the PropertiesView
    */
   private void init(String uiID) {
      // Commenting this out to remove the assumption the property view is used only in a read-only mode.
      // GlobalFunction.waitForRefreshButtonEnable(Timeout.THIRTY_SECONDS.getDuration());
      String propertiesViewContent =
            viewObject.getProperties(PROPERTY_VIEW_PROPERTIES, arrayDelimiter);
      initPropertiesViewElements(propertiesViewContent);
   }

   private void initPropertiesViewElements(String content) {
      String[] properties = content.split(DEFAULT_PROPERTIES_SPLITTER);
      String[] uidProperties = properties[UID_PROPERTIES_POS].split(arraySplitter);
      String[] classNamesProperties =
            properties[CLASS_NAME_PROPERTIES_POS].split(arraySplitter);
      String[] visibleProperties =
            properties[VISIBLE_PROPERTIES_POS].split(arraySplitter);
      String propertyObjectID = null;
      list = new ArrayList<PropertyViewObject>();
      for (int i = 0; i < uidProperties.length; i++) {
         if (!Boolean.parseBoolean(visibleProperties[i])) {
            continue;
         }
         propertyObjectID = uiID + "/uid=" + uidProperties[i];
         list.add(create(
               new DisplayObject(propertyObjectID, BrowserUtil.flashSelenium),
               classNamesProperties[i]));
      }
   }

   /**
    * Create a new PropertyViewObject of the appropriate type, depending on the
    * className property.
    *
    * @param object - the DisplayObject component
    * @param className - the class name of the component
    * @param pos - the y coordinate
    * @return PropertyViewObject
    */
   private PropertyViewObject create(DisplayObject object, String className) {
      PropertyViewObject propertyViewObject;

      if (PROPERTY_VIEW_PROP_VALUE_CLASS_NAME.equals(className)) {
         propertyViewObject = new PropertyViewObject(object, false);
      } else if (PROPERTY_VIEW_PROP_MULTI_VALUE_CLASS_NAME.equals(className)) {
         propertyViewObject = new MultiPropertyViewObject(object);
      } else if (PROPERTY_VIEW_PROP_MESSAGE_CLASS_NAME.equals(className)) {
         propertyViewObject = new PropertyViewMessageObject(object);
      } else if (PROPERTY_VIEW_HEADING_CLASS_NAME.equals(className)) {
         propertyViewObject = new PropertyViewObject(object, true);
      } else {
         // default choice
         propertyViewObject = new PropertyViewObject(object, false);
      }

      return propertyViewObject;
   }

   /**
    * Reloads the data from the properties view
    */
   public void refresh() {
      init(uiID);

   }

   /**
    * Returns the name of the heading on <code>index</code> position in the
    * properties view. The top item( heading or property) in properties view has
    * index 0, the item under it 1 and etc. If the given <code>index</code> is
    * the position of property the method returns null
    *
    * @param index
    * @return String
    */
   public String getHeadingName(int index) {

      if (list.size() < index || !list.get(index).isHeading()) {

         return null;
      }

      return list.get(index).getName();
   }

   /**
    * Returns whether property on <code>index</code> is heading
    *
    * @param index
    * @return boolean
    */
   public boolean isHeading(int index) {

      if (list.size() < index) {

         return false;
      }
      return list.get(index).isHeading();
   }

   /**
    * Returns the count of properties view elements
    */
   public int getPropertiesViewElementsSize() {

      if (list == null) {

         return -1;
      }
      return list.size();
   }

   /**
    * Returns the name of the property on <code>index</code> position in the
    * properties view. The top item( heading or property) in properties view has
    * index 0, the item under it 1 and etc. If the <code>index</code> is the
    * position of heading the method returns null
    *
    * @param index
    * @return String
    */
   public String getPropertyName(int index) {

      if (list.size() < index || !list.get(index).isProperty()) {

         return null;
      }

      return list.get(index).getName();
   }

   /**
    * Returns the value of the property on <code>index</code> position in
    * properties view. The top item( heading if exist or property) in properties
    * view has index 0, the item under it 1 and etc. If the <code>index</code> is the position of heading the method
    * returns null
    *
    * @param index
    * @return String
    */
   public String getPropertyValue(int index) {

      if (list.size() < index || !list.get(index).isProperty()) {

         return null;
      }

      return list.get(index).getValue();
   }

   /**
    * Returns the value of the property of the properties view which has <code>propertyName<code> name
    *
    * @param propertyName
    */
   public String getPropertyValue(String propertyName) {

      for (PropertyViewObject object : list) {
         if (object.isProperty() && object.getName().equals(propertyName)) {
            return object.getValue();
         }

      }
      return null;
   }

   /**
    * Returns the value of the property of the properties view which has
    * <code>propertyName<code> name and is under heading with <code>headingName<code> name
    *
    * @param headingName
    * @param propertyName
    */
   public String getPropertyValue(String headingName, String propertyName) {

      PropertyViewObject heading = null;
      for (PropertyViewObject object : list) {
         if (object.isHeading() && object.getName().equals(headingName)) {
            heading = object;
         }

         if (heading != null && object.isProperty()
               && object.getName().equals(propertyName)) {
            return object.getValue();
         }

      }

      return null;
   }

   /**
    * Verify that properties view has property with <code>propertyName<code> name
    * and <code>propertyValue<code> value
    *
    * @param headingName
    */
   public void verifyPropertyValue(String propertyName, String propertyValue) {

      String value = getPropertyValue(propertyName);
      if (value != null) {
         verifySafely(propertyValue, value, "The properties view has  property with "
               + propertyName + " name with corect value ");
         return;
      } else {
         verifyTrueSafely(false, "Property with " + propertyName
               + " name is not part of the properties view");
      }

   }

   /**
    * Verify that properties view has heading with <code>headingName<code> name
    *
    * @param headingName
    */
   public void verifyHeading(String headingName) {

      verifyTrueSafely(
            hasHeadingWithName(headingName),
            "Properties view has heading with name " + headingName);

   }

   /**
    * Verify that properties view has property with <code>propertyName<code> name
    *
    * @param headingName
    * @param propertyName
    */
   public void verifyPropertyName(String propertyName) {

      verifyTrueSafely(
            hasPropertyWithName(propertyName),
            "Properties view has property with name " + propertyName);

   }

   /**
    * Verify that properties view has property with <code>propertyName<code> name
    * under the heading with <code>headingName</code> name
    *
    * @param headingName
    * @param propertyName
    */
   public void verifyPropertyNameUnderHeading(String headingName, String propertyName) {

      verifyTrueSafely(
            hasPropertyWithName(headingName, propertyName),
            "Properties view has property with name " + propertyName + " under "
                  + headingName + " heading ");
   }

   /**
    * Verify Data for All headings - main heading and each single entry under
    * that main heading:
    *
    * @param allMainHeadingsData: each entry constructs from heading name and
    *            all property names with according values under that heading
    */
   public void verifyPropertyAllHeadings(
         HashMap<String, HashMap<String, String[]>> allMainHeadingsData) {
      for (Entry<String, HashMap<String, String[]>> singleEntry : allMainHeadingsData
            .entrySet()) {
         String headingName = singleEntry.getKey();
         HashMap<String, String[]> expectedPageDataMap = singleEntry.getValue();
         verifyProperty(headingName, expectedPageDataMap);
      }
   }

   /**
    * Verify that heading in Ready To complete wizard page contains property
    * entries with according values - these values are from array type String[]
    *
    * @param headingName - name of the heading
    * @param expectedPageDataMap: each entry constructs from propertyName - name
    *            of property under heading propertyValue[] - value of that
    *            property under heading
    */
   public void verifyProperty(String headingName,
         HashMap<String, String[]> expectedPageDataMap) {
      PropertyViewObject heading = null;
      int headingIndex = 0;
      for (PropertyViewObject object : list) {
         if (object.isHeading() && object.getName().equals(headingName)) {
            heading = object;
            break;
         }
         headingIndex++;
      }

      if (heading == null) {
         verifyTrueSafely(false, headingName + " heading is found ");
         return;
      }

      if (expectedPageDataMap == null || expectedPageDataMap.size() == 0) {
         return;
      }

      String propertyName = null;
      String[] propertyValues = null;

      boolean propertyFound = false;
      for (Entry<String, String[]> expectedPageData : expectedPageDataMap.entrySet()) {
         propertyName = expectedPageData.getKey();
         propertyValues = expectedPageData.getValue();

         propertyFound = false;
         for (int i = headingIndex + 1; i < list.size(); i++) {
            PropertyViewObject object = list.get(i);
            if (object.isHeading()) {
               verifyTrueSafely(false, propertyName + " property " + " under "
                     + headingName + " heading is found ");
               return;
            }

            if (object.getName().equals(propertyName)) {
               if (object instanceof MultiPropertyViewObject) {
                  Set<String> values = ((MultiPropertyViewObject) object).getValues();
                  for (String propertyValue : propertyValues) {

                     verifyTrueSafely(values.contains(propertyValue), "Property "
                           + propertyName + " under " + headingName
                           + " heading has expected " + propertyValue + " value ");
                  }
               } else {
                  logger.info("object.getValue() = " + object.getValue());
                  logger.info("propertyValues[0] = " + propertyValues[0]);
                  verifySafely(object.getValue(), propertyValues[0], "Property "
                        + propertyName + " under " + headingName
                        + " heading has correct value ");
               }
               propertyFound = true;
               break;
            }

         }
         TestBaseUI.verifyTrueSafely(propertyFound, propertyName + " property "
               + " under " + headingName + " heading is found ");
      }
   }

   /**
    * Returns true if property view contain property with <code>propertyName</code> name under heading with
    * <code>headingName</code> name, otherwise false
    *
    * @param headingName
    * @param @param headingName
    * @returns boolean
    */
   public boolean hasPropertyWithName(String headingName, String propertyName) {

      PropertyViewObject heading = null;
      for (PropertyViewObject object : list) {
         if (object.isHeading() && object.getName().equals(headingName)) {
            heading = object;
         }

         if (heading != null && object.isProperty()
               && object.getName().equals(propertyName)) {
            return true;
         }

      }
      return false;
   }

   /**
    * Returns true if property view contain property with <code>propertyName</code> name, otherwise false
    *
    * @param propertyName
    * @returns boolean
    */
   public boolean hasPropertyWithName(String propertyName) {

      for (PropertyViewObject object : list) {
         if (object.isProperty() && object.getName().equals(propertyName)) {
            return true;
         }

      }
      return false;
   }

   /**
    * Returns true if property view contain heading with <code>headingName</code> name, otherwise false
    *
    * @param headingName
    * @returns boolean
    */
   public boolean hasHeadingWithName(String headingName) {

      for (PropertyViewObject object : list) {
         if (object.isHeading() && object.getName().equals(headingName)) {
            return true;
         }

      }
      return false;
   }

   /**
    * Represent PropertyViewHeading and PropertyViewPropValue object part of
    * property view
    */
   class PropertyViewObject {

      protected static final String TITLE_PROPERTY = "title";
      private static final String VALUE_PROPERTY = "value";

      protected DisplayObject object;
      protected String value;
      protected String name;
      private boolean isHeading;

      PropertyViewObject(DisplayObject object, boolean isHeading) {
         this.object = object;
         this.isHeading = isHeading;
      }

      /**
       * Returns true if the object is heading, otherwise false
       */
      public boolean isHeading() {
         return isHeading;
      }

      /**
       * Returns true if the object is a property, false otherwise
       */
      public boolean isProperty() {
         return !isHeading;
      }


      /**
       * Returns property value
       */
      public String getValue() {
         if (value == null) {
            value = object.getProperty(VALUE_PROPERTY);
         }

         return value;
      }

      /**
       * Returns property name
       */
      public String getName() {
         if (name == null) {
            name = object.getProperty(TITLE_PROPERTY);
         }

         return name;
      }

   }

   /**
    * Represent PropertyViewPropMultiValue object part of property view
    */
   class MultiPropertyViewObject extends PropertyViewObject {

      MultiPropertyViewObject(DisplayObject object) {
         super(object, false);
      }

      @Override
      public String getValue() {
         if (value == null) {
            value = object.getProperties(VALUES_PROPERTY);
         }
         return value;
      }

      /**
       * Returns property name
       */
      @Override
      public String getName() {
         if (name == null) {
            name = object.getProperty(TITLE_PROPERTY);
         }
         return name;
      }

      /**
       * Returns property values
       */
      public Set<String> getValues() {
         Set<String> values = new HashSet<String>();
         String multiPropValue = getValue();
         if (multiPropValue != null) {
            String[] array = multiPropValue.split(",");
            for (String value : array) {
               values.add(value);
            }
            return values;
         }
         if (value != null) {
            values.add(value);
         }

         return values;
      }
   }

   /**
    * Represent PropertyViewMessage object part of property view
    */
   class PropertyViewMessageObject extends PropertyViewObject {
      private static final String MESSAGE_PROPERTY = "message";

      PropertyViewMessageObject(DisplayObject object) {
         super(object, false);
      }

      @Override
      public String getValue() {
         if (value == null) {
            value = object.getProperty(MESSAGE_PROPERTY);
         }

         return value;
      }

      @Override
      public String getName() {
         return getValue();
      }
   }
}
