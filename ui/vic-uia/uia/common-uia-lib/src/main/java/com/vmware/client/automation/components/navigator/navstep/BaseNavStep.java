/**
 * Copyright 2012 VMware, Inc. All rights reserved. -- VMware Confidential
 */
package com.vmware.client.automation.components.navigator.navstep;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.vmware.client.automation.components.navigator.spec.LocationSpec;
import com.vmware.suitaf.SUITA;
import com.vmware.suitaf.UIAutomationTool;

/**
 * Base implementation of a <code>NavigationStep</code>.
 *
 * All atomic navigation actions should extend this class. It is
 * also strongly advised that all the standard navigation is implemented
 * through the navigator's model.
 *
 * Specific UI parameters such as control ID should be passed over a
 * constructor.
 *
 * Each <code>NavigationStep</code> should be registered with a corresponding
 * <code>Navigator</code> object.
 */
public abstract class BaseNavStep implements NavigationStep {

   protected static final Logger _logger = LoggerFactory.getLogger(BaseNavStep.class);

   protected static final UIAutomationTool UI = SUITA.Factory.UI_AUTOMATION_TOOL;

   private String _nid;

   /**
    * The constructor initializes the object with a corresponding
    * navigation identified.
    *
    * All <code>NavigationStep</code> should have a navigation
    * identified.
    *
    * @param nid
    *    Test navigation identifier.
    */
   public BaseNavStep(String nid) {
      _nid = nid;
   }

   @Override
   public String getNId() {
      return _nid;
   }

   @Override
   public abstract void doNavigate(LocationSpec locationSpec) throws Exception;
}
