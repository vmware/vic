/* Copyright 2015 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import java.awt.Point;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.flexui.componentframework.UIComponent;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * Represents the vertical scroll bar control. Provides ability to scroll up and down by moving by page.
 */
public class VerticalScrollBar {

   private static final String VERTICAL_POSITION_PROPERTY_NAME = "y";
   private static final String HEIGHT_PROPERTY_NAME = "height";
   private static final String WIDTH_PROPERTY_NAME = "width";
   private static final String SCROLL_THUMB_ID = "/className=ScrollThumb";

   private static final Logger _logger = LoggerFactory.getLogger(VerticalScrollBar.class);

   private String _id;
   private String _scrollThumbId;
   private String _pageNavigationButtonId;
   private UIComponent _scrollBar = null;

   /**
    * Create vertical scroll bar control.
    * @param id
    */
   public VerticalScrollBar(String id) {
      this._id = id;
      this._scrollThumbId = _id + SCROLL_THUMB_ID;
      initScrollerIds();
   }

   /**
    * Move down by a page.
    * To move down the automation click on the scroll bar above the down button.
    * @return true if the scroll is moved down.
    */
   public boolean moveDown() {
      // get the scroller thumb y position
      UIComponent scrollThumb = new UIComponent(_scrollThumbId, BrowserUtil.flashSelenium);
      int scrollThumbYPosition = Integer.parseInt(scrollThumb.getProperties(VERTICAL_POSITION_PROPERTY_NAME));
      // get the scroll bar height
      UIComponent scrollBar = new UIComponent(_pageNavigationButtonId, BrowserUtil.flashSelenium);
      int scrollBarHeight = Integer.parseInt(scrollBar.getProperties(HEIGHT_PROPERTY_NAME));
      int scrollBarWidth = Integer.parseInt(scrollBar.getProperties(WIDTH_PROPERTY_NAME));
      scrollBar.leftMouseClick(false, false, false, new Point(scrollBarWidth/2, scrollBarHeight));
      int newYPosition = Integer.parseInt(scrollThumb.getProperties(VERTICAL_POSITION_PROPERTY_NAME));
      return newYPosition > scrollThumbYPosition;
   }

   /**
    * Move up by a page.
    * To move up the automation click on the scroll bar under the top button.
    * @return true if the scroll might be moved and the scroll button is clicked.
    */
   public boolean moveUp() {
      UIComponent scrollThumb = new UIComponent(_scrollThumbId, BrowserUtil.flashSelenium);
      int scrollThumbYPosition = Integer.parseInt(scrollThumb.getProperties(VERTICAL_POSITION_PROPERTY_NAME));
      UIComponent scrollBar = new UIComponent(_pageNavigationButtonId, BrowserUtil.flashSelenium);
      int scrollBarWidth = Integer.parseInt(scrollBar.getProperties(WIDTH_PROPERTY_NAME));
      scrollBar.leftMouseClick(false, false, false, new Point(scrollBarWidth/2, 0));
      int newYPosition = Integer.parseInt(scrollThumb.getProperties(VERTICAL_POSITION_PROPERTY_NAME));
      return newYPosition < scrollThumbYPosition;
   }

   /**
    * Scroll to the top of the scroll bar moving page by page.
    * @return true when got to the top of the scroll bar.
    */
   public boolean goToTop() {
      int counter = 0;
      while(moveUp()) {
         if(counter++ > 50) {
            _logger.warn("The scroller is not able to get to the top of or is longer than 50 pages.");
            return false;
         }
      }
      return true;
   }

   /**
    * Return true if the scroll bar component is visible on the screen.
    * @return
    */
   public boolean isVisible() {
      return _scrollBar.isVisibleOnPath();
   }

   /**
    * To click up/down in the scroll bar and move page by page the automation clicks not on the scroll bar up and down
    * buttons but on the scroll bar itself.
    * To move up it clicks a position under the up button and to move down clicks upper the down button. The scroll bar
    * structure consists of 4 buttons - up, down scroll bar thumb and the scroll bar pane button from top to down.
    * The method iterates the buttons and get the longest as this is the scroll bar pane button.
    */
   private void initScrollerIds() {
      _scrollBar = new UIComponent(_id, BrowserUtil.flashSelenium);
      // heights of all buttons in the scroll bar
      String heights = _scrollBar.getProperties("[].height");
      String[] scrollHegths = heights.split(",");
      // names of all the buttons in the scroll bar
      String names = _scrollBar.getProperties("[].name");
      String[] scrollerButtons = names.split(",");
      // find out the longest - this is the button to click on to move by one page(up or down)
      int tmpMaxHeigth = Integer.parseInt(scrollHegths[0]);
      int index = 0;
      for (int i = 1; i < scrollHegths.length; i++) {
         if(tmpMaxHeigth < Integer.parseInt(scrollHegths[i])) {
            tmpMaxHeigth = Integer.parseInt(scrollHegths[i]);
            index = i;
         }
      }
      _pageNavigationButtonId = String.format("name=%s", scrollerButtons[index]);
   }
}
