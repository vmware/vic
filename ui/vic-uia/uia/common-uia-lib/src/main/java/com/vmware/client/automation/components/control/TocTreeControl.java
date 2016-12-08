/* Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential */
package com.vmware.client.automation.components.control;

import com.vmware.flexui.componentframework.controls.mx.Tree;
import com.vmware.flexui.selenium.BrowserUtil;

/**
 * This class is an abstraction of TOC tree (Table Of Contents tree). <br>
 * <br>
 * More specifically this class has next capabilities at the moment: <br>
 * 1) executes traversal of the UI tree in search of visible TOC trees (this logic is
 * executed in the class constructors); <br>
 * 2) selects TOC item, found by index or name; <br>
 * <br>
 * Example of usage: <br>
 * <code>
 *    NGCTocTree tocTree = new NGCTocTree();
 *    tocTree.selectItemByIndex(1);
 * </code>
 */
public class TocTreeControl {

   /** Most of the TocTree in NGC has ID="tocTree". */
   private static final String DEFAULT_TOC_TREE_ID = "tocTree";

   /** Reference to <code>Tree</code> element, that comes from the underlying libraries. */
   private Tree _tocTree;

   // Five seconds timeout constant
   private static final long TT_FIVE_SECONDS_CONSTANT = 5000;

   /**
    * Creates instance of the currently visible TocTree component. Default tocTree ID is
    * used.
    */
   public TocTreeControl() {
      this(DEFAULT_TOC_TREE_ID);
   }

   /**
    * Creates instance of the currently visible TocTree component.
    *
    * @param tocTreeId
    *           simple UID of the tocTree component;
    */
   public TocTreeControl(String tocTreeId) {
      _tocTree = new Tree(tocTreeId, BrowserUtil.flashSelenium);
      _tocTree.waitForElementEnable(TT_FIVE_SECONDS_CONSTANT);
   }

   /**
    * This method will select an item in the tocTree. The item is determined by its
    * index.
    *
    * @param index
    *           the index of the target Toc item
    */
   public void selectItemByIndex(int index) {
      _tocTree.selectItemViClient(String.valueOf(index));
   }
}