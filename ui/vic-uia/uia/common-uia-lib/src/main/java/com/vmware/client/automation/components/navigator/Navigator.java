/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.navigator;

import java.util.HashMap;
import java.util.regex.Pattern;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.google.common.base.Strings;
import com.vmware.client.automation.common.view.BaseView;
import com.vmware.client.automation.components.navigator.navstep.BaseNavStep;
import com.vmware.client.automation.components.navigator.navstep.ChildEntityNavStep;
import com.vmware.client.automation.components.navigator.navstep.EntityNavStep;
import com.vmware.client.automation.components.navigator.navstep.HomeViewNavStep;
import com.vmware.client.automation.components.navigator.navstep.NavigationStep;
import com.vmware.client.automation.components.navigator.navstep.TreeEntityNavStep;
import com.vmware.client.automation.components.navigator.navstep.TreeTabNavStep;
import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.flexui.componentframework.controls.mx.custom.InventoryTree;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.SubToolAudit;

/**
 * The class provides support for performing standard navigation
 * in the client application. The latter includes the object
 * navigator, first (primary) and second (secondary) level tabs,
 * and right hand table of content.
 *
 * The navigation sequence is defined as a path of navigation steps.
 *
 * The identifiers of all navigation steps should be defined inside
 * the <code>Navigator</code> or its subclasses.
 */
public class Navigator {

   // Logger utility
   private static final Logger _logger = LoggerFactory.getLogger(Navigator.class);

   public final static String NID_HOME_ROOT = "home.root";

   private static final String LINK_UID_HOME_ROOT = "homeLogo";

   private final HashMap<String, BaseNavStep> _navActions =
         new HashMap<String, BaseNavStep>();

   private static Navigator s_Navigator = null;

   /**
    * Return an instance of the <code>Navigator</code>. Currently
    * a single instance is always provided.
    *
    * @return
    *    Navigator instance.
    */
   public static Navigator getInstance() {
      if (s_Navigator == null) {
         s_Navigator = new Navigator();
      }

      return s_Navigator;
   }

   /**
    * Register the <code>NavigationStep</code> processors for the known
    * navigation steps.
    */
   public Navigator() {
      registerNavStep(
            new HomeViewNavStep(NID_HOME_ROOT, LINK_UID_HOME_ROOT));
   }

   /**
    * Register a <code>NavigationStep</code> processor for a navigation step.
    */
   public void registerNavStep(BaseNavStep navStep) {
      if (navStep == null) {
         throw new IllegalArgumentException(
               "The required navStep parameter is not set.");
      }

      if (!(navStep instanceof BaseNavStep)) {
         throw new IllegalArgumentException(
               "The provided navStep doesn't not inherit BaseNavStep");
      }

      String stepNId = navStep.getNId();

      if (!isRegisteredNavStep(stepNId)) {
         _navActions.put(stepNId, navStep);
      } else {
         _logger.warn(
               "A navStep with the same ID is already registered. This one will be dismissed: "+ stepNId);
      }
   }

   /**
    * Perform a standard client navigation as defined in the
    * <code>LocationSpec</code>.
    *
    * @param locationSpec
    *    Location spec
    *
    * @return
    *    True if the navigation was successful. False - otherwise.
    */
   public boolean navigateTo(LocationSpec locationSpec) {

      // Validate the location spec.
      validateLocationSpec(locationSpec);

      // Get path elements.
      String path = locationSpec.path.get();
      String[] pathElements = getPathElements(path);

      boolean navResult = true;
      boolean isTreeContext = false;

      BaseNavStep navStep;
      for (String pathElement : pathElements) {

         navStep = null;
         _logger.info("Navigate PATH: " + pathElement);

         if (isEntityName(pathElement)) {
            // Build an entity navigation step.
            String entityName = getEntityNameFromPathElement(pathElement);
            if (isTreeContext) {
               navStep = new TreeEntityNavStep(entityName);
            } else {
               navStep = new EntityNavStep(entityName);
            }
         } else if (isChildEntityName(pathElement)) {
            String entityName = getEntityNameFromPathElement(pathElement);
            navStep = buildChildEntityNavStep(entityName);
         } else {
            // Get the existing navigation step.
            navStep = _navActions.get(pathElement);
            if (TreeTabNavStep.class.isInstance(navStep)) {
               isTreeContext = true;
            }
         }

         if (navStep != null) {
            try {
               // Execute the navigation step.
               SUITA.Factory.UI_AUTOMATION_TOOL.logger.info("NAVIGATE: " + navStep.getClass().toString());
               navStep.doNavigate(locationSpec);

               new BaseView().waitForPageToRefresh();
            } catch (Throwable t) {
               SUITA.Factory.UI_AUTOMATION_TOOL.audit.snapshotAppScreen(
                     SubToolAudit.getFPID(), t);
               _logger.error(
                     String.format(
                           "Error performing navigation step: %s\n%s\n%s",
                           pathElement, t.getMessage(), t.getStackTrace()));
               navResult = false;
               break;
            }
         } else {
            _logger.error(
                  String.format(
                        "Cannot find path element: %s. Navigation steps execution is stopped",
                        pathElement));
            navResult = false;
            break;
         }
      }

      return navResult;
   }

   /**
    * Validate the provided location spec is correct and logically complete.
    *
    * @param locationSpec
    *    Location spec.
    */
   private void validateLocationSpec(LocationSpec locationSpec) {
      // Location spec is set.
      if (locationSpec == null) {
         throw new IllegalArgumentException(
               "Required LocationSpec parameter is not passed.");
      }

      // Path is set.
      String path = locationSpec.path.get();

      if (Strings.isNullOrEmpty(path)) {
         throw new IllegalArgumentException(
               "Required location path parameter is not set.");
      }

      String[] pathElements = path.split(LocationSpec.PATH_SEPARATOR);

      for (String pathElement : pathElements) {

         // Path element is valid.
         if (!isValidPathElement(pathElement)) {
            throw new IllegalArgumentException(
                  "Invalid path element found: " + pathElement);
         }

         // Navigation step implementation is registered.
         if (!isEntityName(pathElement) && !isRegisteredNavStep(pathElement)
               && !isChildEntityName(pathElement)) {
            throw new IllegalArgumentException(
                  "The navigation step identifier doesn't have registered implementation: "
                        + pathElement);
         }
      }

      // TODO: Validate the elements sequence is correct in the specified path.
      // This will require defining parent navigation steps for each step. One step
      // might have more than single parent.
   }

   /**
    * Return true if the path element is not empty.
    */
   private boolean isValidPathElement(String pathElement) {
      return !Strings.isNullOrEmpty(pathElement);
   }

   /**
    * Return true if the pathElement starts with the entity identifier as
    * defined in the <code>LocationSpec.ENTITY_IDENTIFIER</code>.
    */
   private boolean isEntityName(String pathElement) {
      return pathElement.startsWith(LocationSpec.ENTITY_IDENTIFIER);
   }

   /**
    * Return true if the pathElement starts with the entity identifier as
    * defined in the <code>LocationSpec.CHILD_ENTITY_IDENTIFIER</code>.
    */
   private boolean isChildEntityName(String pathElement) {
      return pathElement.startsWith(LocationSpec.CHILD_ENTITY_IDENTIFIER);
   }

   /**
    * Return true if there is <code>NavigationStep</code> registered for the
    * specified navigation path element.
    */
   private boolean isRegisteredNavStep(String pathElement) {
      return getNavStepForPathElement(pathElement) != null;
   }

   /**
    * Return a <code>String</code> array of all path elements in the specified
    * path.
    *
    * The separator in the <code>LocationSpec.PATH_SEPARATOR</code> is used.
    */
   private String[] getPathElements(String path) {
      return path.split(LocationSpec.PATH_SEPARATOR);
   }

   /**
    * Return an instance of the registered <code>NavigationStep</code> action
    * for the specified pathElement.
    */
   private NavigationStep getNavStepForPathElement(String pathElement) {
      return _navActions.get(pathElement);
   }

   /**
    * Return the entity name from the specified entity/child entity name pathElement.
    *
    * The entity name annotation specified in
    * <code>LocationSpec.ENTITY_IDENTIFIER/LocationSpec.CHILD_ENTITY_IDENTIFIER</code>
    * is removed.
    */
   private String getEntityNameFromPathElement(String pathElement) {
      if (isEntityName(pathElement)) {
         return pathElement.replaceFirst(LocationSpec.ENTITY_IDENTIFIER, "");
      } else if (isChildEntityName(pathElement)) {
         return pathElement.replaceFirst(LocationSpec.CHILD_ENTITY_IDENTIFIER, "");
      } else {
         return "";
      }
   }

   /** Constructs the child entity's NavStep object based on the content of the entity name <br>
    * If entity name has an embedded gridId, use it to uniquely identify the grid's ID <br>
    * Else, just use the entity name
    *
    * @param entityNameToParse child entity's name, which could have the grid ID embedded in this format
    * parentOfGridObjectID>gridObjectID$EntityName
    * @return ChildEntityNavStep
    */
   private ChildEntityNavStep buildChildEntityNavStep(String entityNameToParse) {
      String escapedPattern = Pattern.quote(LocationSpec.CHILD_ENTITY_GRID_ID_DELIMITER);
      String[] splitStrArray = entityNameToParse.split(escapedPattern);
      if (splitStrArray.length > 1) {
         String gridId = splitStrArray[0].replace(LocationSpec.CHILD_ENTITY_GRID_ID_PATH_SEPARATOR, LocationSpec.PATH_SEPARATOR);
         String entityName = splitStrArray[1];
         return new ChildEntityNavStep(gridId, entityName);
      } else {
         return new ChildEntityNavStep(entityNameToParse);
      }
   }
}