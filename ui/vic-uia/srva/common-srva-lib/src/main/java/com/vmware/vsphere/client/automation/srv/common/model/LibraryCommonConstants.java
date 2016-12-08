package com.vmware.vsphere.client.automation.srv.common.model;

import com.vmware.client.automation.util.SrvLocalizationUtil;

/**
 * TODO: move to library-uia-lib when old fixture implementation is disabled
 * FixtureImpl.java
 */
public class LibraryCommonConstants {
   private static final String SUBSCRIBED_STR =
         SrvLocalizationUtil.getLocalizedString("library.type.subscribed");
   private static final String LOCAL_STR =
         SrvLocalizationUtil.getLocalizedString("library.type.local");

   /**
    * Defines the different types of libraries
    */
   public enum LibraryType {
      // TODO rename LOCAL as it conflicts with the vim type
      SUBSCRIBED(SUBSCRIBED_STR),
      LOCAL(LOCAL_STR);

      private final String type;

      private LibraryType(String libType) {
         this.type = libType;
      }

      /**
       * Gets the type string.
       *
       * @return  The retrieved type string
       */
      public String getTypeString() {
         return type;
      }

      /**
       * Gets Library type for string.
       *
       * @param strLibType    The string for which the Library type will be retrieved
       * @return              The retrieved Library Type
       */
      public static LibraryType toLibraryType(String strLibType) {
         if (strLibType.equals(SUBSCRIBED_STR)) {
            return SUBSCRIBED;
         }
         if (strLibType.equals(LOCAL_STR)) {
            return LOCAL;
         }
         throw new IllegalArgumentException(strLibType + " is NOT a valid Library Type");
      }
   }

   /**
    * Defines the different types of library templates
    */
   public enum LibraryItemType {
      VAPP_TEMPLATE, VM_TEMPLATE
   }

   /**
    * Defines the different types of Summary views from which one can launch the wizard
    */
   public enum LibraryWizardLaunchView {
      DATACENTER_LAUNCH_VIEW,
      DATACENTER_CONTEXT_MENU_LAUNCH_VIEW,
      DATACENTER_RO_VAPP_ICON_LAUNCH_VIEW,
      DATACENTER_RO_VM_ICON_LAUNCH_VIEW,
      HOST_LAUNCH_VIEW,
      HOST_CONTEXT_MENU_LAUNCH_VIEW,
      HOST_RO_VAPP_ICON_LAUNCH_VIEW,
      HOST_RO_VM_ICON_LAUNCH_VIEW,
      CONTENT_LIBRARY_ITEM_LAUNCH_VIEW,
      CONTENT_LIBRARY_ITEM_TEMPLATE_ACTIONS_LAUNCH_VIEW,
      CONTENT_LIBRARY_ITEM_ICON_LAUNCH_VIEW,
      VAPP_LAUNCH_VIEW,
      VAPP_CONTEXT_MENU_LAUNCH_VIEW,
      VAPP_RO_VAPP_ICON_LAUNCH_VIEW,
      VAPP_RO_VM_ICON_LAUNCH_VIEW,
      CLUSTER_LAUNCH_VIEW,
      CLUSTER_CONTEXT_MENU_LAUNCH_VIEW,
      CLUSTER_RO_VAPP_ICON_LAUNCH_VIEW,
      CLUSTER_RO_VM_ICON_LAUNCH_VIEW,
      VM_FOLDER_LAUNCH_VIEW,
      VM_FOLDER_ICON_LAUNCH_VIEW,
      VM_FOLDER_CONTEXT_MENU_LAUNCH_VIEW,
      @Deprecated VC_LAUNCH_VIEW,
      VC_ICON_LAUNCH_VIEW,
      RESPOOL_LAUNCH_VIEW,
      RESPOOL_CONTEXT_MENU_LAUNCH_VIEW,
      RESPOOL_RO_VAPP_ICON_LAUNCH_VIEW,
      RESPOOL_RO_VM_ICON_LAUNCH_VIEW
   }

   /**
    * Defines the different types of legacy resources
    */
   public enum LegacyResources {
      HOST, CLUSTER, VAPP, RP
   }

   public enum DownloadContentType {
      IMMEDIATE, ON_DEMAND;
   }

}
