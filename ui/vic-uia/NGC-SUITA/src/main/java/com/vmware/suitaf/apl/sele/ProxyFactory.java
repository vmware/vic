package com.vmware.suitaf.apl.sele;

import java.lang.reflect.Constructor;

import com.thoughtworks.selenium.FlashSelenium;
import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid;
import com.vmware.flexui.componentframework.controls.mx.Alert;
import com.vmware.flexui.componentframework.controls.mx.Application;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.CheckBox;
import com.vmware.flexui.componentframework.controls.mx.ComboBase;
import com.vmware.flexui.componentframework.controls.mx.ComboBox;
import com.vmware.flexui.componentframework.controls.mx.DataGrid;
import com.vmware.flexui.componentframework.controls.mx.DateField;
import com.vmware.flexui.componentframework.controls.mx.HBox;
import com.vmware.flexui.componentframework.controls.mx.Image;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.componentframework.controls.mx.List;
import com.vmware.flexui.componentframework.controls.mx.Menu;
import com.vmware.flexui.componentframework.controls.mx.NavBar;
import com.vmware.flexui.componentframework.controls.mx.NumericStepper;
import com.vmware.flexui.componentframework.controls.mx.Panel;
import com.vmware.flexui.componentframework.controls.mx.PopUpButton;
import com.vmware.flexui.componentframework.controls.mx.RadioButton;
import com.vmware.flexui.componentframework.controls.mx.TextInput;
import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.componentframework.controls.mx.VBox;
import com.vmware.flexui.componentframework.controls.mx.custom.PermanentTabBar;
import com.vmware.flexui.componentframework.controls.mx.custom.RichComboBox;
import com.vmware.flexui.componentframework.controls.spark.SparkCheckBox;
import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.componentframework.controls.spark.SparkRadioButton;

/**
 * A static factory class that creates proxy UI-component instances
 * @author dkozhuharov
 */
class ProxyFactory {
   private static final String NAME_PN = "name";
   private static final String[] AUX_ID_PN = new String[] {"uid", "text"};

   /**
    * This methods checks if the given Selenium ID is "temporary"
    * @param id - the given Selenium ID
    * @return <b>true</b> if the ID ends with indexing pattern (e.g. button[3])
    */
   public static final boolean isTempId(String id) {
      return id.matches(".+\\[[0-9]+\\]");
   }

   /**
    * This method checks if the given Selenium ID is unique within the current
    * screen rendering.
    * @param apl - an instance of the APL interface implementation
    * @param id - the given Selenium ID
    * @return <b>true</b> if there is no UI component that matches the
    * temporary ID with index [1] (e.g. okButton is unique if there is no
    * okButton[1] component)
    */
   public static final boolean isUnique(SeleAPLImpl apl, String id) {
      if (id == null) {
         return false;
      }
      if (isTempId(id)) {
         return true;
      }
      else {
         return apl.seleHelper.safeGetProperty(
               id + "[1]", true, NAME_PN) == null;
      }
   }

   /**
    * This is a factory method that creates proxy-instances given a particular
    * Selenium ID. The class of the returned instances is of a type "guessed"
    * from the "className" property of the UI component.
    * @param apl - an instance of the APL interface implementation
    * @param id - the given Selenium ID
    * @return a proxy-instance with a base class of {@link DisplayObject}
    */
   public static DisplayObject getInstance(SeleAPLImpl apl, String id) {
      // Retrieval of the name property is used also as a check for existence
      // of a component with such an ID.
      String idPValue = apl.seleHelper.safeGetProperty(id, false, NAME_PN);
      String uid = NAME_PN + "=" + idPValue;

      // If an equivalent unique ID is found - it is put in the "id" var
      if (isUnique(apl, uid)) {
         id = uid;
      }
      else {
         // If the "name" property value is not unique enough
         // - try alternative properties
         for (String idPName : AUX_ID_PN) {
            idPValue = apl.seleHelper.safeGetProperty(id, true, idPName);
            if (idPValue != null) {
               // Check the potential identifier for uniqueness
               uid = idPName + "=" + idPValue;
               if ( isUnique(apl, uid) ) {
                  id = uid;
                  break;
               }
            }
         }
      }

      Class<? extends DisplayObject> cls = getProxyClass(apl, id);
      try {
         Constructor<? extends DisplayObject> wrapperConstructor =
               cls.getConstructor(String.class, FlashSelenium.class);
         return wrapperConstructor.newInstance(
               id, apl.seleLinks.flashSelenium());
      } catch (Exception e) {
         throw new RuntimeException("Failed creation of proxy class <"
               + cls.getName() + "> for ID:" + id, e);
      }
   }

   /**
    * This method "guesses" the type for a proxy-instance from the "className"
    * property of the UI component.
    * @param apl - an instance of the APL interface implementation
    * @param uid - the given Selenium ID
    * @return a {@link Class} instance determining the type of the
    * proxy-instance
    */
   private static Class<? extends DisplayObject> getProxyClass(
         SeleAPLImpl apl, String uid) {
      String className = apl.seleHelper.safeGetProperty(
            uid, true, SeleHelper.PN_CLASS_NAME);
      if (className == null) {
         return DisplayObject.class;
      }
      else if (className.endsWith("Application")) {
         return Application.class;
      }
      else if (className.endsWith("HBox")) {
         return HBox.class;
      }
      else if (className.endsWith("VBox")) {
         return VBox.class;
      }
      else if (className.endsWith("MultilineCheckBox")) {
         return CheckBox.class;
      }
      else if (className.endsWith("CheckBox")) {
         return SparkCheckBox.class;
      }
      else if (className.endsWith("PopUpButton")) {
         return PopUpButton.class;
      }
      else if (className.endsWith("MultilineRadioButton")) {
         return RadioButton.class;
      }
      else if (className.endsWith("RadioButton")) {
         return SparkRadioButton.class;
      }
      else if (className.endsWith("Button")) {
         return Button.class;
      }
      else if (className.endsWith("RichComboBox")) {
          return RichComboBox.class;
      }
      else if (className.endsWith("ComboBox")) {
         return ComboBox.class;
      }
      else if (className.endsWith("DateField")) {
         return DateField.class;
      }
      else if (className.endsWith("ComboBase")) {
         return ComboBase.class;
      }
      else if (className.endsWith("DropDownList")) {
         return SparkDropDownList.class;
      }
      else if (className.endsWith("AdvancedDataGrid")) {
         return AdvancedDataGrid.class;
      }
      else if (className.endsWith("DataGrid")) {
         return DataGrid.class;
      }
      else if (className.endsWith("Image")) {
         return Image.class;
      }
      else if (matchSuffix(className, "TextInput", "Field", "TextInputEx", "RichEditableText", "TextArea", "TextAreaEx")) {
         return TextInput.class;
      }
      else if (matchSuffix(className, "Label", "LabelEx", "Text")) {
         return Label.class;
      }
      else if (className.endsWith("Menu")) {
         return Menu.class;
      }
      else if (className.endsWith("List")) {
         return List.class;
      }
      else if (matchSuffix(className, "Tree", "LeftNavigationControl")) {
         return Tree.class;
      }
      else if (matchSuffix(className, "TabNavigator")) {
         return PermanentTabBar.class;
      }
      else if (matchSuffix(className, "NavBar", "TabBar", "ButtonBar")) {
         return NavBar.class;
      }
      else if (className.endsWith("NumericStepper")) {
         return NumericStepper.class;
      }
      else if (className.endsWith("Alert")) {
         return Alert.class;
      }
      else if (className.endsWith("Panel")) {
         return Panel.class;
      }
      else {
         apl.hostLogger.debug( String.format(
               "ID '%s' points unknown type '%s'." +
                     " Using generic type: %s",
                     uid, className, UIComponent.class.getSimpleName()
               ));
         return UIComponent.class;
      }
   }

   /**
    * A tool method helping the "guessing" of proxy-instance type
    * @param className - the value of the className property of the
    * UI component
    * @param checkSuffixes - an array of characteristic suffixes to be matched
    * @return <b>true</b> if any of the suffixes was matched
    */
   private static boolean matchSuffix(
         String className, String...checkSuffixes) {
      for (String suffix : checkSuffixes) {
         if (className.endsWith(suffix)) {
            return true;
         }
      }
      return false;
   }
}