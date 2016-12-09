package com.vmware.suitaf.apl.sele;

import java.awt.event.KeyEvent;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import org.openqa.selenium.server.commands.NativeCommand;

import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.flexui.componentframework.DisplayObjectContainer;
import com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid;
import com.vmware.flexui.componentframework.controls.mx.Button;
import com.vmware.flexui.componentframework.controls.mx.ComboBase;
import com.vmware.flexui.componentframework.controls.mx.DataGrid;
import com.vmware.flexui.componentframework.controls.mx.Label;
import com.vmware.flexui.componentframework.controls.mx.NavBar;
import com.vmware.flexui.componentframework.controls.mx.NumericStepper;
import com.vmware.flexui.componentframework.controls.mx.TextInput;
import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.componentframework.controls.mx.custom.PermanentTabBar;
import com.vmware.flexui.componentframework.controls.spark.SparkCheckBox;
import com.vmware.flexui.componentframework.controls.spark.SparkDropDownList;
import com.vmware.flexui.componentframework.controls.spark.SparkRadioButton;
import com.vmware.flexui.selenium.MethodNameConstants;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.Key;
import com.vmware.suitaf.apl.Property;
import com.vmware.suitaf.util.FailureCauseHolder;

public  class SeleHelper {
   public static final Class<?>[] USER_INPUT_CLASSES = {
      com.vmware.flexui.componentframework.controls.mx.CheckBox.class,
      com.vmware.flexui.componentframework.controls.mx.PopUpButton.class,
      com.vmware.flexui.componentframework.controls.mx.RadioButton.class,
      com.vmware.flexui.componentframework.controls.mx.Button.class,
      com.vmware.flexui.componentframework.controls.mx.ComboBox.class,
      com.vmware.flexui.componentframework.controls.mx.DateField.class,
      com.vmware.flexui.componentframework.controls.mx.ComboBase.class,
      com.vmware.flexui.componentframework.controls.mx.AdvancedDataGrid.class,
      com.vmware.flexui.componentframework.controls.mx.DataGrid.class,
      com.vmware.flexui.componentframework.controls.mx.Menu.class,
      com.vmware.flexui.componentframework.controls.mx.List.class,
      com.vmware.flexui.componentframework.controls.mx.NavBar.class,
      com.vmware.flexui.componentframework.controls.mx.NumericStepper.class,
      com.vmware.flexui.componentframework.controls.mx.TextInput.class,
      com.vmware.flexui.componentframework.controls.spark.SparkCheckBox.class,
      com.vmware.flexui.componentframework.controls.spark.SparkDropDownList.class,
      com.vmware.flexui.componentframework.controls.spark.SparkRadioButton.class,
   };

   private enum StatusMapper {

      POWERED_OFF ("Stopped, Powered Off"),
      POWERED_ON ("Running, Powered On"),
      RESOLVED ("Ready, Successful"),
      UNRESOLVED (""),
      ;

      private String expectedStatus;
      StatusMapper(String expectedStatus) {
         this.expectedStatus = expectedStatus;
      }

      public String getExpectedStatus() {
         return this.expectedStatus;
      }

   }

   public static final String PN_NUM_CHILDREN = "numChildren";
   public static final String PN_PARENT = "parent";
   public static final String PN_CLASS_NAME = "className";

   //VI CLient flex commands
   //TODO: add dependency from VI Client project at some stage
   public static final String FLEX_GET_FLEX_PROPERTY_TIWO_DIALOG_RAW_CHILDREN_BUTTON ="getFlexPropertyTiwoDialogRawChildrenButton";
   public static final String FLEX_CLICK_TIWO_DIALOG_RAW_CHILDREN_BUTTON ="clickTiwoDialogRawChildrenButton";


   private final SeleAPLImpl apl;
   public SeleHelper(SeleAPLImpl apl) {
      this.apl = apl;
   }



   // =========================
   // Structure logging methods
   // =========================
   /**
    * Method that dumps the component tree given some initial parameters
    * @param id - the abstract ID of the start component for the dumping
    * @param foTreeDepth - maximum tree depth for dumping
    * @param fullyLoggedClasses - array of proxy-classes that should be
    * dumped with the full set of properties. All other classes are dumped with
    * the {@link Category#BRIEF} set of properties.
    * @param loggedPropertyCategories - array of property categories that
    * comprises the full set of dumped properties. The full set of properties
    * will be dumped only for the classes mentioned in parameter
    * <b>fullyLoggedClasses</b>.
    */
   public void logFlexObject(
         DirectID id,
         int foTreeDepth,
         Class<?>[] fullyLoggedClasses,
         Category[] loggedPropertyCategories) {

      List<DisplayObject> foInstance = null;

      if (id == null) {
         // Force Flash Selenium instance creation if not created
         apl.seleLinks.flashSelenium();

         String rootsId = SeleLinks.PATH_TO_FLEX_APP;
         String rootsParenName = apl.seleHelper.safeGetProperty(
               SeleLinks.PATH_TO_FLEX_APP, false, "parent.name");
         if (rootsParenName != null) {
            rootsId = "parent.name" + "=" + rootsParenName;
         }
         foInstance = apl.seleLinks.findAll(
               DirectID.from(apl, IDGroup.toIDGroup(rootsId)));

         if (foInstance == null) {
            apl.hostLogger.error("Failed to track the root components.");
            return;
         }
      }
      else {
         foInstance = apl.seleLinks.findAll(id, false);

         if (foInstance == null) {
            apl.hostLogger.dump("LOGGING Flex Object at:<" +
                  id.getID() + "> was not found.");
            return;
         }
      }

      // Apply default setting: Fully log just the class of the root component
      if (fullyLoggedClasses == null) {
         fullyLoggedClasses = new Class<?>[0];
         fullyLoggedClasses[0] = foInstance.getClass();
      }

      failingProp4Class.clear();

      for (DisplayObject to : foInstance) {
         logFlexObject(to, foTreeDepth, fullyLoggedClasses,
               loggedPropertyCategories);
      }
   }

   /**
    * Method that dumps the component tree given some initial parameters
    * @param foInstance - the proxy-object of the start component for the
    * dumping
    * @param foTreeDepth - maximum tree depth for dumping
    * @param fullyLoggedClasses - array of proxy-classes that should be
    * dumped with the full set of properties. All other classes are dumped with
    * the {@link Category#BRIEF} set of properties.
    * @param loggedPropertyCategories - array of property categories that
    * comprises the full set of dumped properties. The full set of properties
    * will be dumped only for the classes mentioned in parameter
    * <b>fullyLoggedClasses</b>.
    */
   public void logFlexObject(
         DisplayObject foInstance,
         int foTreeDepth,
         Class<?>[] fullyLoggedClasses,
         Category...loggedPropertyCategories) {

      String initialMessage = "LOGGING Flex Object Tree ( ";

      // Apply default setting: Log all tree levels
      if (foTreeDepth < 0) {
         foTreeDepth = 1000;
      }

      if (foInstance != null) {
         initialMessage += "Root:" +
               safeGetProperty(foInstance.getUniqueId(), true, PN_PARENT) + "; ";
      }
      else {
         apl.hostLogger.error("Component dump FAILED because a starting" +
               " instance could not be obtained");
         return;
      }

      // Apply default setting: Logs only identification and info properties
      if (loggedPropertyCategories == null ||
            loggedPropertyCategories.length == 0) {
         loggedPropertyCategories =
               new Category[] { Category.ID, Category.INFO };
      }

      // Apply default setting: Fully log just the Flex input classes
      if (fullyLoggedClasses == null) {

         fullyLoggedClasses = USER_INPUT_CLASSES;
      }

      apl.hostLogger.dump(
            initialMessage + "TreeDepth:" + foTreeDepth + ";)...");
      logTestObjectTree( foInstance, foTreeDepth,
            fullyLoggedClasses, loggedPropertyCategories, 0);
   }

   /**
    * This is a recursive method that crawls the UI component tree and
    * initiates the dumping of each UI component it reaches.
    * @param foInstance - the proxy-object of the current component for the
    * dumping
    * @param foTreeDepth - maximum tree depth for dumping
    * @param fullyLoggedClasses - array of proxy-classes that should be
    * dumped with the full set of properties. All other classes are dumped with
    * the {@link Category#BRIEF} set of properties.
    * @param loggedPropertyCategories - array of property categories that
    * comprises the full set of dumped properties. The full set of properties
    * will be dumped only for the classes mentioned in parameter
    * <b>fullyLoggedClasses</b>.
    * @param foCurrentTreeDepth - the current tree-depth reached through
    * the recursive calls
    */
   private void logTestObjectTree(
         DisplayObject loggedTO,
         int foTreeDepth,
         Class<?>[] fullyLoggedClasses,
         Category[] loggedPropertyCategories,
         int foCurrentTreeDepth) {

      List<DisplayObject> children = null;
      if (foCurrentTreeDepth + 1 < foTreeDepth) {
         children = safeGetChildren(loggedTO);
      }
      else {
         children = new ArrayList<DisplayObject>();
      }

      // IF the screen component is without "className" property
      // AND there are no child components
      // THEN skip this component
      if (loggedTO.getClass().equals(DisplayObject.class)
            && children.size() == 0) {
         return;
      }

      logSingleTestObject(loggedTO, fullyLoggedClasses,
            loggedPropertyCategories, foCurrentTreeDepth);

      foCurrentTreeDepth++;
      for (DisplayObject child : children) {
         logTestObjectTree( child, foTreeDepth, fullyLoggedClasses,
               loggedPropertyCategories, foCurrentTreeDepth);
      }
   }

   /**
    * This method dumps the information of a single UI component
    * @param loggedTO - the proxy-object representing the logged component
    * @param fullyLoggedClasses - array of proxy-classes that should be
    * dumped with the full set of properties. All other classes are dumped with
    * the {@link Category#BRIEF} set of properties.
    * @param loggedPropertyCategories - array of property categories that
    * comprises the full set of dumped properties. The full set of properties
    * will be dumped only for the classes mentioned in parameter
    * <b>fullyLoggedClasses</b>.
    * @param loggingTreeDepth - the tree-depth of the current UI component
    * relative to the starting component of the dumping.
    */
   private void logSingleTestObject(
         DisplayObject loggedTO,
         Class<?>[] fullyLoggedClasses,
         Category[] loggedPropertyCategories,
         int loggingTreeDepth) {

      if (loggedTO == null) {
         return;
      }
      ArrayList<Category> cat = new ArrayList<Category>();
      if (loggedPropertyCategories != null) {
         cat.addAll(Arrays.asList(loggedPropertyCategories));
         if (cat.contains(Category.ANY)) {
            // Move the ANY category at the last position of the list
            cat.clear();
            cat.addAll(Arrays.asList(Category.values()));
            cat.remove(Category.BRIEF);
            cat.remove(Category.REGULAR);
            cat.remove(Category.SYSTEM);
            cat.remove(Category.ANY);
            cat.add(Category.ANY);
         }
      }

      // =================================
      // Pre-extract come Test Object features
      // =================================
      String indent = "                              "
            .substring(0, (loggingTreeDepth < 15 ? loggingTreeDepth * 2 : 30));
      StringBuilder loTitleText = new StringBuilder();

      loTitleText.append(indent);
      loTitleText.append("[").append(loggingTreeDepth).append("] ");

      // If the screen component is without "className" property
      if (loggedTO.getClass().equals(DisplayObject.class)) {
         loTitleText.append("???");
         loTitleText.append("<").append(loggedTO.getUniqueId()).append(">");

         apl.hostLogger.dump(loTitleText.toString());
         return;
      }

      String xy = getComponentInfo(loggedTO, -1, true);
      String simpleClassName = getSingleProperty(
            loggedTO, true, Property.CLASS_NAME);

      loTitleText.append(simpleClassName);
      loTitleText.append("(").append(xy).append(")");

      loTitleText.append("<").append(loggedTO.getUniqueId()).append("> ");

      // Compute if the property logging is muted
      boolean muted = false;
      // Muted if logged object's class is not of the fully logged classes
      muted = muted ||
            ( !isInstanceOf(loggedTO, fullyLoggedClasses) );
      // Muted if no property categories are given for logging
      muted = muted || ( cat.size() == 0 );

      // Retrieve the list of properties supported by the logged class
      HashSet<Property> propertiesToLog =
            new HashSet<Property>(Property.CAT2PROP.get(Category.REGULAR));

      // =================================
      // Logging of test object title line
      // =================================
      if (muted) {
         // If one of the muted test objects - log just the brief info
         cat = new ArrayList<Category>(Arrays.asList(Category.BRIEF));
      }

      // ======================================
      // Full Logging of test object properties
      // ======================================
      for (Category category : cat) {
         if (propertiesToLog.size() == 0) {
            break;
         }
         String propertyText = propertyListToText(
               loggedTO, simpleClassName, Property.CAT2PROP.get(category),
               propertiesToLog, muted);

         if (propertyText.length() > 0 || loTitleText.length() > 0) {
            String pref = indent + "    ";

            if (loTitleText.length() > 0) {
               pref = loTitleText.toString();
               loTitleText.delete(0, loTitleText.length());
            }

            apl.hostLogger.dump( pref +
                  category.toString().toUpperCase() + "->" + propertyText);
         }
      }
   }

   /**
    * Dynamically populated register of the properties not available for
    * UI components with a given "className" property
    */
   private static final HashMap<String, HashSet<Property>> failingProp4Class =
         new HashMap<String, HashSet<Property>>();

   /**
    * This method extracts a list of properties for a given UI component
    * @param to - the proxy-instance representing the UI component
    * @param simpleClassName - the "className" property of the UI component
    * @param propertyList - the list of properties that needed for extraction
    * @param filerSet - a general list of allowed for extraction properties
    * @param briefed - a flag for clipping the property values
    * @return a {@link String} representation of the list of property-value
    * pairs extracted for the given UI component
    */
   private String propertyListToText(
         DisplayObject to, String simpleClassName,
         List<Property> propertyList, Set<Property> filerSet,
         boolean briefed) {
      HashSet<Property> failingProp = null;
      if (failingProp4Class.containsKey(simpleClassName)) {
         failingProp = failingProp4Class.get(simpleClassName);
      }
      else {
         failingProp = new HashSet<Property>();
         failingProp4Class.put(simpleClassName, failingProp);
      }

      ArrayList<String> pStringList = new ArrayList<String>();

      if (to != null && propertyList != null && filerSet != null) {
         for (Property pName : propertyList) {
            if (filerSet.contains(pName)) {
               filerSet.remove(pName);

               if (failingProp.contains(pName)) {
                  continue;
               }

               String propValue = getSingleProperty(to, true, pName);
               if (propValue == null) {
                  List<String> tmp = getListProperty(to, true, pName);
                  if (tmp != null) {
                     propValue = tmp.toString();
                  }
               }
               if (propValue == null) {
                  List<List<String>> tmp = getGridProperty(to,true,pName);
                  if (tmp != null) {
                     propValue = tmp.toString();
                  }
               }
               if (propValue == null) {
                  failingProp.add(pName);
                  continue;
               }

               // Clip the value if in briefed mode
               if (briefed) {
                  if (propValue.indexOf("\n") > -1) {
                     propValue = propValue.substring(
                           0, propValue.indexOf("\n")-1) + " ...";
                  }
               }

               // Add the property and its value
               pStringList.add(pName + "='" + propValue + "'");
            }
         }

      }

      if (pStringList.size() > 0) {
         return pStringList.toString();
      }
      else {
         return "";
      }
   }

   // ===================================================================
   // Safer property extraction method implementations
   // ===================================================================
   /**
    * This method takes care to safely retrieve the child components
    * of a given UI component
    * @param to - the proxy-instance representing the UI component
    * @return a list of proxy-instances representing the child components.
    * Empty list is returned if no children were found.
    */
   public List<DisplayObject> safeGetChildren(DisplayObject to) {
      ArrayList<DisplayObject> children = new ArrayList<DisplayObject>();
      // null safe behavior
      if (to == null || ProxyFactory.isTempId(to.getUniqueId())) {
         return children;
      }

      String root = to.getUniqueId();
      // Retrieve the potential child components
      List<DisplayObject> childs = apl.seleLinks.findAll(
            DirectID.directid(apl, "parent." + root), false);
      if (childs != null) {
         children.addAll(childs);
      }

      return children;
   }

   /**
    * A tool method implementing "instanceof" logic over an array of classes
    * @param instance - an instance to be checked
    * @param ofClasses - array of classes
    * @return - true if the instance is of any of listed classes
    */
   private boolean isInstanceOf(Object instance, Class<?>[] ofClasses) {
      // Deterministic null behavior
      if (instance == null) {
         return false;
      }
      // Check for all classes one by one
      for ( int i=0; i<ofClasses.length; i++) {
         if (ofClasses[i] != null && ofClasses[i].isInstance(instance)) {
            return true;
         }
      }
      // False if none was the super class
      return false;
   }

   private static final String NO_ELEMENT_ERROR_MSG =
         "element '%s' was not found";
   private static final String NO_PROPERTY_ERROR_MSG =
         "no '%s' property for the element";
   private static final String ERROR_MSG = "Error:";

   /**
    * Safe wrappers over Selenium test object related method
    * @param to
    * @param propName
    * @return
    */
   public String safeGetProperty(
         DisplayObject to, boolean safeGet, String...propNames) {
      return safeGetProperty(to.getUniqueId(), safeGet, propNames);
   }

   String safeGetProperty(
         String uniqueID, boolean safeGet, String...propNames) {
      FailureCauseHolder fch = new FailureCauseHolder();
      String value = null;

      // Try if any of the given property names will return a value
      for (String propName: propNames) {
         try {
            value = seleniumFlexAPICall( MethodNameConstants.FLEX_GET_PROPERTY,
                  uniqueID, propName);

            if (value == null) {
               throw new RuntimeException("Null value retrieved.");
            }
            if (value.contains( String.format(
                  NO_ELEMENT_ERROR_MSG, uniqueID))) {
               throw new RuntimeException(value);
            }
            if (value.contains( String.format(
                  NO_PROPERTY_ERROR_MSG, propName))) {
               throw new RuntimeException(value);
            }

            // Return the first successfully retrieved value
            return value;

         } catch (Throwable e) {
            fch.setCause(e);
            value = null;
         }
      }

      // if no property value could be retrieved - return null
      if (!safeGet) {
         fch.escalateCause();
      }
      else {
         this.apl.hostLogger.debug( "PropertyGet <" + uniqueID + "." +
               Arrays.asList(propNames) + ">" + " failed with: " +
               fch.getCause().getMessage());
      }

      return value;
   }

   // ===================================================================
   // Property extraction method implementations
   // ===================================================================
   /**
    * Method for retrieval of single-valued properties. The method maps the
    * general identifiers from {@link Property} to the appropriate native
    * property.
    * @param to - test object instance
    * @param property - property to be retrieved
    * @return - the value of the property
    */
   public String getSingleProperty(DisplayObject to, Property property) {
      return getSingleProperty(to, false, property);
   }
   String getSingleProperty(
         DisplayObject to, boolean safeGet, Property property) {

      String value;
      // Mapping between the general properties and native property names
      switch (property) {
      case AUTOMATION_NAME:
         return safeGetProperty(to, safeGet, "automationName");
      case ID:
         return safeGetProperty(to, safeGet, "id");
      case ENABLED:
         return safeGetProperty(to, safeGet, "enabled");
      case ERRORTEXT:
         return safeGetProperty(to, safeGet, "errorString");
      case FOCUSED:
         return safeGetProperty(to, safeGet, "focusEnabled");
      case VALUE:
         if (to instanceof Button) {
            return safeGetProperty(to, safeGet, "selected");
         }
         else if (to instanceof PermanentTabBar) {
            int index = getSelectedIndex(to, safeGet);
            if (index > -1) {
               String selectedTabId = ((PermanentTabBar) to).
                     getChildIdAtIndex(index + "");
               return (new DisplayObject(selectedTabId, apl.seleLinks.flashSelenium())).getProperty("label");
            } else {
               return null;
            }

         }
         else if (to instanceof NavBar) {
            int index = getSelectedIndex(to, safeGet);
            if (index > -1 && to instanceof DisplayObjectContainer) {
               return ((DisplayObjectContainer) to).getChildIdAtIndex(
                     index + "");
            }
            else {
               return null;
            }
         }
         else if (to instanceof TextInput) {
            return safeGetProperty(to, safeGet, "text");
         }
         else if (to instanceof NumericStepper) {
            return safeGetProperty(to, safeGet, "value");
         }
         else if (to instanceof AdvancedDataGrid) {
            int index = getSelectedIndex(to, safeGet);
            if (index > -1) {
               return getAdvancedDataGridItem(
                     (AdvancedDataGrid) to, index, safeGet);
            }
            else {
               return null;
            }
         }
         else if (to instanceof Tree) {
            int index = getSelectedIndex(to, safeGet);
            if (index > -1) {
               List<String> valueList =
                     getListProperty(to, safeGet, Property.VALUE_LIST);
               return valueList.get(index);
            }
            else {
               return null;
            }
         }
         else if (to instanceof com.vmware.flexui.componentframework.controls.mx.ListBase) {
            int index = getSelectedIndex(to, safeGet);
            if (index > -1) {
               String prefix = "dataProvider." + index + ".";
               return safeGetProperty(to, safeGet,
                     prefix + "text", prefix + "label");
            }
            else {
               return null;
            }
         }
         else if (to instanceof ComboBase) {
            return safeGetProperty(to, safeGet, "text", "selectedLabel", "selectedItem");
         }
         else if (to instanceof Label) {
            return safeGetProperty(to, safeGet, "text");
         }
         else if (to instanceof SparkRadioButton) {
            return safeGetProperty(to, safeGet, "selected");
         }
         else if (to instanceof SparkCheckBox) {
            return safeGetProperty(to, safeGet, "selected");
         }
         else if (to instanceof SparkDropDownList) {
            String dropDownValue = null;
            try {
               dropDownValue = safeGetProperty(to, safeGet, "labelDisplay.text");
            } catch(RuntimeException e) {
               // try to get textInput not labelDisplay
               dropDownValue = safeGetProperty(to, safeGet, "textInput.text");
            }

            return dropDownValue;
         }
         // get value for editable stepper property
         else if (to instanceof com.vmware.flexui.componentframework.UIComponent) {
            return safeGetProperty(to, safeGet, "value");
         }
         else {
            // attempt to obtain scroll position
            return safeGetProperty(to, safeGet, "scrollPosition");
         }
      case VALUE_INDEX:
         return getSelectedIndex(to, safeGet) + "";
      case POS_X:
         value = getComponentInfo(to, 0, safeGet);
         if (value != null) {
            return value;
         }
         break;
      case POS_Y:
         value = getComponentInfo(to, 1, safeGet);
         if (value != null) {
            return value;
         }
         break;
      case POS_WIDTH:
         value = safeGetProperty(to, safeGet, "width");
         if (value != null) {
            return value;
         }
         break;
      case POS_HEIGHT:
         value = safeGetProperty(to, safeGet, "height");
         if (value != null) {
            return value;
         }
         break;
      case TEXT:
         return safeGetProperty(to, safeGet, "text", "label", "title");
      case CAPTION:
         return safeGetProperty(to, safeGet, "label");
      case TITLE:
         return safeGetProperty(to, safeGet, "title");
      case VISIBLE:
         value = null;
         FailureCauseHolder fch = new FailureCauseHolder();
         try {
            // This special method checks not only the visibility of the
            // component itself, but of all its parents
            value = seleniumFlexAPICall( MethodNameConstants.FLEX_GET_VISIBLE_ON_PATH,
                  to.getUniqueId());
            if (value == null) {
               throw new RuntimeException("Null value retrieved.");
            }
            if (value.contains( String.format(
                  NO_ELEMENT_ERROR_MSG, to.getUniqueId()))) {
               throw new RuntimeException(value);
            }
         } catch (Throwable e) {
            fch.setCause(e);
            value = null;
         }
         if (!safeGet) {
            fch.escalateCause();
         }
         return value;
      case TOOLTIP:
         return safeGetProperty(to, safeGet, "toolTip");
      case CLASS_NAME:
         return safeGetProperty(to, safeGet, PN_CLASS_NAME);
      case CONTEXT:
         return null;
      case DIRECT_ID:
         return to.getUniqueId();
      case X_PATH_ID:
         return null;
      case SCR_X:
         String xres = apl.seleLinks.selenium.getEval(NativeCommand.FLEXAPP_X_Y.encode());
         value = getComponentInfo(to, 0, safeGet);
         if (value != null) {
            int result = Integer.valueOf(value) + Integer.valueOf(NativeCommand.unpackArray(xres)[0]);
            return String.valueOf(result);
         }
         break;
      case SCR_Y:
         String yres = apl.seleLinks.selenium.getEval(NativeCommand.FLEXAPP_X_Y.encode());
         value = getComponentInfo(to, 1, safeGet);
         if (value != null) {
            int result = Integer.valueOf(value) + Integer.valueOf(NativeCommand.unpackArray(yres)[1]);
            return String.valueOf(result);
         }
         break;

      case SCROLL_MAX_POS:
         return safeGetProperty(to, safeGet, "maxScrollPosition");

      case SHOW_PROGRESS:
         return safeGetProperty(to, safeGet, "showProgress");

      case SCROLL_MIN_POS:
         return safeGetProperty(to, safeGet, "minScrollPosition");

      case SCROLL_PAGE_SIZE:
         return safeGetProperty(to, safeGet, "pageSize");

      case SCROLL_POSITION:
         return safeGetProperty(to, safeGet, "scrollPosition");

      case CHILDREN_NUMBER:
         return safeGetProperty(to, safeGet, "numChildren");

      case CURRENT_STATE:
         return safeGetProperty(to, safeGet, "currentState");

      case MINIMIZE_DIRECTION:
         return safeGetProperty(to, safeGet, "minimizeDirection");

      default:
         return null;
      }

      if (!safeGet) {
         throw new RuntimeException("Property: " + property +
               " unsupported for: " + to.getUniqueId());
      }
      else {
         return null;
      }
   }

   /**
    * Checks if Raw Component exists on screen
    * @param componentID - DirectID of the Raw component
    */
   public boolean rawComponentExists(DirectID componentID){
      String rawComponentID = getRawProperty(componentID, Property.ID);
      return  rawComponentID != null && rawComponentID.equals(componentID.getRawComponentId());
   }

   /**
    * Returns value of specific Raw Component property
    * @param componentID - DirectID of the Raw component
    * @param property - the property we want to obtain value for
    */
   public String getRawProperty(DirectID componentID, Property property) {
      String value = null;
      FailureCauseHolder fch = new FailureCauseHolder();

      try {
         String parentId = componentID.getRawParentId();
         String componentId = componentID.getRawComponentId();
         switch (property) {
         case AUTOMATION_NAME:
            return safeGetRawProperty(parentId, componentId, "automationName");
         case ID:
            return safeGetRawProperty(parentId, componentId, "id");
         case ENABLED:
            return safeGetRawProperty(parentId, componentId, "enabled");
         case ERRORTEXT:
            return safeGetRawProperty(parentId, componentId, "errorString");
         case FOCUSED:
            return safeGetRawProperty(parentId, componentId, "focusEnabled");
         case TEXT:
            return safeGetRawProperty(parentId, componentId, "text", "label", "title");
         case CAPTION:
            return safeGetRawProperty(parentId, componentId, "label");
         case TITLE:
            return safeGetRawProperty(parentId, componentId,  "title");
         case VALUE:
            return safeGetRawProperty(parentId, componentId,  "value", "text", "selected");

         default:
            throw new Exception("Property (" + property + ") not implemented for raw components");
         }
      } catch (Exception e){
         fch.setCause(e);
      }

      return value;
   }

   /**
    * Attempts to click on a Raw component
    * @param component - the raw component to click on
    */
   public void clickRawComponent(DirectID component){
      if (component.isRawType() && rawComponentExists(component)){
         seleniumFlexAPICall(FLEX_CLICK_TIWO_DIALOG_RAW_CHILDREN_BUTTON,
               component.getRawParentId(), component.getRawComponentId());
      } else {
         String message = "Cannot execute Raw component click: " + component;
         apl.hostLogger.error(message);
         throw new RuntimeException(message);
      }
   }

   String safeGetRawProperty(
         String parentID, String componentID, String...propNames) {
      FailureCauseHolder fch = new FailureCauseHolder();
      String value = null;

      // Try if any of the given property names will return a value
      for (String propName: propNames) {
         try {
            value = seleniumFlexAPICall(
                  FLEX_GET_FLEX_PROPERTY_TIWO_DIALOG_RAW_CHILDREN_BUTTON,
                  parentID, componentID, propName);

            if (value == null) {
               throw new RuntimeException("Null value retrieved.");
            }
            if (value.contains( String.format(
                  NO_ELEMENT_ERROR_MSG, componentID))) {
               throw new RuntimeException(value);
            }
            if (value.contains( String.format(
                  NO_PROPERTY_ERROR_MSG, propName))) {
               throw new RuntimeException(value);
            }

            // Return the first successfully retrieved value
            return value;

         } catch (Throwable e) {
            fch.setCause(e);
            value = null;
         }
      }
      fch.escalateCause();
      return value;
   }



   private int getSelectedIndex(DisplayObject to, boolean safeGet) {
      String indexText = safeGetProperty(to, safeGet, "selectedIndex");
      if (indexText != null) {
         return Integer.parseInt(indexText);
      }
      return -1;
   }

   private final String FLEX_GET_GLOBAL_POSITION = "getFlexGlobalPosition";
   private final String FLEX_GET_GLOBAL_POSITION_EX ="getFlexGlobalPositionEx";
   private String getComponentInfo(
         DisplayObject to, int infoItem, boolean safeGet) {
      String info = null;
      FailureCauseHolder fch = new FailureCauseHolder();

      int retryCnt = 4;
      while (retryCnt > 0 && info == null) {
         try {
            switch (retryCnt) {
            case 4:
               info = seleniumFlexAPICall(FLEX_GET_GLOBAL_POSITION, to.getUniqueId());
               break;
            case 3:
               info = seleniumFlexAPICall(FLEX_GET_GLOBAL_POSITION_EX, to.getUniqueId());
               break;
               // These two calls are for backward compatibility with
               // older versions of the VMwareSeleniumFlexAPI
            case 2:
               info = seleniumFlexAPICall(FLEX_GET_GLOBAL_POSITION, to.getUniqueId(), null);
               break;
            case 1:
               info = seleniumFlexAPICall(FLEX_GET_GLOBAL_POSITION_EX, to.getUniqueId(),null);
               break;
            }
            fch.setCause(null);
         } catch (Throwable e) {
            fch.setCause(e);
         }
         retryCnt--;
      }

      if (!safeGet) {
         fch.escalateCause();
      }

      if (info == null) {
         return null;
      } else if (infoItem == 0) {
         return info.split(",")[0];
      } else if (infoItem == 1) {
         return info.split(",")[1];
      } else {
         return info;
      }
   }

   private String getAdvancedDataGridItem(
         AdvancedDataGrid to, int row, boolean safeGet) {
      int col = gridColumnsCount(to, safeGet);
      List<String> rowItems = new ArrayList<String>();
      for (int i = 0; i < col; i++ ) {
         rowItems.add( to.getAdvancedDataGridCellValue(row + "", i + "") );
      }
      return Property.Convert.rowItemsToString(rowItems);

   }

   /**
    * Use that string as prefix for setting default column name where no name
    * is set.
    */
   private final String COLUMN_NAME_PREFIX = "#noName_";
   /**
    * Method for retrieval of list-valued properties. The method maps the
    * general identifiers from {@link Property} to the appropriate native
    * property.
    * @param to - test object instance
    * @param property - property to be retrieved
    * @return - the value list of the property
    */
   public List<String> getListProperty(DisplayObject to, Property property) {
      return getListProperty(to, false, property);
   }
   List<String> getListProperty(
         DisplayObject to, boolean safeGet, Property property) {
      FailureCauseHolder fch = new FailureCauseHolder();

      // Mapping between the general properties and native property names
      switch (property) {
      case GRID_COLUMNS:
         if (to instanceof DataGrid) {
            List<String> columnNames = new ArrayList<String>();
            int columnCount = gridColumnsCount((DataGrid) to, safeGet);
            for (int i = 0; i < columnCount; i++) {
               String colName = "null";
               try {
                  if (to instanceof AdvancedDataGrid) {
                     colName = ((AdvancedDataGrid) to)
                           .getAdvancedDataGridHeader(Integer.toString(i));
                  }
                  else {
                     colName = ((DataGrid) to).getDataGridHeader(
                           Integer.toString(i));
                  }
               } catch (Throwable e) {
                  fch.setCause(e);
                  if (!safeGet) {
                     fch.escalateCause();
                  }
               }
               columnNames.add(colName);
            }

            for (int i = 0; i < columnNames.size(); i++) {
               String name = columnNames.get(i);
               String fakeColumnName =  COLUMN_NAME_PREFIX + i;
               if(name.equals("null")) {
                  columnNames.set(i, fakeColumnName);
               }
            }
            return columnNames;
         }
         break;
      case VALUE_LIST:
         List<String> values = new ArrayList<String>();

         if (to instanceof AdvancedDataGrid) {
            for (List<String> row :
               getGridProperty(to, safeGet, Property.GRID_DATA)) {
               values.add(Property.Convert.rowItemsToString(row));
            }
            return values;
         }
         else if (to instanceof Tree) {
            List<String> propSplit = getNextTreeItemPrefix(to, null);
            List<String> valueParts = new ArrayList<String>();
            while (propSplit.size() > 0) {
               // There are two elements in the property prefix
               // for each tree level - so divide by 2 to get the level
               int level = propSplit.size() / 2;
               // Make place for the next tree-leaf value
               while (valueParts.size() >= level) {
                  valueParts.remove(valueParts.size() - 1);
               }
               // Extract and add the new tree-leaf value
               valueParts.add( level - 1,
                     safeGetProperty(to, safeGet,
                           join(propSplit) + PROP_LABEL));
               // Convert the current tree path to string and add to values
               values.add(Property.Convert.treePathToString(valueParts));

               // Get the next tree item's property prefix
               propSplit = getNextTreeItemPrefix(to, propSplit);
            }
            return values;
         }
         else {
            try {
               int length = Integer.parseInt(
                     safeGetProperty(to,safeGet, "dataProvider.length"));
               for (int i = 0; i < length; i++) {
                  String prefix = "dataProvider." + i + ".";
                  String value = safeGetProperty(to, true,
                        prefix + "text", prefix + "label");
                  values.add(value);
               }
               return values;
            } catch (Throwable e1) {
               fch.setCause(e1);
            }
         }
         break;

      case CHILD_IDS:
         ArrayList<String> childIDs = new ArrayList<String>();
         // null safe behavior
         if (to == null) {
            return childIDs;
         }

         String numChildren = safeGetProperty(to, safeGet, PN_NUM_CHILDREN);
         if (numChildren == null || numChildren.equals("0")) {
            return childIDs;
         }

         for (int i = 0; i < Integer.valueOf(numChildren); i++) {
            String newID = null;
            try {
               newID = seleniumFlexAPICall( MethodNameConstants.FLEX_GET_CHILD_ID_AT_INDEX,
                     to.getUniqueId(), i + "" ) ;
            } catch (Throwable e) {
               newID = e.getMessage();
               if (!newID.toLowerCase().startsWith("error: ")) {
                  newID = "Error: " + newID;
               }
            }

            childIDs.add(newID);
         }

         return childIDs;
      default:
         return null;
      }

      if (!safeGet) {
         throw new RuntimeException("Property: " + property +
               " unsupported for: " + to.getUniqueId());
      }
      else {
         return null;
      }
   }

   private final String PROP_COLUMNS = "columns";
   private final String PROP_PROVIDER = "dataProvider";
   private final String PROP_LENGTH = "length";
   private final String PROP_LABEL = "label";
   private final String PROP_CHILDREN = "children";
   private final String PROP_STATUS_NAME = "status.name";
   private final String PROP_PENDING_STATUS = "pendingStatus";

   /**
    * safeGetProperty(to, true, "dataProvider.1.children.0.label")
    * safeGetProperty(to, true, "dataProvider.1.children.length")
    *
    * @param to
    * @param oldPrefix
    * @return
    */
   private List<String> getNextTreeItemPrefix(
         DisplayObject to, List<String> prefixSplit){
      // Generate the first property
      if (prefixSplit == null) {
         prefixSplit = new ArrayList<String>();
         prefixSplit.add(PROP_PROVIDER);
         prefixSplit.add("-1");
      }
      else {
         // Try to add another tree level
         prefixSplit.add(PROP_CHILDREN);
         prefixSplit.add("-1");
      }

      while (prefixSplit.size() > 0) {
         // Get current tree level index
         int currLevelIndex =
               Integer.parseInt(prefixSplit.get(prefixSplit.size() - 1));
         // Get tree level length
         String currLevelLengthProperty =
               join(prefixSplit.subList(0, prefixSplit.size() - 1)) +
               PROP_LENGTH;
         // Get tree level length
         String currLevelLength =
               safeGetProperty(to, true, currLevelLengthProperty);
         // If the current level has more elements - point to the next one
         if (currLevelLength != null
               && Integer.parseInt(currLevelLength) > currLevelIndex + 1) {
            prefixSplit.set( prefixSplit.size() - 1,
                  String.valueOf(currLevelIndex + 1));

            break;
         }

         // if no more elements on the current level - remove it
         prefixSplit.remove(prefixSplit.size() - 1);
         prefixSplit.remove(prefixSplit.size() - 1);
      }

      return prefixSplit;
   }
   private String join(List<String> parts) {
      StringBuilder sb = new StringBuilder();
      for (String part : parts) {
         sb.append(part).append(".");
      }
      return sb.toString();
   }
   private Integer gridColumnsCount(DataGrid to, boolean safeGet) {
      FailureCauseHolder fch = new FailureCauseHolder();
      try {
         return to.getColumnCount();
      } catch (Throwable e) {
         fch.setCause(e);
      }
      String tmp = null;
      try {
         safeGetProperty(to, false, PROP_COLUMNS + "." + PROP_LENGTH);
         return Integer.decode(tmp);
      } catch (Throwable e) {
         fch.setCause(e);
      }
      if (!safeGet) {
         fch.escalateCause();
      }
      return null;
   }
   private Integer gridRowsCount(DataGrid to, boolean safeGet) {
      FailureCauseHolder fch = new FailureCauseHolder();
      try {
         return to.getDatagridRowCount();
      } catch (Throwable e) {
         fch.setCause(e);
      }
      String tmp = null;
      try {
         tmp = safeGetProperty(to, false, PROP_PROVIDER + "." + PROP_LENGTH);
         return Integer.decode(tmp);
      } catch (Throwable e) {
         fch.setCause(e);
      }
      if (!safeGet) {
         fch.escalateCause();
      }
      return null;
   }

   /**
    * Method for retrieval of grid-valued properties. The method maps the
    * general identifiers from {@link Property} to the appropriate native
    * property.
    * @param to - test object instance
    * @param property - property to be retrieved
    * @return - the value grid of the property
    */
   public List<List<String>> getGridProperty(DisplayObject to, Property property) {
      return getGridProperty(to, false, property);
   }
   List<List<String>> getGridProperty(
         DisplayObject to, boolean safeGet, Property property) {

      // Mapping between the general properties and native property names
      switch (property) {
      case GRID_DATA:
         List<List<String>> data = new ArrayList<List<String>>();

         if (to instanceof AdvancedDataGrid) {
            AdvancedDataGrid aGrid = (AdvancedDataGrid) to;

            int rowCount =  Integer.parseInt(safeGetProperty(aGrid, safeGet,
                  PROP_PROVIDER + "." + PROP_LENGTH));
            int colCount = gridColumnsCount(aGrid, safeGet);

            for (int i_row=0; i_row < rowCount; i_row++) {
               data.add(new ArrayList<String>());
            }

            for (int i_col=0; i_col < colCount; i_col++) {
               String[] columnValues = aGrid.getAdvancedDataGridEntireColumnValue(i_col + "");
               for (int i_row=0; i_row < columnValues.length; i_row++) {
                  // TODO: Remove the workaround after the Selenium API issue is fixed.
                  // Temporary fix for the bug of the Selenium API that returns more
                  // raws when the content of a cell contains end of line symbol.
                  if(i_row >= data.size()) {
                     break;
                  }

                  // Ugly workaround for obtaining status value for the vCloud grids.
                  if(columnValues[i_row].contains("object") && columnValues[i_row].toLowerCase().contains("status")) {
                     String statusName = safeGetProperty(aGrid, safeGet, PROP_PROVIDER + "." + i_row + "." + PROP_STATUS_NAME);
                     StatusMapper value = StatusMapper.valueOf(statusName);
                     switch (value) {
                     case UNRESOLVED:
                        statusName = safeGetProperty(aGrid, safeGet, PROP_PROVIDER + "." + i_row + "." + PROP_PENDING_STATUS);
                        break;
                     default:
                        statusName = value.getExpectedStatus();
                        break;
                     }
                     data.get(i_row).add("System Alerts, " + statusName);
                  } else {
                     data.get(i_row).add(columnValues[i_row]);
                  }
               }
            }
            return data;
         }
         else if (to instanceof DataGrid) {
            DataGrid dGrid = ((DataGrid) to);

            Integer rowCount = gridRowsCount(dGrid, safeGet);
            Integer colCount = gridColumnsCount(dGrid, safeGet);
            if (rowCount == null || colCount == null) {
               return null;
            }
            for (int i_row=0; i_row < rowCount; i_row++) {
               data.add(new ArrayList<String>());
               for (int i_col=0; i_col < colCount; i_col++) {
                  String cellData = null;
                  try {
                     cellData =
                           dGrid.getItemInCell("" + i_row, "" + i_col);
                  } catch (Throwable e) {
                  }
                  data.get(i_row).add(cellData);
               }
            }

            return data;
         }

         break;
      default:
         return null;
      }

      if (!safeGet) {
         throw new RuntimeException("Property: " + property +
               " unsupported for: " + to.getUniqueId());
      }
      else {
         return null;
      }
   }

   // ===================================================================
   // Methods for generation of internal key-typing string
   // ===================================================================
   /**
    * Method generates {@link KeyItem} sequence from the general format
    * key-typing sequence.
    * @param keySequence - general format key-typing sequence.
    * @return - {@link KeyItem} sequence
    */
   public static ArrayList<KeyItem> getKeyItems(Object...keySequence) {
      ArrayList<KeyItem> keyList = new ArrayList<KeyItem>();

      for (Object keyItem: keySequence) {
         if (keyItem instanceof String) {
            keyList.addAll(KeyItem.from((String)keyItem));
         }
         else if (keyItem instanceof Key) {
            keyList.add(KeyItem.from((Key)keyItem));
         }
         else {
            throw new IllegalArgumentException(
                  "Type <" + keyItem.getClass().getSimpleName() + ">" +
                        " is not allowed for key-typing entries."
                  );
         }
      }

      return keyList;
   }

   public static final class KeyItem {
      public final Boolean isHoldKey;
      public final Boolean isPasteText;
      public final String value;
      private KeyItem(Boolean isHoldKey, Boolean isPasteText, String value) {
         this.isHoldKey = isHoldKey;
         this.isPasteText = isPasteText;
         this.value = value;
      }
      /**
       * Creates a functional {@link KeyItem}
       * @param value the integer code of the functional key
       * @return key item
       */
      private static KeyItem f(int value) {
         return new KeyItem(false, false, "" + value);
      }
      /**
       * Creates a functional hold-down {@link KeyItem}
       * @param value the integer code of the functional key
       * @return key item
       */
      private static KeyItem h(int value) {
         return new KeyItem(true, false, "" + value);
      }

      public static final String TYPEABLE_LO =
            "a-z0-9 \\-\\/\\.\\,\\;\\=\\[\\]\\\\";
      public static final String TYPEABLE_UP =
            "A-Z";
      public static final String UNTYPEABLE =
            "^" + TYPEABLE_UP + TYPEABLE_LO;
      public static ArrayList<KeyItem> from(String text) {
         ArrayList<KeyItem> res = new ArrayList<KeyItem>();
         if (text.matches(".*[" + UNTYPEABLE + "].*")) {
            res.add(new KeyItem( false, true, text));
            res.add(KeyItem.from(Key.CTRL));
            res.add(new KeyItem(false, false, "" + "V".codePointAt(0)));
         }
         else {
            for (int i=0; i < text.length(); i++) {
               String letter = text.substring(i,i+1);
               if (letter.matches("[" + TYPEABLE_UP + "]")) {
                  res.add(KeyItem.from(Key.SHIFT));
               }

               res.add(new KeyItem(
                     false, false,
                     "" + letter.toUpperCase().codePointAt(0)
                     ));
            }
         }
         return res;
      }
      public static KeyItem from(Key key) {
         switch (key) {
         case ALT: return h(KeyEvent.VK_ALT);
         //        case APPLICATION: return f(-1);
         case BACKSPACE: return f(KeyEvent.VK_BACK_SPACE);
         //        case BREAK: return "" + -1;
         case CANCEL: return f(KeyEvent.VK_CANCEL);
         case CAPSLOCK: return f(KeyEvent.VK_CAPS_LOCK);
         case CLEAR: return f(KeyEvent.VK_CLEAR);
         case CTRL: return h(KeyEvent.VK_CONTROL);
         case DELETE: return f(KeyEvent.VK_DELETE);
         case DOWN: return f(KeyEvent.VK_DOWN);
         case END: return f(KeyEvent.VK_END);
         case ENTER: return f(KeyEvent.VK_ENTER);
         case ESC: return f(KeyEvent.VK_ESCAPE);
         //        case EXECUTE: return f(-1);
         case F1: return f(KeyEvent.VK_F1);
         case F2: return f(KeyEvent.VK_F2);
         case F3: return f(KeyEvent.VK_F3);
         case F4: return f(KeyEvent.VK_F4);
         case F5: return f(KeyEvent.VK_F5);
         case F6: return f(KeyEvent.VK_F6);
         case F7: return f(KeyEvent.VK_F7);
         case F8: return f(KeyEvent.VK_F8);
         case F9: return f(KeyEvent.VK_F9);
         case F10: return f(KeyEvent.VK_F10);
         case F11: return f(KeyEvent.VK_F11);
         case F12: return f(KeyEvent.VK_F12);
         case GREATER_THAN: return f(KeyEvent.VK_GREATER);
         case HELP: return f(KeyEvent.VK_HELP);
         case HOME: return f(KeyEvent.VK_HOME);
         case INSERT: return f(KeyEvent.VK_INSERT);
         //        case LEFT_ALT: return f(-1);
         //        case LEFT_CTRL: return f(-1);
         //        case LEFT_SHIFT: return f(-1);
         //        case LEFT_WIN: return f(-1);
         case LESS_THAN: return f(KeyEvent.VK_LESS);
         //        case MENU: return f(-1);
         //        case MODE: return f(-1);
         case PAUSE: return f(KeyEvent.VK_PAUSE);
         case PGDN: return f(KeyEvent.VK_PAGE_DOWN);
         case PGUP: return f(KeyEvent.VK_PAGE_UP);
         case PRTSCR: return f(KeyEvent.VK_PRINTSCREEN);
         //        case REDO: return f(-1);
         case RIGHT: return f(KeyEvent.VK_RIGHT);
         //        case RIGHT_ALT: return f(-1;
         //        case RIGHT_CTRL: return f(-1;
         //        case RIGHT_OPTION: return f(-1);
         //        case RIGHT_SHIFT: return f(-1);
         //        case RIGHT_WIN: return f(-1);
         case SCROLL_LOCK: return f(KeyEvent.VK_SCROLL_LOCK);
         case SHIFT: return h(KeyEvent.VK_SHIFT);
         case SPACE: return f(KeyEvent.VK_SPACE);
         case TAB: return f(KeyEvent.VK_TAB);
         //        case UNDO: return f(-1);
         case UP: return f(KeyEvent.VK_UP);
         }
         throw new java.lang.IllegalArgumentException("Unmapped Key: " + key);
      }
   }

   /**
    * This method safely invokes remote Selenium Flex API call
    * @param operation - the name of the operation to be invoked
    * @param parameters - the parameters to the operation. For more information look
    *                     at selenium flex api
    * @return - result from the remote operation if any
    */

   public String seleniumFlexAPICall(String operation, String ... parameters){
      String value = null;
      FailureCauseHolder fch = new FailureCauseHolder();

      try {
         Integer numParam = parameters.length;

         if (numParam >= 0){
            if (numParam == 0) {
               value = apl.seleLinks.flashSelenium().call(operation);
            } else {
               value = apl.seleLinks.flashSelenium().call(operation, parameters);
            }
         } else {
            throw new RuntimeException("Selenium flex API call invokation failure.");
         }

         if (value == null) {
            throw new RuntimeException("Null value retrieved.");
         }

         if (value.contains(ERROR_MSG)) {
            throw new RuntimeException(value);
         }

         // Return the first successfully retrieved value
         return value;

      } catch (Throwable e) {
         fch.setCause(e);
         fch.escalateCause();
      }
      return value;
   }

   /**
    * Returns the existing count of a raw component
    * @param componentId the direct component ID
    * @return the existing count
    */
   public int getExistingCount(DirectID componentId) {
      if (rawComponentExists(componentId)) {
         return 1;
      } else {
         return 0;
      }
   }

}
