/* Copyright 2014 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.controls.mx.custom.Portlet;
import com.vmware.flexui.selenium.BrowserUtil;
import com.vmware.suitaf.SUITA;

/**
 * This class is an abstraction of Portlet component. <br>
 * <br>
 * The class encapsulates basic actions related to portlets - expand, collapse (minimize, maximize when need arises)
 * Example of usage: <br>
 * <code>
 *    PortletControl portlet = new PortletControl();
 *    portlet.expand();
 * </code>
 */
public class PortletControl {

   /** Reference to <code>Portlet</code> element, that comes from the underlying libraries. */
   private Portlet _portlet;

   private static final String PORTLET_COLLAPSED_PROPERTY = "collapsed";


   /**
    * Creates instance of the currently visible Portlet component.
    *
    * @param portletId
    */
   public PortletControl(String portletId) {
      _portlet = new Portlet(portletId, BrowserUtil.flashSelenium);
      _portlet.waitForElementEnable(SUITA.Environment.getUIOperationTimeout());
   }

   /**
    * This method will expand the portlet in case it's collapsed
    * @return true if this method has expanded the portlet
    */
   public boolean expand() {
      if (getPortletCollapseState()) {
         _portlet.expandCollapse(SUITA.Environment.getUIOperationTimeout());
         return true;
      }
      return false;
   }

   /**
    * This method will collapse the portlet in case it's expanded
    * @return true if this method has collapsed the portlet
    */
   public boolean collapse() {
      if (!getPortletCollapseState()) {
         _portlet.expandCollapse(SUITA.Environment.getUIOperationTimeout());
         return true;
      }
      return false;
   }

   /**
    * Gets the collapsed state of the Portlet. If it is expanded should return false.
    *
    * @return state of the portlet - if collapsed should be true, false otherwise
    */
   public boolean getPortletCollapseState() {
      return Boolean.parseBoolean(_portlet.getProperty(PORTLET_COLLAPSED_PROPERTY));
   }
}