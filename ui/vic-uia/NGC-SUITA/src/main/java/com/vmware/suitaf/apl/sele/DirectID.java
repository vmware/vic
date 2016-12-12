package com.vmware.suitaf.apl.sele;

import java.util.Arrays;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Map;
import java.util.Set;

import com.vmware.flexui.componentframework.DisplayObject;
import com.vmware.suitaf.apl.AutomationPlatformLink;
import com.vmware.suitaf.apl.Category;
import com.vmware.suitaf.apl.IDGroup;
import com.vmware.suitaf.apl.IDPart;
import com.vmware.suitaf.apl.Property;

public class DirectID {
   private final IDGroup idGroup;
   private final SeleAPLImpl apl;
   private final String directID;
   private final String grapchicalID;
   private static final String SEPARATOR_RAW_TYPE = "$";

   private DirectID(SeleAPLImpl apl, IDGroup idGroup) {
      if (idGroup == null) {
         throw new IllegalArgumentException(
               "Parameter IDGroup must not be null.");
      }
      this.apl = apl;
      this.idGroup = idGroup;
      this.grapchicalID = idGroup.getValue(Property.GRAPHICAL_ID);

      String id = null;
      if (this.grapchicalID == null) {
         id = idGroup.getValue(Property.DIRECT_ID);
         if (id == null) {
            id = seleIdFromXPath(idGroup.getValue(Property.X_PATH_ID));
         }
         if (id == null) {
            apl.hostLogger.warn(
                  "No usable " +
                  Property.DIRECT_ID + " in: " + idGroup);
         }
      }
      directID = id;
   }

   /**
    * Factory method for creation of {@link DirectID} instances from a single
    * {@link IDGroup}.
    *
    * @param apl - instance of the main class that implements the
    * {@link AutomationPlatformLink} interface
    * @param idGroup - the id group instance to be wrapped
    *
    * @return the new instance
    */
   public static final DirectID from(
         SeleAPLImpl apl, IDGroup idGroup) {
      return new DirectID(apl, idGroup);
   }

   /**
    * Factory method for creation of {@link XPathID} instances from a direct
    * identifier and an array of {@link IDPart} elements.
    *
    * @param apl - instance of the main class that implements the
    * {@link AutomationPlatformLink} interface
    * @param xpath - a xpath component identifier
    * @param parts - optional {@link IDPart} elements
    *
    * @return the new instance
    */
   public static final DirectID directid(
         SeleAPLImpl apl, String directID, IDPart... parts) {
      if (directID == null) {
         return null;
      }
      return new DirectID(apl, IDGroup.toIDGroup(directID, parts));
   }

   public String getID() {
      if (directID == null) {
         throw new RuntimeException("No DirectID in IDGroup: " + idGroup);
      }
      return directID;
   }

   public boolean hasIndex() {
      return directID.contains("[") && directID.contains("]");
   }

   /**
    * @return - Returns true if the ID is not null.
    */
   public boolean hasID() {
      if (directID == null) {
         return false;
      }
      return true;
   }

   /**
    * @return - Returns true if the ID is raw type ID Raw IDs
    * Raw IDs are those ids, which contain the parent dialog
    * id concatenated with component id
    * Example of Raw Id:
    * tiwoDialog$okButton
    */
   public boolean isRawType() {
      if (directID == null) {
         return false;
      }
      return directID.contains(this.SEPARATOR_RAW_TYPE);
   }

   /**
    * @return - Returns the parent dialog ID of the raw component ID
    */
   public String getRawParentId() {
      if (!isRawType()) {
         throw new RuntimeException(
               "Cannot get parent unless ID is RawType: " + idGroup);
      }
      return directID.split("\\" + this.SEPARATOR_RAW_TYPE)[0];
   }

   /**
    * @return - Returns the raw component ID
    */
   public String getRawComponentId() {
      if (!isRawType()) {
         throw new RuntimeException(
               "Cannot get component unless ID is RawType: " + idGroup);
      }
      return directID.split("\\" + this.SEPARATOR_RAW_TYPE)[1];
   }

   /**
    * This method is used to return the graphicalID parameter as a String.
    * It throws an exception if the item is null. In case the exception
    * is not desired a pre-check if the parameter exists can be created
    * using  hasGrapchicalID() method.
    *
    * @return Returns the graphicalID as a String.
    */
   public String getGraphicalID() {
      if (grapchicalID == null) {
         throw new RuntimeException("No GraphicalID in IDGroup: " + idGroup);
      }
      return grapchicalID;
   }

   /**
    * Return the HTML id as a String.
    * @return
    */
   public String getHtmlID() {
      return idGroup.getValue(Property.H5_ID);
   }

   /**
    * This method can be used to check if the ID has a grapchicalID property.
    *
    * @return Returns false if the ID does not have a graphicalID.
    */
   public boolean hasGraphicalID() {
      if (grapchicalID == null) {
         return false;
      }
      return true;
   }

   /**
    * This method returns the Screen X coordinate of the group. If it is not
    * set it will return 0 in order to state the search area to be the entire
    * screen.
    *
    * @return Returns the X screen coordinate property of the group.
    */
   public int getScreenX() {
      int x = 0;
      try {
         x = Integer.parseInt(idGroup.getValue(Property.SCR_X));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "Screen X coordinage not set! getScreenX() will " + "return 0!");
      }

      return x;
   }

   /**
    * This method returns the Screen Y coordinate of the group. If it is not
    * set it will return 0 in order to state the search area to be the entire
    * screen.
    *
    * @return Returns the Y screen coordinate property of the group.
    */
   public int getScreenY() {
      int y = 0;
      try {
         y = Integer.parseInt(idGroup.getValue(Property.SCR_Y));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "Screen Y coordinage not set! getScreenY() will " + "return 0!");
      }
      return y;
   }

   /**
    * This method returns the Application X coordinate of the group. If it is
    * not
    * set it will return 0 in order to state the search area to be the entire
    * screen.
    *
    * @return Returns the X screen coordinate property of the group.
    */
   public int getPosX() {
      int x = 0;
      try {
         x = Integer.parseInt(idGroup.getValue(Property.POS_X));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "Application X coordinage not set! getPosX() will " +
               "return 0!");
      }
      return x;
   }

   /**
    * This method returns the Application Y coordinate of the group. If it is
    * not
    * set it will return 0 in order to state the search area to be the entire
    * screen.
    *
    * @return Returns the Y screen coordinate property of the group.
    */
   public int getPosY() {
      int y = 0;
      try {
         y = Integer.parseInt(idGroup.getValue(Property.POS_Y));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "Application Y coordinage not set! getPosY() will " +
               "return 0!");
      }
      return y;
   }

   /**
    * This method returns the Height coordinate of the group. If it is not
    * set it will return the monitor Height detected from Sikuli in order to
    * state the search area to be the entire screen.
    *
    * @return Returns the Height coordinate property of the group.
    */
   public int getHeight() {
      int height = apl.sikuliXHelper.getScreenDimensions().height;
      try {
         height = Integer.parseInt(idGroup.getValue(Property.POS_HEIGHT));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "POS_HEIGHT not set! getHeight() will return " +
               "monitor height!");
      }
      return height;
   }

   /**
    * This method returns the Width coordinate of the group. If it is not
    * set it will return the monitor Width detected from Sikuli in order to
    * state the search area to be the entire screen.
    *
    * @return Returns the Width coordinate property of the group.
    */
   public int getWidth() {
      int width = apl.sikuliXHelper.getScreenDimensions().width;
      try {
         width = Integer.parseInt(idGroup.getValue(Property.POS_WIDTH));
      } catch (Exception e) {
         apl.hostLogger.debug(
               "POS_WIDTH not set! getWidth() will return " + "monitor width!");
      }
      return width;
   }

   /**
    * Match the properties of a {@link DisplayObject} instance against the
    * id-parts from the wrapped {@link IDGroup}.
    *
    * @param to
    * @param properties - restricting list of property identifiers to be
    * compared
    *
    * @return <b>true</b> if all the compared properties match
    */
   public boolean matchProperties(DisplayObject to, Property... properties) {
      return matchProperties(
            to,
            new HashSet<Property>(Arrays.asList(properties)),
            new HashSet<Category>());
   }

   /**
    * Match the properties of a to object instance against the
    * id-parts from the wrapped {@link IDGroup}.
    *
    * @param to
    * @param categories - restricting list of property categories whose
    * properties must be compared
    *
    * @return <b>true</b> if all the compared properties match
    */
   public boolean matchProperties(DisplayObject to, Category... categories) {
      return matchProperties(
            to, new HashSet<Property>(), new HashSet<Category>(
            Arrays.asList(
                  categories)));
   }

   private boolean matchProperties(
         DisplayObject to,
         HashSet<Property> p,
         HashSet<Category> c) {
      for (Property propToMatch : idGroup.getProperties()) {
         if (p.contains(propToMatch) || haveIntersection(
               c,
               propToMatch.category)) {
            if (idGroup.matchComponentValue(
                  apl.seleHelper.getSingleProperty(to, true, propToMatch),
                  propToMatch)) {
               continue;
            }
            if (idGroup.matchComponentValueList(
                  apl.seleHelper.getListProperty(to, true, propToMatch),
                  propToMatch)) {
               continue;
            }
            if (idGroup.matchComponentValueGrid(
                  apl.seleHelper.getGridProperty(to, true, propToMatch),
                  propToMatch)) {
               continue;
            }
            return false;
         }
      }

      return true;
   }

   private boolean haveIntersection(Set<?> s1, Set<?> s2) {
      if (s1 == null || s2 == null) {
         return false;
      }

      for (Object element : s1)
         if (s2.contains(element)) {
            return true;
         }

      return false;
   }

   @Override
   public String toString() {
      return idGroup.toString();
   }

   // ====================================================================
   // XPath helper functions
   // ====================================================================

   /**
    * XPath parsing helper function. Splits a XPath on three parts:
    * <li> The rightmost filter section
    * <li> The rightmost component class name
    * <li> The root of the XPath
    * Here are two example XPaths splitted:
    * <br>{@code //TitleWindow[@caption='Open']//}
    * <br>&nbsp;{@code FlexButton}
    * <br>&nbsp;&nbsp;{@code [@id='okBtn'][2]}
    * <br><br>
    *
    * @param xPath - the XPath to be split
    *
    * @return - array with three elements: root / type / filter
    */
   static String[] splitRootTypeFilter(String xPath) {
      String[] rootTypeFilter = new String[]{"", "", ""};

      while (xPath.endsWith("]")) {
         int i1 = xPath.lastIndexOf("[");
         if (i1 >= 0) {
            rootTypeFilter[2] = xPath.substring(i1) + rootTypeFilter[2];
            xPath = xPath.substring(0, i1);
         }
      }

      int i2 = xPath.lastIndexOf("/");
      if (i2 >= 0) {
         rootTypeFilter[1] = xPath.substring(i2 + 1);
         xPath = xPath.substring(0, i2 + 1);
      }

      rootTypeFilter[0] = xPath;

      return rootTypeFilter;
   }

   /**
    * Parses an XPath component identifier and from attribute restrictions
    * section extracts attributeName/AttributeValue pairs. The attribute
    * filtering clause is expected to have the structure:<br>
    * {@code [@attr1 = 'attr1value' and @attr2='attr2value' and ...]}
    *
    * @param xpath - the XPath to parse
    *
    * @return Hash map with the name/value pairs
    */
   public static HashMap<String, String> getXPathAttribs(String filter) {
      HashMap<String, String> attribs = new HashMap<String, String>();
      if (filter != null) {
         String[] filterAttribs = ("-" + filter).split("@");

         for (int i = 1; i < filterAttribs.length; i++) {
            int j = filterAttribs[i].indexOf("=");
            if (j < 0) {
               continue;
            }

            int v1 = filterAttribs[i].indexOf("'", j + 1);
            if (v1 < 0) {
               continue;
            }

            int v2 = filterAttribs[i].indexOf("'", v1 + 1);
            if (v2 < 0) {
               continue;
            }

            String attribName = filterAttribs[i].substring(0, j).trim();
            String attribValue = filterAttribs[i].substring(v1 + 1, v2).trim();

            attribs.put(attribName, attribValue);
         }
      }

      return attribs;
   }

   protected static final String SELE_IDATTRIB_ID = "id";
   protected static final String SELE_IDATTRIB_WINDOWID = "windowid";
   protected static final String SELE_IDATTRIB_AUTOMATIONNAME =
         "automationName";
   protected static final String SELE_IDATTRIB_CAPTION = "caption";
   protected static final String[] FLEX_TYPE_PREF =
         {"Flex", "SparkApplication"};

   /**
    * This is a specialized translation method that generates direct-id format
    * from a xpath format identifier.
    *
    * @param xpath - xpath format identifier
    *
    * @return direct-id format identifier
    */
   public static String seleIdFromXPath(String xpath) {
      String DirectID = "";
      while (xpath != null && xpath.contains("/")) {
         String[] xpathSplit = splitRootTypeFilter(xpath);

         // Check if the identifier's component type is Flex type
         boolean isFlexType = false;
         for (String prefix : FLEX_TYPE_PREF) {
            if (xpathSplit[1].startsWith(prefix)) {
               isFlexType = true;
               break;
            }
         }

         if (isFlexType) {
            HashMap<String, String> attr = getXPathAttribs(xpathSplit[2]);
            String id = seleIdFromAttrib(attr);
            if (id != null) {
               DirectID = id + ((DirectID.length() > 0) ? "/" : "") + DirectID;
            }
         }
         xpath = xpathSplit[0];
         while (xpath.endsWith("/")) {
            xpath = xpath.substring(0, xpath.length() - 1);
         }
      }
      return (DirectID.length() == 0) ? null : DirectID;
   }

   public static String seleIdFromAttrib(Map<String, String> attr) {

      if (attr.containsKey(SELE_IDATTRIB_ID) && isIdUsable(
            attr.get(
                  SELE_IDATTRIB_ID))) {
         return attr.get(SELE_IDATTRIB_ID);
      } else if (attr.containsKey(SELE_IDATTRIB_WINDOWID) && attr.get(
            SELE_IDATTRIB_WINDOWID) != null) {
         return attr.get(SELE_IDATTRIB_WINDOWID);
      } else if (attr.containsKey(SELE_IDATTRIB_AUTOMATIONNAME) && attr.get(
            SELE_IDATTRIB_AUTOMATIONNAME) != null) {
         return SELE_IDATTRIB_AUTOMATIONNAME + "=" + attr.get(
               SELE_IDATTRIB_AUTOMATIONNAME);
      } else if (attr.containsKey(SELE_IDATTRIB_CAPTION) && attr.get(
            SELE_IDATTRIB_CAPTION) != null) {
         return SELE_IDATTRIB_AUTOMATIONNAME + "=" + attr.get(
               SELE_IDATTRIB_CAPTION);
      }
      return null;
   }

   private static boolean isIdUsable(String anValue) {
      if (anValue != null && !anValue.isEmpty() && !anValue.equals("null") &&
          anValue.trim().replaceAll(
                "([\\d\\p{javaUpperCase}\\p{javaLowerCase}_])+", "").length() ==
          0) {
         return true;
      }
      return false;
   }
}
